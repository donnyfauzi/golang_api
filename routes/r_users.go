package routes

import (
	"golang_api/controller"

	"github.com/gorilla/mux"
)

// RegisterRoutes untuk menangani rute-rute autentikasi
func RegisterRoutes(r *mux.Router) {
	// Rute untuk registrasi
	r.HandleFunc("/register", controller.Register).Methods("POST")
}
