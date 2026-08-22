package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"coding-assistant/internal/config"
	"coding-assistant/internal/tools"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const MaxSteps = 5

type Agent struct {
	client   openai.Client
	registry *tools.Registry
	model    string
}

func NewAgent(config config.Config) Agent {
	return Agent{
		registry: tools.NewRegistry(),
		model:    config.AgentModel,
		client: openai.NewClient(
			option.WithBaseURL(config.ApiBaseUrl),
			option.WithAPIKey(config.ApiKey),
		),
	}
}

func (a Agent) Prompt(prompt string) (string, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a senior software engineer. You help with code analyzing. You have access to tools that can read files from the codebase. When you need to analyze a file, use the read_file tool to get its content."),
		openai.UserMessage(prompt),
	}

	for _ = range MaxSteps {
		completion, err := a.client.Chat.Completions.New(
			context.Background(),
			openai.ChatCompletionNewParams{
				Model:    a.model,
				Messages: messages,
				Tools:    ToOpenAITools(a.registry.AsSlice()),
			},
		)
		if err != nil {
			return "", fmt.Errorf("openai api error: %w", err)
		}
		message := completion.Choices[0].Message

		messages = append(messages, message.ToParam())

		if len(message.ToolCalls) == 0 {
			break
		}

		for _, call := range message.ToolCalls {
			tool, err := a.registry.ResolveTool(call.Function.Name)
			if err != nil {
				continue
			}

			var args map[string]any
			err = json.Unmarshal([]byte(call.Function.Arguments), &args)
			if err != nil {
				return "", fmt.Errorf("json unmarshal error: %w", err)
			}

			content, err := tool.Execute(args)
			if err != nil {
				messages = append(messages, openai.ToolMessage("Error executing tool: "+err.Error(), call.ID))
				continue
			}

			messages = append(messages, openai.ToolMessage(content, call.ID))
		}
	}

	return messages[len(messages)-1].OfAssistant.Content.OfString.Value, nil
}
