package entities

import (
	"sprint/go/pkg/config"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	VideoID     string `json:"videoId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishedAt int    `json:"publishedAt"`
}

func (b *Video) CreateVideo() *Video {
	config.GetDB().Create(&b)
	return b
}

func GetAllVideos() []Video {
	var Videos []Video
	config.GetDB().Find(&Videos)
	return Videos
}
