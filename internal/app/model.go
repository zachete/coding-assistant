package app

import (
	"coding-assistant/internal/agent"
	"coding-assistant/internal/ui/confirmation"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/glamour/v2/styles"
	"charm.land/lipgloss/v2"
)

func initialModel(agent agent.Agent) *model {
	const (
		width  = 80
		height = 20
	)
	greeting := `**Hallo, Hallo!**`

	vp := viewport.New(viewport.WithWidth(width), viewport.WithHeight(height))
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9c7cff")).
		PaddingRight(2)
	vp.KeyMap.PageDown.Unbind()

	ti := textinput.New()
	ti.Placeholder = "Ask something"
	ti.Focus()
	ti.Prompt = "> "

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5fef"))

	confirm := confirmation.New()

	const glamourGutter = 3
	glamourRenderWidth := width - vp.Style.GetHorizontalFrameSize() - glamourGutter
	style := styles.DarkStyleConfig
	glam, _ := glamour.NewTermRenderer(
		glamour.WithStyles(style),
		glamour.WithWordWrap(glamourRenderWidth),
	)

	render, _ := glam.Render(greeting)
	history := render

	vp.SetContent(history)

	go agent.Run()

	return &model{
		textInput:    ti,
		spinner:      s,
		viewPort:     vp,
		agent:        agent,
		confirmation: confirm,
		glamRender:   glam,
		history:      history,
		pending:      false,
	}
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textInput.SetWidth(msg.Width)
		m.viewPort.SetWidth(msg.Width)
		m.viewPort.SetHeight(msg.Height - 3)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if !m.confirmation.Active {
				message := "-> " + m.textInput.Value()
				m.history = appendToHistory(m.history, message)

				m.viewPort.GotoBottom()
				m.viewPort.SetContent(m.history)

				cmds = append(cmds, m.spinner.Tick)
				m.agent.CmdChan <- agent.Command{Kind: agent.PromptCommand, Text: m.textInput.Value()}

				m.textInput.SetValue("")
			}
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	case agent.Event:
		switch msg.Kind {
		case agent.AgentResponse:
			renderedText, _ := m.glamRender.Render(msg.Text)
			m.history = appendToHistory(m.history, renderedText)
			m.viewPort.SetContent(m.history)
			m.viewPort.GotoBottom()
			m.pending = false
		case agent.AgentStart:
			m.pending = true
		case agent.AgentToolConfirm:
			m.pending = false
			m.confirmation.Active = true
			m.confirmation.SetCallID(msg.ToolCall.CallID)
			m.confirmation.SetPrompt(fmt.Sprintf("Call tool (%s?)", msg.ToolCall.Name))
		case agent.AgentError:
			m.history = appendToHistory(m.history, errorHistoryStyle(msg.Error.Error()))
			m.pending = false
			m.viewPort.SetContent(m.history)
			m.viewPort.GotoBottom()
		}
	case confirmation.ConfirmYesMsg:
		m.agent.CmdChan <- agent.Command{Kind: agent.ConfirmCommand, ToolCallID: m.confirmation.CallID}
		m.confirmation.Active = false
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	m.viewPort, cmd = m.viewPort.Update(msg)

	if m.confirmation.Active {
		m.confirmation, cmd = m.confirmation.Update(msg)
		cmds = append(cmds, cmd)
	}

	cmd = m.listenAgent()
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() tea.View {
	var sb strings.Builder
	sb.WriteString(m.viewPort.View())
	sb.WriteString("\n")

	if m.pending {
		sb.WriteString(lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).Height(1).
			BorderForeground(lipgloss.Color("#ff5fef")).Width(m.viewPort.Width()).Render(m.spinner.View()))

	} else if m.confirmation.Active {
		sb.WriteString(lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).Height(1).
			BorderForeground(lipgloss.Color("#8888bb")).Width(m.viewPort.Width()).Render(m.confirmation.View().Content))

	} else {
		sb.WriteString(lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).Height(1).
			BorderForeground(lipgloss.Color("#ff5fef")).Width(m.viewPort.Width()).Render(m.textInput.View()))
	}

	v := tea.NewView(sb.String())
	v.AltScreen = true
	return v
}
