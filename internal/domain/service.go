package domain

import (
	"errors"
	"tagger/internal/mbfetcher"
)

type Repo interface {
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

	SelectLabelByMBID(id string) (Label, error)
}

type Service struct {
	repo   Repo
	client *mbfetcher.BrainzClient
}

func NewService(r Repo, c *mbfetcher.BrainzClient) *Service {
	return &Service{repo: r, client: c}
}

func (s *Service) UpsertLabel(in mbfetcher.LabelInfo) (Label, error) {
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
