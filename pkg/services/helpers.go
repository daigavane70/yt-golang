package services

import "fmt"

func GetSearchUrl(publishedAfter, publishedBefore string) string {
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
