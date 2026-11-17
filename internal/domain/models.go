package domain

import "time"

type Album struct {
	ID                int64
	MBReleaseID       string
	MBReleaseGroupID  string
	Name              string
	AlbumArtistCredit string // Display name as string with join phrases from mb
	AlbumArtists      []Artist
	ReleaseDate       time.Time
	Labels            []Label
	TotalDiscs        int
	Genres            []Genre
	Tracks            []Track
}

type Artist struct {
	ID   int64
	MBID string
	Name string
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
}

type Label struct {
	ID   int64
	MBID string
	Name string
}
