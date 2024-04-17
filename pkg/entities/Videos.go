package entities

import (
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/config"
	"sprint/go/pkg/models"
	"strings"

	"gorm.io/gorm"
)

var videoDB *gorm.DB

type Video struct {
	Id          int    `json:"id"`
	VideoID     string `json:"videoId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	PublishedAt int    `json:"publishedAt"`
}

func init() {
	if config.GetDB() == nil {
		config.ConnectDB()
	}
	videoDB = config.GetDB()
}

func (b *Video) CreateVideo() *Video {
	videoDB.Create(&b)
	return b
}

func CreateVideos(videos []Video) error {
	if len(videos) == 0 {
		return nil // No videos to create
	}
	result := videoDB.Create(&videos)

	if result.Error != nil {
		logger.Error("[CreateVideos] Error creating videos: ", result.Error)
		return result.Error
	}
	return nil
}

func GetAllVideos() []Video {
	var videos []Video
	videoDB.Where("title = ?", "Hello").Find(&videos)
	return videos
}

func SearchVideosByKeyword(keyword string, pageSize int, page int) ([]Video, models.MetaData, error) {
	var videos []Video
	offset := (page - 1) * pageSize
	var totalCount int64

	keywords := strings.Fields(keyword)

	query := videoDB

	// Loop through each word and adding a condition to search for it in title or description
	for _, word := range keywords {
		query = query.Where("videos.title LIKE ?", "%"+word+"%").Or("videos.description LIKE ?", "%"+word+"%")
	}

	// Get total count without limit
	query.Find(&videos).Count(&totalCount)

	// Apply pagination and fetch the videos
	result := query.Offset(offset).Limit(pageSize).Find(&videos)
	if result.Error != nil {
		return nil, models.MetaData{}, result.Error
	}

	metaData := models.MetaData{TotalResults: int(totalCount), PageSize: pageSize, Page: page}

	return videos, metaData, nil
}

func GetAllExistingVideoIds() []string {
	var videos []Video

	result := videoDB.Select("video_id").Find(&videos)

	if result.Error != nil {
		logger.Error("[GetAllExistingVideoIds] Unable to fetch the existing videoIds, error: ", result.Error)
		return []string{}
	}

	videoIds := []string{}

	for _, video := range videos {
		videoIds = append(videoIds, video.VideoID)
	}

	return videoIds
}

func GetLastPublishedVideoTime() (int64, error) {
	var latestTime int
	var err error

	// fetch the maximum value of PublishedAt
	result := videoDB.Table("videos").Select("MAX(published_at)").Scan(&latestTime)

	if result.Error != nil {
		// Handle the error if the query fails
		err = result.Error
		return 0, err
	}

	return int64(latestTime), err
}
