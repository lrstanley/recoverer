// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package recoverer provides an http middleware to catch and log panics,
// and optionally display a text (or html) page with the details (useful
// when one is debugging or has a debug flag enabled). recoverer will also
// show exported expvar variables in html mode.
package recoverer

import (
	"expvar"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

// LoggerWriter can be used to convert a *log.Logger to an io.Writer.
type LoggerWriter struct{ *log.Logger }

func (w LoggerWriter) Write(b []byte) (int, error) {
	w.Printf("%s", b)
	return len(b), nil
}

// Options is the configuration which you can pass to the recoverer, to allow
// showing/hiding stack, etc.
type Options struct {
	// Logger is an optional io.Writer which the panic error and stack trace
	// are written to.
	Logger io.Writer

	// Show renders the panic and stack trace from the handler.
	Show bool

	// Simple renders the panic and stack trace in plain text, rather than
	// the default of HTML.
	Simple bool

	// Fn is a function which is called before the panic is written back to
	// the connection (useful if you want, for example, to increment the
	// total number of exceptions for a service, or if you want to only show
	// the error to a specific set of IP's). If the returned error is not nil,
	// the recoverer will not show the error back to the end user, and the
	// error will be logged. recoverer WILL NOT prevent this function from
	// panicing.
	Fn func(req *http.Request, err interface{}, file string, line int) error
}

type recoverer struct {
	options Options
	next    http.Handler

	Stack []byte
	File  string
	Line  int
	Err   interface{}

	ExpVars map[string]expvar.Var
}

func (rec *recoverer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec.Err = recover(); rec.Err != nil {
			rec.Stack = debug.Stack()
			if rec.options.Logger != nil {
				fmt.Fprintf(rec.options.Logger, "panic: %+v\n%s", rec.Err, rec.Stack)
			}

			_, rec.File, rec.Line, _ = runtime.Caller(3)

			if rec.options.Fn != nil {
				err := rec.options.Fn(r, rec.Err, rec.File, rec.Line)
				if err != nil {
					if rec.options.Logger != nil {
						fmt.Fprintf(rec.options.Logger, "panic: recover function error occurred: %+v\n", err)
					}

					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			}

			if !rec.options.Show {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			expvar.Do(func(kv expvar.KeyValue) {
				rec.ExpVars[kv.Key] = kv.Value
			})

			w.WriteHeader(http.StatusInternalServerError)

			// Check if they accept text/html, otherwise provide basic output
			// using the simple format.
			if accept := r.Header.Get("Accept"); accept == "" || !strings.Contains(accept, "text/html") {
				rec.simple(w)
				return
			}

			if rec.options.Simple {
				rec.simple(w)
				return
			}

			rec.html(w)
		}
	}()

	rec.next.ServeHTTP(w, r)
}

// New creates a new recoverer handler with specific options. See
// DefaultRecoverer() for a sane set of default options.
func New(options Options) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &recoverer{
			options: options, next: next,
			ExpVars: make(map[string]expvar.Var),
		}
	}
}

// DefaultRecoverer provides sane defaults to catch and log the panic to
// stderr, and throwing a generic 500 error back to the http client. See
// New() to specify custom options.
func DefaultRecoverer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &recoverer{
			options: Options{Logger: os.Stderr, Show: false},
			next:    next,
			ExpVars: make(map[string]expvar.Var),
		}
	}
}

func (rec *recoverer) simple(w io.Writer) {
	fmt.Fprintf(w, "panic: %+v\n", rec.Err)
	fmt.Fprintf(w, "in: %s:%d\n\n", rec.File, rec.Line)
	fmt.Fprint(w, "stack at time of panic:\n")
	w.Write(rec.Stack)
}

func (rec *recoverer) html(w io.Writer) {
	tmpl.Execute(w, rec)
}

var tmpl = template.Must(template.New("main").Parse(`<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<meta name="description" content="Internal Exception Occurred">
		<title>panic: {{ html (printf "%+v" .Err) }}</title>
	</head>

	<body>
		<div class="header">
			<span>
				<img src="https://blog.golang.org/gopher/gopher.png">
				<h2>500 -- AN EXCEPTION OCCURRED</h2>
			</span>
		</div>

		<div class="main">
			<h2>panic: <span class="panic-text">{{ html (printf "%+v" .Err) }}</span></h2>
			<hr>
			<p>file: <code class="inline">{{ html .File }}</code> (line {{ .Line }})</p>

			<pre class="panic"><code>{{ html (printf "%s" .Stack) }}</code></pre>

			<h2 style="margin-top: 30px">exported variables (with expvars)</h2>
			<hr>
			<table>
				<tr>
					<th>Key</th>
					<th>Value</th>
				</tr>

				{{ range $key, $val := .ExpVars }}
				<tr>
					<td>{{ $key }}</td>
					<td class="skip-padding"><pre><code>{{ $val }}</code></pre></td>
				</tr>
				{{ end }}
			</table>
		</div>

		<style type="text/css">
		* { font-family: "Helvetica Neue", Helvetica, Arial, sans-serif; }

		html, body {
			margin: 0;
			padding: 0;
			width: 100%;
			height: 100%;
		}

		h1, h2, h3, h4, h5, h6 {
			margin-top: 5px;
			margin-bottom: .5rem;
			font-weight: 300;
		}

		h1, h2, h3, h4, h5, h6 {
		    margin-bottom: 0.5rem;
		    font-family: inherit;
		    font-weight: 500;
		    line-height: 1.1;
		    color: inherit;
		}

		h1 { font-size: 2.5rem; }
		h2 { font-size: 2rem; }
		h3 { font-size: 1.75rem; }
		h4 { font-size: 1.5rem; }
		h5 { font-size: 1.25rem; }
		h6 { font-size: 1rem; }

		hr {
			box-sizing: content-box;
			height: 0;
			overflow: visible;
		}

		.header {
			background-color: #EE605E;
			color: white;
			font-size: 35px;
		}

		.header > span {
			display: flex;
			padding: 10px;
		}
		.header > span > img {
			margin-top: 5px;
			height: 55px;
			display: inline-block;
		}
		.header > span > h2 {
			display: inline-block;
			padding: 15px 30px;
			margin: 0;
		}

		.main { padding: 20px;  }
		.panic-text { color: #AAAAAA; }

		code.inline {
			background-color: #e6afaf;
			color: #383838;
			padding: 4px;
			border-radius: 3px;
		}

		pre.panic {
			background-color: #383838;
			color: white;
			padding: 15px;
			border-radius: 6px;
		}

		pre.panic > code {
			font-family: "Courier New", Courier, monospace;
			white-space: pre-wrap;
			word-wrap: break-word;
		}

		table {
			border-collapse: collapse;
			table-layout: fixed;
			width: 100%;
		}

		td:not(.skip-padding), th:not(.skip-padding) {
			border: 1px solid #DDDDDD;
			padding: 8px;
		}

		th, td { width: 100px; text-align: center; }
		td+td, th+th { width: auto; text-align: left; }

		td+td {
			max-height: 200px;
			overflow-y: auto;
		}

		tr:nth-child(even) {
			background-color: #DDDDDD;
		}

		td > pre {
			background-color: #383838;
			color: white;
			padding: 10px;
			margin: 0;
			text-align: left;
			max-height: 150px;
			overflow-x: auto;
		}

		td > pre > code {
			font-family: "Courier New", Courier, monospace;
			white-space: pre-wrap;
			word-wrap: break-word;
		}

		::-webkit-scrollbar {
			width: 10px;
			height: 6px;
		}

		::-webkit-scrollbar-track-piece {
			background-color: #F5F5F5;
			background-clip: padding-box;
		}

		::-webkit-scrollbar-thumb {
			background-color: #EE605E;
			background-clip: padding-box;
			border: 2px solid #FFFFFF;
			border-radius: 6px;
		}

		::-webkit-scrollbar-thumb:window-inactive {
			background-color: #BBBBBB;
		}
		</style>
	</body>
</html>`))
