package entities

import (
	"sprint/go/pkg/config"

	"gorm.io/gorm"
)

type Config struct {
	gorm.Model
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (c *Config) GetValueByKey(key string) string {
	var configValue Config
	result := config.GetDB().Where("key = ?", key).First(&configValue)

	if result.Error != nil {
		return ""
	}

	return configValue.Value
}
