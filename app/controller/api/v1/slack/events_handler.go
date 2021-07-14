package slack

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack/slackevents"
)

func (h handlerImpl) SlackEvents(c echo.Context) error {
	ctx := c.Request().Context()
	r := c.Request()
	w := c.Response()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	// 認証
	err = h.eventVerify(r.Header, body, w, eventsAPIEvent.Type)
	if err != nil {
		return err
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		err := h.slackRouter.EventsRoute(ctx, innerEvent)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Print(err)
			return err
		}
	}
	return nil
}
