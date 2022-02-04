package main

import (
	"log"
	"os"

	fiber "github.com/gofiber/fiber/v2"
	dialgoflowsrv "github.com/snehalyelmati/telegram-bot-golang/internal/core/services/dialogflowsrv"
	"github.com/snehalyelmati/telegram-bot-golang/internal/core/services/telegramsrv"
	"github.com/snehalyelmati/telegram-bot-golang/internal/handlers/telegramhdl"
)

func main() {
	logger := log.New(os.Stdout, "go-api:", log.LstdFlags)
	PROJECT_ID := os.Getenv("PROJECT_ID")
	LANGUAGE := os.Getenv("LANGUAGE")

	dialogflowService := dialgoflowsrv.New(logger)
	telegramService := telegramsrv.New(logger, dialogflowService)
	telegramHandler := telegramhdl.NewHTTPHandler(logger, telegramService, PROJECT_ID, LANGUAGE)

	app := fiber.New()
	app.Get("/healthcheck", telegramHandler.HealthCheck)

	TOKEN := os.Getenv("TOKEN")
	URI := "/webhook" + TOKEN
	app.Post(URI, telegramHandler.SendMessage)

	app.Listen(":3000")
	logger.Println("Server running on port 3000")
}
