package slackcontoller

import (
	"context"
	"fmt"
	"regexp"

	"github.com/shiv3/slackube/app/controller/slackcontoller/test"

	"github.com/shiv3/slackube/app/controller/slackcontoller/updateimage"

	"github.com/shiv3/slackube/app/controller/slackcontoller/get"

	"github.com/slack-go/slack"

	"github.com/shiv3/slackube/app/usecase"

	"github.com/slack-go/slack/slackevents"
)

type SlackRouter struct {
	mentionEventHandlers map[*regexp.Regexp]func(ctx context.Context, ev *slackevents.AppMentionEvent) error
	slackActionsHandlers map[string]func(ctx context.Context, actionCallBack *slack.InteractionCallback) error
}

func NewSlackRouter(usecases *usecase.UsecasesImpl, slackClient *slack.Client) *SlackRouter {
	return &SlackRouter{
		mentionEventHandlers: map[*regexp.Regexp]func(ctx context.Context, ev *slackevents.AppMentionEvent) error{
			regexp.MustCompile(`test`):         test.NewHandler(usecases, slackClient).Test,
			regexp.MustCompile(`get ns`):       get.NewGetHandler(usecases, slackClient).GetNs,
			regexp.MustCompile(`update image`): updateimage.NewUpdateImageHandler(usecases, slackClient).Start,
		},
		slackActionsHandlers: map[string]func(ctx context.Context, actionCallBack *slack.InteractionCallback) error{
			"update-image-1": updateimage.NewUpdateImageHandler(usecases, slackClient).SelectDeployment,
			"update-image-2": updateimage.NewUpdateImageHandler(usecases, slackClient).GetImageTag,
			"update-image-3": updateimage.NewUpdateImageHandler(usecases, slackClient).UpdateImageTag,
		},
	}
}

func (r *SlackRouter) EventsRoute(ctx context.Context, event slackevents.EventsAPIInnerEvent) error {
	switch ev := event.Data.(type) {
	case *slackevents.AppMentionEvent:
		for k, f := range r.mentionEventHandlers {
			if k.MatchString(ev.Text) {
				err := f(ctx, ev)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *SlackRouter) ActionsRoute(ctx context.Context, actionCallBack *slack.InteractionCallback) error {
	switch actionCallBack.Type {
	case slack.InteractionTypeBlockActions:
		for _, blockAction := range actionCallBack.ActionCallback.BlockActions {
			for k, f := range r.slackActionsHandlers {
				if k == blockAction.BlockID {
					return f(ctx, actionCallBack)
				}
			}
		}
		return fmt.Errorf("not matched")
	}
	return nil
}
