package domain

import (
	"errors"
	"tagger/internal/mbfetcher"
)

type Repo interface {
	SelectAlbumByID(id int64) (Album, error)
	InsertAlbum(in Album) (Album, error)
	// UpdateAlbum(id int64) error
	// DeleteAlbum(id int64) error

	InsertAlbumArtistLinks(in Album) error
	InsertAlbumGenreLinks(in Album) error
	InsertAlbumLabelLinks(in Album) error

	SelectArtistByID(id int64) (Artist, error)
	SelectArtistByMBID(id string) (Artist, error)
	SelectArtistByName(name string) (Artist, error)
	PopulateAlbumArtists(album *Album) error
	// ListArtists(query string) ([]Artist, error)
	InsertArtist(in Artist) (Artist, error)
	// UpdateArtist(id int64) error
	// DeleteArtist(id int64) error

	// SelectTrack(id int64) (Track, error)
	// ListTracks(albumID int64) ([]Track, error)
	InsertTracks(in []Track) ([]Track, error)
	PopulateTracks(album *Album) error
	// UpdateTrack(id int64) error
	// DeleteTrack(id int64) error
	InsertTrackArtistLinks(in []Track) error
	PopulateTrackArtists(album *Album) error

	// SelectGenre(id int64) (Genre, error)
	ListGenres() ([]Genre, error)
	InsertGenre(in Genre) (Genre, error)
	PopulateGenres(album *Album) error
	// UpdateGenre(id int64) error
	// DeleteGenre(id int64) error

	InsertLabel(in Label) (Label, error)
	PopulateLabels(album *Album) error
	SelectLabelByMBID(id string) (Label, error)
}

type Service struct {
	repo   Repo
	client *mbfetcher.BrainzClient
}

func NewService(r Repo, c *mbfetcher.BrainzClient) *Service {
	return &Service{repo: r, client: c}
}

func (s *Service) CreateAlbum(release mbfetcher.Album, genres []Genre) (int64, error) {
	if err := s.createMissingGenres(genres); err != nil {
		return 0, err
	}

	labels := ParseLabels(release.Label)
	if err := s.resolveLabels(labels); err != nil {
		return 0, err
	}

	artists := ParseArtists(release.AlbumArtists)
	if err := s.resolveArtists(artists); err != nil {
		return 0, err
	}

	album := ParseAlbum(release)
	album, err := s.repo.InsertAlbum(album)
	if err != nil {
		return 0, err
	}

	album.AlbumArtists = artists
	album.Genres = genres
	album.Labels = labels
	if err := s.repo.InsertAlbumArtistLinks(album); err != nil {
		return 0, err
	}
	if err := s.repo.InsertAlbumGenreLinks(album); err != nil {
		return 0, err
	}
	if err := s.repo.InsertAlbumLabelLinks(album); err != nil {
		return 0, err
	}

	tracks := ParseTracks(release.Tracks, album.ID)
	tracks, err = s.repo.InsertTracks(tracks)

	if err := s.resolveTrackArtists(tracks); err != nil {
		return 0, err
	}

	if err := s.repo.InsertTrackArtistLinks(tracks); err != nil {
		return 0, err
	}

	return album.ID, nil
}

func (s *Service) GetAlbumTags(id int64) ([]map[string][]string, error) {
	album, err := s.GetAlbum(id)
	if err != nil {
		return nil, err
	}

	tags := AlbumToTagger(album)

	return tags, nil
}

func (s *Service) GetAlbum(id int64) (Album, error) {
	album, err := s.repo.SelectAlbumByID(id)
	if err != nil {
		return album, err
	}

	if err := s.repo.PopulateAlbumArtists(&album); err != nil {
		return album, err
	}
	if err := s.repo.PopulateGenres(&album); err != nil {
		return album, err
	}
	if err := s.repo.PopulateLabels(&album); err != nil {
		return album, err
	}
	if err := s.repo.PopulateTracks(&album); err != nil {
		return album, err
	}
	if err := s.repo.PopulateTrackArtists(&album); err != nil {
		return album, err
	}

	return album, nil
}

func (s *Service) createMissingGenres(genres []Genre) error {
	for i := range genres {
		if genres[i].ID == 0 {
			genre, err := s.repo.InsertGenre(genres[i])
			if err != nil {
				return err
			}

			genres[i] = genre
		}
	}

	return nil
}

func getOrCreate[T any](items []T, selectFn func(T) (T, error), insertFn func(T) (T, error)) error {
	for i := range items {
		item, err := selectFn(items[i])

		if errors.Is(err, ErrNoRows) {
			item, err = insertFn(items[i])
		}

		if err != nil {
			return err
		}

		items[i] = item
	}

	return nil
}

func (s *Service) resolveLabels(labels []Label) error {
	return getOrCreate(
		labels,
		func(label Label) (Label, error) {
			return s.repo.SelectLabelByMBID(label.MBID)
		},
		s.repo.InsertLabel,
	)
}

func (s *Service) resolveArtists(artists []Artist) error {
	return getOrCreate(
		artists,
		func(artist Artist) (Artist, error) {
			return s.repo.SelectArtistByMBID(artist.MBID)
		},
		s.repo.InsertArtist,
	)
}

// in domain model, artist should appear once per tracklist with array of indices linking to tracks they appear on
// artist order is handled by artistCredit string
func (s *Service) resolveTrackArtists(tracks []Track) error {
	for i := range tracks {
		if err := s.resolveArtists(tracks[i].Artists); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetGenres() ([]Genre, error) {
	return s.repo.ListGenres()
}
