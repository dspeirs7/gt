package router

import (
	"fmt"
	"net/http"
	"path/filepath"
	"slices"
)

type (
	Mux struct {
		mux         *http.ServeMux
		basePath    string
		middlewares []Middleware
	}

	MuxOption func(*Mux)
)

func getRoute(method, basePath, path string) string {
	return fmt.Sprintf("%s %s", method, filepath.Join(basePath, path))
}

func applyMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) <= 0 {
		return handler
	}

	next := handler

	for _, middleware := range slices.Backward(middlewares) {
		next = middleware(next)
	}

	return next
}

func copyMiddlewares(mws []Middleware) []Middleware {
	copied := make([]Middleware, len(mws))

	copy(copied, mws)

	return copied
}

func (r *Mux) setRoute(method, path string, handler http.HandlerFunc) Router {
	route := getRoute(method, r.basePath, path)

	wrappedHandler := applyMiddleware(handler, r.middlewares...)

	r.mux.Handle(route, wrappedHandler)

	return r
}

func (r *Mux) Get(path string, handler http.HandlerFunc) Router {
	return r.setRoute(http.MethodGet, path, handler)
}

func (r *Mux) Post(path string, handler http.HandlerFunc) Router {
	return r.setRoute(http.MethodPost, path, handler)
}

func (r *Mux) Put(path string, handler http.HandlerFunc) Router {
	return r.setRoute(http.MethodPut, path, handler)
}

func (r *Mux) Delete(path string, handler http.HandlerFunc) Router {
	return r.setRoute(http.MethodDelete, path, handler)
}

func (r *Mux) Use(middlewares ...Middleware) Router {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *Mux) Group(basePath string) Router {
	rCopy := r

	rCopy.basePath = basePath
	rCopy.middlewares = copyMiddlewares(r.middlewares)

	return rCopy
}

func (r *Mux) SubGroup(path string) Router {
	rCopy := r

	rCopy.basePath = filepath.Join(rCopy.basePath, path)
	rCopy.middlewares = copyMiddlewares(r.middlewares)

	return rCopy
}

func WithBasePath(basePath string) MuxOption {
	return func(m *Mux) {
		m.basePath = basePath
	}
}

func WithMiddleware(middlewares ...Middleware) MuxOption {
	return func(m *Mux) {
		m.middlewares = append(m.middlewares, middlewares...)
	}
}
