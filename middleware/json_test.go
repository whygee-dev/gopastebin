package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJsonContentTypeMiddleware(t *testing.T) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    middlewareHandler := JsonContentTypeMiddleware(handler)

    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    middlewareHandler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
        t.Errorf("handler returned wrong content type: got %v want %v",
            contentType, "application/json")
    }

    expected := "Hello, World!"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}
