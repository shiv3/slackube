package test

import (
	"context"

	"github.com/shiv3/slackube/app/adapter/slacksender"

	"github.com/shiv3/slackube/app/view/slack/common"

	"github.com/shiv3/slackube/app/usecase"

	"github.com/slack-go/slack/slackevents"

	"github.com/slack-go/slack"
)

type Handler struct {
	usecases      *usecase.UsecasesImpl
	slackBotToken string
	slackClient   *slack.Client
}

func NewHandler(usecase *usecase.UsecasesImpl, slackClient *slack.Client) *Handler {
	return &Handler{
		usecases: usecase, slackClient: slackClient}
}

func (h Handler) Test(ctx context.Context, ev *slackevents.AppMentionEvent) error {
	res, err := h.usecases.ListNameSpace(ctx)
	if err != nil {
		return err
	}
	return slacksender.NewSender(h.slackClient, ev.Channel).
		PostBlocks([]slack.Block{common.CodeBlockListBlock(res)})
}
