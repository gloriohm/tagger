package tag_loader

import (
	"fmt"

	"go.senan.xyz/taglib"
)

func parseTrack(path string, tags map[string][]string, props taglib.Properties) Track {
	track := Track{
		Path:                      path,
		Title:                     firstTag(tags, taglib.Title),
		Artists:                   tags[taglib.Artist],
		AlbumArtists:              tags[taglib.AlbumArtist],
		Album:                     firstTag(tags, taglib.Album),
		Date:                      firstTag(tags, taglib.Date),
		TrackNumber:               firstTag(tags, taglib.TrackNumber),
		DiscNumber:                firstTag(tags, taglib.DiscNumber),
		Label:                     firstTag(tags, taglib.Label),
		BPM:                       firstTag(tags, taglib.BPM),
		CatalogNumber:             firstTag(tags, taglib.CatalogNumber),
		MusicBrainzTrackID:        firstTag(tags, taglib.MusicBrainzTrackID),
		MusicBrainzReleaseTrackID: firstTag(tags, taglib.MusicBrainzReleaseTrackID),
		MusicBrainzAlbumID:        firstTag(tags, taglib.MusicBrainzAlbumID),
		MusicBrainzReleaseGroupID: firstTag(tags, taglib.MusicBrainzReleaseGroupID),
		MusicBrainzArtistID:       firstTag(tags, taglib.MusicBrainzArtistID),
		MusicBrainzAlbumArtistID:  firstTag(tags, taglib.MusicBrainzAlbumArtistID),
		Duration:                  props.Length,
	}
	return track
}

func firstTag(tags map[string][]string, key string) string {
	values := tags[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func parseAlbumIdentifiers(t Track) (AlbumIdentifiers, error) {
	if len(t.AlbumArtists) == 0 {
		return AlbumIdentifiers{}, fmt.Errorf("%w: no album artist", ErrNoValue)
	}

	return AlbumIdentifiers{
		ReleaseID:      t.MusicBrainzAlbumID,
		ReleaseGroupID: t.MusicBrainzReleaseGroupID,
		AlbumName:      t.Album,
		AlbumArtist:    t.AlbumArtists[0],
	}, nil
}
