<!-- template:begin:header -->
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- template:end:toc -->

## :grey_question: Why

recoverer is a simple Go http middleware to catch (and optionally display when
debugging) panics, and attempt to gracefully recover them. recoverer also has
the ability to display such errors (and exported expvar variables) via a clean
and simple html generated error page (shown below).

## Examples

<!-- template:begin:goget -->
<!-- template:end:goget -->

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

<!-- template:begin:support -->
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- template:end:license -->
