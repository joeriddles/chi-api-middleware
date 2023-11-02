package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	chi "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func Test_Main_RootPath(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(ApiMiddleware(r))
	addRoutes(r)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	// Act
	actual := captureOutput(func() {
		r.ServeHTTP(w, req)
	})

	// Assert
	expected := "api=/\n"
	if actual != expected {
		t.Fatalf("expected output to equal '%s' but it was '%s'", expected, actual)
	}
}

func Test_Main_NestedPath(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(ApiMiddleware(r))
	addRoutes(r)

	req, err := http.NewRequest("GET", "/hello", nil)
	require.NoError(t, err)

	// Act
	actual := captureOutput(func() {
		r.ServeHTTP(w, req)
	})

	// Assert
	expected := "api=/hello\n"
	if actual != expected {
		t.Fatalf("expected output to equal '%s' but it was '%s'", expected, actual)
	}
}

func Test_Main_NestedPathWithVariable(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(ApiMiddleware(r))
	addRoutes(r)

	req, err := http.NewRequest("GET", "/hello/john", nil)
	require.NoError(t, err)

	// Act
	actual := captureOutput(func() {
		r.ServeHTTP(w, req)
	})

	// Assert
	expected := "api=/hello/{name}\n"
	if actual != expected {
		t.Fatalf("expected output to equal '%s' but it was '%s'", expected, actual)
	}
}

// Taken from https://gist.githubusercontent.com/hauxe/e935a7f9012bf2649710cf75af323dbf/raw/ecf0ab41dcf0a743e3af70c9997c7d8a8a155bc7/output_capturing_full.go
// and https://medium.com/@hau12a1/golang-capturing-log-println-and-fmt-println-output-770209c791b4
func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}
