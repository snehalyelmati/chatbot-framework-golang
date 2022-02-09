package transcriptsrepo

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"github.com/snehalyelmati/telegram-bot-golang/internal/core/domain"
)

type FirestoreRepo struct {
	firebaseApp *firebase.App
	l           *log.Logger
}

func NewFirestoreRepo(firebaseApp *firebase.App, logger *log.Logger) *FirestoreRepo {
	return &FirestoreRepo{
		firebaseApp: firebaseApp,
		l:           logger,
	}
}

func (fr *FirestoreRepo) Save(transcript domain.Transcript) error {
	ctx := context.Background()
	firestore, err := fr.firebaseApp.Firestore(ctx)
	if err != nil {
		fr.l.Println("Couldn't initialize firestore instance to save transcripts")
		return err
	}

	if _, err := firestore.Collection("transcripts").Doc(transcript.MessageID).Set(ctx, transcript); err != nil {
		fr.l.Println("Couldn't save transcript")
		return err
	}

	fr.l.Println("Saved transcript successfully")
	fr.l.Println(transcript)
	return nil
}
