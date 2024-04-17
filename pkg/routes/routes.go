package routes

import (
	"sprint/go/pkg/controllers"

	"github.com/gorilla/mux"
)

var RegisterRoutes = func(router *mux.Router) {
	router.HandleFunc("/ping", controllers.Test).Methods("GET")
	router.HandleFunc("/video", controllers.SearchVideos).Methods("GET")
	router.HandleFunc("/jobs/once", controllers.RunJobOnce).Methods("GET")
	router.HandleFunc("/jobs/start", controllers.StartJob).Methods("GET")
}
