package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	chi "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func Test_Middleware_Order(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Use(recoverMiddleware())

	pattern := "/middleware"
	r.With(otelMiddleware(pattern, nil)).Get(pattern, func(w http.ResponseWriter, r *http.Request) {
		panic("oh no!")
	})

	req, err := http.NewRequest("GET", pattern, nil)
	require.NoError(t, err)

	// Act
	actual := captureOutput(func() {
		r.ServeHTTP(w, req)
	})

	// Assert
	expected := []string{
		"pattern=/middleware",
		"recovered from error with value=oh no!",
	}
	assertLogs(t, actual, expected)
}

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
	expected := []string{"api=/"}
	assertLogs(t, actual, expected)
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
	expected := []string{"api=/hello"}
	assertLogs(t, actual, expected)
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
	expected := []string{"api=/hello/{name}"}
	assertLogs(t, actual, expected)
}

// Taken from https://gist.githubusercontent.com/hauxe/e935a7f9012bf2649710cf75af323dbf/raw/ecf0ab41dcf0a743e3af70c9997c7d8a8a155bc7/output_capturing_full.go
// and https://medium.com/@hau12a1/golang-capturing-log-println-and-fmt-println-output-770209c791b4
func captureOutput(f func()) []string {
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

	cap := <-out
	lines := strings.Split(cap, "\n")
	return lines
}

func assertLogs(t *testing.T, actual []string, expected []string) {
	if len(actual) != len(expected) {
		t.Errorf("Expected to have %v logs but only had %v", len(expected), len(actual))
	}

	for i, expectedLog := range expected {
		actualLog := actual[i]
		if actualLog != expectedLog {
			t.Errorf("Expected log at %v to be \"%v\" but it was \"%v\"", i, expectedLog, actualLog)
		}
	}
}
