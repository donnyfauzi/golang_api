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
	CreatedAt string `json:"created_at"` 
	UpdatedAt string `json:"updated_at"` 
}

func CreateChart(user_id, product_id, quantity int) (ChartItem, error) {
	var chart ChartItem

	now := time.Now().Format("2006-01-02 15:04:05") 
	var existingQuantity int
	var existingID int
	var createdAt string

	err := database.DB.QueryRow("SELECT id, quantity, created_at FROM chart WHERE user_id = ? AND product_id = ?", user_id, product_id).Scan(&existingID, &existingQuantity, &createdAt)

	if err == nil {
		// kalo data sudah ada -> Update quantity
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
			CreatedAt: createdAt, 
			UpdatedAt: now,       
		}

	} else if err == sql.ErrNoRows {
		// kalo data belum ada -> Insert baru
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
		return chart, err
	}

	return chart, nil
}

// Hapus semua chart nerdasarkan ID
func DeleteChart(userID, productID int) error {
	_, err := database.DB.Exec("DELETE FROM chart WHERE user_id = ? AND product_id = ?", userID, productID)
	if err != nil {
		return err
	}
	return CreateChartEvent(userID, productID, 0, "deleted_all")
}

func DeleteChartQuantity(userID, productID, quantityToDelete int) error {
	var currentQty int
	err := database.DB.QueryRow("SELECT quantity FROM chart WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&currentQty)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	if quantityToDelete >= currentQty {
		return DeleteChart(userID, productID)
	}

	// Update quantity
	newQty := currentQty - quantityToDelete
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err = database.DB.Exec(`UPDATE chart SET quantity = ?, updated_at = ? WHERE user_id = ? AND product_id = ?`,
		newQty, now, userID, productID)
	if err != nil {
		return err
	}

	return CreateChartEvent(userID, productID, quantityToDelete, "deleted")
}

func CreateChartEvent(user_id, product_id, quantity int, action string) error {

	validActions := map[string]bool{
		"added":       true,
		"deleted":     true,
		"deleted_all": true,
	}

	if !validActions[action] {
		log.Printf("Action tidak valid untuk chart_event: %s", action)
	}

	now := time.Now().Format("2006-01-02 15:04:05")

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




