# hc

[![Build
Status](https://travis-ci.org/skamenetskiy/hc.svg?branch=master)](https://travis-ci.org/skamenetskiy/hc)

hc stands for "http client". It is built on top of
[fasthttp](https://github.com/valyala/fasthttp) and provides some handy
shortcuts to http client methods.

## Usage

```go
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"

	DefaultReadTimeout  = time.Duration(0)
	DefaultWriteTimeout = time.Duration(0)
)
```

#### func  Delete

```go
func Delete(url string, body []byte, headers Headers) (int, []byte, error)
```
Delete makes a DELETE request to url, sends body in request body, sends headers
and returns http status code, response body, error

#### func  Get

```go
func Get(url string) (int, []byte, error)
```
Get makes a GET request to url and returns: http status code, response body,
error

#### func  MakeRawRequest

```go
func MakeRawRequest(req *Request, res *Response) error
```
MakeRequest is doing a request using req and res and returns an error

#### func  Post

```go
func Post(url string, body []byte, headers Headers) (int, []byte, error)
```
Post makes a POST request to url, sends body in request body, sends headers and
returns http status code, response body, error

#### func  Put

```go
func Put(url string, body []byte, headers Headers) (int, []byte, error)
```
Put makes a PUT request to url, sends body in request body, sends headers and
returns http status code, response body, error

#### func  SetReadTimeout

```go
func SetReadTimeout(t time.Duration)
```
SetReadTimeout sets the client ReadTimeout

#### func  SetWriteTimeout

```go
func SetWriteTimeout(t time.Duration)
```
SetWriteTimeout sets the client WriteTimeout

#### type Headers

```go
type Headers map[string]string
```

Headers key/value map if headers

#### func (Headers) Add

```go
func (h Headers) Add(k string, v string)
```
Add adds a header to Headers

#### type Request

```go
type Request struct {
	*fasthttp.Request
}
```

Request request object

#### func  AcquireRequest

```go
func AcquireRequest() *Request
```
AcquireRequest returns a new Request

#### func (*Request) WriteJSON

```go
func (r *Request) WriteJSON(v interface{}) error
```
WriteJSON will write v (as json) to Request body

#### type Response

```go
type Response struct {
	*fasthttp.Response
}
```

Response response object

#### func  AcquireResponse

```go
func AcquireResponse() *Response
```
AcquireResponse returns a new response

#### func  MakeRequest

```go
func MakeRequest(method string, url string, body []byte, headers Headers) (*Response, error)
```
MakeRequest makes a new request using method to url, sends body in post body,
sends headers in headers and returns Response object and/or error

#### func (*Response) ReadJSON

```go
func (r *Response) ReadJSON(v interface{}) error
```
ReadJSON will read json into v from Response body
