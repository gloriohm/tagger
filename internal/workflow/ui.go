package workflow

import (
	"tagger/internal/compare"
	"tagger/internal/mbfetcher"
)

type ReleaseSelection struct {
	ReleaseID string
}

type Decision struct {
	Approved  bool
	ReleaseID string
}

type UserSelector interface {
	SelectRelease(releases []mbfetcher.ReleaseInfo) (ReleaseSelection, error)
	ReviewRelease(comp compare.Comparision) (Decision, error)
}
