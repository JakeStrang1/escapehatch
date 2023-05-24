package app

import (
	"log"

	"github.com/JakeStrang1/saas-template/db"
	"github.com/JakeStrang1/saas-template/email"
	"github.com/gin-gonic/gin"
)

/****************************************************************************************
 * app.go
 *
 * This file is intended to:
 * - Instantiate app dependencies (ideally as singletons)
 * - Create an API router
 * - Define any start up and shutdown procedures
 ****************************************************************************************/

type App struct {
	router *gin.Engine
}

type Config struct {
	MongoHost         string
	MongoDatabaseName string
	CORSAllowOrigin   string
	UseSendGrid       string
	SendGridAPIKey    string
	SendGridFromEmail string
}

func NewApp(config Config) App {
	err := db.Setup(config.MongoHost, config.MongoDatabaseName)
	if err != nil {
		panic(err)
	}
	log.Printf("Using DB: %s/%s\n", config.MongoHost, config.MongoDatabaseName)

	if config.UseSendGrid == "true" {
		email.SetupSendGridMailer(email.SendGridConfig{
			APIKey:      config.SendGridAPIKey,
			FromAddress: config.SendGridFromEmail,
		})
		log.Printf("Email is enabled from %s\n", config.SendGridFromEmail)
	} else {
		email.SetupMockMailer() // Mock mailer implementation
		log.Printf("Email is disabled\n")
	}

	return App{
		router: Router(config),
	}
}

func (a *App) Run() {
	a.router.Run()
}

func (a *App) Router() *gin.Engine {
	return a.router
}

func (a *App) Close() {
	db.Close()
}
