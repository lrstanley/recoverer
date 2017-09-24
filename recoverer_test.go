// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package recoverer

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const panicError = "this should never actually panic"

func getTestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic(panicError)
	}
}

func TestNew(t *testing.T) {
	at := assert.New(t)

	tests := []struct {
		description   string
		url           string
		acceptHeaders string
		options       Options

		wantBodyGlob string
		wantTypeGlob string
	}{
		{
			description: "html show", url: "/",
			acceptHeaders: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			options:       Options{Show: true, Simple: false},
			wantBodyGlob:  "(?s)<html .*</html>",
		},
		{
			description: "html show (url 2)", url: "/another/url",
			acceptHeaders: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			options:       Options{Show: true, Simple: false},
			wantBodyGlob:  "(?s)<html .*</html>",
		},
		{
			description: "html show, but request plaintext", url: "/another/url",
			acceptHeaders: "text/plain",
			options:       Options{Show: true, Simple: false},
			wantBodyGlob:  "(?s)panic: .*in: .*stack at time of panic:",
		},
		{
			description: "simple show", url: "/",
			acceptHeaders: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			options:       Options{Show: true, Simple: true},
			wantBodyGlob:  "(?s)panic: .*in: .*stack at time of panic:",
		},
		{
			description: "hide", url: "/",
			acceptHeaders: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			options:       Options{Show: false},
			wantBodyGlob:  http.StatusText(http.StatusInternalServerError),
		},
	}

	var logger bytes.Buffer

	for _, tt := range tests {
		// Set our custom logger.
		logger.Reset()
		tt.options.Logger = &logger
		server := httptest.NewServer(New(tt.options)(getTestHandler()))

		client := &http.Client{}
		req, err := http.NewRequest("GET", server.URL+tt.url, nil)
		at.NoError(err, tt.description)

		req.Header.Set("Accept", tt.acceptHeaders)
		resp, err := client.Do(req)
		at.NoError(err, tt.description)

		// Status code should always be an internal service error.
		at.Equal(resp.StatusCode, 500, tt.description)

		b, err := ioutil.ReadAll(resp.Body)
		at.NoError(err, tt.description)

		at.Regexp(tt.wantBodyGlob, string(b), tt.description)

		// Verify that the logged output is actually logged.
		line, err := logger.ReadString(0xa)
		at.NoError(err, tt.description)

		at.Equal("panic: "+panicError+"\n", line, tt.description)

		server.Close()
	}
}
