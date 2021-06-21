package slackevents

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/shiv3/slackube/app/usecase/list"

	"github.com/slack-go/slack"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack/slackevents"
)

const SlackEventsEndpoint = "/eventsapi"

type (
	Handler interface {
		SlackEvents(c echo.Context) error
	}

	handlerImpl struct {
		signingSecret string
		slackBotToken string

		CommandHandlers map[*regexp.Regexp]func()

		ListUseCaseImpl list.ListUseCaseImpl
	}
)

func NewHandlerImpl(signingSecret string, slackBotToken string) *handlerImpl {
	return &handlerImpl{
		signingSecret: signingSecret,
		slackBotToken: slackBotToken,
	}
}

var api = slack.New("")

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
	err = h.verify(r.Header, body, w, eventsAPIEvent.Type)
	if err != nil {
		return err
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			r := regexp.MustCompile(`get ns`)
			if r.MatchString(ev.Text) {
				//api.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Name Space List : %v",ns), false))

				o, err := h.ListUseCaseImpl.ListNameSpace(ctx)
				if err != nil {
					return err
				}
				if _, _, err := api.PostMessage(ev.Channel, o); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}

				//if _, err := api.PostEphemeral(ev.Channel, ev.User, fallbackText, blocks); err != nil {
				//	w.WriteHeader(http.StatusInternalServerError)
				//}
			}
		}
	}
	return nil
}
