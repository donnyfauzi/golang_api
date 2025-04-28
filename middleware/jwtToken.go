package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Mendefinisikan tipe baru untuk kunci context
type ContextKey string

const UserIDKey ContextKey = "user_id" 

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

		// Ambil klaim dari token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Klaim token tidak valid", http.StatusForbidden)
			return
		}

		// Ambil user_id dari klaim
		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "user_id tidak ditemukan dalam token", http.StatusForbidden)
			return
		}

		// Tambahkan userID ke dalam context request
		ctx := context.WithValue(r.Context(), UserIDKey, int(userID))
		r = r.WithContext(ctx)

		// Token valid, lanjutkan ke handler
		next.ServeHTTP(w, r)
	})
}
