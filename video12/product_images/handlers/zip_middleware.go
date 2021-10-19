package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipHandler struct {
}

// gorilla mux definition for middleware
// main coding for gzip still remains
func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// this statement returns true or false obviously
		if strings.Contains(r.Header.Get("Accept-encoding"), "gzip") {
			rw.Write([]byte("hello"))
			//create a gzipped response
			wrw := NewWrappedResponseWriter(&rw)
			wrw.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(wrw, r)
			defer wrw.Flush()
			return
		}

		// handle normal, i.e performs normal execution
		next.ServeHTTP(rw, r)
	})
}

// created a struct
type WrappedResponseWriter struct {
	rw http.ResponseWriter
	w  *gzip.Writer
}

// entry point to the whole zipping operation
func NewWrappedResponseWriter(rw *http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(*rw) // default way to use gzip package
	return &WrappedResponseWriter{rw: *rw, w: gw}
}

// implementing the header interface
func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.rw.Header() // here we call and return the inbuilt package function which returns the header
	// I know it seems why to do this
	// this is the only way to do it, i.e implementing the interface
}

func (wr *WrappedResponseWriter) Write(d []byte) (int, error) {
	return wr.w.Write(d) // same whihc is exmplained above in Header() is done here
}

func (wr *WrappedResponseWriter) WriteHeader(statuscode int) {
	wr.rw.WriteHeader(statuscode)
}

// might take some time to sink in
func (wr *WrappedResponseWriter) Flush() {
	wr.w.Flush()
	wr.w.Close()
}
