package routes

import (
	"sprint/go/pkg/controllers"

	"github.com/gorilla/mux"
)

var RegisterRoutes = func(router *mux.Router) {
	router.HandleFunc("/", controllers.Test).Methods("GET")
}
