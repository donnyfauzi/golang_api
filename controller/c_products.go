package controller

import (
	"golang_api/helper"
	"golang_api/model"
	"net/http"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := model.GetAllProducts()

	if err != nil {
		helper.JSONError(w, http.StatusInternalServerError, "Gagal menampilkan produk")
		return
	}
	
	response := map[string]any{
		"message": "Berhasil menampilkan produk",
		"products": products,
	}

	helper.JSONResponse(w, http.StatusOK, response)
}