package controllers

import (
	"net/http"
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/common/utils"
	"sprint/go/pkg/entities"
	"sprint/go/pkg/models"
	"strconv"
)

// SearchVideos handles the API endpoint for searching videos by keyword.
var SearchVideos = func(w http.ResponseWriter, r *http.Request) {
	// Get search parameters from URL query
	searchQuery := r.URL.Query().Get("searchQuery")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1 // Default page if not provided
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 5 // Default page size if not provided
	}

	logger.Info("[SearchVideos] Searching videos for searchQuery: ", searchQuery, ", page: ", page, ", pageSize: ", pageSize)

	// Call function to search videos by keyword
	videos, metaData, _ := entities.SearchVideosByKeyword(searchQuery, pageSize, page)
	logger.Success("[SearchVideos] returning ", len(videos), " records")
	response := models.CreateCommonSuccessWithMetaDataResponse(videos, metaData)
	utils.SendJSONResponse(w, http.StatusOK, response)
}
