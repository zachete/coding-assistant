package agent

import (
	"coding-assistant/internal/config"
	"coding-assistant/internal/tools"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func New(config config.Config) Agent {
	return Agent{
		CmdChan:   make(chan Command, 10),
		EventChan: make(chan Event, 10),
		registry:  tools.NewRegistry(),
		model:     config.AgentModel,
		client: openai.NewClient(
			option.WithBaseURL(config.ApiBaseUrl),
			option.WithAPIKey(config.ApiKey),
		),
		state: State{
			messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are a senior software engineer. You help with code analyzing. You have access to tools that can read files from the codebase. When you need to analyze a file, use the read_file tool to get its content."),
			},
			toolCalls: make(map[string]ToolCall),
		},
	}
}

func (a *Agent) Run() {
	for {
		cmd := <-a.CmdChan
		switch cmd.Kind {
		case PromptCommand:
			a.EventChan <- a.sendPrompt(cmd.Text)
		case ToolCallCommand:
			call, ok := a.state.toolCalls[cmd.ToolCall.CallID]
			if !ok {
				a.EventChan <- Event{Kind: AgentError, Error: fmt.Errorf("Can't find a tool")}
				continue
			}

			res, err := call.Tool.Execute()
			delete(a.state.toolCalls, cmd.ToolCall.CallID)

			if err != nil {
				a.EventChan <- Event{Kind: AgentError, Error: err}
				continue
			}

			a.state.messages = append(a.state.messages, openai.ToolMessage(res, cmd.ToolCall.CallID))

			a.EventChan <- a.send()
		}
	}
}
