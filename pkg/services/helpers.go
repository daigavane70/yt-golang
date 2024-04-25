package services

import (
	"fmt"
	"sprint/go/pkg/common/utils"
	"sprint/go/pkg/models"
)

func GetSearchUrl(publishedAfter int64, publishedBefore int64, pageToken string) string {
	const (
		part        = "snippet"
		maxCounts   = 50
		searchQuery = "ipl"
		contentType = "video"
		dateOrder   = "date"
		timeFormat  = "2006-01-02T15:04:05Z"
	)

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=%s&q=%s&maxResults=%d&order=%s&publishedAfter=%s&type=%s&key=%s&pageToken=%s", part, searchQuery, maxCounts, dateOrder, utils.EpochToRFC3339(publishedAfter), contentType, apiKey, pageToken)

	if publishedBefore != 0 {
		url += fmt.Sprintf("&publishedBefore=%s", utils.EpochToRFC3339(publishedBefore))
	}

	return url
}

func filterTheExistingVideoIds(newVideos []models.YouTubeVideoItem) []models.YouTubeVideoItem {
	var filteredVideos []models.YouTubeVideoItem
	for _, video := range newVideos {
		if !videoIdsMap[video.ID.VideoID] {
			filteredVideos = append(filteredVideos, video)
		}
	}
	return filteredVideos
}
