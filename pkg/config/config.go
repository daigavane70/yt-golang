package config

import (
	"fmt"
	"os"
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/constants/env"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	loadEnvFile()
}

func loadEnvFile() {
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file:", err)
	}
}

func GetPortAndHost() (host, port string) {
	host = getEnvOrDefault(env.SERVER_HOST, "localhost")
	port = getEnvOrDefault(env.SERVER_PORT, "8080")
	return
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetApiKey(keyNumber int) string {
	return os.Getenv(fmt.Sprintf("%s%d", env.API_KEY_PREFIX, keyNumber))
}

func ConnectDB() {
	loadEnvFile()

	dbUsername := os.Getenv(env.DB_USERNAME)
	dbPassword := os.Getenv(env.DB_PASSWORD)
	dbHost := os.Getenv(env.DB_HOST)
	dbPort := os.Getenv(env.DB_PORT)
	dbName := os.Getenv(env.DB_NAME)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUsername, dbPassword, dbHost, dbPort, dbName)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		logger.Error("Error connecting to database:", err)
		return
	}
	logger.Success("Successfully connected to the database at port:", dbPort)
}

func GetDB() *gorm.DB {
	return db
}
