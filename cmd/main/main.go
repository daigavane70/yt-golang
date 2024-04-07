package main

import (
	"fmt"
	"net/http"
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/config"
	"sprint/go/pkg/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	routes.RegisterRoutes(r)

	http.Handle("/", r)

	host, port := config.GetPortAndHost()
	serverUrl := fmt.Sprintf("%s:%s", host, port)

	config.ConnectDB()

	logger.Error(http.ListenAndServe(serverUrl, r))
}
