package audiosservice

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

type soundCloudAudioRepo struct {
	client *http.Client
}

// Added ID and Title so we can verify if the track unmarshaled at all, even if StreamURL is missing
type soundCloudTrack struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	StreamURL string `json:"stream_url"`
	Duration  int    `json:"duration"`
}

type soundCloudSearchResult struct {
	Collection []soundCloudTrack `json:"collection"`
	NextHref   string            `json:"next_href"`
}

func NewSoundCloudAudioRepo(client *http.Client) *soundCloudAudioRepo {
	return &soundCloudAudioRepo{client: client}
}

// 1. Search for the track and get the protected stream URL
func (r *soundCloudAudioRepo) SearchTrackStreamURL(ctx context.Context, query string) (string, int, error) {
	endpoint := "https://api.soundcloud.com/tracks?q=" + url.QueryEscape(query) + "&limit=1&offset=0&linked_partitioning=true"

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	resp, err := r.client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	// Read the raw body to log it exactly as it came from SoundCloud
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	var result soundCloudSearchResult

	// Attempt 1: Unmarshal as an Object (linked_partitioning=true format)
	err = json.Unmarshal(body, &result)

	// Attempt 2: If collection is empty, SoundCloud might have returned a raw array
	if err != nil || len(result.Collection) == 0 {
		var rawArray []soundCloudTrack
		if errArray := json.Unmarshal(body, &rawArray); errArray == nil && len(rawArray) > 0 {
			log.Printf("Fallback: Successfully parsed response as a raw JSON array instead of a collection object.")
			result.Collection = rawArray
		}
	}

	if len(result.Collection) == 0 {
		return "", 0, errors.New("no tracks found on soundcloud for this query (or unmarshal completely failed)")
	}

	// 3. Check if the track exists but lacks a stream URL
	if result.Collection[0].StreamURL == "" {
		// Log the raw JSON of just this track to see what fields actually exist
		var debugTrack interface{}
		json.Unmarshal(body, &debugTrack)
		log.Printf("Track is missing 'stream_url'. Full track JSON: %+v", debugTrack)
		return "", 0, errors.New("track found, but 'stream_url' is missing. The token might lack permissions or the track is not streamable")
	}

	return result.Collection[0].StreamURL, result.Collection[0].Duration, nil
}

// 2. Intercept the Redirect to get the raw MP3 URL
func (r *soundCloudAudioRepo) ResolvePlayableCDNURL(ctx context.Context, streamURL string) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, streamURL, nil)

	redirectErr := errors.New("stop redirect")
	tempClient := &http.Client{
		Transport: r.client.Transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return redirectErr
		},
	}

	resp, err := tempClient.Do(req)

	if err != nil && !errors.Is(err, redirectErr) && resp == nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound {
		cdnURL := resp.Header.Get("Location")
		return cdnURL, nil
	}

	return "", errors.New("failed to resolve CDN url from soundcloud")
}
