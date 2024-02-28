package handlers

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// Acts as an adapter for `http.ResponseWriter` type to store response status
// code and size.
type ResponseWriter struct {
	http.ResponseWriter

	code int
	size int
}

// Returns a new `ResponseWriter` type by decorating `http.ResponseWriter` type.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
	}
}

// Overrides `http.ResponseWriter` type.
func (r *ResponseWriter) WriteHeader(code int) {
	if r.Code() == 0 {
		r.code = code
		r.ResponseWriter.WriteHeader(code)
	}
}

// Overrides `http.ResponseWriter` type.
func (r *ResponseWriter) Write(body []byte) (int, error) {
	if r.Code() == 0 {
		r.WriteHeader(http.StatusOK)
	}

	var err error
	r.size, err = r.ResponseWriter.Write(body)

	return r.size, err
}

// Overrides `http.Flusher` type.
func (r *ResponseWriter) Flush() {
	if fl, ok := r.ResponseWriter.(http.Flusher); ok {
		if r.Code() == 0 {
			r.WriteHeader(http.StatusOK)
		}

		fl.Flush()
	}
}

// Overrides `http.Hijacker` type.
func (r *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the hijacker interface is not supported")
	}

	return hj.Hijack()
}

// Returns response status code.
func (r *ResponseWriter) Code() int {
	return r.code
}

// Returns response size.
func (r *ResponseWriter) Size() int {
	return r.size
}
