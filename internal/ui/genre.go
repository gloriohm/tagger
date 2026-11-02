package ui

import "github.com/charmbracelet/huh"

func SelectGenres(genres []string) ([]string, error) {
	var selected []string

	options := make([]huh.Option[string], 0, len(genres))
	for _, genre := range genres {
		options = append(options, huh.NewOption(genre, genre))
	}

	err := huh.NewMultiSelect[string]().
		Title("Select album genres").
		Description("Press / to search, Space to select").
		Options(options...).
		Filterable(true).
		Height(10).
		Value(&selected).
		Run()

	return selected, err
}
