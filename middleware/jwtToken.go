package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func JwtVerify(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization") // Format: Bearer {token}
		if tokenHeader == "" {
			http.Error(w, "Token tidak tersedia", http.StatusForbidden)
			return
		}

		splitToken := strings.Split(tokenHeader, " ")
		if len(splitToken) != 2 {
			http.Error(w, "Format token salah", http.StatusForbidden)
			return
		}

		tokenPart := splitToken[1]
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenPart, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("method token tidak valid")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token tidak valid", http.StatusForbidden)
			return
		}

		// Token valid, lanjutkan ke handler
		next.ServeHTTP(w, r)
	})
}
