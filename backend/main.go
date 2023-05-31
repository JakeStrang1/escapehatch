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
	err := godotenv.Overload() // Will overwrite existing env vars
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Access production secrets
	if os.Getenv("PRODUCTION") == "true" {
		log.Println("PRODUCTION=true")
		log.Println("Loading production secrets from Google Secret Manager...")
		loadProdSecrets()
	} else {
		log.Println("PRODUCTION=false")
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

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, "MONGO_HOST"),
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}
	os.Setenv("MONGO_HOST", string(result.Payload.Data))

	// Build the request.
	accessRequest = &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, "SENDGRID_API_KEY"),
	}

	// Call the API.
	result, err = client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}
	os.Setenv("SENDGRID_API_KEY", string(result.Payload.Data))
}
