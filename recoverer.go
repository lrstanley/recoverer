// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package recoverer provides an http middleware to catch and log panics,
// and optionally
package recoverer

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
)

// TODO:
//  - context keys?
//  - change based on requested headers (e.g. plaintext)?

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
}

type recoverer struct {
	options Options
	next    http.Handler

	Stack []byte
	File  string
	Line  int
	Err   interface{}
}

func (rec *recoverer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec.Err = recover(); rec.Err != nil {
			rec.Stack = debug.Stack()
			if rec.options.Logger != nil {
				fmt.Fprintf(rec.options.Logger, "panic: %+v\n%s", rec.Err, rec.Stack)
			}

			if !rec.options.Show {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			_, rec.File, rec.Line, _ = runtime.Caller(3)

			w.WriteHeader(http.StatusInternalServerError)

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
		return &recoverer{options: options, next: next}
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
		}
	}
}

func (rec *recoverer) simple(w io.Writer) {
	fmt.Fprintf(w, "panic: %+v:\n", rec.Err)
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
		<title>panic: {{printf "%+v" .Err}}</title>
	</head>

	<body>
		<div class="header">
			<span>
				<img src="https://blog.golang.org/gopher/gopher.png">
				<h2>500 -- AN EXCEPTION OCCURRED</h2>
			</span>
		</div>

		<div class="main">
			<h1>panic: <span class="panic-text">{{printf "%+v" .Err}}</span></h1>
			<hr>
			<p>file: <code class="inline">{{.File}}</code> (line {{.Line}})</p>

			<pre class="panic"><code>{{printf "%s" .Stack}}</code></pre>
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
			margin-top: 0;
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
			font-size: 40px;
		}

		.header > span {
			display: flex;
			padding: 20px;
		}
		.header > span > img {
			height: 65px;
			display: inline-block;
		}
		.header > span > h2 {
			display: inline-block;
			padding: 20px 30px;
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
		</style>
	</body>
</html>`))
