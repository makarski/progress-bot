package slack

import (
	"fmt"

	"github.com/slack-go/slack"

	"github.com/makarski/progress-bot/config"
	"github.com/makarski/progress-bot/jira"
)

var issueKindEmoji = map[string]string{
	"Bug":         ":beetle:",
	"Story":       ":book:",
	"Improvement": ":hammer_and_wrench:",
	"Task":        ":memo:",
	"Sub-task":    ":memo:",
	"Epic":        ":eight_pointed_black_star:",
}

const defaultItemEmoji = ":memo:"
const dailyReportTitle = "Daily Team Progress Report :chart_with_upwards_trend:"
const dateFormat = "02 Jan 2006"

type SlackMessenger struct {
	client        *slack.Client
	statusesEmoji []config.StatusEmoji
}

func NewSlackMessenger(token string, statusesEmoji []config.StatusEmoji) *SlackMessenger {
	return &SlackMessenger{
		client:        slack.New(token),
		statusesEmoji: statusesEmoji,
	}
}

func (sm *SlackMessenger) SendMessage(
	channel string,
	issues jira.IssueReport,
	impactMap map[string]string,
) error {
	msg := slack.NewBlockMessage()
	msg = slack.AddBlockMessage(msg, slack.NewHeaderBlock(slack.NewTextBlockObject(slack.PlainTextType, dailyReportTitle, false, false)))

	msg = slack.AddBlockMessage(msg, slack.NewContextBlock("", slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("*%s - %s*", issues.From.Format(dateFormat), issues.To.Format(dateFormat)), false, false)))

	msg = slack.AddBlockMessage(msg, slack.NewDividerBlock())
	msg = slack.AddBlockMessage(msg, slack.NewSectionBlock(slack.NewTextBlockObject(slack.MarkdownType, impactMap["summary"], false, false), nil, nil))

	for _, status := range sm.statusesEmoji {
		issueItems, ok := issues.Changes[status.Name]
		if !ok {
			continue
		}
		msg = slack.AddBlockMessage(msg, slack.NewDividerBlock())
		msg = slack.AddBlockMessage(msg, slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, status.Emoji+" *"+status.Name+"*", false, false),
			nil,
			nil,
		))

		for i, item := range issueItems {
			issueKindEmoji, ok := issueKindEmoji[item.Kind]
			if !ok {
				issueKindEmoji = defaultItemEmoji
			}

			markdownText := fmt.Sprintf("%d. *<%s|%s>: %s*\n\t%s %s: %s\n\t",
				i+1,
				item.Link,
				item.IssueKey,
				item.Summary,
				issueKindEmoji,
				item.Kind,
				impactMap[item.IssueKey],
			)

			msg = slack.AddBlockMessage(msg, slack.NewSectionBlock(
				slack.NewTextBlockObject(slack.MarkdownType, markdownText, false, false),
				nil,
				nil,
			))

			msg = slack.AddBlockMessage(msg, slack.NewContextBlock("",
				slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("*Reporter*: %s *Assignee*: %s", item.Reporter, item.Reporter), false, false),
			))
		}
	}

	_, _, err := sm.client.PostMessage(channel, slack.MsgOptionBlocks(msg.Blocks.BlockSet...))
	return err
}
