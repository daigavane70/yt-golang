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
	"sprint/go/pkg/entities"
	"sprint/go/pkg/models"
)

var (
	downtimeSync      sync.WaitGroup
	apiKey            string
	usingSecondaryKey bool
)

func init() {
	downtimeSync.Add(1)
	logger.Info("Fetching videos uploaded during downtime")

	latestVideoTime, err := entities.GetLastPublishedVideoTime()
	if err != nil {
		logger.Error("Unable to fetch last video publish time before downtime:", err)
	}
	fetchDataFromYoutube(latestVideoTime, utils.FormatToRFC3339(time.Now()))
	downtimeSync.Done()
}

func useSecondaryKey() {
	secondaryKey, err := entities.GetValueByKey(entities.SecondaryApiKey)
	if err != nil {
		logger.Error("Error fetching secondary key:", err)
		return
	}
	apiKey = secondaryKey
	usingSecondaryKey = true
}

func getLastPublishedValue() string {
	val, err := entities.GetValueByKey(entities.LastFetchedAtKey)
	if err != nil || val == "" {
		return utils.FormatToRFC3339(time.Now().Add(-1 * time.Hour))
	}

	epoch, _ := strconv.ParseInt(val, 10, 64)
	unixTime := utils.FormatToRFC3339(time.Unix(epoch, 0))
	return unixTime
}

func getSearchUrl(publishedAfter, publishedBefore string) string {
	const (
		part         = "snippet"
		maxCounts    = 20
		searchQuery  = "cricket"
		contentType  = "video"
		dateOrder    = "date"
		timeFormat   = "2006-01-02T15:04:05Z"
		apiKeyPrefix = "&key="
	)

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=%s&q=%s&maxResults=%d&order=%s&publishedAfter=%s&type=%s%s%s", part, searchQuery, maxCounts, dateOrder, publishedAfter, contentType, apiKeyPrefix, apiKey)

	if publishedBefore != "" {
		url += fmt.Sprintf("&publishedBefore=%s", publishedBefore)
	}

	return url
}

func fetchDataFromYoutube(publishedAfter, publishedBefore string) {
	client := &http.Client{}
	url := getSearchUrl(publishedAfter, publishedBefore)
	logger.Info("URL:", url)

	res, err := client.Get(url)
	if err != nil {
		logger.Error("Error fetching data from YouTube:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		if usingSecondaryKey {
			logger.Error("Both API keys expired:", res.Status)
			return
		}
		logger.Error("API key expired:", res.Status)
		useSecondaryKey()
		fetchDataFromYoutube(publishedAfter, publishedBefore)
		return
	}

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

	logger.Info("[StartFetchVideosJob] Waiting for downtime sync")
	downtimeSync.Wait()
	logger.Info("[StartFetchVideosJob] Downtime sync completed")

	for {
		publishedAfter := getLastPublishedValue()
		go fetchDataFromYoutube(publishedAfter, "")
		time.Sleep(30 * time.Second)
	}
}
