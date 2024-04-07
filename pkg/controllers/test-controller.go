package controllers

import (
	"encoding/json"
	"net/http"
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/models"
)

var Test = func(res http.ResponseWriter, req *http.Request) {
	logger.Success("Testing the api")
	response := models.CreateCommonSuccessResponse("Status ok")
	jsonResponse, _ := json.Marshal(response)
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonResponse)
}
