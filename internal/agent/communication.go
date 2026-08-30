package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
)

func (a *Agent) send() Event {
	a.EventChan <- Event{Kind: AgentStart}

	completion, err := a.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model:    a.model,
		Messages: a.state.messages,
		Tools:    ToOpenAITools(a.registry.AsSlice()),
	})

	if err != nil {
		return Event{Kind: AgentError, Error: fmt.Errorf("openai api error: %w", err)}
	}

	message := completion.Choices[0].Message
	a.state.messages = append(a.state.messages, message.ToParam())

	if len(message.ToolCalls) == 0 {
		return Event{Kind: AgentResponse, Text: message.Content}
	}

	for _, call := range message.ToolCalls {
		tool, err := a.registry.ResolveTool(call.Function.Name)
		if err != nil {
			continue
		}

		var args map[string]any
		err = json.Unmarshal([]byte(call.Function.Arguments), &args)
		if err != nil {
			return Event{Kind: AgentError, Error: fmt.Errorf("json unmarshal error: %w", err)}
		}

		tool.SetArgs(args)

		a.state.toolCalls[call.ID] = ToolCall{
			CallID: call.ID,
			Tool:   tool,
		}

		if tool.NeedConfirm() {
			a.EventChan <- Event{Kind: AgentToolConfirm, ToolCall: a.state.toolCalls[call.ID]}
		} else {
			a.EventChan <- Event{Kind: AgentToolCall, ToolCall: a.state.toolCalls[call.ID]}
			a.CmdChan <- Command{Kind: ToolCallCommand, ToolCall: a.state.toolCalls[call.ID]}
		}

	}

	return Event{Kind: AgentFinish}
}

func (a *Agent) sendPrompt(prompt string) Event {
	a.state.messages = append(a.state.messages, openai.UserMessage(prompt))
	return a.send()
}
