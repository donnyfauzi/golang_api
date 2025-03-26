package middleware

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword menerima password plain text dan mengembalikan hash bcrypt-nya
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return "", err
	}
	return string(hash), nil
}
