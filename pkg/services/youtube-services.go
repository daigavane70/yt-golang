package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/common/utils"
	"sprint/go/pkg/config"
	"sprint/go/pkg/entities"
	"sprint/go/pkg/models"
	"strconv"
	"sync"
	"time"
)

var downtimeSync sync.WaitGroup

// fetch all the videos that were uploaded during downtime
func init() {
	downtimeSync.Add(1)
	logger.Info("Fetching the videos that were uploaded during downtime")
	latestVideoTime, err := entities.GetLastPublishedVideoTime()
	if err != nil {
		logger.Error("Unable to fetch the last video publish time before downtime, error: ", err)
	}
	fetchDataFromYoutube(latestVideoTime, utils.FormatToRFC3339(time.Now()))
	downtimeSync.Done()
}

func getLastPublishedValue() string {
	// Get the last fetched value from the database
	val, err := entities.GetValueByKey(entities.LastFetchedAtKey)
	if err != nil || val == "" {
		// If there's an error or the value is empty, return the current time minus 1 hour in RFC3339 format
		return utils.FormatToRFC3339(time.Now().Add(-1 * time.Hour))
	}

	// Convert the fetched value (epoch time) to int64
	epoch, _ := strconv.ParseInt(val, 10, 64)

	// Convert the epoch time to UTC time string
	unixTime := utils.FormatToRFC3339(time.Unix(epoch, 0))
	return unixTime
}

func getSearchUrl(publishedAfter string, publishedBefore string) string {
	const (
		part         = "snippet"
		maxCounts    = 20
		searchQuery  = "cricket"
		contentType  = "video"
		dateOrder    = "date"
		timeFormat   = "2006-01-02T15:04:05Z"
		apiKeyPrefix = "&key="
	)

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=%s&q=%s&maxResults=%d&order=%s&publishedAfter=%s&type=%s%s%s", part, searchQuery, maxCounts, dateOrder, publishedAfter, contentType, apiKeyPrefix, config.GetGoogleApiKey())

	if publishedBefore != "" {
		url = fmt.Sprintf("%s&publishedBefore=%s", url, publishedBefore)
	}

	return url
}

func fetchDataFromYoutube(publishedAfter string, publishedBefore string) {

	client := &http.Client{}

	url := getSearchUrl(publishedAfter, publishedBefore)

	logger.Info("url: ", url)

	res, err := client.Get(url)
	if err != nil {
		logger.Error("Error fetching data from YouTube:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Error("Non-OK status code received from YouTube:", res.Status)
		return
	}

	var response models.YouTubeVideoList
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		logger.Error("Error decoding JSON response from YouTube:", err)
		return
	}

	items := response.Items

	logger.Info("Fetched total videos:", len(items))

	for _, item := range items {
		newVideo := entities.Video{
			VideoID:     item.ID.VideoID,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: int(item.Snippet.PublishTime.Unix()),
		}
		newVideo.CreateVideo()
	}

	if len(items) > 0 {
		newPublishedAtValue := items[0].Snippet.PublishTime.UTC().Unix()
		updatedLastPublishConfig := entities.Config{Key: entities.LastFetchedAtKey, Value: fmt.Sprint(newPublishedAtValue)}
		entities.UpdateValue(updatedLastPublishConfig)
	}
}

func StartFetchVideosJob() {
	defer logger.Info("[StartFetchVideosJob] Completed")

	// Wait for the downtime sync to finish
	logger.Info("[StartFetchVideosJob] Waiting for downtime sync")
	downtimeSync.Wait()
	logger.Info("[StartFetchVideosJob] Downtime sync completed")

	for {
		publishedAfter := getLastPublishedValue()
		go fetchDataFromYoutube(publishedAfter, "")
		time.Sleep(30 * time.Second)
	}
}
