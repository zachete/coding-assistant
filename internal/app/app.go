package app

import (
	"coding-assistant/internal/agent"
	"coding-assistant/internal/config"

	tea "charm.land/bubbletea/v2"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	config := config.NewConfig()
	agent := agent.New(config)
	p := tea.NewProgram(initialModel(agent))
	p.Run()
}
