package common

import (
	"github.com/slack-go/slack"
)

func CodeBlockListBlock(in []string) slack.Block {
	d := "```"
	for _, ns := range in {
		d += "- " + ns + "\n"
	}
	d += "```"

	return slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, d, false, false), nil, nil)
}
