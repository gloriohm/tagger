package tagger

import "errors"

type Orchestrator struct {
	store Store
}

type Store interface {
	ResumeImport() error
}

type UserAction int

const (
	NoAction UserAction = iota
	ProvideMatch
	EnterGenres
)

type ImportStage int

const (
	PathStored ImportStage = iota
	MetaValidated
	MetaMatched
	DomainModelCreated
	DomainModelValidated
	FilesTagged
)

type ActionObject struct {
	Action      UserAction
	Description []string
}

func (o *Orchestrator) Run() error {
	for {
		err := o.store.ResumeImport()
		if errors.Is(err, ErrNoRows) {
			// no albums awaiting import
			return nil
		} else if err != nil {
			return err
		}

	}
}

func (s *Service) AdvanceAlbum(album StagedAlbum) error {
	switch album.Stage {
	case StageImportStarted:
		return s.GetMusicBrainzMatch(album)
	}
}

func (o *Orchestrator) ()