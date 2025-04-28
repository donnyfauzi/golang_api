package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"golang_api/helper"

	"github.com/golang-jwt/jwt"
)

// Mendefinisikan tipe baru untuk kunci context
type ContextKey string

const UserIDKey ContextKey = "user_id" 

func JwtVerify(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization") // Format: Bearer {token}
		if tokenHeader == "" {
			helper.JSONError(w, http.StatusUnauthorized, "Token tidak tersedia")
			return
		}

		splitToken := strings.Split(tokenHeader, " ")
		if len(splitToken) != 2 {
			helper.JSONError(w, http.StatusUnauthorized, "Format token salah")
			return
		}

		tokenPart := splitToken[1]
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenPart, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				helper.JSONError(w, http.StatusUnauthorized, "Method token tidak valid")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			helper.JSONError(w, http.StatusUnauthorized, "Token tidak valid")
			return
		}

		// Ambil klaim dari token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helper.JSONError(w, http.StatusUnauthorized, "Klaim token tidak valid")
			return
		}

		// Ambil user_id dari klaim
		userID, ok := claims["user_id"].(float64)
		if !ok {
			helper.JSONError(w, http.StatusUnauthorized, "User_id tidak ditemukan dalam token")
			return
		}

		// Tambahkan userID ke dalam context request
		ctx := context.WithValue(r.Context(), UserIDKey, int(userID))
		r = r.WithContext(ctx)

		// Token valid, lanjutkan ke handler
		next.ServeHTTP(w, r)
	})
}
