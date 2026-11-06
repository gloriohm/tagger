package domain

import (
	"tagger/internal/mbfetcher"
)

func ParseAlbum(in mbfetcher.Album) Album {
	var a Album
	a.MBReleaseID = in.MBReleaseID
	a.MBReleaseGroupID = in.MBReleaseGroupID
	a.Name = in.Name

	// check this against zero value and request user input if true
	a.ReleaseDate = in.ReleaseDate

}
