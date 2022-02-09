package domain

import (
	"time"
)

type Transcript struct {
	SessionID       string
	MessageID       string
	FirstName       string
	Utterance       string
	FulfillmentText string
	Request         TelegramReq
	CreatedTime     time.Time
}

func NewTranscript(sessionID, messageID, firstName, utterance, fulfillmentText string, request TelegramReq) *Transcript {
	return &Transcript{
		SessionID:       sessionID,
		MessageID:       messageID,
		FirstName:       firstName,
		Utterance:       utterance,
		FulfillmentText: fulfillmentText,
		Request:         request,
		CreatedTime:     time.Now(),
	}
}
