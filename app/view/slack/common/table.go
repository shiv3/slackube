package common

import (
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/slack-go/slack"
)

func TableBlock(header []string, in [][]string) slack.Block {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, codeBlock(header, in), false, false), nil, nil)
}

func codeBlock(header []string, in [][]string) string {
	out := &strings.Builder{}
	table := tablewriter.NewWriter(out)
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(in)
	table.Render()

	d := "```\n"
	d += out.String()
	d += "```"
	return d
}
