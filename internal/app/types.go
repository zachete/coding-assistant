package app

import (
	"coding-assistant/internal/agent"
	"coding-assistant/internal/ui/confirmation"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	"charm.land/glamour/v2"
)

type model struct {
	agent        agent.Agent
	textInput    textinput.Model
	spinner      spinner.Model
	confirmation confirmation.Model
	viewPort     viewport.Model
	glamRender   *glamour.TermRenderer
	// messages     []string
	history string
	pending bool
}
