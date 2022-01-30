package models

type TelegramMessageRes struct {
	ChatID      string                        `json:"chat_id"`
	Text        string                        `json:"text"`
	ReplyMarkup TelegramMessageResReplyMarkup `json:"reply_markup"`
}

type TelegramMessageResReplyMarkup struct {
	OneTimeKeyboard bool                                      `json:"one_time_keyboard"`
	ResizeKeyboard  bool                                      `json:"resize_keyboard"`
	Keyboard        [][]TelegramMessageResReplyMarkupKeyboard `json:"keyboard"`
}

type TelegramMessageResReplyMarkupKeyboard struct {
	Text string `json:"text"`
}

