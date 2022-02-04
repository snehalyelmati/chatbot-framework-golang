package ports

type TelegramService interface {
	SendMessage(utterance string, chatID int, projectID, language, telegramAPI string) error
	HealthCheck() string
}

type DialogflowService interface {
	DetectIntentText(projectID, sessionID, text, languageCode string) (string, string, error)
}
