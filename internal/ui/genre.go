package ui

import (
	"errors"
	"fmt"
	"strings"
	"tagger/internal/compare"
	"tagger/internal/mbfetcher"

	"github.com/charmbracelet/huh"
)

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

func SelectRelease(releases []mbfetcher.ReleaseInfo) (string, error) {
	const manual = "__manual__"
	var selected string

	options := make([]huh.Option[string], 0, len(releases)+1)

	for _, release := range releases {
		options = append(options, huh.NewOption(release.Quality, release.ID))
	}

	options = append(options, huh.NewOption("Enter release ID manually", manual))

	err := huh.NewSelect[string]().
		Title("Select a release").
		Options(options...).
		Value(&selected).
		Run()

	if selected == manual {
		return InputReleaseID()
	}

	return selected, err
}

func InputReleaseID() (string, error) {
	var input string

	err := huh.NewInput().
		Title("Input MusicBrainz Release ID").
		Value(&input).
		Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return errors.New("release ID cannot be empty")
			}
			return nil
		}).
		Run()

	return strings.TrimSpace(input), err
}

func ReviewMatch(c compare.Comparision) (bool, error) {
	var approved bool
	var details strings.Builder

	fmt.Fprintf(&details,
		"Match: %.1f%%\n\n",
		(1-c.Distance)*100,
	)

	for _, diff := range c.Differences {
		fmt.Fprintf(&details,
			"%s: %.2f\n",
			diff.Field,
			diff.Distance,
		)
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Release discrepancies").
				Description(details.String()),

			huh.NewConfirm().
				Title("Use this release?").
				Affirmative("Approve").
				Negative("Reject").
				Value(&approved),
		),
	).Run()

	return approved, err
}
