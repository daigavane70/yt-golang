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
	downtimeSync      sync.WaitGroup
	apiKey            string
	usingSecondaryKey bool
)

func fetchDownTimeData() {
	downtimeSync.Add(1)
	logger.Info("[fetchDownTimeData] Fetching videos uploaded during downtime")

	latestVideoTime, err := entities.GetLastPublishedVideoTime()
	if err != nil {
		logger.Error("[fetchDownTimeData] Unable to fetch last video publish time before downtime:", err)
	}
	fetchDataFromYoutube(latestVideoTime, utils.FormatToRFC3339(time.Now()), nil)
	downtimeSync.Done()
}

func init() {
	apiKey = config.GetApiKey(false)
	fetchDownTimeData()
}

func useSecondaryKey() {
	apiKey = config.GetApiKey(false)
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

func fetchDataFromYoutube(publishedAfter, publishedBefore string, ch chan string) {
	client := &http.Client{}

	url := GetSearchUrl(publishedAfter, publishedBefore)

	res, err := client.Get(url)
	if err != nil {
		logger.Error("Error fetching data from YouTube:", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		if usingSecondaryKey {
			logger.Error("Both API keys expired:", res.Status)
		}
		logger.Error("API key expired:", res.Status)
		useSecondaryKey()
		fetchDataFromYoutube(publishedAfter, publishedBefore, ch)
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("Non-OK status code received from YouTube: ", res.Status)
	}

	var response models.YouTubeVideoList
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		logger.Error("Error decoding JSON response from YouTube: ", err)
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
	logger.Info("[StartFetchVideosJob] Waiting for downtime sync")
	downtimeSync.Wait()
	logger.Info("[StartFetchVideosJob] Downtime sync completed")

	ch := make(chan string)

	for {
		publishedAfter := getLastPublishedValue()
		go fetchDataFromYoutube(publishedAfter, "", ch)
		time.Sleep(60 * time.Second)
	}
}
