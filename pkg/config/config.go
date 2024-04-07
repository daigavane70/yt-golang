package config

import (
	"fmt"
	"os"
	"sprint/go/pkg/common/logger"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// init initializes the configuration by loading environment variables from a .env file.
// It logs an error if the .env file is not found or cannot be loaded.
func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file, using default values")
	}
}

// GetPortAndHost retrieves the port and host from environment variables or uses default values.
// If PORT or HOST environment variables are not set, it defaults to port 8080 and host localhost.
func GetPortAndHost() (host string, port string) {
	// Get environment variables or use default values
	port = os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Default port if not specified in .env
	}

	host = os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost" // Default host if not specified in .env
	}
	return
}

// GetGoogleApiKey retrieves the Google API key from environment variables.
// It returns an empty string if API_KEY environment variable is not set.
func GetGoogleApiKey() (apiKey string) {
	apiKey = os.Getenv("API_KEY")
	return
}

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbUsername, dbPassword, dbHost, dbPort, dbName)

	logger.Info("db source name: ", dsn)

	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Error("Error connecting to database, ", err)
	}

	// logger.Success("Successfully connect to the database at port: ", dbPort)
}

func GetDB() *gorm.DB {
	return db
}
