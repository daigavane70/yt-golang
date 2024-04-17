package services

import (
	"fmt"
	"sprint/go/pkg/common/utils"
	"sprint/go/pkg/models"
)

func GetSearchUrl(publishedAfter, publishedBefore int64) string {
	const (
		part         = "snippet"
		maxCounts    = 20
		searchQuery  = "cricket"
		contentType  = "video"
		dateOrder    = "date"
		timeFormat   = "2006-01-02T15:04:05Z"
		apiKeyPrefix = "&key="
	)

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=%s&q=%s&maxResults=%d&order=%s&publishedAfter=%s&type=%s%s%s", part, searchQuery, maxCounts, dateOrder, utils.EpochToRFC3339(publishedAfter), contentType, apiKeyPrefix, apiKey)

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
