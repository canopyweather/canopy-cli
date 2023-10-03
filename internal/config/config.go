package config

import (
	"log"
	"os"
)

type Config struct {
	ApiKey string
	Url    string
}

var config *Config

func Setup() {
	canopyKey := os.Getenv("CANOPY_API_KEY")

	var environment string

	if os.Getenv("CANOPY_ENVIRONMENT") == "" {
		environment = "production"
	} else {
		environment = os.Getenv("CANOPY_ENVIRONMENT")

		switch environment {
		case "staging":
			fallthrough
		case "production":
			break
		default:
			log.Fatal("Invalid environment. Can be either 'production' or 'staging'")
		}
	}

	url := ""

	if environment == "production" {
		url = "https://api.canopyweather.com"
	} else {
		url = "https://api.canopyweather.net"
	}

	config = &Config{
		ApiKey: canopyKey,
		Url:    url,
	}
}

func GetConfig() Config {
	return *config
}
