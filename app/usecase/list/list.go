package list

import (
	"context"

	"github.com/shiv3/slackube/app/adapter/k8s"
	"github.com/slack-go/slack"

	"github.com/shiv3/slackube/app/adapter/k8s/list"
)

type ListUseCaseImpl struct {
	ListAdapter list.ListAdapterInterface
}

func (u ListUseCaseImpl) ListNameSpace(ctx context.Context) (slack.MsgOption, error) {

	ks, err := k8s.NewK8SClientClient()
	if err != nil {
		return nil, err
	}

	nsList, err := ks.ListNs(ctx)
	if err != nil {
		return nil, err
	}

	options := make([]*slack.OptionBlockObject, 0, len(nsList.Items))
	for _, ns := range nsList.Items {
		optionText := slack.NewTextBlockObject(slack.PlainTextType, ns.Name, false, false)
		options = append(options, slack.NewOptionBlockObject(ns.Name, optionText, nil))
	}

	placeholder := slack.NewTextBlockObject(slack.PlainTextType, "Select NameSpace", false, false)
	selectMenu := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, placeholder, "", options...)
	actionBlock := slack.NewActionBlock("select-namespace", selectMenu)
	text := slack.NewTextBlockObject(slack.MarkdownType, "Please select *namespace*.", false, false)
	textSection := slack.NewSectionBlock(text, nil, nil)
	//fallbackText := slack.MsgOptionText("This client is not supported.", false)
	blocks := slack.MsgOptionBlocks(textSection, actionBlock)
	return blocks, nil
}
