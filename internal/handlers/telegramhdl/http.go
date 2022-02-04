package telegramhdl

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/snehalyelmati/telegram-bot-golang/internal/core/ports"
)

// implement telegram handler and required functions here
type HTTPHandler struct {
	telegramService ports.TelegramService
	l               *log.Logger
	ProjectID       string
	Language        string
}

func NewHTTPHandler(l *log.Logger, telegramService ports.TelegramService, projectID, language string) *HTTPHandler {
	return &HTTPHandler{
		l:               l,
		telegramService: telegramService,
		ProjectID:       projectID,
		Language:        language,
	}
}

func (hdl *HTTPHandler) HealthCheck(c *fiber.Ctx) error {
	res := hdl.telegramService.HealthCheck()
	return c.SendString(res)
}

func (hdl *HTTPHandler) SendMessage(c *fiber.Ctx) error {
	// user services and do return the response
	req := NewTelegramReq()
	json.Unmarshal(c.Body(), req)

	// initialize stuff
	inputMessage := req.Message.Text
	if inputMessage == "/start" {
		// save the user data
		inputMessage = "hi"
	}

	hdl.telegramService.SendMessage(inputMessage, req.Message.Chat.ID, hdl.ProjectID, hdl.Language)

	return c.SendStatus(fiber.StatusOK)
}
