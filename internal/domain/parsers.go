package domain

import (
	"strconv"
	"tagger/internal/mbfetcher"

	"go.senan.xyz/taglib"
)

func ParseAlbum(in mbfetcher.Album) Album {
	var a Album
	a.MBReleaseID = in.MBReleaseID
	a.MBReleaseGroupID = in.MBReleaseGroupID
	a.Name = in.Name
	a.ReleaseDate = in.ReleaseDate
	a.TotalDiscs = in.TotalDiscs

	return a
}

func ParseLabels(in []mbfetcher.LabelInfo) []Label {
	labels := make([]Label, len(in))

	for _, l := range in {
		label := Label{
			MBID: l.MBID,
			Name: l.Name,
		}

		labels = append(labels, label)
	}

	return labels
}

func ParseArtists(in []mbfetcher.Artist) []Artist {
	artists := make([]Artist, len(in))

	for _, a := range in {
		artist := Artist{
			MBID: a.MBID,
			Name: a.Name,
		}

		artists = append(artists, artist)
	}

	return artists
}

func ParseTracks(in []mbfetcher.Track, album_id int64) []Track {
	tracks := make([]Track, len(in))

	for _, t := range in {
		track := Track{
			MBTrackID:     t.MBTrackID,
			MBRecordingID: t.MBRecordingID,
			AlbumID:       album_id,
			Name:          t.Name,
			Artists:       ParseArtists(t.Artists),
			ArtistCredit:  t.ArtistCredit,
			TrackNumber:   t.TrackNumber,
			DiscNumber:    t.DiscNumber,
		}

		tracks = append(tracks, track)
	}

	return tracks
}

func AlbumToTagger(album Album) []map[string][]string {
	out := make([]map[string][]string, len(album.Tracks))

	for _, track := range album.Tracks {
		tags := map[string][]string{
			// track
			taglib.Title:                 {track.Name},
			taglib.TrackNumber:           {strconv.Itoa(track.TrackNumber)},
			taglib.DiscNumber:            {strconv.Itoa(track.DiscNumber)},
			taglib.Artist:                {track.ArtistCredit},
			taglib.Artists:               extractArtistsName(track.Artists),
			"MUSICBRAINZ_ARTISTID":       extractArtistsMBID(track.Artists),
			"MUSICBRAINZ_TRACKID":        {track.MBRecordingID},
			"MUSICBRAINZ_RELEASETRACKID": {track.MBTrackID},
			// album
			"MUSICBRAINZ_ALBUMID":        {album.MBReleaseID},
			"MUSICBRAINZ_RELEASEGROUPID": {album.MBReleaseGroupID},
			taglib.AlbumArtist:           {album.AlbumArtistCredit},
			"ALBUMARTISTS":               extractArtistsName(album.AlbumArtists),
			"MUSICBRAINZ_ALBUMARTISTID":  extractArtistsMBID(album.AlbumArtists),
			taglib.Genre:                 extractGenreName(album.Genres),
			taglib.Label:                 extractLabelName(album.Labels),
			taglib.ReleaseDate:           {album.ReleaseDate.Format("2006-01-02")},
			taglib.OriginalDate:          {album.ReleaseDate.Format("2006-01-02")},
			"DISCTOTAL":                  {strconv.Itoa(album.TotalDiscs)},
		}

		out = append(out, tags)
	}

	return out
}

func extractArtistsName(in []Artist) []string {
	out := make([]string, len(in))
	for _, a := range in {
		out = append(out, a.Name)
	}

	return out
}

func extractArtistsMBID(in []Artist) []string {
	out := make([]string, len(in))
	for _, a := range in {
		out = append(out, a.MBID)
	}

	return out
}

func extractGenreName(in []Genre) []string {
	out := make([]string, len(in))
	for _, g := range in {
		out = append(out, g.Name)
	}

	return out
}

func extractLabelName(in []Label) []string {
	out := make([]string, len(in))
	for _, l := range in {
		out = append(out, l.Name)
	}

	return out
}
