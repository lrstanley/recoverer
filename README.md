# recoverer

<p align="center">recoverer -- simple Go http middleware to catch panics</p>
<p align="center">
  <a href="https://pkg.go.dev/github.com/lrstanley/recoverer"><img src="https://pkg.go.dev/badge/github.com/lrstanley/recoverer" alt="pkg.go.dev"></a>
  <a href="https://github.com/lrstanley/recoverer/actions"><img src="https://github.com/lrstanley/recoverer/workflows/test/badge.svg" alt="test status"></a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/recoverer"><img src="https://goreportcard.com/badge/github.com/lrstanley/recoverer" alt="goreportcard"></a>
  <a href="https://gocover.io/github.com/lrstanley/recoverer"><img src="http://gocover.io/_badge/github.com/lrstanley/recoverer" alt="gocover"></a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/Community-Chat%20with%20us-green.svg" alt="Community Chat"></a>
</p>

recoverer is a simple Go http middleware to catch (and optionally display when
debugging) panics, and attempt to gracefully recover them. recoverer also has
the ability to display such errors (and exported expvar variables) via a clean
and simple html generated error page (shown below).

## Examples

### Using net/http's default ServeMux

```go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/lrstanley/recoverer"
)

func hello(w http.ResponseWriter, r *http.Request) {
	panic("uhoh.. things happened.")

	w.Write([]byte("Hello World!\n"))
}

func main() {
	rec := recoverer.New(recoverer.Options{
		Logger: os.Stderr, Show: true, Simple: false,
	})

	http.Handle("/", rec(http.HandlerFunc(hello)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Using go-chi

```go
package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lrstanley/recoverer"
)

func main() {
	r := chi.NewRouter()

	r.Use(recoverer.New(recoverer.Options{Logger: os.Stderr, Show: true, Simple: false}))
	r.Use(middleware.Logger)
	r.Use(middleware.DefaultCompress)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!\n"))

		panic("uhoh.. things happened.")
	})

    log.Fatal(http.ListenAndServe(":8080", r))
}
```

## Screenshot

![Example screenshot](https://i.imgur.com/TF0Y7gV.png)

## Contributing

Please review the [CONTRIBUTING](https://github.com/lrstanley/recoverer/blob/master/CONTRIBUTING.md)
doc for submitting issues/a guide on submitting pull requests and helping out.

## License

    LICENSE: The MIT License (MIT)
    Copyright (c) 2017 Liam Stanley <me@liamstanley.io>

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
