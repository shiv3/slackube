package common

import (
	"fmt"

	"github.com/slack-go/slack"
)

func ErrorViewMsgOption(err error, blocks ...slack.Block) slack.MsgOption {
	var errMsg string
	errMsg = fmt.Sprintf("err: ```%s``` ", err)
	if blocks != nil {
		blocks := &slack.Blocks{BlockSet: blocks}
		b, _ := blocks.MarshalJSON()
		errMsg += fmt.Sprintf("request block: ```%s``` ", b)
	}

	return slack.MsgOptionAttachments(
		slack.Attachment{
			Title: "error",
			Text:  errMsg,
			Color: "#FF0000",
		},
	)
}
