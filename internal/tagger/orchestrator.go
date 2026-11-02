package tagger

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
	ValidateFolder
	KillProgram
	RetryMusicBrainzReleaseGroup
	ValidateMusicBrainzResult
	DisplayResult
	InputMusicBrainzReleaseID
)

type ActionObject struct {
	Action      UserAction
	Description []string
}

func (o *Orchestrator) RunImportService(path string) error {

}

func (o *Orchestrator) AdvanceJob(job Job) (Signal, error) {
	switch job.Stage {
	case StageImportStarted:
		return o.svc.GetMusicBrainzMatch(job)
	}
}

func (o *Orchestrator) Run() error {
	for {
		signal, err := o.AdvanceJob(job)
	}
}

type Job struct {
	ID    int
	Stage ImportStage
	Path  string
}

type Signal struct {
	Job        Job
	UserAction UserAction
}
