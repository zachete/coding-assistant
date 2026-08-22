package agent

import (
	"context"
	"encoding/json"
	"os"

	"coding-assistant/internal/tools"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const DefaultModel = "gemini/gemini-3.1-flash-lite"
const MaxSteps = 5

func Prompt(prompt string) string {
	apiKey := os.Getenv("API_KEY")
	baseURL := os.Getenv("API_BASE_URL")
	model := os.Getenv("AGENT_MODEL")

	if apiKey == "" {
		panic("API_KEY environment variable is not set")
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if model == "" {
		model = DefaultModel
	}

	registry := tools.NewRegistry()
	toolsList := ToOpenAITools(registry.AsSlice())
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
	)
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a senior software engineer. You help with code analyzing. You have access to tools that can read files from the codebase. When you need to analyze a file, use the read_file tool to get its content."),
		openai.UserMessage(prompt),
	}

	for _ = range 5 {
		completion, err := client.Chat.Completions.New(
			context.Background(),
			openai.ChatCompletionNewParams{
				Model:    model,
				Messages: messages,
				Tools:    toolsList,
			},
		)
		if err != nil {
			panic(err)
		}
		message := completion.Choices[0].Message

		messages = append(messages, message.ToParam())

		if len(message.ToolCalls) == 0 {
			break
		}

		for _, call := range message.ToolCalls {
			tool, err := registry.ResolveTool(call.Function.Name)
			if err != nil {
				continue
			}

			var args map[string]any
			err = json.Unmarshal([]byte(call.Function.Arguments), &args)
			if err != nil {
				panic(err)
			}

			content, err := tool.Execute(args)
			if err != nil {
				messages = append(messages, openai.ToolMessage("Error executing tool: "+err.Error(), call.ID))
				continue
			}

			messages = append(messages, openai.ToolMessage(content, call.ID))
		}
	}

	return messages[len(messages)-1].OfAssistant.Content.OfString.Value
}
