package controller

import (
	"encoding/json"
	"golang_api/helper"
	"golang_api/kafka"
	"golang_api/middleware"
	"golang_api/model"
	"net/http"
)

func AddToChart(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		helper.JSONError(w, http.StatusUnauthorized, "Pengguna tidak terautentikasi")
		return
	}

	var req model.ChartItem

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helper.JSONError(w, http.StatusBadRequest, "Request tidak valid")
		return
	}

	req.UserID = userID

	message, _ := json.Marshal(req)

	err = kafka.ProduceMessage(message)
	if err != nil {
		helper.JSONError(w, http.StatusInternalServerError, "Gagal mengirim data ke Kafka")
		return
	}

	response := map[string]any{
		"message": "Produk berhasil ditambahkan ke chart & chart_event",
	}

	helper.JSONResponse(w, http.StatusOK, response)
}

func DeleteChartHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		helper.JSONError(w, http.StatusUnauthorized, "Pengguna tidak terautentikasi")
		return
	}

	var req struct {
		ProductID int `json:"product_id"`
	}
	
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	productID := req.ProductID
	

	err = model.DeleteChart(userID, productID)
	if err != nil {
		helper.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.JSONResponse(w, http.StatusOK, map[string]string{"message": "Data dihapus seluruhnya"})
}

type PartialDeleteRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func DeleteChartQuantityHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		helper.JSONError(w, http.StatusUnauthorized, "Pengguna tidak terautentikasi")
		return
	}

	var req PartialDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.JSONError(w, http.StatusBadRequest, "Permintaan tidak valid")
		return
	}

	if req.Quantity <= 0 {
		helper.JSONError(w, http.StatusBadRequest, "Quantity harus lebih dari 0")
		return
	}

	err := model.DeleteChartQuantity(userID, req.ProductID, req.Quantity)
	if err != nil {
		helper.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.JSONResponse(w, http.StatusOK, map[string]string{"message": "Quantity berhasil dikurangi"})
}

