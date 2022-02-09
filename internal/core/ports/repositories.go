package ports

import (
	"github.com/snehalyelmati/telegram-bot-golang/internal/core/domain"
)

type TranscriptsRepository interface {
	Save(domain.Transcript) error
	SaveUser(domain.From) error
}
