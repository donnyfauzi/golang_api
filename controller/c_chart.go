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
