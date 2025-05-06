// =====================================================================================================================
// == LICENSE:       Copyright (c) 2025 Kevin De Coninck
// ==
// ==                Permission is hereby granted, free of charge, to any person
// ==                obtaining a copy of this software and associated documentation
// ==                files (the "Software"), to deal in the Software without
// ==                restriction, including without limitation the rights to use,
// ==                copy, modify, merge, publish, distribute, sublicense, and/or sell
// ==                copies of the Software, and to permit persons to whom the
// ==                Software is furnished to do so, subject to the following
// ==                conditions:
// ==
// ==                The above copyright notice and this permission notice shall be
// ==                included in all copies or substantial portions of the Software.
// ==
// ==                THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// ==                EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// ==                OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// ==                NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// ==                HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// ==                WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// ==                FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// ==                OTHER DEALINGS IN THE SOFTWARE.
// =====================================================================================================================

// Quality assurance: Verify (and measure the performance) of the public API of the "rapi" package.
package rapi_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-essentials/assert"
	"github.com/go-essentials/rapi"
	"github.com/go-essentials/tstsrv"
)

// UT: Make an HTTP POST request.
func TestPOST(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	t.Run("When the host is NOT available.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// ARRANGE.
		var got string

		request := rapi.POSTRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: "http://xyz.local/",
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
			},
		}

		// ACT.
		err := request.POST(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the endpoint isn't available.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	t.Run("When there's a custom handler for the received status code.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusUnauthorized},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got any

		request := rapi.POSTRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				HttpStatusCodeHandlers: map[int]func() error{
					http.StatusUnauthorized: func() error {
						return errors.New("error raised from the custom handler")
					},
				},
			},
		}

		// ACT.
		err := request.POST(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  The custom handler is for the received status code is invoked.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equal(t, err.Error(), "error raised from the custom handler", "", "\n\n"+
			"UT Name:  The custom handler is for the received status code is invoked.\n"+
			"\033[32mExpected: error raised from the custom handler\033[0m\n"+
			"\033[31mActual:   %s\033[0m\n\n", err.Error())
	})

	t.Run("When the endpoint is NOT available.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{})

		defer srvFake.Close()

		// ARRANGE.
		var got string

		request := rapi.POSTRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
			},
		}

		// ACT.
		err := request.POST(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the endpoint isn't available.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	t.Run("When the status code of the response is different from the 'OK' status code.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusBadRequest},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got string

		request := rapi.POSTRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				OkStatusCode: http.StatusOK,
			},
		}

		// ACT.
		err := request.POST(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the response is different from the 'OK' status code.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equal(t, err.Error(), "status code 400", "", "\n\n"+
			"UT Name:  An 'error' is returned when the response is different from the 'OK' status code.\n"+
			"\033[32mExpected: status code 400\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err.Error())
	})

	t.Run("When the HTTP response of can't be read.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusOK, DropConnection: true},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got any

		request := rapi.POSTRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				OkStatusCode: http.StatusOK,
			},
		}

		// ACT.
		err := request.POST(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the HTTP response can't be read.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	t.Run("When the HTTP response doesn't contain valid JSON.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusOK, Body: "invalid json"},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got any

		request := rapi.POSTRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				OkStatusCode: http.StatusOK,
			},
		}

		// ACT.
		err := request.POST(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the HTTP response doesn't contain valid JSON.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	t.Run("When the HTTP response does contain valid JSON.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusOK, Body: `{"id":"0"}`},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		type Response struct {
			Id string `json:"Id"`
		}

		var got Response
		var want Response = Response{Id: "0"}

		request := rapi.POSTRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				OkStatusCode: http.StatusOK,
			},
		}

		// ACT.
		err := request.POST(http.DefaultClient, &got)

		// ASSERT.
		assert.Nil(t, err, "", "\n\n"+
			"UT Name:  NO 'error' is returned when the HTTP response contains valid JSON.\n"+
			"\033[32mExpected: <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equal(t, got, want, "", "\n\n"+
			"UT Name:  The deserialized HTTP response is returned.\n"+
			"\033[32mExpected: %v\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", got, want)
	})
}

// UT: Make an HTTP GET request.
func TestGET(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	t.Run("When the host is NOT available.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// ARRANGE.
		var got string

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: "http://xyz.local/",
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
			},
		}

		// ACT.
		err := request.GET(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the endpoint isn't available.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	t.Run("When there's a custom handler for the received status code.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusUnauthorized},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got any

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				HttpStatusCodeHandlers: map[int]func() error{
					http.StatusUnauthorized: func() error {
						return errors.New("error raised from the custom handler")
					},
				},
			},
		}

		// ACT.
		err := request.GET(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  The custom handler is for the received status code is invoked.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equal(t, err.Error(), "error raised from the custom handler", "", "\n\n"+
			"UT Name:  The custom handler is for the received status code is invoked.\n"+
			"\033[32mExpected: error raised from the custom handler\033[0m\n"+
			"\033[31mActual:   %s\033[0m\n\n", err.Error())
	})

	t.Run("When the endpoint is NOT available.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{})

		defer srvFake.Close()

		// ARRANGE.
		var got string

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
			},
		}

		// ACT.
		err := request.GET(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the endpoint isn't available.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	t.Run("When the status code of the response is different from the 'OK' status code.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusBadRequest},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got string

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				OkStatusCode: http.StatusOK,
			},
		}

		// ACT.
		err := request.GET(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the response is different from the 'OK' status code.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equal(t, err.Error(), "status code 400", "", "\n\n"+
			"UT Name:  An 'error' is returned when the response is different from the 'OK' status code.\n"+
			"\033[32mExpected: status code 400\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err.Error())
	})

	t.Run("When the HTTP response of can't be read.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusOK, DropConnection: true},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got any

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				OkStatusCode: http.StatusOK,
			},
		}

		// ACT.
		err := request.GET(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the HTTP response can't be read.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	t.Run("When the HTTP response doesn't contain valid JSON.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusOK, Body: "invalid json"},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got any

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				OkStatusCode: http.StatusOK,
			},
		}

		// ACT.
		err := request.GET(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the HTTP response doesn't contain valid JSON.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	t.Run("When the HTTP response does contain valid JSON.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusOK, Body: `{"id":"0"}`},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		type Response struct {
			Id string `json:"Id"`
		}

		var got Response
		var want Response = Response{Id: "0"}

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				OkStatusCode: http.StatusOK,
			},
		}

		// ACT.
		err := request.GET(http.DefaultClient, &got)

		// ASSERT.
		assert.Nil(t, err, "", "\n\n"+
			"UT Name:  NO 'error' is returned when the HTTP response contains valid JSON.\n"+
			"\033[32mExpected: <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equal(t, got, want, "", "\n\n"+
			"UT Name:  The deserialized HTTP response is returned.\n"+
			"\033[32mExpected: %v\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", got, want)
	})
}

// UT: Make a "plain" HTTP GET request.
func TestGETPlain(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	t.Run("When the host is NOT available.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// ARRANGE.
		var got string

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: "http://xyz.local/",
			},
		}

		// ACT.
		err := request.GETPlain(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the endpoint isn't available.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	t.Run("When there's a custom handler for the received status code.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusUnauthorized},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got string

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				HttpStatusCodeHandlers: map[int]func() error{
					http.StatusUnauthorized: func() error {
						return errors.New("error raised from the custom handler")
					},
				},
			},
		}

		// ACT.
		err := request.GETPlain(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  The custom handler is for the received status code is invoked.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equal(t, err.Error(), "error raised from the custom handler", "", "\n\n"+
			"UT Name:  The custom handler is for the received status code is invoked.\n"+
			"\033[32mExpected: error raised from the custom handler\033[0m\n"+
			"\033[31mActual:   %s\033[0m\n\n", err.Error())
	})

	t.Run("When the endpoint is NOT available.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{})

		defer srvFake.Close()

		// ARRANGE.
		var got string

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
			},
		}

		// ACT.
		err := request.GETPlain(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the endpoint isn't available.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	t.Run("When the status code of the response is different from the 'OK' status code.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusBadRequest},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got string

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				OkStatusCode: http.StatusOK,
			},
		}

		// ACT.
		err := request.GETPlain(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the response is different from the 'OK' status code.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equal(t, err.Error(), "status code 400", "", "\n\n"+
			"UT Name:  An 'error' is returned when the response is different from the 'OK' status code.\n"+
			"\033[32mExpected: status code 400\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err.Error())
	})

	t.Run("When the HTTP response of can't be read.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusOK, DropConnection: true},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got string

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				OkStatusCode: http.StatusOK,
			},
		}

		// ACT.
		err := request.GETPlain(http.DefaultClient, &got)

		// ASSERT.
		assert.NotNil(t, err, "", "\n\n"+
			"UT Name:  An 'error' is returned when the HTTP response can't be read.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	t.Run("When the HTTP response does contain data.", func(t *testing.T) {
		t.Parallel() // Enable parallel execution.

		// FAKE SETUP.
		srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
			"/": {
				Responses: []tstsrv.Response{
					{StatusCode: http.StatusOK, Body: "HELLO, WORLD!"},
				},
			},
		})

		defer srvFake.Close()

		// ARRANGE.
		var got string

		request := rapi.GETRequestMsg{
			BaseRequest: rapi.BaseRequest{
				Endpoint: srvFake.URL(),
				HttpHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				OkStatusCode: http.StatusOK,
			},
		}

		// ACT.
		err := request.GETPlain(http.DefaultClient, &got)

		// ASSERT.
		assert.Nil(t, err, "", "\n\n"+
			"UT Name:  NO 'error' is returned when the HTTP response contains valid JSON.\n"+
			"\033[32mExpected: <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equal(t, got, "HELLO, WORLD!", "", "\n\n"+
			"UT Name:  The deserialized HTTP response is returned.\n"+
			"\033[32mExpected: HELLO, WORLD!\033[0m\n"+
			"\033[31mActual:   %s\033[0m\n\n", got)
	})
}
