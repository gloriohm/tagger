package tagger

import "time"

type Album struct {
	ID               int64
	MBReleaseID      string
	MBReleaseGroupID string
	Name             string
	AlbumArtists     []Artist
	ReleaseDate      time.Time
	Label            Label
	CatalogueNumber  string
	TotalDiscs       int
	Genres           []Genre
	Tracks           []Track
}

type Artist struct {
	ID   int64
	MBID string
	Name string
	Slug string
}

type Track struct {
	ID            int64
	MBTrackID     string
	MBRecordingID string
	AlbumID       int64
	Name          string
	Artists       []Artist
	ArtistCredit  string
	TrackNumber   int
	DiscNumber    int
}

type Genre struct {
	ID   int64
	Name string
	Slug string
}

type Label struct {
	ID   int64
	Name string
	Slug string
}
