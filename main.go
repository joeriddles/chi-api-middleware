package main

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(ApiMiddleware)

	r.Use(middleware.RequestID)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("hello %s", r.URL.Query().Get("name"))))
	})

	http.ListenAndServe(":8000", r)
}

func ApiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api := chi.RouteContext(r.Context()).RoutePattern()
		fmt.Printf("before: api=%s\n", api)

		next.ServeHTTP(w, r)

		api = chi.RouteContext(r.Context()).RoutePattern()
		fmt.Printf("after: api=%s\n", api)
	})
}
