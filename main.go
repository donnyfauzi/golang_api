package main

import (
	"fmt"
	"golang_api/database"
	"golang_api/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Menghubungkan ke database
	database.ConnectDB()

	// Menyiapkan router
	r := mux.NewRouter()

	// Daftarkan rute untuk autentikasi
	routes.RegisterRoutes(r)

	// Menjalankan server
	fmt.Println("Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
