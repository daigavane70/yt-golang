package routes

import (
	"sprint/go/pkg/controllers"

	"github.com/gorilla/mux"
)

var RegisterRoutes = func(router *mux.Router) {
	router.HandleFunc("/ping", controllers.Test).Methods("GET")
	router.HandleFunc("/video", controllers.GetAllVideos).Methods("GET")
}
