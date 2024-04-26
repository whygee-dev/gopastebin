package middleware

import (
	"gopastebin/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
)

func TestAuthMiddlewareHappyPath(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value("token").(*jwt.Token)

		if token != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Token: " + token.Raw ))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("No token found\n"))
		}
	})

	middlewareHandler := AuthMiddleware(handler)

	token, err := service.CreateToken()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer " + token)
	w := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expectedBody := "Token: " + token 
	if w.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, w.Body.String())
	}
}

func TestAuthMiddlewareUnauthorized(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value("token").(*jwt.Token)

		if token != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Token: " + token.Raw ))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("No token found\n"))
		}
	})

	middlewareHandler := AuthMiddleware(handler)

	token := "invalid"
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer " + token)
	w := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddlewareIgnorePublicSignupRoute(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middlewareHandler := AuthMiddleware(handler)

	req := httptest.NewRequest("POST", "/user/signup", nil)
	w := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddlewareIgnorePublicLoginRoute(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middlewareHandler := AuthMiddleware(handler)

	req := httptest.NewRequest("GET", "/user/login", nil)
	w := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddlewareIgnoreNoBearer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middlewareHandler := AuthMiddleware(handler)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bear" )
	w := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
