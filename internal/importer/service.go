package importer

import (
	"database/sql"
	"errors"
	"fmt"
	"tagger/internal/fetcher"

	"github.com/google/uuid"
)

type Service struct {
	repo Repo
	job  uuid.UUID
}

func NewService(r Repo) *Service {
	return &Service{repo: r}
}

type Repo interface {
	SelectArtistIDByMBID(mbid string) (int64, error)
	InsertArtist(in Artist) (int64, error)

	SelectLabelByMBID(in string) (int64, error)
	InsertLabel(in Label) (int64, error)

	InsertTrack(in Track) (int64, error)

	InsertAlbum(in Album) (int64, error)

	InsertAlbumArtists(in AlbumArtistMap) error
	InsertTrackArtists(in []TrackArtistMap) error
}

func (s *Service) RunCreateAlbumJob(in fetcher.Album) error {
	albumArtistsIDs, err := s.upsertArtists(in.AlbumArtists)
	if err != nil {
		return err
	}

	albumID, err := s.insertAlbum(in)
	if err != nil {
		return err
	}

	albumMap := AlbumArtistMap{AlbumID: albumID, ArtistIDs: albumArtistsIDs}
	if err := s.repo.InsertAlbumArtists(albumMap); err != nil {
		return err
	}

	trackMap, err := s.insertTracks(albumID, in.Tracks)
	if err != nil {
		return err
	}

	if err := s.repo.InsertTrackArtists(trackMap); err != nil {
		return err
	}

	return nil
}

func (s *Service) insertTracks(albumID int64, tracks []fetcher.Track) ([]TrackArtistMap, error) {
	var trackMap []TrackArtistMap
	for _, t := range tracks {
		trackID, err := s.repo.InsertTrack(parseTrack(t, albumID))
		if err != nil {
			return nil, err
		}

		artistIDs, err := s.upsertArtists(t.Artists)
		if err != nil {
			return nil, err
		}

		trackArtist := TrackArtistMap{TrackID: trackID, ArtistIDs: artistIDs}
		trackMap = append(trackMap, trackArtist)
	}

	return trackMap, nil
}

func (s *Service) upsertArtists(in []fetcher.Artist) ([]int64, error) {
	var artistIDs []int64
	for _, a := range in {
		id, err := s.repo.SelectArtistIDByMBID(a.MBID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				id, err := s.repo.InsertArtist(parseArtist(a))
				if err != nil {
					return artistIDs, fmt.Errorf("failed to insert artist: %w", err)
				} else {
					artistIDs = append(artistIDs, id)
				}
			}
		} else {
			artistIDs = append(artistIDs, id)
		}
	}

	return artistIDs, nil
}

func (s *Service) upsertLabel(in fetcher.LabelInfo) (int64, error) {
	id, err := s.repo.SelectLabelByMBID(in.MBID)
	if errors.Is(err, sql.ErrNoRows) {
		return s.repo.InsertLabel(parseLabel(in))
	}

	if err != nil {
		return 0, fmt.Errorf("failed to upsert label: %w", err)
	}

	return id, nil
}

func (s *Service) insertAlbum(in fetcher.Album) (int64, error) {
	album := parseAlbum(in)

	return s.repo.InsertAlbum(album)
}
