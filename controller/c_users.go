package controller

import (
	"encoding/json"
	"golang_api/middleware"
	"golang_api/model"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Register menangani registrasi pengguna
func Register(w http.ResponseWriter, r *http.Request) {
	// Membaca body request sebagai JSON
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Gagal membaca data request", http.StatusBadRequest)
		return
	}

	// Hash password menggunakan bcrypt
	hashedPassword, err := middleware.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Gagal memproses password", http.StatusInternalServerError)
		return
	}

	// Perbarui user dengan password yang telah di-hash
	user.Password = hashedPassword

	// Simpan user baru
	user, err = model.CreateUser(user.Name, user.Email, user.Password)
	if err != nil {
		log.Println("Error menyimpan data pengguna:", err)
		http.Error(w, "Gagal registrasi", http.StatusInternalServerError)
		return
	}

	// Menampilkan response sukses dengan data pengguna dalam format JSON
	response := map[string]interface{}{
		"message": "Registrasi berhasil!",
		"user":    user,
	}

	// Set header content type dan kirim response dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login menangani proses login dan generate JWT token
func Login(w http.ResponseWriter, r *http.Request) {
	var creds model.User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Format data tidak valid", http.StatusBadRequest)
		return
	}

	// Cari user berdasarkan email
	user, err := model.FindUserByEmail(creds.Email)
	if err != nil {
		http.Error(w, "Email tidak ditemukan", http.StatusUnauthorized)
		return
	}

	// Cek password
	if !middleware.CheckPasswordHash(creds.Password, user.Password) {
		http.Error(w, "Password salah", http.StatusUnauthorized)
		return
	}

	// Buat JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	})

	// Ambil secret dari .env
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		http.Error(w, "JWT secret tidak tersedia", http.StatusInternalServerError)
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Gagal membuat token", http.StatusInternalServerError)
		return
	}

	// Kirim token ke client
	response := map[string]interface{}{
		"message": "Login berhasil",
		"token":   tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
