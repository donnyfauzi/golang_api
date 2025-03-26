package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB instance
var DB *sql.DB

// ConnectDB untuk menghubungkan ke database MySQL
func ConnectDB() {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang_api") // Ubah sesuai kredensial DB
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	// Cek koneksi
	err = DB.Ping()
	if err != nil {
		log.Fatal("Database tidak bisa dijangkau:", err)
	}

	fmt.Println("Database MySql terhubung!")
}
