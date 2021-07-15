package get

import (
	"context"

	"github.com/shiv3/slackube/app/adapter/slacksender"
	"github.com/shiv3/slackube/app/view/slack/common"

	"github.com/slack-go/slack/slackevents"

	"github.com/slack-go/slack"
)

func (h GetHandler) GetDp(ctx context.Context, ev *slackevents.AppMentionEvent) error {
	res, err := h.usecases.ListNameSpace(ctx)
	if err != nil {
		return err
	}
	return slacksender.NewSender(h.slackClient, ev.Channel).
		PostBlocks([]slack.Block{common.CodeBlockListBlock(res)})
}
