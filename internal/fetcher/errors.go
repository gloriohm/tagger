package fetcher

import (
	"errors"

	"go.uploadedlobster.com/musicbrainzws2"
)

var (
	ErrNoDate    = errors.New("date is empty")
	ErrNotFound  = errors.New("no results from metadata provider")
	ErrRateLimit = errors.New("rate limited")
)

func handleResponseError(err error) error {
	var clientErr *musicbrainzws2.ClientError

	if errors.As(err, &clientErr) {
		switch clientErr.StatusCode {
		case 404:
			return ErrNotFound
		case 429:
			return ErrRateLimit
		default:
			return err
		}
	}

	return err
}
