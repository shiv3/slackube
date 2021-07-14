package get_ns

import (
	"context"

	"github.com/shiv3/slackube/app/usecase"

	"github.com/slack-go/slack/slackevents"

	"github.com/slack-go/slack"
)

type ListHandler struct {
	usecases      *usecase.UsecasesImpl
	slackBotToken string
	slackClient   *slack.Client
}

func NewListHandler(usecase *usecase.UsecasesImpl, slackClient *slack.Client) *ListHandler {
	return &ListHandler{
		usecases: usecase, slackClient: slackClient}
}

func (h ListHandler) GetNs(ctx context.Context, ev *slackevents.AppMentionEvent) error {
	res, err := h.usecases.ListNameSpace(ctx)
	if err != nil {
		return err
	}

	options := make([]*slack.OptionBlockObject, 0, len(res))
	for _, ns := range res {
		optionText := slack.NewTextBlockObject(slack.PlainTextType, ns, false, false)
		options = append(options, slack.NewOptionBlockObject(ns, optionText, nil))
	}

	placeholder := slack.NewTextBlockObject(slack.PlainTextType, "Select NameSpace", false, false)
	selectMenu := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, placeholder, "", options...)
	actionBlock := slack.NewActionBlock("select-namespace", selectMenu)

	text := slack.NewTextBlockObject(slack.MarkdownType, "Please select *namespace*.", false, false)
	textSection := slack.NewSectionBlock(text, nil, nil)
	fallbackText := slack.MsgOptionText("This client is not supported.", false)
	blocks := slack.MsgOptionBlocks(textSection, actionBlock)

	if _, _, err := h.slackClient.PostMessage(ev.Channel, fallbackText, blocks); err != nil {
		return err
	}
	return nil
}
