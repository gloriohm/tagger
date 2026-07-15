package filer

type Track struct {
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
}

type Album []Track

type Identifiers struct {
	ReleaseID      string
	ReleaseGroupID string
	AlbumName      string
	AlbumArtists   []string
}
