package ui

import (
	"tagger/internal/tagger"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	app     app
	screen  screen
	cursor  int
	choices []choice
}

type choice struct {
	label  string
	action action
}

type action int

const (
	selectImport action = iota
	selectSearch
	selectQuit
)

type screen int

const (
	homeScreen screen = iota
)

type app interface {
	ImportAlbum(tagger.ImportStatus) error
}

func NewModel(app app) Model {
	return Model{
		app: app,
	}
}

func Run(app app) error {
	model := NewModel(app)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
