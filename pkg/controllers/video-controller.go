package controllers

import (
	"encoding/json"
	"net/http"
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/entities"
	"sprint/go/pkg/models"
	"strconv"
	"time"
)

var NewVideo entities.Video

var SearchVideos = func(w http.ResponseWriter, r *http.Request) {
	searchKeyword := r.URL.Query().Get("searchKeyword")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 5
	}
	if searchKeyword == "" {
		// If searchKeyword is empty, return a bad request response
		response := models.CreateCommonErrorResponse("Search keyword is required")
		res, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	logger.Info("Fetching all the videos")
	videos, metaData, _ := entities.SearchVideosByKeyword(searchKeyword, pageSize, page)
	response := models.CreateCommonSuccessWithMetaDataResponse(videos, metaData)
	res, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

var CreateVideo = func(w http.ResponseWriter, r *http.Request) {
	logger.Info("Creating a new video")
	newVideo := entities.Video{VideoID: "hello", Title: "Ipl20202", Description: "Sample description", PublishedAt: int(time.Now().UTC().Unix())}
	newVideo = *newVideo.CreateVideo()
	response := models.CreateCommonSuccessResponse(newVideo)
	res, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
