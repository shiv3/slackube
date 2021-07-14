package updateimage

import (
	"context"
	"fmt"
	"strings"

	"github.com/shiv3/slackube/app/usecase"

	"github.com/slack-go/slack/slackevents"

	"github.com/slack-go/slack"
)

type UpdateImageHandler struct {
	usecases      *usecase.UsecasesImpl
	slackBotToken string
	slackClient   *slack.Client
}

func NewUpdateImageHandler(usecase *usecase.UsecasesImpl, slackClient *slack.Client) *UpdateImageHandler {
	return &UpdateImageHandler{
		usecases: usecase, slackClient: slackClient}
}

func (h UpdateImageHandler) Start(ctx context.Context, ev *slackevents.AppMentionEvent) error {
	res, err := h.usecases.ListNameSpace(ctx)
	if err != nil {
		return err
	}

	options := make([]*slack.OptionBlockObject, 0, len(res))
	for _, ns := range res {
		optionText := slack.NewTextBlockObject(slack.PlainTextType, ns, false, false)
		options = append(options, slack.NewOptionBlockObject(ns, optionText, nil))
	}

	placeholder := slack.NewTextBlockObject(slack.PlainTextType, "Select NameSpace", false, false)
	selectMenu := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, placeholder, "", options...)
	actionBlock := slack.NewActionBlock("update-image-1", selectMenu)

	text := slack.NewTextBlockObject(slack.MarkdownType, "Please select *namespace*.", false, false)
	textSection := slack.NewSectionBlock(text, nil, nil)
	fallbackText := slack.MsgOptionText("This client is not supported.", false)
	blocks := slack.MsgOptionBlocks(textSection, actionBlock)

	if _, _, err := h.slackClient.PostMessage(ev.Channel, fallbackText, blocks); err != nil {
		return err
	}
	return nil
}

func (h UpdateImageHandler) SelectDeployment(ctx context.Context, actionCallBack *slack.InteractionCallback) error {
	if len(actionCallBack.BlockActionState.Values) < 1 {
		return fmt.Errorf("value < 1")
	}
	var namespaceValue string
	for _, v := range actionCallBack.BlockActionState.Values {
		for _, action := range v {
			namespaceValue = action.SelectedOption.Value
		}
	}

	res, err := h.usecases.ListDeployment(ctx, namespaceValue)
	if err != nil {
		return err
	}

	optionGroups := make([]*slack.OptionGroupBlockObject, 0, len(res))
	for _, dp := range res {
		options := make([]*slack.OptionBlockObject, 0, len(res))
		for container, image := range dp.ContainerImages {
			optionText := slack.NewTextBlockObject(slack.PlainTextType, container, false, false)
			optionDeployment := slack.NewTextBlockObject(slack.PlainTextType, dp.Deployment, false, false)
			options = append(options, slack.NewOptionBlockObject(image, optionText, optionDeployment))
		}
		optionGroupText := slack.NewTextBlockObject(slack.PlainTextType, dp.Deployment, false, false)
		optionGroups = append(optionGroups, slack.NewOptionGroupBlockElement(optionGroupText, options...))
	}

	placeholder := slack.NewTextBlockObject(slack.PlainTextType, "Select Container", false, false)
	selectMenu := slack.NewOptionsGroupSelectBlockElement(slack.OptTypeStatic, placeholder, "", optionGroups...)
	actionBlock := slack.NewActionBlock("update-image-2", selectMenu)

	blocks := slack.MsgOptionBlocks(
		slack.NewSectionBlock(slack.NewTextBlockObject(slack.MarkdownType, "Please select *deployment*.", false, false),
			[]*slack.TextBlockObject{
				slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("namespace: %s", namespaceValue), true, false),
			}, nil),
		actionBlock,
	)

	if _, _, _, err := h.slackClient.UpdateMessage(actionCallBack.Channel.ID, actionCallBack.Message.Msg.Timestamp,
		slack.MsgOptionText("This client is not supported.", false),
		blocks); err != nil {
		return err
	}

	return nil
}

func (h UpdateImageHandler) GetImageTag(ctx context.Context, actionCallBack *slack.InteractionCallback) error {
	var imageValue, deploymentValue, containerName string
	for _, v := range actionCallBack.BlockActionState.Values {
		for _, action := range v {
			imageValue = action.SelectedOption.Value
			containerName = action.SelectedOption.Text.Text
			deploymentValue = action.SelectedOption.Description.Text
		}
	}
	var namespaceValue string
	for _, v := range actionCallBack.Message.Msg.Blocks.BlockSet {
		if b, t := v.(*slack.SectionBlock); t {
			for _, field := range b.Fields {
				ns := strings.Split(field.Text, "namespace: ")
				if len(ns) > 1 {
					namespaceValue = ns[1]
				}
			}
		}
	}
	tags, err := h.usecases.ListImageTag(ctx, strings.Split(strings.Split(imageValue, ":")[0], "@")[0])
	if err != nil {
		return err
	}

	if len(tags) < 1 {
		return fmt.Errorf("failed get List tag")
	}

	if len(tags) > 20 {
		tags = tags[0:20]
	}
	options := make([]*slack.OptionBlockObject, 0, len(tags))
	for _, tag := range tags {
		digest := strings.Split(tag.Digest, "sha256:")[1][0:12]
		tags := strings.Join(append([]string{digest}, tag.Tags...), " ")
		optionText := slack.NewTextBlockObject(slack.PlainTextType, tags, false, false)
		descriptionText := slack.NewTextBlockObject(slack.PlainTextType, deploymentValue, false, false)
		options = append(options, slack.NewOptionBlockObject(tag.Digest, optionText, descriptionText))
	}

	placeholder := slack.NewTextBlockObject(slack.PlainTextType, "Select Image Tag", false, false)
	selectMenu := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, placeholder, "", options...)
	actionBlock := slack.NewActionBlock("update-image-3", selectMenu)

	blocks := slack.MsgOptionBlocks(
		slack.NewSectionBlock(slack.NewTextBlockObject(slack.MarkdownType, "Please select *image*.", false, false),
			[]*slack.TextBlockObject{
				slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("namespace: %s", namespaceValue), true, false),
				slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("deployment: %s", deploymentValue), true, false),
				slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("container: %s", containerName), true, false),
				slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("image: %s", imageValue), true, false),
			}, nil),
		actionBlock,
	)
	if _, _, _, err := h.slackClient.UpdateMessage(actionCallBack.Channel.ID, actionCallBack.Message.Msg.Timestamp,
		slack.MsgOptionText("This client is not supported.", false), blocks); err != nil {
		return err
	}
	return nil
}

func (h UpdateImageHandler) UpdateImageTag(ctx context.Context, actionCallBack *slack.InteractionCallback) error {
	var tags, digest string
	for _, v := range actionCallBack.BlockActionState.Values {
		for _, action := range v {
			digest = action.SelectedOption.Value
			tags = action.SelectedOption.Text.Text
		}
	}

	var namespace, deployment, container, image string
	for _, v := range actionCallBack.Message.Msg.Blocks.BlockSet {
		if b, t := v.(*slack.SectionBlock); t {
			for _, field := range b.Fields {
				ns := strings.Split(field.Text, "namespace: ")
				if len(ns) > 1 {
					namespace = ns[1]
				}

				dp := strings.Split(field.Text, "deployment: ")
				if len(dp) > 1 {
					deployment = dp[1]
				}

				cn := strings.Split(field.Text, "container: ")
				if len(cn) > 1 {
					container = cn[1]
				}

				img := strings.Split(field.Text, "image: ")
				if len(img) > 1 {
					image = img[1]
				}
				tg := strings.Split(field.Text, "tags: ")
				if len(tg) > 1 {
					tags = tg[1]
				}
			}
		}
	}
	_, err := h.usecases.UpdateImage(ctx, namespace, deployment, container, digest)
	if err != nil {
		return err
	}

	blocks := slack.MsgOptionBlocks(
		slack.NewSectionBlock(slack.NewTextBlockObject(slack.MarkdownType, "*Deploy Done*.", false, false),
			[]*slack.TextBlockObject{
				slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("namespace: %s", namespace), true, false),
				slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("deployment: %s", deployment), true, false),
				slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("container: %s", container), true, false),
				slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("image: %s", image), true, false),
				slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("digest: %s", digest), true, false),
				slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("tags: %s", tags), true, false),
			}, nil),
	)
	if _, _, _, err := h.slackClient.UpdateMessage(actionCallBack.Channel.ID, actionCallBack.Message.Msg.Timestamp,
		slack.MsgOptionText("This client is not supported.", false), blocks); err != nil {
		return err
	}

	return nil
}
