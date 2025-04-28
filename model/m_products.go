package model

import (
	"golang_api/database"
)

type Product struct {
	Id				int    		`json:"id"`
	Name			string 		`json:"name"`
	Price			float64 	`json:"price"`	
	Stock 			int 		`json:"stock"`
	Description 	string 		`json:"description"`
	Image_url 		string 		`json:"image_url"`
	CreatedAt   	string  	`json:"created_at"`
	UpdatedAt   	string  	`json:"updated_at"`
}

func GetAllProducts() ([]Product, error) {
	var products []Product

	query := "SELECT * FROM products"
	rows, err := database.DB.Query(query)

	if err != nil {
		return products, err
	}
	
	defer rows.Close()

	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.Description,
			&product.Image_url,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return products, err
		}
		products = append(products, product)
	}

	return products, nil
}

