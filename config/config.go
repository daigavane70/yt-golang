package config

import (
	"os"
	"sprint/go/common/logger"

	"github.com/joho/godotenv"
)

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
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified in .env
	}

	host = os.Getenv("HOST")
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
