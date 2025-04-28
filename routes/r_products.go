package routes

import (
	"golang_api/controller"
	"golang_api/middleware"

	"github.com/gorilla/mux"
)

func ProductsRoutes(r *mux.Router) {
	r.Handle("/GetAllProducts", middleware.JwtVerify(controller.GetProducts)).Methods("GET")
}