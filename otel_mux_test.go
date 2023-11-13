package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	chi "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestOtelMux(t *testing.T) {
	t.Parallel()

	noop := func(w http.ResponseWriter, r *http.Request) {}

	tests := []struct {
		method     string
		pattern    string
		addHandler func(mx *OtelMux)
	}{
		// HTTP method-specific funcs
		{
			"GET",
			"/",
			func(mx *OtelMux) { mx.Get("/", noop) },
		},
		{
			"POST",
			"/",
			func(mx *OtelMux) { mx.Post("/", noop) },
		},
		{
			"CONNECT",
			"/",
			func(mx *OtelMux) { mx.Connect("/", noop) },
		},
		{
			"DELETE",
			"/",
			func(mx *OtelMux) { mx.Delete("/", noop) },
		},
		{
			"HEAD",
			"/",
			func(mx *OtelMux) { mx.Head("/", noop) },
		},
		{
			"OPTIONS",
			"/",
			func(mx *OtelMux) { mx.Options("/", noop) },
		},
		{
			"PATCH",
			"/",
			func(mx *OtelMux) { mx.Patch("/", noop) },
		},
		{
			"PUT",
			"/",
			func(mx *OtelMux) { mx.Put("/", noop) },
		},
		{
			"TRACE",
			"/",
			func(mx *OtelMux) { mx.Trace("/", noop) },
		},
		// custom handler funcs
		{
			"GET",
			"/",
			func(mx *OtelMux) { mx.Handle("/", chi.NewMux()) },
		},
		{
			"GET",
			"/",
			func(mx *OtelMux) { mx.HandleFunc("/", noop) },
		},
		{
			"GET",
			"/",
			func(mx *OtelMux) { mx.Method("GET", "/", chi.NewMux()) },
		},
		{
			"POST",
			"/",
			func(mx *OtelMux) { mx.Method("POST", "/", chi.NewMux()) },
		},
		// sub routes
		{
			"GET",
			"/",
			func(mx *OtelMux) { mx.Group(func(chi.Router) {}).Get("/", noop) },
		},
		{
			"GET",
			"/hello/world",
			func(mx *OtelMux) {
				mx.Route("/hello", func(r chi.Router) {
					r.Get("/world", noop)
				})
			},
		},
		{
			"GET",
			"/",
			func(mx *OtelMux) {
				mx.With(func(handler http.Handler) http.Handler { return handler }).Get("/", noop)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			t.Cleanup(func() {
				otelRequests = []*http.Request{}
			})

			mx := NewOtelMux(nil)
			test.addHandler(mx)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(test.method, test.pattern, nil)
			require.NoError(t, err)
			mx.ServeHTTP(w, req)

			require.Equal(t, len(otelRequests), 1, "no otel middleware recorded for %s", test.method)
		})
	}

}

func TestOtelMux_Group_CreatesOtelMux(t *testing.T) {
	mx := NewOtelMux(nil)
	r := mx.Group(func(r chi.Router) {})
	switch r.(type) {
	case *OtelMux:
		return
	default:
		t.Errorf("expected r.type to be ... but it was %T", r)
	}
}

func TestOtelMux_Group_CreatesNewInstance(t *testing.T) {
	mx := NewOtelMux(nil)
	r := mx.Group(func(r chi.Router) {})
	if mx == r {
		t.Errorf("expected r to be a new instance")
	}
}
