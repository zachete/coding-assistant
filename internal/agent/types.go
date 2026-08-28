package agent

import (
	"coding-assistant/internal/tools"

	"github.com/openai/openai-go/v3"
)

type ToolCall struct {
	Name   string
	CallID string
	tool   tools.Tool
	args   map[string]any
}

type State struct {
	step      int
	messages  []openai.ChatCompletionMessageParamUnion
	toolCalls map[string]ToolCall
}

type Event struct {
	Kind     EventKind
	ToolCall ToolCall
	Text     string
	Error    error
}

type EventKind int

const (
	AgentStart EventKind = iota
	AgentResponse
	AgentToolConfirm
	AgentFinish
	AgentError
)

type Command struct {
	Kind       CommandKind
	Text       string
	ToolCallID string
}

type CommandKind int

const (
	PromptCommand CommandKind = iota
	ConfirmCommand
)

type Agent struct {
	CmdChan   chan Command
	EventChan chan Event
	queue     []string
	registry  *tools.Registry
	client    openai.Client
	model     string
	state     State
}
