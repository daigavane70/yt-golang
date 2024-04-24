package entities

import (
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/config"

	"gorm.io/gorm"
)

var (
	configDb *gorm.DB
)

const (
	LastFetchedAtKey = "lastFetchedAt"
)

type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func init() {
	if config.GetDB() == nil {
		config.ConnectDB()
	}
	configDb = config.GetDB()
}

func GetValueByKey(key string) (string, error) {
	var configValue Config

	result := configDb.Where("configs.key = ?", key).First(&configValue)

	if result.Error != nil {
		logger.Error("[GetValueByKey] failed to fetch the config for key: ", key, ", error: ", result.Error)
		return "", result.Error
	}

	return configValue.Value, nil
}

func UpdateValue(updatedConfig Config) Config {
	saveResult := configDb.Where("configs.key = ?", updatedConfig.Key).Save(&updatedConfig)
	if saveResult.Error != nil {
		logger.Error("[UpdateValueByKey] Error while updating config key: ", updatedConfig.Key, ", error: ", saveResult.Error)
	}
	return updatedConfig
}
