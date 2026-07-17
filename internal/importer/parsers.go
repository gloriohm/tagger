package importer

import "tagger/internal/fetcher"

func parseTrack(in fetcher.Track, albumID int64) Track {
	return Track{
		MBTrackID:     in.MBTrackID,
		MBRecordingID: in.MBRecordingID,
		AlbumID:       albumID,
		Name:          in.Name,
		ArtistCredit:  in.ArtistCredit,
		TrackNumber:   in.TrackNumber,
		DiscNumber:    in.DiscNumber,
	}
}

func parseArtists(in []fetcher.Artist) []Artist {
	var artists []Artist
	for _, a := range in {
		artists = append(artists, parseArtist(a))
	}

	return artists
}

func parseArtist(in fetcher.Artist) Artist {
	return Artist{
		MBID: in.MBID,
		Name: in.Name,
	}
}

func parseLabel(in fetcher.LabelInfo) Label {
	return Label{
		Name: in.Name,
		MBID: in.MBID,
	}
}

func parseAlbum(in fetcher.Album) Album {
	return Album{
		MBReleaseID:      in.MBReleaseID,
		MBReleaseGroupID: in.MBReleaseGroupID,
		Name:             in.Name,
		ReleaseDate:      in.ReleaseDate,
		CatalogueNumber:  in.Label[0].CatalogNumber,
		TotalDiscs:       in.TotalDiscs,
	}
}
