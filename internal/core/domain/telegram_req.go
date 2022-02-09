package domain

type TelegramReq struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type From struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Type      string `json:"type"`
}

func NewTelegramReq() *TelegramReq {
	return new(TelegramReq)
}
