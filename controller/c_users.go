package controller

import (
	"encoding/json"
	"golang_api/helper"
	"golang_api/middleware"
	"golang_api/model"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Register menangani registrasi pengguna
func Register(w http.ResponseWriter, r *http.Request) {
	// Membaca body request sebagai JSON
	var user model.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	if user.Name == "" || user.Email == "" || user.Password == "" {
		helper.JSONError(w, http.StatusBadRequest, "Nama, email, dan password wajib diisi")
		return
	}

	existingUser,err := model.FindUserByEmail(user.Email)
	if err == nil && existingUser.Id != 0 {  // Jika email ditemukan
		helper.JSONError(w, http.StatusBadRequest, "Email sudah terdaftar")
		return
	}

	// Hash password menggunakan bcrypt
	user.Password, _ = middleware.HashPassword(user.Password)

	// Simpan user baru
	user, err = model.CreateUser(user.Name, user.Email, user.Password)
	if err != nil {
		helper.JSONError(w, http.StatusInternalServerError, "Gagal registrasi pengguna")
		return
	}

	// Menampilkan response sukses dengan data pengguna dalam format JSON
	response := map[string]any{
		"message": "Registrasi berhasil!",
		"user":    user,
	}

	// Set header content type dan kirim response dalam format JSON
	helper.JSONResponse(w, http.StatusCreated, response)
}

// Login menangani proses login dan generate JWT token
func Login(w http.ResponseWriter, r *http.Request) {
	var creds model.User
	_ = json.NewDecoder(r.Body).Decode(&creds)

	// Cari user berdasarkan email
	user, err := model.FindUserByEmail(creds.Email)
	if err != nil {
		helper.JSONError(w, http.StatusUnauthorized, "Email tidak ditemukan")
		return
	}

	// Cek password
	if !middleware.CheckPasswordHash(creds.Password, user.Password) {
		helper.JSONError(w, http.StatusUnauthorized, "Password salah")
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
		helper.JSONError(w, http.StatusInternalServerError, "JWT secret tidak tersedia")
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		helper.JSONError(w, http.StatusInternalServerError, "Gagal membuat token")
		return
	}

	// Kirim token ke client
	response := map[string]any{
		"message": "Login berhasil",
		"token":   tokenString,
	}

	helper.JSONResponse(w, http.StatusOK, response)
}
