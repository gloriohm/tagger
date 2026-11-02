package tag_loader

import "time"

type Track struct {
	Path                      string
	Title                     string
	Artists                   []string
	AlbumArtists              []string
	Album                     string
	Date                      string
	TrackNumber               string
	DiscNumber                string
	Label                     string
	BPM                       string
	CatalogNumber             string
	MusicBrainzTrackID        string
	MusicBrainzReleaseTrackID string
	MusicBrainzAlbumID        string
	MusicBrainzReleaseGroupID string
	MusicBrainzArtistID       string
	MusicBrainzAlbumArtistID  string
	Duration                  time.Duration
}

type Album []Track

type AlbumIdentifiers struct {
	ReleaseID      string
	ReleaseGroupID string
	AlbumName      string
	AlbumArtist    string
	TotalDiscs     int
	AlbumTags      Album
}
