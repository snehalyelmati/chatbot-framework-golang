package telegramhdl

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/snehalyelmati/telegram-bot-golang/internal/core/ports"
)

type HTTPHandler struct {
	telegramService ports.TelegramService
	l               *log.Logger
	ProjectID       string
	Language        string
	TelegramAPI     string
}

func NewHTTPHandler(l *log.Logger, telegramService ports.TelegramService, projectID, language, telegramAPI string) *HTTPHandler {
	return &HTTPHandler{
		l:               l,
		telegramService: telegramService,
		ProjectID:       projectID,
		Language:        language,
		TelegramAPI:     telegramAPI,
	}
}

func (hdl *HTTPHandler) HealthCheck(c *fiber.Ctx) error {
	res := hdl.telegramService.HealthCheck()
	return c.SendString(res)
}

func (hdl *HTTPHandler) SendMessage(c *fiber.Ctx) error {
	req := NewTelegramReq()
	json.Unmarshal(c.Body(), req)

	chatID := req.Message.Chat.ID
	utterance := req.Message.Text

	err := hdl.telegramService.SendMessage(utterance, chatID, hdl.ProjectID, hdl.Language, hdl.TelegramAPI)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
