package get

import (
	"github.com/shiv3/slackube/app/usecase"
	"github.com/slack-go/slack"
)

type GetHandler struct {
	usecases    *usecase.UsecasesImpl
	slackClient *slack.Client
}

func NewGetHandler(usecase *usecase.UsecasesImpl, slackClient *slack.Client) *GetHandler {
	return &GetHandler{
		usecases: usecase, slackClient: slackClient}
}
