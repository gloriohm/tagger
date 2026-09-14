package ui

import tea "charm.land/bubbletea/v2"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", "space":

		}
	}
	return m, nil
}

func (m Model) handleEnter(selected action) (tea.Model, tea.Cmd) {
	switch selected {
	case selectQuit:
		return m, tea.Quit
	}

	return m, nil
}
