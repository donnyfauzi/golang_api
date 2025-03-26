package model

import (
	"golang_api/database"
	"log"
)

// User representasi dari pengguna dalam sistem
type User struct {
	ID       int    `json:"ID"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

// CreateUser menyimpan data pengguna baru ke dalam database
func CreateUser(user User) (User, error) {
	// Query untuk memasukkan data pengguna baru
	result, err := database.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		log.Println("Error menambahkan user:", err)
		return user, err
	}

	// Mengambil ID pengguna yang baru saja dimasukkan
	userID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error mendapatkan ID pengguna:", err)
		return user, err
	}

	// Menetapkan ID yang baru diperoleh ke objek user
	user.ID = int(userID)

	return user, nil
}
