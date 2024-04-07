package controllers

import (
	"encoding/json"
	"net/http"
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/entities"
	"sprint/go/pkg/models"
)

var NewVideo entities.Video

var GetAllVideos = func(w http.ResponseWriter, r *http.Request) {
	logger.Info("Fetching all the videos")
	videos := entities.GetAllVideos()
	response := models.CreateCommonSuccessResponse(videos)
	res, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
