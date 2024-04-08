package config

import (
	"fmt"
	"os"
	"sprint/go/pkg/common/logger"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	host = getEnvOrDefault("SERVER_HOST", "localhost")
	port = getEnvOrDefault("SERVER_PORT", "8080")
	return
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetGoogleApiKey() string {
	return os.Getenv("API_KEY")
}

func ConnectDB() {
	loadEnvFile()

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUsername, dbPassword, dbHost, dbPort, dbName)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Error connecting to database:", err)
		return
	}
	logger.Success("Successfully connected to the database at port:", dbPort)
}

func GetDB() *gorm.DB {
	return db
}
