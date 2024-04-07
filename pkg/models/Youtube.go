package models

import "time"

type YouTubeVideoList struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []YouTubeVideoItem `json:"items"`
}

type YouTubeVideoItem struct {
	Kind    string              `json:"kind"`
	Etag    string              `json:"etag"`
	ID      YouTubeVideoID      `json:"id"`
	Snippet YouTubeVideoSnippet `json:"snippet"`
}

type YouTubeVideoID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

type YouTubeVideoSnippet struct {
	PublishedAt          time.Time              `json:"publishedAt"`
	ChannelID            string                 `json:"channelId"`
	Title                string                 `json:"title"`
	Description          string                 `json:"description"`
	Thumbnails           YouTubeVideoThumbnails `json:"thumbnails"`
	ChannelTitle         string                 `json:"channelTitle"`
	LiveBroadcastContent string                 `json:"liveBroadcastContent"`
	PublishTime          time.Time              `json:"publishTime"`
}

type YouTubeVideoThumbnails struct {
	Default YouTubeVideoThumbnail `json:"default"`
	Medium  YouTubeVideoThumbnail `json:"medium"`
	High    YouTubeVideoThumbnail `json:"high"`
}

type YouTubeVideoThumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
