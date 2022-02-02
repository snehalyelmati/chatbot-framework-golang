package domain

type TelegramRes struct {
	ChatID      string              `json:"chat_id"`
	Text        string              `json:"text"`
	ReplyMarkup TelegramReplyMarkup `json:"reply_markup"`
}

func NewTelegramRes(chatId, text string, oneTimeKeyboard, resizeKeyboard bool) TelegramRes {
	return TelegramRes{
		ChatID: chatId,
		Text:   text,
		ReplyMarkup: TelegramReplyMarkup{
			OneTimeKeyboard: oneTimeKeyboard,
			ResizeKeyboard:  resizeKeyboard,
			Keyboard:        [][]TelegramKeyboard{},
		},
	}
}

type TelegramReplyMarkup struct {
	OneTimeKeyboard bool                 `json:"one_time_keyboard"`
	ResizeKeyboard  bool                 `json:"resize_keyboard"`
	Keyboard        [][]TelegramKeyboard `json:"keyboard"`
}

type TelegramKeyboard struct {
	Text string `json:"text"`
}
