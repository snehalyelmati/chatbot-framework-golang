package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/snehalyelmati/telegram-bot-golang/models"
)

func main() {
	l := log.New(os.Stdout, "go-api:", log.LstdFlags)

	TOKEN := os.Getenv("TOKEN")
	SERVER_URL := os.Getenv("SERVER_URL")

	TELEGRAM_API := "https://api.telegram.org/bot" + TOKEN
	URI := "/webhook" + TOKEN
	WEBHOOK_URL := SERVER_URL + URI

	// l.Println(TELEGRAM_API + "/setWebhook?url=" + WEBHOOK_URL)
	res, err := http.Get(TELEGRAM_API + "/setWebhook?url=" + WEBHOOK_URL)
	if err != nil {
		l.Println(err)
	}
	l.Println("Response from the initial webhook:", res)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		l.Println("Get request")
		return c.SendString("Hello, world!")
	})

	app.Post(URI, func(c *fiber.Ctx) error {
		l.Println("Post request")
		reqBody := new(models.MessageReq)
		json.Unmarshal(c.Body(), reqBody)
		l.Printf("%v", reqBody)

		data, err := json.Marshal(map[string]string{
			"chat_id": strconv.Itoa(reqBody.Message.Chat.ID),
			"text":    reqBody.Message.Text,
		})
		if err != nil {
			l.Println(err)
		}

		_, err = http.Post(TELEGRAM_API+"/sendMessage", "application/json", bytes.NewBuffer(data))
		if err != nil {
			l.Println(err)
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":3000")
	l.Println("Server running on port 3000")
}
