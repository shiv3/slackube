package get

import (
	"context"
	"fmt"

	"github.com/shiv3/slackube/app/adapter/slacksender"
	"github.com/shiv3/slackube/app/view/slack/common"

	"github.com/slack-go/slack"
)

type ActionFunc struct {
	ActionID string
	Func     func(ctx context.Context, actionCallBack *slack.InteractionCallback) error
	NextFunc func(ctx context.Context, actionCallBack *slack.InteractionCallback) error
}

func (h GetHandler) NewGetDeploymentActionFunc() ActionFunc {
	return ActionFunc{
		Func: h.getDeployment,
	}
}

func (h GetHandler) getDeployment(ctx context.Context, actionCallBack *slack.InteractionCallback) error {
	if len(actionCallBack.BlockActionState.Values) < 1 {
		return fmt.Errorf("value < 1")
	}
	var namespaceValue string
	for _, v := range actionCallBack.BlockActionState.Values {
		for _, action := range v {
			namespaceValue = action.SelectedOption.Value
		}
	}
	res, err := h.usecases.ListDeployment(ctx, namespaceValue)
	if err != nil {
		return err
	}
	if len(res) > 20 {
		res = res[0:20]
	}
	var list []string
	for _, re := range res {
		list = append(list, re.Deployment)
	}
	return slacksender.NewSender(h.slackClient, actionCallBack.Channel.ID).
		PostBlocks([]slack.Block{common.CodeBlockListBlock(list)})
}
