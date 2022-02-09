package telegramhdl

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/snehalyelmati/telegram-bot-golang/internal/core/domain"
	"github.com/snehalyelmati/telegram-bot-golang/internal/core/ports"
)

type HTTPHandler struct {
	telegramService       ports.TelegramService
	transcriptsRepository ports.TranscriptsRepository
	l                     *log.Logger
	ProjectID             string
	Language              string
	TelegramAPI           string
}

func NewHTTPHandler(l *log.Logger, telegramService ports.TelegramService, transcriptsRepository ports.TranscriptsRepository,
	projectID, language, telegramAPI string) *HTTPHandler {
	return &HTTPHandler{
		l:                     l,
		telegramService:       telegramService,
		transcriptsRepository: transcriptsRepository,
		ProjectID:             projectID,
		Language:              language,
		TelegramAPI:           telegramAPI,
	}
}

func (hdl *HTTPHandler) HealthCheck(c *fiber.Ctx) error {
	res := hdl.telegramService.HealthCheck()
	return c.SendString(res)
}

func (hdl *HTTPHandler) SendMessage(c *fiber.Ctx) error {
	req := domain.NewTelegramReq()
	json.Unmarshal(c.Body(), req) // nolint

	chatID := req.Message.Chat.ID
	utterance := req.Message.Text
	messageID := req.Message.MessageID
	firstName := req.Message.Chat.FirstName

	// initialize stuff
	if utterance == "/start" {
		// save the user data
		utterance = "hi"
		err := hdl.transcriptsRepository.SaveUser(req.Message.From)
		if err != nil {
			return err
		}
	}

	dialogflowResponse, err := hdl.telegramService.SendMessage(utterance, chatID, hdl.ProjectID, hdl.Language, hdl.TelegramAPI)
	if err != nil {
		return err
	}

	// saving transcripts
	err = hdl.transcriptsRepository.Save(*domain.NewTranscript(strconv.Itoa(chatID), strconv.Itoa(messageID), firstName, utterance, dialogflowResponse, *req))
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
