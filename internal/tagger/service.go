package tagger

import (
	"errors"
	"fmt"
	"tagger/internal/fetcher"
	"tagger/internal/filer"
)

type Repo interface {
	InsertPaths(paths []string) error //new album paths to directory table
	SelectPath(status ImportStatus) (*string, error)

	// SelectAlbum(id int64) (Album, error)
	// InsertAlbum(in Album) (int64, error)
	// UpdateAlbum(id int64) error
	// DeleteAlbum(id int64) error

	SelectArtistByID(id int64) (*Artist, error)
	SelectArtistByName(name string) (*Artist, error)
	// ListArtists(query string) ([]Artist, error)
	// InsertArtist(in Artist) (int64, error)
	// UpdateArtist(id int64) error
	// DeleteArtist(id int64) error

	// SelectTrack(id int64) (Track, error)
	// ListTracks(albumID int64) ([]Track, error)
	// InsertTrack(in Track) (int64, error)
	// UpdateTrack(id int64) error
	// DeleteTrack(id int64) error

	// SelectGenre(id int64) (Genre, error)
	// ListGenres(query string) ([]Genre, error)
	// InsertGenre(in Genre) (int64, error)
	// UpdateGenre(id int64) error
	// DeleteGenre(id int64) error
}

type Service struct {
	repo   Repo
	client *fetcher.BrainzClient
}

func NewService(r Repo, c *fetcher.BrainzClient) *Service {
	return &Service{repo: r, client: c}
}

func (s *Service) InitializeDirectory(root string) error {
	folders, err := findAllAlbumFolders(root)
	if err != nil {
		return err
	}

	if err := s.repo.InsertPaths(folders); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetMusicBrainzMatch(album StagedAlbum) error {

}

func (s *Service) ResumeImport() error {
	// first try albums that are not fully imported
	trying := ImportStatusPendingInput
	err := s.ImportAlbum(trying)
	if !errors.Is(err, ErrNoRows) {
		return err
	}

	// then try albums that have not been attempted imported yet
	trying = ImportStatusInitialized
	err = s.ImportAlbum(trying)
	if !errors.Is(err, ErrNoRows) {
		return err
	}

	return ErrNoRows
}

func (s *Service) ImportAlbum(status ImportStatus) error {
	path, err := s.repo.SelectPath(status)
	if errors.Is(err, ErrNoRows) {
		return ErrNoRows
	} else if err != nil {
		return err
	}

	if err := filer.LoadAlbum(*path); err != nil {
		return fmt.Errorf("could not load album at path %s: %w", *path, err)
	}

	return nil
}

func (s *Service) GetRelease(id string) (fetcher.Album, error) {
	if id == "" {
		return fetcher.Album{}, errors.New("id was empty")
	}

	return s.client.GetReleaseByID(id)
}

func (s *Service) GetReleaseGroup(id filer.Identifiers) ([]fetcher.ReleaseInfo, error) {
	if id.ReleaseGroupID != "" {
		releases, err := s.client.GetReleaseGroup(id.ReleaseGroupID)
		if len(releases) > 0 {
			return releases, nil
		}
		if err != nil && !errors.Is(err, fetcher.ErrNotFound) {
			return nil, err
		}
	}

	if id.AlbumName != "" && id.AlbumArtist != "" {
		return s.client.SearchReleaseGroup(id.AlbumName, id.AlbumArtist)
	}

	return nil, errors.New("release_group_id and album/artist was empty")
}

func ParseAlbum(in fetcher.Album) Album {
	var a Album
	a.MBReleaseID = in.MBReleaseID
	a.MBReleaseGroupID = in.MBReleaseGroupID
	a.Name = in.Name

	// check this against zero value and request user input if true
	a.ReleaseDate = in.ReleaseDate

}

func (s *Service) parseLabel(in fetcher.LabelInfo) (Label, error) {
	var out Label
	out, err := s.repo.SelectLabelByMBID(in.MBID)
	if errors.Is(err, ErrNoRows) {
		out.MBID = in.MBID
		out.Name = in.Name
		out.Slug = toSlug(in.Name)
		id, err := s.repo.InsertLabel(out)
		if err != nil {
			return out, err
		}

		out.ID = id
	}

	return out, nil
}
