package main

import (
	"fmt"
	"golang_api/database"
	"golang_api/kafka"
	"golang_api/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Menghubungkan ke database
	database.ConnectDB()

	// Connect ke Kafka
	kafka.InitKafka()
	go kafka.StartConsumer()

	// Menyiapkan router
	r := mux.NewRouter()

	// Daftarkan rute untuk autentikasi
	routes.UserRoutes(r)
	routes.ProductsRoutes(r)
	routes.ChartRoutes(r)

	// Menjalankan server
	fmt.Println("Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
