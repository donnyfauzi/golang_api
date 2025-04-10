package routes

import (
	"golang_api/controller"

	"github.com/gorilla/mux"
)

// RegisterRoutes untuk menangani rute-rute autentikasi
func UserRoutes(r *mux.Router) {

	// Rute untuk registrasi
	r.HandleFunc("/register", controller.Register).Methods("POST")

	// Rute untuk login
	r.HandleFunc("/login", controller.Login).Methods("POST")
}
