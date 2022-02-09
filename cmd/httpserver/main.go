package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	firebase "firebase.google.com/go"
	fiber "github.com/gofiber/fiber/v2"
	dialgoflowsrv "github.com/snehalyelmati/telegram-bot-golang/internal/core/services/dialogflowsrv"
	"github.com/snehalyelmati/telegram-bot-golang/internal/core/services/telegramsrv"
	"github.com/snehalyelmati/telegram-bot-golang/internal/handlers/telegramhdl"
	"github.com/snehalyelmati/telegram-bot-golang/internal/respositories/transcriptsrepo"
)

func main() {
	logger := log.New(os.Stdout, "go-api:", log.LstdFlags)
	PROJECT_ID := os.Getenv("PROJECT_ID")
	LANGUAGE := os.Getenv("LANGUAGE")
	TOKEN := os.Getenv("TOKEN")
	SERVER_URL := os.Getenv("SERVER_URL")
	TELEGRAM_API := "https://api.telegram.org/bot" + TOKEN
	URI := "/webhook" + TOKEN
	WEBHOOK_URL := SERVER_URL + URI

	// initialize firebase
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: PROJECT_ID}
	firebaseApp, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	transcriptsFirestore := transcriptsrepo.NewFirestoreRepo(firebaseApp, logger)
	dialogflowService := dialgoflowsrv.New(logger)
	telegramService := telegramsrv.New(logger, dialogflowService)
	telegramHandler := telegramhdl.NewHTTPHandler(logger, telegramService, transcriptsFirestore, PROJECT_ID, LANGUAGE, TELEGRAM_API)

	res, err := http.Get(TELEGRAM_API + "/setWebhook?url=" + WEBHOOK_URL)
	if err != nil {
		logger.Println(err)
	}
	logger.Println("Response from the initial webhook:", res)

	app := fiber.New()
	app.Get("/healthcheck", telegramHandler.HealthCheck)

	app.Post(URI, telegramHandler.SendMessage)

	PORT := 3000
	err = app.Listen(":" + strconv.Itoa(PORT))
	if err != nil {
		logger.Fatal("Cannot run server on port:", PORT)
	}
	logger.Println("Server running on port:", PORT)
}
