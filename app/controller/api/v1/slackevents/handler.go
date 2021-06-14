package slackevents

import (
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack/slackevents"
)

const SlackEventsEndpoint = "/events-endpoint"

type (
	Handler interface {
		SlackEvents(c echo.Context) error
	}

	handlerImpl struct {
		signingSecret string
	}
)

func NewHandlerImpl() *handlerImpl {
	return &handlerImpl{}
}

var api = slack.New("TOKEN")

func (h handlerImpl) SlackEvents(c echo.Context) error {

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
	err = h.verify(r.Header, body, w, eventsAPIEvent.Type)
	if err != nil {
		return err
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
		}
	}
	return nil
}
