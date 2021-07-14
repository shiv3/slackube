package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/slack-go/slack"

	"github.com/labstack/echo/v4"
)

func (h handlerImpl) SlackActions(c echo.Context) error {
	r := c.Request()
	w := c.Response()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	if err := h.auth(r.Header, body, c.Response()); err != nil {
		fmt.Println(err)
		return err
	}

	str, _ := url.QueryUnescape(string(body))
	payload := strings.Replace(str, "payload=", "", 1)

	var actionCallBack *slack.InteractionCallback
	if err := json.Unmarshal([]byte(payload), &actionCallBack); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	ctx := c.Request().Context()
	if err := h.slackRouter.ActionsRoute(ctx, actionCallBack); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return nil
}
