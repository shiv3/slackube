package v1

import (
	"github.com/shiv3/slackube/app/controller/api/v1/slack"

	"github.com/labstack/echo/v4"
	"github.com/shiv3/slackube/app/controller/api/v1/ping"
)

const (
	V1Prefix = "/v1"
)

type Router interface {
	Dispatch(e *echo.Echo) error
}

// RouterImpl v1パスに来た際のRouteを記述する。
type RouterImpl struct {
	pingEndpoint        string
	slackEventsEndpoint string
	slackActionEndpoint string

	slackEventsHandler slack.Handler
	pingHandler        ping.Handler
}

// NewRouter V1Routerの作成
func NewV1Router(
	pingEndpoint string,
	slackEventsEndpoint string,
	slackActionEndpoint string,
	slackHandler slack.Handler,
	pingHandler ping.Handler,
) RouterImpl {
	return RouterImpl{
		pingEndpoint:        pingEndpoint,
		slackEventsEndpoint: slackEventsEndpoint,
		slackActionEndpoint: slackActionEndpoint,
		slackEventsHandler:  slackHandler,
		pingHandler:         pingHandler,
	}
}

// Dispatch V1RouterへHandlerを登録する。
func (r RouterImpl) Dispatch(e *echo.Echo) error {
	group := e.Group(V1Prefix)
	group.GET(r.pingEndpoint, r.pingHandler.GetPing)
	group.POST(r.slackEventsEndpoint, r.slackEventsHandler.SlackEvents)
	group.POST(r.slackActionEndpoint, r.slackEventsHandler.SlackActions)
	return nil
}
