package main

import (
	"context"
	"fmt"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
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
	// Load env vars
	if os.Getenv("PRODUCTION") == "true" {
		log.Println("PRODUCTION=true")
		log.Println("Loading production secrets from Google Secret Manager...")
		loadProdSecrets()
		// Non-secret env vars are already loaded from app.yaml file
	} else {
		log.Println("PRODUCTION=false")
		log.Println("Loading env vars from .env file...")
		err := godotenv.Overload() // Will overwrite existing env vars
		if err != nil {
			log.Println("Warning: Error loading .env file")
		}
	}

	config := app.Config{
		MongoHost:         os.Getenv("MONGO_HOST"),
		MongoDatabaseName: os.Getenv("MONGO_DB_NAME"),
		CORSAllowOrigin:   os.Getenv("ORIGIN"),
		UseSendGrid:       os.Getenv("USE_SENDGRID"),
		SendGridAPIKey:    os.Getenv("SENDGRID_API_KEY"),
		SendGridFromEmail: os.Getenv("SENDGRID_FROM_EMAIL"),
		UseGCS:            os.Getenv("USE_GCS"),
		GCSBucketName:     os.Getenv("GCS_BUCKET_NAME"),
		StaticURLRoot:     os.Getenv("STATIC_URL_ROOT"),
		UseAtlasSearch:    os.Getenv("USE_ATLAS_SEARCH"),
	}
	mainApp := app.NewApp(config)
	defer mainApp.Close()

	mainApp.Run()
}

// loadProdSecrets uses Google Cloud Application Default Credentials (ADC) to authenticate with the Google Secret Manager and
// load secrets into environment variables
// Source: https://cloud.google.com/secret-manager/docs/reference/libraries#client-libraries-install-go
func loadProdSecrets() {
	projectID := os.Getenv("PROJECT_ID")

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	// MONGO_HOST
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, "MONGO_HOST"),
	}
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}
	os.Setenv("MONGO_HOST", string(result.Payload.Data))

	// SENDGRID_API_KEY
	accessRequest = &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, "SENDGRID_API_KEY"),
	}
	result, err = client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}
	os.Setenv("SENDGRID_API_KEY", string(result.Payload.Data))
}
