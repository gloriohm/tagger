package fetcher

import "time"

type Album struct {
	MBReleaseID      string
	MBReleaseGroupID string
	Name             string
	ReleaseDate      time.Time
	Label            []LabelInfo
	TotalDiscs       int
	AlbumArtists     []Artist
	Tracks           []Track
}

type Track struct {
	MBTrackID     string
	MBRecordingID string
	ArtistCredit  string
	Artists       []Artist
	Name          string
	DiscNumber    int
	TrackNumber   int
}

type Artist struct {
	MBID string
	Name string
}

type LabelInfo struct {
	Name          string
	MBID          string
	CatalogNumber string
}

type ReleaseInfo struct {
	ID      string
	Quality string
	Country string
	Barcode string
	Date    time.Time
	Media   MediaInfo
}

type Medium struct {
	Format string
	Tracks int
}

type MediaInfo struct {
	Discs       []Medium
	TotalTracks int
}
