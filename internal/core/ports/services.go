package ports

import ()

type TelegramService interface {
	SendMessage(utterance string, chatID int, projectID, language string)
	HealthCheck() string
}

type DialogflowService interface {
	DetectIntentText(projectID, sessionID, text, languageCode string) (string, string, error)
}
