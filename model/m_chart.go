package model

import (
	"database/sql"
	"golang_api/database"
	"log"
	"time"
)

type ChartItem struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	CreatedAt string `json:"created_at"` // ubah jadi string
	UpdatedAt string `json:"updated_at"` // ubah jadi string
}

func CreateChart(user_id, product_id, quantity int) (ChartItem, error) {
	var chart ChartItem

	now := time.Now().Format("2006-01-02 15:04:05") // format waktu standar
	var existingQuantity int
	var existingID int
	var createdAt string

	err := database.DB.QueryRow("SELECT id, quantity, created_at FROM chart WHERE user_id = ? AND product_id = ?", user_id, product_id).Scan(&existingID, &existingQuantity, &createdAt)

	if err == nil {
		// Data sudah ada -> Update quantity
		newQuantity := existingQuantity + quantity

		log.Printf("Updating chart: user_id=%d, product_id=%d, old_quantity=%d, new_quantity=%d", user_id, product_id, existingQuantity, newQuantity)

		_, err = database.DB.Exec(`
			UPDATE chart 
			SET quantity = ?, updated_at = ? 
			WHERE id = ?
		`, newQuantity, now, existingID)

		if err != nil {
			return chart, err
		}

		chart = ChartItem{
			ID:        existingID,
			UserID:    user_id,
			ProductID: product_id,
			Quantity:  newQuantity,
			CreatedAt: createdAt, // pakai createdAt lama
			UpdatedAt: now,       // update updated_at
		}

	} else if err == sql.ErrNoRows {
		// Data belum ada -> Insert baru
		query := "INSERT INTO chart (user_id, product_id, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
		result, err := database.DB.Exec(query, user_id, product_id, quantity, now, now)
		if err != nil {
			return chart, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return chart, err
		}

		chart = ChartItem{
			ID:        int(id),
			UserID:    user_id,
			ProductID: product_id,
			Quantity:  quantity,
			CreatedAt: now,
			UpdatedAt: now,
		}

	} else {
		// Error lain (bukan no rows)
		return chart, err
	}

	return chart, nil
}

// CreateChartEvent menyimpan event chart ke tabel chart_event
func CreateChartEvent(user_id, product_id, quantity int, action string) error {
	now := time.Now().Format("2006-01-02 15:04:05") // format waktu standar

	query := `
		INSERT INTO chart_event (user_id, product_id, quantity, action, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := database.DB.Exec(query, user_id, product_id, quantity, action, now)
	if err != nil {
		log.Printf("Gagal menyimpan chart_event: %v", err)
		return err
	}

	log.Printf("Berhasil menyimpan chart_event untuk user_id=%d, product_id=%d, action=%s", user_id, product_id, action)
	return nil
}

