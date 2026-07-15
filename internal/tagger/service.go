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
