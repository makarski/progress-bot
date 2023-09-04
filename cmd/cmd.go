package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/makarski/progress-bot/config"
	"github.com/makarski/progress-bot/jira"
	"github.com/makarski/progress-bot/openai"
	"github.com/makarski/progress-bot/slack"
)

// Run executes the command.
func Run() error {
	cfg, err := config.LoadConfig(config.DefaultFileName)
	if err != nil {
		return err
	}

	jv, err := jira.NewJiraViewer(cfg.Jira)
	if err != nil {
		return err
	}

	sinceDuration, err := cfg.SinceDuration()
	if err != nil {
		return err
	}

	issues, err := jv.ListUpdatedSince(cfg.Jira.Project, sinceDuration)
	if err != nil {
		return err
	}

	compareToDate := time.Now().Add(-sinceDuration)
	filtered, err := jv.FilterWithStatusChange(issues, compareToDate)
	if err != nil {
		return err
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Date.After(filtered[j].Date)
	})

	issueReport := jira.GroupByTargetStatus(filtered)

	b, err := json.MarshalIndent(issueReport, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal issue report: %s", err)
	}

	reporter := openai.NewOpenAI(cfg.OpenAI.Token, cfg.OpenAI.GPTModel)
	report, err := reporter.CompileReport(context.Background(), b)
	if err != nil {
		return err
	}

	mappedImpact := make(map[string]string, len(issueReport.Changes))
	if err := json.Unmarshal([]byte(report), &mappedImpact); err != nil {
		return fmt.Errorf("failed to unmarshal ai generated report: %s", err)
	}

	msg := slack.NewSlackMessenger(cfg.Slack.Token, cfg.Slack.StatusesEmoji)
	return msg.SendMessage(cfg.Slack.Channel, issueReport, mappedImpact)
}
