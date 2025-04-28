// routes/chartRoutes.go
package routes

import (
	"golang_api/controller"
	"golang_api/middleware"

	"github.com/gorilla/mux"
)

func ChartRoutes(r *mux.Router) {
	r.Handle("/AddToChart", middleware.JwtVerify(controller.AddToChart)).Methods("POST")
}
