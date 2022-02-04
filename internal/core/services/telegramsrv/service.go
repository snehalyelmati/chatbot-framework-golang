package telegramsrv

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/snehalyelmati/telegram-bot-golang/internal/core/ports"
	telegramhdl "github.com/snehalyelmati/telegram-bot-golang/internal/handlers/telegramhdl"
)

type service struct {
	dialogflowService ports.DialogflowService
	l                 *log.Logger
}

func New(l *log.Logger, dialogflowService ports.DialogflowService) *service {
	return &service{
		l:                 l,
		dialogflowService: dialogflowService,
	}
}

func (srv *service) HealthCheck() string {
	res := "Server running, healthcheck successful!"
	srv.l.Println(res)
	return res
}

func (srv *service) SendMessage(utterance string, chatID int, projectID, language, telegramAPI string) error {
	// initialize stuff
	if utterance == "/start" {
		// TODO: save the user data
		utterance = "hi"
	}

	// get the fulfillment text from the dialogflow service
	dialogflowResponse, queryResult, err := srv.dialogflowService.DetectIntentText(projectID, strconv.Itoa(chatID), utterance, language)
	if err != nil {
		srv.l.Println(err)
		return err
	}
	srv.l.Println("Query result:", queryResult)

	// prepare the response from all the parameters available and post it to the telegram API
	srv.sendResponsesToTelegram(dialogflowResponse, strconv.Itoa(chatID), telegramAPI)

	return nil
}

func (srv *service) sendResponsesToTelegram(dialogflowResponse, chatID, telegramAPI string) {
	statements := strings.Split(dialogflowResponse, "\n")

	for _, statement := range statements {
		// prepare a response for telegram
		data, err := json.Marshal(map[string]string{
			"chat_id": chatID,
			"text":    statement,
		})
		if err != nil {
			srv.l.Println(err)
		}

		response := telegramhdl.NewTelegramRes(chatID, "", true, true)
		// response.ReplyMarkup.Keyboard = make([][]telegramhdl.TelegramKeyboard, 0)

		// get all the Quick Replies from the text
		// prepare a map with label, postBack and qrType
		regEx := `<(?P<qrType>[\sa-zA-Z0-9]*),(?P<postBack>[\sa-zA-Z0-9]*),(?P<label>[\sa-zA-Z0-9]*)>`
		quickReplies, statement := getParams(regEx, statement)
		response.Text = statement
		srv.l.Println("Quick replies:", quickReplies)

		// append available quick replies
		for _, qr := range quickReplies {
			if qr["qrType"] == "OPT" {
				// append to response.ReplyMarkup.Keyboard
				response.ReplyMarkup.Keyboard = append(response.ReplyMarkup.Keyboard, []telegramhdl.TelegramKeyboard{{Text: qr["label"]}})
			} else if qr["qrType"] == "SUGT" {
				// TODO: add buttons for suggestions
			}
		}
		srv.l.Println("Response ReplyMarkup Keyboard:", response.ReplyMarkup.Keyboard)

		if len(response.ReplyMarkup.Keyboard) > 0 {
			response.ReplyMarkup.OneTimeKeyboard = true
			response.ReplyMarkup.ResizeKeyboard = true
		}

		// prepare a json object from the response object
		data, err = json.Marshal(response)
		if err != nil {
			srv.l.Println(err)
		}

		// call the Telegram API
		_, err = http.Post(telegramAPI+"/sendMessage", "application/json", bytes.NewBuffer(data))
		if err != nil {
			srv.l.Println(err)
		}
	}
}

/**
 * Parses string with the given regular expression and returns the
 * group values defined in the expression.
 *
 */
func getParams(regEx, str string) ([]map[string]string, string) {

	var compRegEx = regexp.MustCompile(regEx)
	matchArr := compRegEx.FindAllStringSubmatch(str, -1)
	if len(matchArr) > 0 {
		str = compRegEx.ReplaceAllString(str, "")
	}

	res := make([]map[string]string, len(matchArr)) // map array

	for _, match := range matchArr {
		paramsMap := make(map[string]string)
		for i, name := range compRegEx.SubexpNames() {
			if i > 0 && i <= len(match) {
				paramsMap[name] = match[i]
			}
		}
		res = append(res, paramsMap)
	}
	return res, str
}
