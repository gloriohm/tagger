package compare

import (
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/agnivade/levenshtein"
)

type Album struct {
	Title      string
	Artist     string
	TotalDiscs int
	Tracks     []Track
}

type Track struct {
	Title       string
	TrackNumber int
	DiscNumber  int
	Duration    time.Duration
}

type Difference struct {
	Field    string
	Distance float64
	Weight   float64
}

type Comparision struct {
	Distance    float64
	Differences []Difference
}

func CompareAlbum(local, mb Album) Comparision {
	var result Comparision
	var weightedSum, totalWeight float64

	add := func(field string, distance, weight float64) {
		weightedSum += distance * weight
		totalWeight += weight

		if distance > 0 {
			result.Differences = append(
				result.Differences,
				Difference{field, distance, weight},
			)
		}
	}

	add("album", stringDistance(local.Title, mb.Title), 3)
	add("artist", stringDistance(local.Artist, mb.Artist), 3)

	// penalize different track counts
	maxTracks := max(len(local.Tracks), len(mb.Tracks))
	if maxTracks > 0 {
		diff := math.Abs(float64(
			len(local.Tracks) - len(mb.Tracks),
		))
		add("track_count", diff/float64(maxTracks), 3)
	}

	// penalize different disc counts
	maxDiscs := max(local.TotalDiscs, mb.TotalDiscs)
	if maxDiscs > 0 {
		diff := math.Abs(float64(local.TotalDiscs - mb.TotalDiscs))
		add("disc_count", diff/float64(maxDiscs), 3)
	}

	// assumes tracks are in the same order
	for i := 0; i < min(len(local.Tracks), len(mb.Tracks)); i++ {
		a, b := local.Tracks[i], mb.Tracks[i]

		add("track_title", stringDistance(a.Title, b.Title), 3)
		add("track_length", durationDistance(a.Duration, b.Duration), 2)
	}

	if totalWeight > 0 {
		result.Distance = weightedSum / totalWeight
	}

	return result
}

func stringDistance(a, b string) float64 {
	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))

	longest := max(utf8.RuneCountInString(a), utf8.RuneCountInString(b))
	if longest == 0 {
		return 0
	}

	return float64(levenshtein.ComputeDistance(a, b)) / float64(longest)
}

func durationDistance(a, b time.Duration) float64 {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}

	tolerance := 2 * time.Second      // ignore diff less than 2s
	maxDifference := 30 * time.Second // 30s or more is max penalty

	if diff <= tolerance {
		return 0
	}

	return min(float64(diff-tolerance)/float64(maxDifference-tolerance), 1.0)
}
