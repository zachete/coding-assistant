package app

import "charm.land/lipgloss/v2"

func errorHistoryStyle(str string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Red).Render(str)
}

func toolCallHistoryStyle(str string) string {
	var greenPoint = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("•")
	return greenPoint + " " + str
}

func appendToHistory(input string, str string) string {
	return input + "\n" + str + "\n"
}
