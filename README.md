<!-- template:define:options
{
  "nodescription": true
}
-->
![logo](https://liam.sh/-/gh/svg/lrstanley/recoverer?icon=ic%3Abaseline-error&icon.height=80&layout=left&icon.color=rgba%28222%2C+63%2C+65%2C+1%29&bgcolor=rgba%2830%2C+0%2C+0%2C+1%29)

<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<p align="center">
  <a href="https://github.com/lrstanley/recoverer/tags">
    <img title="Latest Semver Tag" src="https://img.shields.io/github/v/tag/lrstanley/recoverer?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/recoverer/commits/master">
    <img title="Last commit" src="https://img.shields.io/github/last-commit/lrstanley/recoverer?style=flat-square">
  </a>


  <a href="https://github.com/lrstanley/recoverer/actions?query=workflow%3Atest+event%3Apush">
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/workflow/status/lrstanley/recoverer/test/master?label=test&style=flat-square&event=push">
  </a>

  <a href="https://codecov.io/gh/lrstanley/recoverer">
    <img title="Code Coverage" src="https://img.shields.io/codecov/c/github/lrstanley/recoverer/master?style=flat-square">
  </a>

  <a href="https://pkg.go.dev/github.com/lrstanley/recoverer">
    <img title="Go Documentation" src="https://pkg.go.dev/badge/github.com/lrstanley/recoverer?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/recoverer">
    <img title="Go Report Card" src="https://goreportcard.com/badge/github.com/lrstanley/recoverer?style=flat-square">
  </a>
</p>
<p align="center">
  <a href="https://github.com/lrstanley/recoverer/issues?q=is:open+is:issue+label:bug">
    <img title="Bug reports" src="https://img.shields.io/github/issues/lrstanley/recoverer/bug?label=issues&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/recoverer/issues?q=is:open+is:issue+label:enhancement">
    <img title="Feature requests" src="https://img.shields.io/github/issues/lrstanley/recoverer/enhancement?label=feature%20requests&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/recoverer/pulls">
    <img title="Open Pull Requests" src="https://img.shields.io/github/issues-pr/lrstanley/recoverer?label=prs&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/recoverer/discussions/new?category=q-a">
    <img title="Ask a Question" src="https://img.shields.io/badge/support-ask_a_question!-blue?style=flat-square">
  </a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/discord-bytecord-blue.svg?style=flat-square" title="Discord Chat"></a>
</p>
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :link: Table of Contents

  - [Why](#grey_question-why)
  - [Examples](#examples)
    - [Using net/http's default ServeMux](#using-nethttps-default-servemux)
    - [Using go-chi](#using-go-chi)
  - [Screenshot](#screenshot)
  - [Support &amp; Assistance](#raising_hand_man-support--assistance)
  - [Contributing](#handshake-contributing)
  - [License](#balance_scale-license)
<!-- template:end:toc -->

## :grey_question: Why

recoverer is a simple Go http middleware to catch (and optionally display when
debugging) panics, and attempt to gracefully recover them. recoverer also has
the ability to display such errors (and exported expvar variables) via a clean
and simple html generated error page (shown below).

## Examples

<!-- template:begin:goget -->
<!-- do not edit anything in this "template" block, its auto-generated -->
```console
go get -u github.com/lrstanley/recoverer@latest
```
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

	w.Write([]byte("Hello World!
"))
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
		w.Write([]byte("Hello World!
"))

		panic("uhoh.. things happened.")
	})

    log.Fatal(http.ListenAndServe(":8080", r))
}
```

## Screenshot

![Example screenshot](https://i.imgur.com/TF0Y7gV.png)

<!-- template:begin:support -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :raising_hand_man: Support & Assistance

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for
     guidelines on ensuring everyone has the best experience interacting with
     the community.
* :raising_hand_man: Take a look at the [support](.github/SUPPORT.md) document on
     guidelines for tips on how to ask the right questions.
* :lady_beetle: For all features/bugs/issues/questions/etc, [head over here](https://github.com/lrstanley/recoverer/issues/new/choose).
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :handshake: Contributing

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for guidelines
     on ensuring everyone has the best experience interacting with the
    community.
* :clipboard: Please review the [contributing](.github/CONTRIBUTING.md) doc for submitting
     issues/a guide on submitting pull requests and helping out.
* :old_key: For anything security related, please review this repositories [security policy](https://github.com/lrstanley/recoverer/security/policy).
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :balance_scale: License

```
MIT License

Copyright (c) 2017 Liam Stanley <me@liamstanley.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

_Also located [here](LICENSE)_
<!-- template:end:license -->
