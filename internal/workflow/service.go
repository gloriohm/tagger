package workflow

import (
	"errors"
	"tagger/internal/compare"
	"tagger/internal/domain"
	"tagger/internal/mbfetcher"
	"tagger/internal/tag_loader"
	"tagger/internal/ui"
)

type Workflow struct {
	mbsvc     mbfetcher.BrainzClient
	domainsvc domain.Service
	user      UserSelector
}

type Job struct {
	path            string
	stage           ImportStage
	tags            tag_loader.AlbumIdentifiers
	release         mbfetcher.Album
	releaseGroup    []mbfetcher.ReleaseInfo
	comparision     compare.Comparision
	genres          []domain.Genre
	InternalAlbumID int64
}

type ImportStage int

const (
	LoadTags ImportStage = iota
	TryReleaseID
	TryReleaseGroup
	UserSelectRelease
	ValidateMusicBrainzResult
	UserReviewRelease
	UserSelectGenres
	CreateDomainModel
	TagFiles
	MoveFolder
	ImportComplete
)

func NewWorkflow(mbclient mbfetcher.BrainzClient) *Workflow {
	return &Workflow{mbsvc: mbclient}
}

func (wf *Workflow) NewImportJob(job *Job) error {
	for job.stage != ImportComplete {
		if err := wf.advanceJob(job); err != nil {
			return err
		}
	}

	return nil
}

func (wf *Workflow) advanceJob(job *Job) error {
	switch job.stage {
	case LoadTags:
		return wf.loadTags(job)
	case TryReleaseID:
		return wf.tryReleaseID(job)
	case TryReleaseGroup:
		return wf.tryReleaseGroup(job)
	case UserSelectRelease:
		return wf.userSelectRelease(job)
	case ValidateMusicBrainzResult:
		return wf.validateMusicBrainzResult(job)
	case UserReviewRelease:
		return wf.userReviewRelease(job)
	case CreateDomainModel:
		return wf.createDomainModel(job)
	case TagFiles:
		return wf.tagFiles(job)
	case MoveFolder:
		return wf.moveFolder(job)
	default:
		return ErrUnknownStage
	}
}

func (wf *Workflow) loadTags(job *Job) error {
	tags, err := tag_loader.LoadAlbum(job.path)
	if err != nil {
		return err
	}

	// consider checking most specific id and assigning stage directly
	// can avoid the fallback logic entirely

	job.tags = tags
	job.stage = TryReleaseID
	return nil
}

func (wf *Workflow) tryReleaseID(job *Job) error {
	release, err := wf.mbsvc.GetReleaseByID(job.tags.ReleaseID)
	if err != nil {
		if errors.Is(err, mbfetcher.ErrNotFound) {
			job.stage = TryReleaseGroup
			return nil
		} else {
			return err
		}
	}

	job.release = release
	job.stage = ValidateMusicBrainzResult

	return nil
}

func (wf *Workflow) tryReleaseGroup(job *Job) error {
	releaseGroup, err := wf.mbsvc.GetReleaseGroup(job.tags.ReleaseGroupID)
	if err != nil && !errors.Is(err, mbfetcher.ErrNotFound) {
		return err
	}

	if len(releaseGroup) > 0 {
		job.releaseGroup = releaseGroup
		job.stage = UserSelectRelease

		return nil // return early if release group was found by id
	}

	releaseGroup, err = wf.mbsvc.SearchReleaseGroup(job.tags.AlbumArtist, job.tags.AlbumName)
	if err != nil && !errors.Is(err, mbfetcher.ErrNotFound) {
		return err
	}

	job.releaseGroup = releaseGroup
	job.stage = UserSelectRelease

	return nil
}

func (wf *Workflow) userSelectRelease(job *Job) error {
	releaseID, err := ui.SelectRelease(job.releaseGroup)
	if err != nil {
		return err
	}

	release, err := wf.mbsvc.GetReleaseByID(releaseID)
	if err != nil {
		if errors.Is(err, mbfetcher.ErrNotFound) {
			// repeats step if not found. means user input invalid id
			// should be updated to include user feedback
			return nil
		}

		return err
	}

	job.release = release
	job.stage = ValidateMusicBrainzResult

	return nil
}

func (wf *Workflow) validateMusicBrainzResult(job *Job) error {
	localAlbum := localAlbumToCompare(job.tags)
	mbAlbum := mbAlbumToCompare(job.release)

	comparision := compare.CompareAlbum(localAlbum, mbAlbum)

	if comparision.Distance < 0.05 {
		job.stage = UserSelectGenres
	} else {
		job.comparision = comparision
		job.stage = UserReviewRelease
	}

	return nil
}

func (wf *Workflow) userReviewRelease(job *Job) error {
	approved, err := ui.ReviewMatch(job.comparision)
	if err != nil {
		return err
	}

	if approved {
		job.stage = UserSelectGenres
	} else {
		job.stage = UserSelectRelease
	}

	return nil
}

func (wf *Workflow) userInputGenre(job *Job) error {
	genres, err := wf.domainsvc.GetGenres()
	if err != nil {
		return err
	}
	selectedGenres, err := ui.SelectGenres(genres)
	if err != nil {
		return err
	}

	job.genres = selectedGenres
	job.stage = CreateDomainModel

	return nil
}

func (wf *Workflow) createDomainModel(job *Job) error {
	id, err := wf.domainsvc.CreateAlbum(job.release, job.genres)
	if err != nil {
		return err
	}

	job.InternalAlbumID = id
	job.stage = TagFiles

	return nil
}

func (wf *Workflow) tagFiles(job *Job) error {
	albumTags, err := wf.domainsvc.GetAlbumTags(job.InternalAlbumID)
	if err != nil {
		return err
	}

	if err := tag_loader.TagAlbum(job.path, albumTags); err != nil {
		return err
	}

	job.stage = MoveFolder

	return nil
}

func (wf *Workflow) moveFolder(job *Job) error {
	return nil
}
