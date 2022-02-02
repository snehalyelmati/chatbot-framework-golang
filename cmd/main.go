package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	. "github.com/snehalyelmati/telegram-bot-golang/models"
	DialogFlow "github.com/snehalyelmati/telegram-bot-golang/services"
)

func main() {
	l := log.New(os.Stdout, "go-api:", log.LstdFlags)

	TOKEN := os.Getenv("TOKEN")
	SERVER_URL := os.Getenv("SERVER_URL")
	PROJECT_ID := os.Getenv("PROJECT_ID")
	LANGUAGE := os.Getenv("LANGUAGE")

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
		reqBody := new(TelgramMessageReq)
		json.Unmarshal(c.Body(), reqBody)
		l.Printf("%v", reqBody)

		// initialize stuff
		inputMessage := reqBody.Message.Text
		if inputMessage == "/start" {
			// save the user data
			inputMessage = "hi"
		}

		// mimic the request
		// data, err := json.Marshal(map[string]string{
		// 	"chat_id": strconv.Itoa(reqBody.Message.Chat.ID),
		// 	"text":    reqBody.Message.Text,
		// })

		// get response from dialogflow

		// sample format for text and quick replies:
		// statement1
		// statement2
		// <qrType1, qrName1, qrText1><qrType2, qrName2, qrText2>...

		sessionID := strconv.Itoa(reqBody.Message.Chat.ID)
		response, queryResult, err := DialogFlow.DetectIntentText(PROJECT_ID, sessionID, inputMessage, LANGUAGE)
		if err != nil {
			l.Println(err)
		}
		l.Println("Query result:", queryResult)

		statements := strings.Split(response, "\n")

		for _, statement := range statements {
			// prepare a response for telegram
			data, err := json.Marshal(map[string]string{
				"chat_id": strconv.Itoa(reqBody.Message.Chat.ID),
				"text":    statement,
			})
			if err != nil {
				l.Println(err)
			}

			response := new(TelegramMessageRes)
			response.ChatID = strconv.Itoa(reqBody.Message.Chat.ID)
			response.ReplyMarkup.Keyboard = make([][]TelegramMessageResReplyMarkupKeyboard, 0)

			// get all the Quick Replies from the text
			// prepare a map with label, postBack and qrType
			regEx := `<(?P<qrType>[\sa-zA-Z0-9]*),(?P<postBack>[\sa-zA-Z0-9]*),(?P<label>[\sa-zA-Z0-9]*)>`
			quickReplies, statement := getParams(regEx, statement)
			response.Text = statement
			l.Println("Quick replies:", quickReplies)

			// initialize keyboard obj
			response.ReplyMarkup.Keyboard = append(response.ReplyMarkup.Keyboard, []TelegramMessageResReplyMarkupKeyboard{})

			// append available quick replies
			for _, qr := range quickReplies {
				if qr["qrType"] == "OPT" {
					// append to response.ReplyMarkup.Keyboard
					response.ReplyMarkup.Keyboard = append(response.ReplyMarkup.Keyboard, []TelegramMessageResReplyMarkupKeyboard{{Text: qr["label"]}})
				} else if qr["qrType"] == "SUGT" {
					// TODO: add buttons for suggestions
				}
			}
			l.Println("Response ReplyMarkup Keyboard:", response.ReplyMarkup.Keyboard)

			if len(response.ReplyMarkup.Keyboard) > 0 {
				response.ReplyMarkup.OneTimeKeyboard = true
				response.ReplyMarkup.ResizeKeyboard = true
			}

			// prepare a json object from the response object
			data, err = json.Marshal(response)
			if err != nil {
				l.Println(err)
			}

			// call the Telegram API
			_, err = http.Post(TELEGRAM_API+"/sendMessage", "application/json", bytes.NewBuffer(data))
			if err != nil {
				l.Println(err)
			}
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":3000")
	l.Println("Server running on port 3000")
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
