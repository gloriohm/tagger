package importer

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	JobID            uuid.UUID
	MBReleaseID      string
	MBReleaseGroupID string
	Name             string
	ReleaseDate      time.Time
	LabelID          int64
	CatalogueNumber  string
	TotalDiscs       int
}

type Track struct {
	JobID         uuid.UUID
	MBTrackID     string
	MBRecordingID string
	AlbumID       int64
	Name          string
	ArtistCredit  string
	TrackNumber   int
	DiscNumber    int
}

type Artist struct {
	JobID uuid.UUID
	MBID  string
	Name  string
	Slug  string
}

type Genre struct {
	JobID uuid.UUID
	Name  string
	Slug  string
}

type Label struct {
	JobID uuid.UUID
	MBID  string
	Name  string
	Slug  string
}

type TrackArtistMap struct {
	JobID     uuid.UUID
	TrackID   int64
	ArtistIDs []int64
}

type AlbumArtistMap struct {
	JobID     uuid.UUID
	AlbumID   int64
	ArtistIDs []int64
}

type AlbumGenreMap struct {
	JobID    uuid.UUID
	AlbumID  int64
	GenreIDs []int64
}
