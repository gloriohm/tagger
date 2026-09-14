package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

func (m Model) View() tea.View {
	switch m.screen {
	case homeScreen:
		return m.renderHome()
	}

	return tea.NewView("This view should never be reached...\n\n")
}

var homeChoices = [...]choice{{label: "Import", action: selectImport}, {label: "Search", action: selectSearch}, {label: "Quit", action: selectQuit}}

func (m Model) renderHome() tea.View {
	content := "Welcome to Tagger\n\n"

	for i, choice := range homeChoices {
		cursorStyle := " "
		if m.cursor == i {
			cursorStyle = ">"
		}

		content += fmt.Sprintf("%s %s", cursorStyle, choice.label)
	}

	return tea.NewView(content)
}
