package fetcher

import (
	"context"
	"fmt"
	"time"

	"go.uploadedlobster.com/mbtypes"
	"go.uploadedlobster.com/musicbrainzws2"
)

type BrainzClient struct {
	client *musicbrainzws2.Client
}

func NewBrainzClient(c *musicbrainzws2.Client) *BrainzClient {
	return &BrainzClient{client: c}
}

func (bc BrainzClient) GetReleaseByID(id string) (Album, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := musicbrainzws2.IncludesFilter{Includes: []string{"recordings", "labels", "release-groups", "artist-credits"}}
	release, err := bc.client.LookupRelease(ctx, mbtypes.MBID(id), filter)
	if err != nil {
		return Album{}, handleResponseError(err)
	}

	album := ReleaseToAlbum(release)
	return album, nil
}

func (bc BrainzClient) GetReleaseGroup(id string) ([]ReleaseInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := musicbrainzws2.IncludesFilter{Includes: []string{"releases", "media"}}

	releaseGroup, err := bc.client.LookupReleaseGroup(ctx, mbtypes.MBID(id), filter)
	if err != nil {
		return nil, handleResponseError(err)
	}

	var releases []ReleaseInfo
	for _, r := range releaseGroup.Releases {
		if r.Status == "Official" && r.Quality != "low" {
			release := ReleaseInfo{
				ID:      string(r.ID),
				Quality: r.Quality,
				Country: string(r.CountryCode),
				Date:    extractReleaseDate(r.Date),
				Barcode: string(r.Barcode),
				Media:   extractMediaInfo(r.Media),
			}

			releases = append(releases, release)
		}
	}

	return releases, nil
}

func (bc BrainzClient) SearchReleaseGroup(album, artist string) ([]ReleaseInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := musicbrainzws2.SearchFilter{
		Query:  fmt.Sprintf(`releasegroup:%q AND artist:%q`, album, artist),
		Dismax: false,
	}

	paginator := musicbrainzws2.Paginator{
		Limit:  10,
		Offset: 0,
	}

	hits, err := bc.client.SearchReleaseGroups(ctx, filter, paginator)
	if err != nil {
		return nil, handleResponseError(err)
	}

	var releaseGroupID string
	for _, rg := range hits.ReleaseGroups {
		if NormalizeString(rg.Title) == NormalizeString(album) {
			for _, ac := range rg.ArtistCredit {
				if NormalizeString(ac.Name) == NormalizeString(artist) {
					releaseGroupID = string(rg.ID)
					break
				}
			}
		}
	}

	if releaseGroupID == "" {
		return nil, ErrNotFound
	}

	return bc.GetReleaseGroup(releaseGroupID)
}
