package slack

import (
	"encoding/json"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func (h handlerImpl) auth(header http.Header, body []byte, w http.ResponseWriter) error {
	sv, err := slack.NewSecretsVerifier(header, h.signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}
	return nil
}

func (h handlerImpl) eventVerify(header http.Header, body []byte, w http.ResponseWriter, eventType string) error {
	if err := h.auth(header, body, w); err != nil {
		return err
	}
	if eventType == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	return nil
}
