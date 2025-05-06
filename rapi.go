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

// Package rapi provides an easy API for making HTTP requests.
package rapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// BaseRequest describes the "base" structure of an HTTP request.
type BaseRequest struct {
	Endpoint               string               // The URL to send the request to.
	HttpHeaders            map[string]string    // The HTTP headers to include in the request.
	HttpStatusCodeHandlers map[int]func() error // Map containing the HTTP status codes and their corresponding handlers.
	OkStatusCode           int                  // The HTTP status code that indicates a successful request.
}

// POSTRequestMsg describes an HTTP POST request.
type POSTRequestMsg struct {
	BaseRequest        // The "base" HTTP request.
	Payload     string // The payload of the request.
}

// GETRequestMsg describes an HTTP GET request.
type GETRequestMsg struct {
	BaseRequest // The "base" HTTP request.
}

// POST uses client to make an HTTP POST request described by req and updates result.
// It return an error if any error occurs or <nil> when no error was returned.
func (req *POSTRequestMsg) POST(client *http.Client, result any) error {
	requestBytes := bytes.NewBuffer([]byte(req.Payload))
	request, _ := http.NewRequest("POST", req.Endpoint, requestBytes)

	for key, value := range req.HttpHeaders {
		request.Header.Add(key, value)
	}

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if handler, found := req.HttpStatusCodeHandlers[response.StatusCode]; found {
		return handler()
	}

	if response.StatusCode == http.StatusNotImplemented {
		return errors.New("not implemented")
	}

	if response.StatusCode != req.OkStatusCode {
		return fmt.Errorf("status code %d", response.StatusCode)
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(responseData, &result); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// GET uses client to make an HTTP GET request described by req and updates result.
// It return an error if any error occurs or <nil> when no error was returned.
func (req *GETRequestMsg) GET(client *http.Client, result any) error {
	request, _ := http.NewRequest("GET", req.Endpoint, nil)

	for key, value := range req.HttpHeaders {
		request.Header.Add(key, value)
	}

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if handler, found := req.HttpStatusCodeHandlers[response.StatusCode]; found {
		return handler()
	}

	if response.StatusCode == http.StatusNotImplemented {
		return errors.New("not implemented")
	}

	if response.StatusCode != req.OkStatusCode {
		return fmt.Errorf("status code %d", response.StatusCode)
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(responseData, &result); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// GETPlain uses client to make an HTTP GET request described by req and updates result.
// It return an error if any error occurs or <nil> when no error was returned.
func (req *GETRequestMsg) GETPlain(client *http.Client, result *string) error {
	request, _ := http.NewRequest("GET", req.Endpoint, nil)

	for key, value := range req.HttpHeaders {
		request.Header.Add(key, value)
	}

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if handler, found := req.HttpStatusCodeHandlers[response.StatusCode]; found {
		return handler()
	}

	if response.StatusCode == http.StatusNotImplemented {
		return errors.New("not implemented")
	}

	if response.StatusCode != req.OkStatusCode {
		return fmt.Errorf("status code %d", response.StatusCode)
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	*result = string(responseData)

	return nil
}
