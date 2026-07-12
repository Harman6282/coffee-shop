package middleware

import (
	"net/http"

	"github.com/Harman6282/coffee-shop/shared/env"
	"github.com/golang-jwt/jwt/v5"
)


func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		jwtSecret := env.GetString("JWT_SECRET", "myjwtsecretforauth")

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(
			cookie.Value,
			claims,
			func(token *jwt.Token) (any, error) {
				return []byte(jwtSecret), nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		)
		if err != nil || !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.Subject != "" {
			r.Header.Set("X-User-Email", claims.Subject)
		}

		next.ServeHTTP(w, r)
	})
}
