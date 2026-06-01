package router

import "net/http"

type (
	Router interface {
		Get(path string, handler http.HandlerFunc) Router
		Post(path string, handler http.HandlerFunc) Router
		Put(path string, handler http.HandlerFunc) Router
		Delete(path string, handler http.HandlerFunc) Router

		Use(middleware ...Middleware) Router

		Group(basePath string) Router
		SubGroup(path string) Router
	}

	Middleware func(http.Handler) http.Handler
)

func NewRouter(mux *http.ServeMux, options ...MuxOption) *Mux {
	m := &Mux{
		mux:         mux,
		middlewares: []Middleware{},
	}

	for _, opt := range options {
		opt(m)
	}

	return m
}
