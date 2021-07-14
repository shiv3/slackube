package slack

import (
	"github.com/shiv3/slackube/app/controller/slackcontoller"

	"github.com/slack-go/slack"

	"github.com/shiv3/slackube/app/usecase"

	"github.com/labstack/echo/v4"
)

type (
	Handler interface {
		SlackEvents(c echo.Context) error
		SlackActions(c echo.Context) error
	}

	handlerImpl struct {
		signingSecret string
		slackBotToken string
		slackRouter   *slackcontoller.SlackRouter
	}
)

func NewHandlerImpl(signingSecret string, slackBotToken string) (*handlerImpl, error) {
	u, err := usecase.NewUsecasesImpl()
	if err != nil {
		return nil, err
	}

	return &handlerImpl{
		signingSecret: signingSecret,
		slackBotToken: slackBotToken,
		slackRouter:   slackcontoller.NewSlackRouter(u, slack.New(slackBotToken)),
	}, nil
}
