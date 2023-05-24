package main

import (
	"log"
	"os"

	"github.com/JakeStrang1/saas-template/app"
	"github.com/joho/godotenv"
)

/****************************************************************************************
 * main.go
 *
 * This file is intended to:
 * - Gather configuration inputs
 * - Construct the app
 * - Run the app and manage app lifecycle
 ****************************************************************************************/

func main() {
	err := godotenv.Overload() // Will overwrite existing env vars
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	config := app.Config{
		MongoHost:         os.Getenv("MONGO_HOST"),
		MongoDatabaseName: os.Getenv("MONGO_DB_NAME"),
		CORSAllowOrigin:   os.Getenv("ORIGIN"),
		UseSendGrid:       os.Getenv("USE_SENDGRID"),
		SendGridAPIKey:    os.Getenv("SENDGRID_API_KEY"),
		SendGridFromEmail: os.Getenv("SENDGRID_FROM_EMAIL"),
	}
	mainApp := app.NewApp(config)
	defer mainApp.Close()

	mainApp.Run()
}
