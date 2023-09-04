package openai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// dailySystemPrompt is the prompt used to generate the daily report.
const dailySystemPrompt = `Please generate a daily report that serves as a communication on behalf of the Engineering Team Lead. Your task is to read a JSON-serialized summary that outlines the team's progress over a period of time and to estimate the areas impacted by each issue or task.

Instructions:
	
1. For each issue or task, try to identify the areas that it impacts (e.g., security, compliance, user experience, etc.).
2. Format your response as a JSON object.
3. Use the key 'issue_key' for each issue or task (e.g., "LI-607", "LI-617").
4. For each issue or task, provide a brief description of the guessed impact as the value.
Include a 'summary' key in the JSON object to provide an overall summary of the team's activities and the areas impacted.

Example JSON Response:
{
	"summary": "Over the reported period, the team completed 3 tasks and resolved 1 bug. The areas impacted are: security, compliance, user access, and permissions.",
	"LI-607": "This task involved managing roles in the identity management system and primarily impacts the areas of user access and permissions.",
	"LI-617": "This task focused on tracking login activities, affecting the system's security and compliance."
}`

// OpenAI is a wrapper around the OpenAI API.
type OpenAI struct {
	client   *openai.Client
	gptModel string
}

// NewOpenAI returns a new OpenAI instance.
func NewOpenAI(apiKey, gptModel string) *OpenAI {
	return &OpenAI{
		client:   openai.NewClient(apiKey),
		gptModel: gptModel,
	}
}

// CompileReport compiles a daily report from the given JSON summary.
func (oai *OpenAI) CompileReport(ctx context.Context, jsonb []byte) (string, error) {
	prompt := dailySystemPrompt

	completions := make([]openai.ChatCompletionMessage, 1, 2)
	completions[0] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: prompt,
	}

	completions = append(completions, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: string(jsonb),
	})

	resp, err := oai.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    oai.gptModel,
			Messages: completions,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %s", err)
	}

	return resp.Choices[0].Message.Content, nil
}
