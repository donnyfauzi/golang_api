package controller

import (
	"encoding/json"
	"golang_api/middleware"
	"golang_api/model"
	"log"
	"net/http"
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
	user, err = model.CreateUser(user)
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
