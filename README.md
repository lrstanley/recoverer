# recoverer [![status](https://travis-ci.org/lrstanley/recoverer.svg?branch=master)](https://travis-ci.org/lrstanley/recoverer) [![godoc](https://godoc.org/github.com/lrstanley/recoverer?status.png)](https://godoc.org/github.com/lrstanley/recoverer) [![goreport](https://goreportcard.com/badge/github.com/lrstanley/recoverer)](https://goreportcard.com/report/github.com/lrstanley/recoverer)

recoverer is a simple Go http middleware to catch (and optionally display when
debugging) panics, and attempt to gracefully recover them. recoverer also has
the ability to display such errors (and exported expvar variables) via a clean
and simple html generated error page (shown below).

## Example

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
