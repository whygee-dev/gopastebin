package middleware

import (
	"context"
	"gopastebin/consts"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")

		publicRoutes := consts.GetPublicRoutes()

		if slices.Contains(publicRoutes, r.URL.Path) {
			next.ServeHTTP(w, r)
		}

		if bearer == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

		bearerSplit := strings.Split(bearer, " ")

		if len(bearerSplit) != 2 || bearerSplit[0] != "Bearer" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		token, err := jwt.ParseWithClaims(bearerSplit[1], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return consts.GetSecret(), nil
		})

		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, "token", token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}