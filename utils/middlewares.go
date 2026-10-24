package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type UserInfo struct {
	User_id  string `json:"User_id"`
	Is_admin bool   `json:"Is_admin"`
	Email    string `json:"Email"`
}
type JWTClaims struct {
	jwt.RegisteredClaims
	UserInfo
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretCookie, err := r.Cookie("access")
		if err != nil {
			http.Error(w, fmt.Errorf("").Error(), http.StatusUnauthorized)
			return
		}
		cookie := secretCookie.Value
		JWT_SECRET := os.Getenv("JWT_SECRET")
		var claims JWTClaims
		token, err := jwt.ParseWithClaims(cookie, &claims, func(t *jwt.Token) (any, error) {
			return []byte(JWT_SECRET), nil
		}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("nanochemshop"), jwt.WithExpirationRequired())
		if err != nil {
			log.Println(err.Error())
			http.Error(w, fmt.Errorf("").Error(), http.StatusInternalServerError)
			return
		}
		if !claims.Is_admin || !token.Valid {
			http.Error(w, fmt.Errorf("Just Admins").Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", &claims.UserInfo)
		next.ServeHTTP(w, r.WithContext(ctx))

		// 3. Logic to run AFTER the handler has finished
		//    (e.g., timing the request, cleaning up)
	})
}
func AdminOrReadOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			next.ServeHTTP(w, r)
			return
		}
		secretCookie, err := r.Cookie("access")
		if err != nil {
			http.Error(w, fmt.Errorf("").Error(), http.StatusUnauthorized)
			return
		}
		cookie := secretCookie.Value
		JWT_SECRET := os.Getenv("JWT_SECRET")
		var claims JWTClaims
		token, err := jwt.ParseWithClaims(cookie, &claims, func(t *jwt.Token) (any, error) {
			return []byte(JWT_SECRET), nil
		}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("nanochemshop"), jwt.WithExpirationRequired())
		if err != nil {
			http.Error(w, fmt.Errorf("").Error(), http.StatusInternalServerError)
			return
		}
		if !claims.Is_admin || !token.Valid {
			http.Error(w, fmt.Errorf("Just Admins").Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", &claims.UserInfo)
		next.ServeHTTP(w, r.WithContext(ctx))

		// 3. Logic to run AFTER the handler has finished
		//    (e.g., timing the request, cleaning up)
	})
}

func AuthenticatedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretCookie, err := r.Cookie("access")
		if err != nil {
			http.Error(w, fmt.Errorf("").Error(), http.StatusUnauthorized)
			return
		}
		cookie := secretCookie.Value
		JWT_SECRET := os.Getenv("JWT_SECRET")
		var claims JWTClaims
		token, err := jwt.ParseWithClaims(cookie, &claims, func(t *jwt.Token) (any, error) {
			return []byte(JWT_SECRET), nil
		}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("nanochemshop"), jwt.WithExpirationRequired())
		if err != nil {
			log.Println(err.Error())
			http.Error(w, fmt.Errorf("").Error(), http.StatusInternalServerError)
			return
		}
		if !token.Valid {
			http.Error(w, fmt.Errorf("Just Admins").Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", &claims.UserInfo)
		next.ServeHTTP(w, r.WithContext(ctx))

		// 3. Logic to run AFTER the handler has finished
		//    (e.g., timing the request, cleaning up)
	})
}
