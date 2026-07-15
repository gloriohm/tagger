package fetcher

import (
	"strings"
	"time"

	"go.uploadedlobster.com/mbtypes"
	"go.uploadedlobster.com/musicbrainzws2"
)

func ReleaseToAlbum(r musicbrainzws2.Release) Album {
	var tracks []Track
	for _, m := range r.Media {
		for _, t := range m.Tracks {
			track := Track{
				MBTrackID:     string(t.ID),
				MBRecordingID: string(t.Recording.ID),
				Name:          t.Title,
				TrackNumber:   t.Position,
				DiscNumber:    m.Position,
				ArtistCredit:  extractArtistCredit(t.Recording.ArtistCredit),
				Artists:       extractArtists(t.Recording.ArtistCredit),
			}

			tracks = append(tracks, track)
		}
	}

	var albumArtists []Artist
	for _, ac := range r.ArtistCredit {
		artist := Artist{
			Name: ac.Artist.Name,
			MBID: string(ac.Artist.ID),
		}
		albumArtists = append(albumArtists, artist)
	}

	var labels []LabelInfo
	for _, l := range r.LabelInfo {
		label := LabelInfo{
			Name:          l.Label.Name,
			MBID:          string(l.Label.ID),
			CatalogNumber: l.CatalogNumber,
		}
		labels = append(labels, label)
	}

	album := Album{
		MBReleaseID:      string(r.ID),
		MBReleaseGroupID: string(r.ReleaseGroup.ID),
		Name:             r.Title,
		ReleaseDate:      extractReleaseDate(r.ReleaseGroup.FirstReleaseDate),
		Label:            labels,
		TotalDiscs:       extractDiscCount(r.Media),
		AlbumArtists:     albumArtists,
		Tracks:           tracks,
	}

	return album
}

func extractArtistCredit(credits musicbrainzws2.ArtistCredit) string {
	// ignores joinphrase and returns early if there is only one credited artist
	if len(credits) == 1 {
		return credits[0].Name
	}

	var out strings.Builder
	for _, ac := range credits {
		out.WriteString(ac.Name)
		out.WriteString(ac.JoinPhrase)
	}

	return out.String()
}

func extractArtists(credits musicbrainzws2.ArtistCredit) []Artist {
	var artists []Artist
	for _, ac := range credits {
		artist := Artist{
			MBID: string(ac.Artist.ID),
			Name: ac.Artist.Name,
		}

		artists = append(artists, artist)
	}

	return artists
}

func extractReleaseDate(date mbtypes.PartialDate) time.Time {
	if date.IsEmpty() {
		return time.Time{}
	}

	if date.Month == 0 {
		date.Month = 1
	}
	if date.Day == 0 {
		date.Day = 1
	}

	releaseDate := time.Date(date.Year, time.Month(date.Month), date.Day, 0, 0, 0, 0, time.UTC)
	return releaseDate
}

func extractDiscCount(media musicbrainzws2.MediaList) int {
	return len(media)
}

func extractMediaInfo(medialist musicbrainzws2.MediaList) MediaInfo {
	var media []Medium
	var totalTracks int
	for _, m := range medialist {
		media = append(media, Medium{Format: m.Format, Tracks: m.TrackCount})
		totalTracks += m.TrackCount
	}

	return MediaInfo{Discs: media, TotalTracks: totalTracks}
}
