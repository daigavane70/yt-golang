package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/common/utils"
	"sprint/go/pkg/config"
	"sprint/go/pkg/entities"
	"sprint/go/pkg/models"
)

var (
	stopJobExecution  = false
	preProcessing     sync.WaitGroup
	apiKey            = config.GetApiKey(false)
	usingSecondaryKey bool
	videoIdsMap       = make(map[string]bool)
)

func init() {
	storeTheExistingVideoIds()
}

func fetchDownTimeData() {
	preProcessing.Add(1)
	logger.Info("[fetchDownTimeData] Fetching videos uploaded during downtime")

	latestVideoTime, err := entities.GetLastPublishedVideoTime()

	if err != nil {
		logger.Error("[fetchDownTimeData] Unable to fetch last video publish time before downtime:", err)
	}
	fetchDataFromYoutube(latestVideoTime, time.Now().Unix())
	preProcessing.Done()
}

func storeTheExistingVideoIds() {
	videoIds := entities.GetAllExistingVideoIds()
	logger.Info("[storeTheExistingVideoIds] storing total: ", len(videoIds), " existing video ids")
	for _, id := range videoIds {
		videoIdsMap[id] = true
	}
}

func useSecondaryKey() {
	apiKey = config.GetApiKey(false)
	usingSecondaryKey = true
}

func getLastPublishedValue() int64 {
	defaultTime := time.Now().Add(-1 * time.Hour).Unix()

	// get last published data from config;
	val, configErr := entities.GetValueByKey(entities.LastFetchedAtKey)

	// if error, then fetch most recently published video time
	if configErr != nil || val == "" {
		epoch, videoErr := entities.GetLastPublishedVideoTime()
		// if error while
		if videoErr != nil {
			return defaultTime
		}
		return epoch
	}

	epoch, err := strconv.ParseInt(val, 10, 64)

	if err != nil {
		return defaultTime
	}

	return epoch
}

func fetchDataFromYoutube(publishedAfter int64, publishedBefore int64) {
	client := &http.Client{}

	url := GetSearchUrl(publishedAfter, publishedBefore)

	res, err := client.Get(url)
	if err != nil {
		logger.Error("Error fetching data from YouTube:", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest || res.StatusCode == http.StatusForbidden {
		if usingSecondaryKey {
			logger.Error("Both API keys expired: ", res.Status)
			stopJobExecution = true
			return
		}
		logger.Error("API key expired: ", res.Status)
		useSecondaryKey()
		fetchDataFromYoutube(publishedAfter, publishedBefore)
		return
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("Non-OK status code received from YouTube: ", res.Status)
		return
	}

	var response models.YouTubeVideoList
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		logger.Error("Error decoding JSON response from YouTube: ", err)
		return
	}

	newVideos := response.Items
	logger.Info("Fetched ", len(newVideos), " for interval: ", utils.FormatAsReadableTime(publishedAfter), " to ", utils.FormatAsReadableTime(publishedBefore))

	filteredVideos := filterTheExistingVideoIds(newVideos)
	logger.Info("Filtered ", len(filteredVideos), " videos to store in db.")

	// storing the videos in database
	for _, item := range filteredVideos {
		thumbnails, err := json.Marshal(item.Snippet.Thumbnails)
		if err != nil {
			logger.Error("Unable to marshal thumbnails")
		}

		newVideo := entities.Video{
			Id:          0,
			VideoID:     item.ID.VideoID,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: int(item.Snippet.PublishTime.Unix()),
			Thumbnail:   string(thumbnails),
		}
		videoIdsMap[item.ID.VideoID] = true
		newVideo.CreateVideo()
	}

	// update the last published at value
	if len(filteredVideos) > 0 {
		videoPublishTime := filteredVideos[0].Snippet.PublishTime.UTC().Unix()
		logger.Info("Updating the lastPublishedValue: ", utils.FormatAsReadableTime(videoPublishTime))
		entities.UpdateValue(entities.Config{Key: entities.LastFetchedAtKey, Value: fmt.Sprint(videoPublishTime)})
	}

	logger.Success("Completed job execution for job at ", utils.FormatAsReadableTime(time.Now().Unix()))
}

func StartFetchVideosJob() {
	logger.Info("[StartFetchVideosJob] Waiting for downtime sync")
	fetchDownTimeData()
	preProcessing.Wait()
	logger.Success("[StartFetchVideosJob] Downtime sync completed")

	for {
		if stopJobExecution {
			break
		}
		publishedAfter := getLastPublishedValue()
		go fetchDataFromYoutube(publishedAfter, 0)
		time.Sleep(60 * time.Second)
	}
}

func RunTheJobOnce() {
	publishedAfter := getLastPublishedValue()
	fetchDataFromYoutube(publishedAfter, 0)
}
