package app

import (
	"coding-assistant/internal/agent"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	textInput textinput.Model
	spinner   spinner.Model
	history   []string
	viewPort  viewport.Model
	pending   bool
}

type responseMsg string

func prompt(prompt string) tea.Cmd {
	return func() tea.Msg {
		res := agent.Prompt(prompt)
		return responseMsg(res)
	}
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Ask something"
	ti.Focus()
	ti.Prompt = "> "

	history := []string{lipgloss.NewStyle().Foreground(lipgloss.Yellow).Render("Hallo, Hallo!")}
	vp := viewport.New(viewport.WithWidth(10), viewport.WithHeight(10))
	vp.SetContent(history[0])
	vp.SoftWrap = true

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{textInput: ti, spinner: s, history: history, viewPort: vp, pending: false}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case responseMsg:
		m.history = append(m.history, string(msg))
		m.pending = false
		m.viewPort.SetContent(strings.Join(m.history, "\n\n") + "\n")
		m.viewPort.GotoBottom()
	case tea.WindowSizeMsg:
		m.textInput.SetWidth(msg.Width)
		m.viewPort.SetWidth(msg.Width)
		m.viewPort.SetHeight(msg.Height - 3)
		m.textInput.SetWidth(msg.Width)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.history = append(m.history, "-> "+m.textInput.Value())

			m.viewPort.SetContent(strings.Join(m.history, "\n\n") + "\n")
			m.viewPort.GotoBottom()

			cmds = append(cmds, m.spinner.Tick)
			cmds = append(cmds, prompt(m.textInput.Value()))

			m.pending = true
			m.textInput.SetValue("")
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	m.viewPort, cmd = m.viewPort.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var sb strings.Builder
	sb.WriteString(m.viewPort.View())
	sb.WriteString("\n")

	var inputOrLoader string
	if m.pending {
		inputOrLoader = m.spinner.View()
	} else {
		inputOrLoader = m.textInput.View()
	}

	sb.WriteString(lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).Height(1).
		BorderForeground(lipgloss.Color("205")).Width(m.viewPort.Width()).Render(inputOrLoader))

	v := tea.NewView(sb.String())
	v.AltScreen = true
	return v
}
