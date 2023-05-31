package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JakeStrang1/escapehatch/app"
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

	fmt.Printf("Envs: MONGO_DB_NAME:%s, ORIGIN:%s, USE_SENDGRID:%s, SENDGRID_FROM_EMAIL:%s\n", os.Getenv("MONGO_DB_NAME"), os.Getenv("ORIGIN"), os.Getenv("USE_SENDGRID"), os.Getenv("SENDGRID_FROM_EMAIL"))
	fmt.Printf("Secrets: MONGO_HOST:%.5s..., SENDGRID_API_KEY:%.5s...\n", os.Getenv("MONGO_HOST"), os.Getenv("SENDGRID_API_KEY"))

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
