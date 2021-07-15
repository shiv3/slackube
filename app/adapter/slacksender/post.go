package slacksender

import (
	"github.com/shiv3/slackube/app/view/slack/common"
	"github.com/slack-go/slack"
)

type Client struct {
	client  *slack.Client
	channel string
}

func NewSender(client *slack.Client, channel string) *Client {
	return &Client{client: client, channel: channel}
}

func (s *Client) PostError(err error, blocks ...slack.Block) error {
	if _, _, err2 := s.client.PostMessage(s.channel, common.ErrorViewMsgOption(err, blocks...)); err2 != nil {
		return err
	}
	return err
}

func (s *Client) PostBlocks(msgBlock []slack.Block) error {
	if _, _, err := s.client.PostMessage(s.channel, slack.MsgOptionBlocks(msgBlock...)); err != nil {
		return s.PostError(err, msgBlock...)
	}
	return nil
}

func (s *Client) PostMsgOptions(msgOptions []slack.MsgOption) error {
	if _, _, err := s.client.PostMessage(s.channel, msgOptions...); err != nil {
		return s.PostError(err)
	}
	return nil
}
