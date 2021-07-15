package get

import (
	"context"

	"github.com/shiv3/slackube/app/adapter/slacksender"
	"github.com/shiv3/slackube/app/view/slack/common"

	"github.com/slack-go/slack/slackevents"

	"github.com/slack-go/slack"
)

func (h GetHandler) GetNs(ctx context.Context, ev *slackevents.AppMentionEvent) error {
	res, err := h.usecases.GetNameSpace(ctx)
	if err != nil {
		return err
	}
	var table [][]string
	for _, re := range res {
		table = append(table, []string{re.Name, re.Status, re.Age})
	}
	//var list []string
	//for _, re := range res {
	//	list = append(list, re.Name)
	//}

	return slacksender.NewSender(h.slackClient, ev.Channel).
		PostBlocks([]slack.Block{common.TableBlock([]string{"Name", "Status", "Age"}, table)})
}
