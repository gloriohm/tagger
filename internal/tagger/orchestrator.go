package tagger

import "errors"

type Orchestrator struct {
	store Store
	svc   *Service
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

type ImportState int

const (
	StateImportStarted ImportState = iota
	StateInvalidAlbum
	StateMetaValidated
	StateMetaMatched
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

func (o *Orchestrator) AdvanceAlbum(jobID int) error {
	state, err := o.svc.GetJobState
	if err != nil {
		return err
	}
	switch state {
	case StateImportStarted:
		return o.svc.GetMusicBrainzMatch(album)
	}
}
