package slackevents

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/shiv3/slackube/app/usecase/list"

	"github.com/shiv3/slackube/app/adapter/k8s"

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

		CommandHandlers map[*regexp.Regexp]func()

		ListUseCaseImpl list.ListUseCaseImpl
	}
)

func NewHandlerImpl(signingSecret string) *handlerImpl {
	return &handlerImpl{signingSecret: signingSecret}
}

var api = slack.New("")

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
			r := regexp.MustCompile(`get ns`)
			if r.MatchString(ev.Text) {
				//api.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Name Space List : %v",ns), false))

				ks, err := k8s.NewK8SClientClient()
				if err != nil {
					return err
				}

				nsList, err := ks.NsList()
				if err != nil {
					return err
				}

				options := make([]*slack.OptionBlockObject, 0, len(nsList.Items))
				for _, ns := range nsList.Items {
					optionText := slack.NewTextBlockObject(slack.PlainTextType, ns.Name, false, false)
					options = append(options, slack.NewOptionBlockObject(ns.Name, optionText, nil))
				}

				placeholder := slack.NewTextBlockObject(slack.PlainTextType, "Select NameSpace", false, false)
				selectMenu := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, placeholder, "", options...)
				actionBlock := slack.NewActionBlock("select-namespace", selectMenu)
				text := slack.NewTextBlockObject(slack.MarkdownType, "Please select *namespace*.", false, false)
				textSection := slack.NewSectionBlock(text, nil, nil)
				//fallbackText := slack.MsgOptionText("This client is not supported.", false)
				blocks := slack.MsgOptionBlocks(textSection, actionBlock)

				if _, _, err := api.PostMessage(ev.Channel, blocks); err != nil {
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
