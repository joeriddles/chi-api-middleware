package main

import (
	"fmt"
	"net/http"
)

type Client = interface{}

var otelRequests []*http.Request = []*http.Request{}

// noop otel middleware
func otelMiddleware(pattern string, otelClient Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			otelRequests = append(otelRequests, r)
			fmt.Printf("pattern=%s\n", pattern)
			next.ServeHTTP(w, r)
		})
	}
}

func recoverMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if errVal := recover(); errVal != nil {
					fmt.Printf("recovered from error with value=%v", errVal)
					w.WriteHeader(500)
					w.Write([]byte("Internal server error"))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
