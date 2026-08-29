package confirmation

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ConfirmYesMsg struct{}
type ConfirmNoMsg struct{}

type Model struct {
	Prompt string
	Meta   any
	Yes    bool
	Active bool
}

func New() Model {
	return Model{Prompt: "Confirm?", Yes: true, Active: false}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetPrompt(prompt string) {
	m.Prompt = prompt
}

func (m *Model) SetMeta(meta any) {
	m.Meta = meta
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			return m, func() tea.Msg { return ConfirmYesMsg{} }
		case "n":
			return m, func() tea.Msg { return ConfirmNoMsg{} }
		case "right":
			m.Yes = false
		case "left":
			m.Yes = true
		case "enter":
			return m, func() tea.Msg {
				if m.Yes {
					return ConfirmYesMsg{}
				} else {
					return ConfirmNoMsg{}
				}
			}
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	var buttons string
	if m.Yes {
		buttons = lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("#666666")).Render("yes") + " " +
			lipgloss.NewStyle().Padding(0, 1).Render("no")
	} else {
		buttons = lipgloss.NewStyle().Padding(0, 1).Render("yes") + " " +
			lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("#666666")).Render("no")
	}

	return tea.NewView(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(
		m.Prompt,
	) + "  " + buttons)
}
