package agent

import (
	"coding-assistant/internal/tools"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

func ToOpenAITools(toolsList []tools.Tool) []openai.ChatCompletionToolUnionParam {
	var openAITools []openai.ChatCompletionToolUnionParam
	for _, tool := range toolsList {
		openAITools = append(openAITools, openai.ChatCompletionToolUnionParam{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: shared.FunctionDefinitionParam{
					Name:        tool.Name(),
					Description: openai.String(tool.Description()),
					Parameters:  tool.Params(),
					Strict:      openai.Bool(false),
				},
			},
		})
	}
	return openAITools
}
