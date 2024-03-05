package models

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter

	code int
	size int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
	}
}

func (r *ResponseWriter) WriteHeader(code int) {
	if r.Code() == 0 {
		r.code = code
		r.ResponseWriter.WriteHeader(code)
	}
}

func (r *ResponseWriter) Write(body []byte) (int, error) {
	if r.Code() == 0 {
		r.WriteHeader(http.StatusOK)
	}

	var err error
	r.size, err = r.ResponseWriter.Write(body)

	return r.size, err
}

func (r *ResponseWriter) Flush() {
	if fl, ok := r.ResponseWriter.(http.Flusher); ok {
		if r.Code() == 0 {
			r.WriteHeader(http.StatusOK)
		}

		fl.Flush()
	}
}

func (r *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the hijacker interface is not supported")
	}

	return hj.Hijack()
}

func (r *ResponseWriter) Code() int {
	return r.code
}

func (r *ResponseWriter) Size() int {
	return r.size
}
