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
	"sprint/go/pkg/constants"
	"sprint/go/pkg/entities"
	"sprint/go/pkg/models"
)

var (
	stopJobExecution = false
	preProcessing    sync.WaitGroup
	keyNumber        = 1
	apiKey           = config.GetApiKey(keyNumber)
	videoIdsMap      = make(map[string]bool)
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
	apiKey = config.GetApiKey(keyNumber + 1)
	keyNumber = keyNumber + 1
}

func getLastPublishedValue() int64 {
	defaultTime := time.Now().Add(-24 * time.Hour).Unix()

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

func getDataFromYouTube(publishedAfter int64, publishedBefore int64, pageToken string, videoList *[]models.YouTubeVideoItem) {
	url := GetSearchUrl(publishedAfter, publishedBefore, pageToken)

	logger.Success("URL: ", url)

	res, err := http.Get(url)
	if err != nil {
		logger.Error("Error fetching data from YouTube:", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		if keyNumber == constants.MAX_API_KEYS {
			logger.Error("All API keys expired: ", res.Status)
			stopJobExecution = true
			return
		}
		logger.Error("API key expired: ", apiKey, res.Status)
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

	if len(response.Items) == 0 {
		return
	}

	*videoList = append(*videoList, response.Items...)

	if response.NextPageToken != "" {
		logger.Info("Making the call for page: ", response.NextPageToken)
		getDataFromYouTube(publishedAfter, publishedBefore, response.NextPageToken, videoList)
	}
}

func fetchDataFromYoutube(publishedAfter int64, publishedBefore int64) {

	var newVideos []models.YouTubeVideoItem

	getDataFromYouTube(publishedAfter, publishedBefore, "", &newVideos)

	logger.Info("Fetched ", len(newVideos), " for interval: ", utils.FormatAsReadableTime(publishedAfter), " to ", utils.FormatAsReadableTime(publishedBefore))

	filteredVideos := filterTheExistingVideoIds(newVideos)
	logger.Info("Filtered ", len(filteredVideos), " videos to store in db.")

	var videos []entities.Video

	// storing the videos in database
	for _, item := range filteredVideos {
		thumbnailString, _ := json.Marshal(item.Snippet.Thumbnails)

		newVideo := entities.Video{
			VideoID:     item.ID.VideoID,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: int(item.Snippet.PublishTime.Unix()),
			Thumbnail:   thumbnailString,
		}

		videoIdsMap[item.ID.VideoID] = true
		videos = append(videos, newVideo)
	}

	entities.CreateVideos(videos)

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
