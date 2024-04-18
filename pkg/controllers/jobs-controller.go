package controllers

import (
	"net/http"
	"sprint/go/pkg/common/utils"
	"sprint/go/pkg/models"
	"sprint/go/pkg/services"
)

var RunJobOnce = func(w http.ResponseWriter, r *http.Request) {
	services.RunTheJobOnce()
	utils.SendJSONResponse(w, http.StatusAccepted, models.CreateCommonSuccessResponse("Triggered job for 1 execution"))
	// go services.StartFetchVideosJob()
}

var StartJob = func(w http.ResponseWriter, r *http.Request) {
	go services.StartFetchVideosJob()
	utils.SendJSONResponse(w, http.StatusAccepted, models.CreateCommonSuccessResponse("Triggered job"))
}
