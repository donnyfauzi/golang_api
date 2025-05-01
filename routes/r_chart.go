// routes/chartRoutes.go
package routes

import (
	"golang_api/controller"
	"golang_api/middleware"

	"github.com/gorilla/mux"
)

func ChartRoutes(r *mux.Router) {
	r.Handle("/AddToChart", middleware.JwtVerify(controller.AddToChart)).Methods("POST")
	r.HandleFunc("/chart/delete", middleware.JwtVerify(controller.DeleteChartHandler)).Methods("DELETE") 
	r.HandleFunc("/chart/delete-quantity", middleware.JwtVerify(controller.DeleteChartQuantityHandler)).Methods("POST") 
}
