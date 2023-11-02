package main

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(ApiMiddleware(r))

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	addRoutes(r)

	http.ListenAndServe(":8000", r)
}

func addRoutes(r chi.Router) chi.Router {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("hello %s", r.URL.Query().Get("name"))))
	})

	return r
}

func ApiMiddleware(router chi.Router) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rctx := chi.NewRouteContext()
			path := r.URL.Path
			op := r.Method
			api := router.Find(rctx, op, path)

			fmt.Printf("api=%s\n", api)

			next.ServeHTTP(w, r)
		})
	}
}
