// # hc
//
// [![Build Status](https://travis-ci.org/skamenetskiy/hc.svg?branch=master)](https://travis-ci.org/skamenetskiy/hc)
//
// hc stands for "http client". It is built on top of [fasthttp](https://github.com/valyala/fasthttp) and provides
// some handy shortcuts to http client methods.
package hc

import (
	"encoding/json"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"

	DefaultReadTimeout  = time.Duration(0)
	DefaultWriteTimeout = time.Duration(0)
)

var client = &fasthttp.Client{
	ReadTimeout:  DefaultReadTimeout,
	WriteTimeout: DefaultWriteTimeout,
}

// Get makes a GET request to url and returns:
// http status code, response body, error
func Get(url string) (int, []byte, error) {
	return makeRequest(MethodGet, url, nil, nil)
}

// Post makes a POST request to url, sends body in request
// body, sends headers and returns http status code,
// response body, error
func Post(url string, body []byte, headers Headers) (int, []byte, error) {
	return makeRequest(MethodPost, url, body, headers)
}

// Put makes a PUT request to url, sends body in request
// body, sends headers and returns http status code,
// response body, error
func Put(url string, body []byte, headers Headers) (int, []byte, error) {
	return makeRequest(MethodPut, url, body, headers)
}

// Delete makes a DELETE request to url, sends body in request
// body, sends headers and returns http status code,
// response body, error
func Delete(url string, body []byte, headers Headers) (int, []byte, error) {
	return makeRequest(MethodDelete, url, body, headers)
}

// AcquireRequest returns a new Request
func AcquireRequest() *Request {
	return &Request{fasthttp.AcquireRequest()}
}

// AcquireResponse returns a new response
func AcquireResponse() *Response {
	return &Response{fasthttp.AcquireResponse()}
}

// MakeRequest makes a new request using method to url, sends
// body in post body, sends headers in headers and returns
// Response object and/or error
func MakeRequest(method string, url string, body []byte, headers Headers) (*Response, error) {
	// create request and response
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	// configure request
	req.Header.SetMethod(method)
	req.SetRequestURI(url)
	req.SetBody(body)
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// do request
	err := client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return &Response{res}, nil
}

// MakeRequest is doing a request using req and res
// and returns an error
func MakeRawRequest(req *Request, res *Response) error {
	return client.Do(req.Request, res.Response)
}

// SetReadTimeout sets the client ReadTimeout
func SetReadTimeout(t time.Duration) {
	client.ReadTimeout = t
}

// SetWriteTimeout sets the client WriteTimeout
func SetWriteTimeout(t time.Duration) {
	client.WriteTimeout = t
}

func makeRequest(method string, url string, body []byte, headers Headers) (int, []byte, error) {
	res, err := MakeRequest(method, url, body, headers)
	if err != nil {
		return 0, nil, err
	}
	return res.StatusCode(), res.Body(), nil
}

// Request request object
type Request struct {
	*fasthttp.Request
}

// WriteJSON will write v (as json) to Request body
func (r *Request) WriteJSON(v interface{}) error {
	return json.NewEncoder(r.BodyWriter()).Encode(v)
}

// Response response object
type Response struct {
	*fasthttp.Response
}

// ReadJSON will read json into v from Response body
func (r *Response) ReadJSON(v interface{}) error {
	return json.Unmarshal(r.Body(), v)
}

// Headers key/value map if headers
type Headers map[string]string

// Add adds a header to Headers
func (h Headers) Add(k string, v string) {
	h[k] = v
}
