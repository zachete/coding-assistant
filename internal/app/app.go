package app

import (
	tea "charm.land/bubbletea/v2"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	p := tea.NewProgram(initialModel())
	p.Run()
}
