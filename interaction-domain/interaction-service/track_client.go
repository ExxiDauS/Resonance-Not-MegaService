package interactionservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TrackClient interface {
	GetRandomTrack(ctx context.Context) (*TrackResponse, error)
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

func (c *HTTPTrackClient) GetRandomTrack(ctx context.Context) (*TrackResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+"/tracks/random",
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Track service response:", string(bodyBytes))

	var trackRes TrackResponse
	if err := json.Unmarshal(bodyBytes, &trackRes); err != nil {
		return nil, err
	}

	return &trackRes, nil
}
