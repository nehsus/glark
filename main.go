package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Debug().Msg("Error loading .env file")
	}
	appID = os.Getenv("APP_ID")
	appSecret = os.Getenv("APP_SECRET")
	grafanaChatName = os.Getenv("GRAFANA_CHAT_NAME")
	mode := os.Getenv("ENVIRONMENT")
	if mode == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func main() {
	lambda.Start(handleRequest)
}
