package jira

import (
	"fmt"
	"sort"
	"time"

	"github.com/andygrunwald/go-jira"

	"github.com/makarski/progress-bot/config"
)

type (
	IssueReport struct {
		From    time.Time                       `json:"from"`
		To      time.Time                       `json:"to"`
		Changes map[string][]*IssueStatusChange `json:"changes"`
	}

	IssueStatusChange struct {
		Summary     string              `json:"summary"`
		Description string              `json:"description"`
		Reporter    string              `json:"reporter"`
		Assignee    string              `json:"assignee"`
		Components  []map[string]string `json:"components"`
		IssueKey    string              `json:"issue_key"`
		Link        string              `json:"link"`
		Kind        string              `json:"kind"`
		Date        time.Time           `json:"date"`
		From        string              `json:"from"`
		To          string              `json:"to"`
	}
)

// JiraViewer is an API wrapper for JIRA
type JiraViewer struct {
	jiraClient *jira.Client
	baseURL    string
}

// NewJiraViewer returns a new JiraViewer instance.
func NewJiraViewer(cfg *config.Jira) (*JiraViewer, error) {
	tp := jira.BasicAuthTransport{
		Username: cfg.User,
		Password: cfg.Token,
	}

	jiraClient, err := jira.NewClient(tp.Client(), cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to init jira client: %s", err)
	}

	return &JiraViewer{jiraClient, cfg.BaseURL}, nil
}

// ListUpdatedSince returns a list of issues updated in the last `since` duration.
func (jv *JiraViewer) ListUpdatedSince(project string, since time.Duration) ([]jira.Issue, error) {
	sinceStr := fmt.Sprintf("%dd", int(since.Hours()/24))
	jql := fmt.Sprintf(`project="%s"&"updatedDate" >= -%s ORDER BY "updated" DESC`, project, sinceStr)

	issues, _, err := jv.jiraClient.Issue.Search(jql, &jira.SearchOptions{Expand: "changelog"})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues for project: %s. %s", project, err)
	}

	return issues, nil
}

// GroupByTargetStatus groups the issues by the target status.
func GroupByTargetStatus(issues []*IssueStatusChange) IssueReport {
	grouped := make(map[string][]*IssueStatusChange, len(issues))

	earliestDate := issues[len(issues)-1].Date
	latestDate := issues[0].Date

	for _, issue := range issues {
		grouped[issue.To] = append(grouped[issue.To], issue)
	}

	return IssueReport{
		From:    earliestDate,
		To:      latestDate,
		Changes: grouped,
	}
}

// FilterWithStatusChange filters the issues with status change.
func (jv *JiraViewer) FilterWithStatusChange(issues []jira.Issue, compareToDate time.Time) ([]*IssueStatusChange, error) {
	filtered := make([]*IssueStatusChange, 0, len(issues)/2)
	filteredMap := make(map[string]*IssueStatusChange, len(issues)/2)
	const statusField = "status"

	for _, issue := range issues {
		if len(issue.Changelog.Histories) < 1 {
			continue
		}

		histories := issue.Changelog.Histories
		sort.Slice(histories, func(i, j int) bool {
			return histories[i].Created < histories[j].Created
		})

		for _, history := range histories {
			if len(history.Items) < 1 {
				continue
			}

			for _, item := range history.Items {
				if item.Field != statusField {
					continue
				}

				changeDate, err := time.Parse("2006-01-02T15:04:05-0700", history.Created)
				if err != nil {
					return nil, fmt.Errorf("failed to parse time: %s. %s. %s", history.Created, err, issue.Key)
				}

				if changeDate.Before(compareToDate) {
					continue
				}

				assignee := "-"
				if issue.Fields.Assignee != nil {
					assignee = issue.Fields.Assignee.DisplayName
				}

				components := make([]map[string]string, 0, len(issue.Fields.Components))
				for _, component := range issue.Fields.Components {
					components = append(components, map[string]string{
						"name":        component.Name,
						"description": component.Description,
					})
				}

				issueStatusChange, ok := filteredMap[issue.Key]
				if !ok {
					statusChange := &IssueStatusChange{
						Summary:     issue.Fields.Summary,
						Description: issue.Fields.Description,
						Components:  components,
						Reporter:    issue.Fields.Reporter.DisplayName,
						Assignee:    assignee,
						IssueKey:    issue.Key,
						Link:        fmt.Sprintf("%s/browse/%s", jv.baseURL, issue.Key),
						Kind:        issue.Fields.Type.Name,
						Date:        changeDate,
						From:        item.FromString,
						To:          item.ToString,
					}

					filteredMap[issue.Key] = statusChange
					filtered = append(filtered, statusChange)
					continue
				}

				issueStatusChange.Summary = issue.Fields.Summary
				issueStatusChange.Description = issue.Fields.Description
				issueStatusChange.Date = changeDate
				issueStatusChange.To = item.ToString
				issueStatusChange.Reporter = issue.Fields.Reporter.DisplayName
				issueStatusChange.Assignee = assignee
				issueStatusChange.Kind = issue.Fields.Type.Name
				issueStatusChange.Components = components
			}
		}
	}

	return filtered, nil
}
