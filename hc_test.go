package hc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

func TestStart(t *testing.T) {
	suite.Run(t, new(HCTestSuite))
}

type HCTestSuite struct {
	suite.Suite
}

func (t *HCTestSuite) TestGet() {
	str := []byte("some string")
	s1 := t.server(func(rs http.ResponseWriter, rq *http.Request) {
		rs.Write(str)
		t.Equal(MethodGet, rq.Method)
	})
	defer s1.Close()
	c1, b1, err := Get(s1.URL)
	t.NoError(err)
	t.Equal(fasthttp.StatusOK, c1)
	t.Equal(str, b1)

	s2 := t.server(func(rs http.ResponseWriter, rq *http.Request) {
		rs.WriteHeader(fasthttp.StatusInternalServerError)
	})
	defer s2.Close()
	c2, b2, err := Get(s2.URL)
	t.NoError(err)
	t.Equal(fasthttp.StatusInternalServerError, c2)
	t.Empty(b2)

	c3, b3, err := Get("")
	t.Error(err)
	t.Empty(c3)
	t.Empty(b3)
}

func (t *HCTestSuite) TestPost() {
	str1 := []byte("some string")
	str2 := []byte("other string")
	s1 := t.server(func(rs http.ResponseWriter, rq *http.Request) {
		b, err := ioutil.ReadAll(rq.Body)
		t.NoError(err)
		t.Equal(MethodPost, rq.Method)
		t.NotEmpty(b)
		t.Equal(str1, b)
		rs.WriteHeader(fasthttp.StatusMethodNotAllowed)
		rs.Write(str2)
	})
	c1, b1, err := Post(s1.URL, str1, nil)
	t.NoError(err)
	t.Equal(str2, b1)
	t.Equal(fasthttp.StatusMethodNotAllowed, c1)
}

func (t *HCTestSuite) TestPut() {
	str1 := []byte("some string")
	str2 := []byte("other string")
	s1 := t.server(func(rs http.ResponseWriter, rq *http.Request) {
		b, err := ioutil.ReadAll(rq.Body)
		t.NoError(err)
		t.Equal(MethodPut, rq.Method)
		t.NotEmpty(b)
		t.Equal(str1, b)
		rs.WriteHeader(fasthttp.StatusNotFound)
		rs.Write(str2)
	})
	c1, b1, err := Put(s1.URL, str1, nil)
	t.NoError(err)
	t.Equal(str2, b1)
	t.Equal(fasthttp.StatusNotFound, c1)
}

func (t *HCTestSuite) TestDelete() {
	str1 := []byte("some string")
	str2 := []byte("other string")
	s1 := t.server(func(rs http.ResponseWriter, rq *http.Request) {
		b, err := ioutil.ReadAll(rq.Body)
		t.NoError(err)
		t.Equal(MethodDelete, rq.Method)
		t.NotEmpty(b)
		t.Equal(str1, b)
		rs.WriteHeader(fasthttp.StatusBadRequest)
		rs.Write(str2)
	})
	c1, b1, err := Delete(s1.URL, str1, nil)
	t.NoError(err)
	t.Equal(str2, b1)
	t.Equal(fasthttp.StatusBadRequest, c1)
}

func (t *HCTestSuite) TestAcquireRequest() {
	r := AcquireRequest()
	t.IsType(&Request{}, r)
	t.IsType(&fasthttp.Request{}, r.Request)
}

func (t *HCTestSuite) TestAcquireResponse() {
	r := AcquireResponse()
	t.IsType(&Response{}, r)
	t.IsType(&fasthttp.Response{}, r.Response)
}

func (t *HCTestSuite) TestHeaders() {
	s := t.server(func(rs http.ResponseWriter, rq *http.Request) {
		t.Equal(rq.Header.Get("Some-Header"), "Value")
		t.Equal(rq.Header.Get("Some-Other-Header"), "Value2")
	})
	h1 := Headers{}
	h1.Add("Some-Header", "Value")
	h1.Add("Some-Other-Header", "Value2")
	Post(s.URL, nil, h1)
	Post(s.URL, nil, map[string]string{
		"Some-Header":       "Value",
		"Some-Other-Header": "Value2",
	})
}

func (t *HCTestSuite) TestSetReadTimeout() {
	rt := time.Second * 60
	SetReadTimeout(rt)
	t.Equal(rt, client.ReadTimeout)
}

func (t *HCTestSuite) TestSetWriteTimeout() {
	rt := time.Second * 60
	SetWriteTimeout(rt)
	t.Equal(rt, client.WriteTimeout)
}

func (t *HCTestSuite) TestMakeRawRequest() {
	rq := AcquireRequest()
	rs := AcquireResponse()
	s := t.server(func(rs http.ResponseWriter, rq *http.Request) {
		rs.WriteHeader(fasthttp.StatusBadGateway)
	})
	rq.SetRequestURI(s.URL)
	err := MakeRawRequest(rq, rs)
	t.NoError(err)
	t.Equal(fasthttp.StatusBadGateway, rs.StatusCode())
}

func (t *HCTestSuite) TestWriteJSON() {
	type tt struct {
		Data string `json:"data"`
	}

	tt1 := tt{
		Data: "good data",
	}
	tt2 := tt{}
	s := t.server(func(rs http.ResponseWriter, rq *http.Request) {
		err := json.NewDecoder(rq.Body).Decode(&tt2)
		t.NoError(err)
		t.Equal(tt1, tt2)
		t.Equal(tt2.Data, "good data")
	})
	rq := AcquireRequest()
	rs := AcquireResponse()
	rq.SetRequestURI(s.URL)
	rq.Header.SetMethod(MethodPost)
	err := rq.WriteJSON(tt1)
	t.NoError(err)
	err = MakeRawRequest(rq, rs)
	t.NoError(err)
}

func (t *HCTestSuite) TestReadJSON() {
	s := t.server(func(rs http.ResponseWriter, rq *http.Request) {
		rs.Write([]byte(`{"d1":"data1","d2":"data2"}`))
	})
	rq := AcquireRequest()
	rs := AcquireResponse()
	rq.SetRequestURI(s.URL)
	err := MakeRawRequest(rq, rs)
	t.NoError(err)
	type tt struct {
		Data1 string `json:"d1"`
		Data2 string `json:"d2"`
	}
	tto := &tt{}
	err = rs.ReadJSON(tto)
	t.NoError(err)
	t.Equal("data1", tto.Data1)
	t.Equal("data2", tto.Data2)
}

func (t *HCTestSuite) server(f func(rs http.ResponseWriter, rq *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(f))
}
