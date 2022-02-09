package dialgoflowsrv

import (
	"context"
	"fmt"
	"log"
	"strings"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type service struct {
	l *log.Logger
}

func New(l *log.Logger) *service {
	return &service{
		l: l,
	}
}

func (srv *service) DetectIntentText(projectID, sessionID, text, languageCode string) (string, string, error) {
	ctx := context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return "", "", err
	}
	defer sessionClient.Close()

	if projectID == "" || sessionID == "" {
		return "", "", fmt.Errorf("Received empty project (%s) or session (%s)", projectID, sessionID)
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", "", err
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	srv.l.Println("Fulfillment text from dialogflow:", strings.ReplaceAll(fulfillmentText, "\r\n", ""))

	return fulfillmentText, queryResult.String(), nil
}
