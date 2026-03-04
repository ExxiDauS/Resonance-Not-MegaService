package interactionservice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type TrackClient interface {
	GetRandomTrack(ctx context.Context) (*Track, error)
}

type HTTPTrackClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPTrackClient(baseURL string) *HTTPTrackClient {
	return &HTTPTrackClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *HTTPTrackClient) GetRandomTrack(ctx context.Context) (*Track, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+"/random-track",
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var track Track
	if err := json.NewDecoder(resp.Body).Decode(&track); err != nil {
		return nil, err
	}

	return &track, nil
}
