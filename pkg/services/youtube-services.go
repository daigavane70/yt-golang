package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/config"
	"sprint/go/pkg/entities"
	"sprint/go/pkg/models"
	"strconv"
	"sync"
	"time"
)

var downtimeSync sync.WaitGroup

func init() {
	downtimeSync.Add(1)
	logger.Info("Fetching the videos that were uploaded during downtime")
	downtimeSync.Done()
}

func getLastPublishedValue() string {
	// Get the last fetched value from the database
	val, err := entities.GetValueByKey(entities.LastFetchedAtKey)
	if err != nil || val == "" {
		// If there's an error or the value is empty, return the current time minus 1 hour in RFC3339 format
		return time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	}

	// Convert the fetched value (epoch time) to int64
	epoch, _ := strconv.ParseInt(val, 10, 64)

	// Convert the epoch time to UTC time string
	unixTime := time.Unix(epoch, 0).UTC().Format(time.RFC3339)
	return unixTime
}

func getSearchUrl(publishedAfter string, publishedBefore *string) string {
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

	if publishedBefore != nil {
		url = fmt.Sprintf("%s&publishedBefore=%s", url, *publishedBefore)
	}

	return url
}

func fetchDataFromYoutube() {
	client := &http.Client{}

	publishedAfter := getLastPublishedValue()

	url := getSearchUrl(publishedAfter, nil)

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

	if len(items) > 0 {
		newPublishedAtValue := items[0].Snippet.PublishTime.UTC().Unix()
		updatedLastPublishConfig := entities.Config{Key: entities.LastFetchedAtKey, Value: fmt.Sprint(newPublishedAtValue)}
		entities.UpdateValueByKey(updatedLastPublishConfig)
	}

	for _, item := range items {
		newVideo := entities.Video{
			VideoID:     item.ID.VideoID,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: int(item.Snippet.PublishTime.Unix()),
		}
		newVideo.CreateVideo()
	}
}

func StartFetchVideosJob() {
	defer logger.Info("[StartFetchVideosJob] Completed")

	// Wait for the downtime sync to finish
	logger.Info("[StartFetchVideosJob] Waiting for downtime sync")
	downtimeSync.Wait()
	logger.Info("[StartFetchVideosJob] Downtime sync completed")

	for {
		go fetchDataFromYoutube()
		time.Sleep(30 * time.Second)
	}
}
