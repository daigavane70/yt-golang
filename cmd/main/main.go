package main

import (
	"fmt"
	"net/http"
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/routes"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	routes.RegisterRoutes(r)

	http.Handle("/", r)

	serverUrl := fmt.Sprintf("%s:%s", "localhost", "8080")

	logger.Error(http.ListenAndServe(serverUrl, r))
}
