package app

import tea "charm.land/bubbletea/v2"

func (m *model) listenAgent() tea.Cmd {
	return func() tea.Msg {
		evt := <-m.agent.EventChan
		return evt
	}
}
