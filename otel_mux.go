package main

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

var _ chi.Router = &OtelMux{}

// Wrapper around base Mux that adds otel metric reporting for APIs
type OtelMux struct {
	chi.Router

	otelClient Client
}

func NewOtelMux(otelClient Client) *OtelMux {
	mx := chi.NewMux()
	return &OtelMux{Router: mx, otelClient: otelClient}
}

func (mx *OtelMux) withRouter(router chi.Router) *OtelMux {
	return &OtelMux{Router: router, otelClient: mx.otelClient}
}

func (mx *OtelMux) Handle(pattern string, handler http.Handler) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).Handle(pattern, handler)
}

func (mx *OtelMux) HandleFunc(pattern string, handlerFn http.HandlerFunc) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).HandleFunc(pattern, handlerFn)
}

func (mx *OtelMux) Method(method, pattern string, handler http.Handler) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).Method(method, pattern, handler)
}

func (mx *OtelMux) Connect(pattern string, handlerFn http.HandlerFunc) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).Connect(pattern, handlerFn)
}

func (mx *OtelMux) Delete(pattern string, handlerFn http.HandlerFunc) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).Delete(pattern, handlerFn)
}

func (mx *OtelMux) Get(pattern string, handlerFn http.HandlerFunc) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).Get(pattern, handlerFn)
}

func (mx *OtelMux) Head(pattern string, handlerFn http.HandlerFunc) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).Head(pattern, handlerFn)
}

func (mx *OtelMux) Options(pattern string, handlerFn http.HandlerFunc) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).Options(pattern, handlerFn)
}

func (mx *OtelMux) Patch(pattern string, handlerFn http.HandlerFunc) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).Patch(pattern, handlerFn)
}

func (mx *OtelMux) Post(pattern string, handlerFn http.HandlerFunc) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).Post(pattern, handlerFn)
}

func (mx *OtelMux) Put(pattern string, handlerFn http.HandlerFunc) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).Put(pattern, handlerFn)
}

func (mx *OtelMux) Trace(pattern string, handlerFn http.HandlerFunc) {
	mx.With(otelMiddleware(pattern, mx.otelClient)).Trace(pattern, handlerFn)
}

func (mx *OtelMux) Group(fn func(r chi.Router)) chi.Router {
	r := mx.Router.Group(fn)
	return mx.withRouter(r)
}

func (mx *OtelMux) Route(pattern string, fn func(r chi.Router)) chi.Router {
	// override entire Mux.Route method because there is no way to inject a
	// difference kind of Router
	if fn == nil {
		panic(fmt.Sprintf("chi: attempting to Route() a nil subrouter on '%s'", pattern))
	}
	subRouter := NewOtelMux(mx.otelClient)
	fn(subRouter)
	mx.Mount(pattern, subRouter)
	return subRouter
}

// TODO: fix this
func (mx *OtelMux) With(middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := mx.Router.With(middlewares...)
	return mx.withRouter(r)
}

// Note: this will result in a stack overflow error
// func (mx *OtelMux) With(middlewares ...func(http.Handler) http.Handler) chi.Router {
// 	mx = NewOtelMux(mx.otelClient)
// 	mx.Router = mx.Router.With(middlewares...)
// 	return mx
// }
