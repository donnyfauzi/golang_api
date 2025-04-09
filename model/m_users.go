package model

import (
	"database/sql"
	"golang_api/database"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser menyimpan user baru dan mengembalikan data lengkapnya
func CreateUser(name, email, hashedPassword string) (User, error) {
	var user User

	query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
	result, err := database.DB.Exec(query, name, email, hashedPassword)
	if err != nil {
		return user, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return user, err
	}

	user = User{
		Id:       int(id),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	return user, nil
}

// FindUserByEmail mencari user berdasarkan email
func FindUserByEmail(email string) (User, error) {
	var user User

	query := "SELECT id, name, email, password FROM users WHERE email = ?"
	err := database.DB.QueryRow(query, email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password,
	)

	if err == sql.ErrNoRows {
		return user, nil // atau return empty + error khusus kalau kamu ingin handle "tidak ditemukan"
	}

	return user, err
}
