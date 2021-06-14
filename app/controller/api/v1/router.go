package v1

import (
	"fmt"

	"github.com/shiv3/slackube/app/controller/api/v1/slackevents"

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
	slackEventsHandler slackevents.Handler
	pingHandler        ping.Handler
}

// NewRouter V1Routerの作成
func NewV1Router(
	slackHandler slackevents.Handler,
	pingHandler ping.Handler,
) RouterImpl {
	return RouterImpl{
		slackEventsHandler: slackHandler,
		pingHandler:        pingHandler,
	}
}

// Dispatch V1RouterへHandlerを登録する。
func (r RouterImpl) Dispatch(e *echo.Echo) error {
	group := e.Group(V1Prefix)
	group.GET(ping.PingEndpoint, r.pingHandler.GetPing)
	group.POST(slackevents.SlackEventsEndpoint, r.slackEventsHandler.SlackEvents)
	return nil
}

func GetPingPath() string {
	return fmt.Sprintf("%s%s", V1Prefix, ping.PingEndpoint)
}
