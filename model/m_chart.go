package model

import (
	"database/sql"
	"golang_api/database"
	"log"
)

type ChartItem struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

func CreateChart(user_id, product_id, quantity, created_at, updated_at int) (ChartItem, error) {
	var chart ChartItem

	// Cek apakah sudah ada data untuk user_id dan product_id ini
	var existingQuantity int
	var existingID int
	err := database.DB.QueryRow("SELECT id, quantity FROM chart WHERE user_id = ? AND product_id = ?", user_id, product_id).Scan(&existingID, &existingQuantity)

	if err == nil {
		// Data sudah ada, maka update quantity yang ada
		newQuantity := existingQuantity + quantity

		// Debug log untuk melihat jumlah quantity
		log.Printf("Updating Quantity: user_id = %d, product_id = %d, old_quantity = %d, new_quantity = %d", user_id, product_id, existingQuantity, newQuantity)

		_, err = database.DB.Exec(`
			UPDATE chart 
			SET quantity = ?, updated_at = ? 
			WHERE id = ?
		`, newQuantity, updated_at, existingID) // Gunakan existingID untuk update

		if err != nil {
			return chart, err
		}

		// Log jika berhasil
		log.Printf("Data berhasil diupdate: user_id = %d, product_id = %d, quantity = %d", user_id, product_id, newQuantity)

		// Set data ChartItem dengan quantity baru
		chart = ChartItem{
			ID:        existingID,   // Gunakan existingID untuk update
			UserID:    user_id,
			ProductID: product_id,
			Quantity:  newQuantity,  // Update quantity dengan jumlah yang baru
			CreatedAt: created_at,   // Bisa mengambil created_at lama jika perlu
			UpdatedAt: updated_at,   // Waktu update
		}

	} else if err == sql.ErrNoRows {
		// Jika tidak ada data, lakukan insert baru
		query := "INSERT INTO chart (user_id, product_id, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
		result, err := database.DB.Exec(query, user_id, product_id, quantity, created_at, updated_at)
		if err != nil {
			return chart, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return chart, err
		}

		// Log ID hasil insert
		log.Printf("Data berhasil disimpan: ID = %d, user_id = %d, product_id = %d, quantity = %d", id, user_id, product_id, quantity)

		chart = ChartItem{
			ID:        int(id),        // ID hasil insert
			UserID:    user_id,
			ProductID: product_id,
			Quantity:  quantity,      // Set quantity yang diinput
			CreatedAt: created_at,    // Timestamp waktu create
			UpdatedAt: updated_at,    // Timestamp waktu update
		}
	} else {
		// Error lainnya
		return chart, err
	}

	return chart, nil
}





