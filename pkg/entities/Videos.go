package entities

import (
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/common/utils"
	"sprint/go/pkg/config"
	"sprint/go/pkg/models"

	"gorm.io/gorm"
)

var videoDB *gorm.DB

type Video struct {
	VideoID     string `json:"videoId"`
	Title       string `json:"title"`
	Description string `json:"description"`
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

	query := videoDB.Where("videos.title LIKE ?", "%"+keyword+"%").Or("videos.description LIKE ?", "%"+keyword+"%")

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

func GetLastPublishedVideoTime() (string, error) {
	var latestTime int
	var err error

	// Fetch the maximum value of PublishedAt using GORM
	result := videoDB.Table("videos").Select("MAX(published_at)").Scan(&latestTime)

	logger.Info("result", result)

	if result.Error != nil {
		// Handle the error if the query fails
		err = result.Error
		return "", err
	}

	return utils.EpochToUTC(int64(latestTime)), err
}
