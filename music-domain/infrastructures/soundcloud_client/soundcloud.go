package soundcloud_client

import (
	"context"
	"log"
	"music-domain/configs"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// SoundCloudTransport fetches the OAuth token and sets it as "OAuth <token>"
// instead of the default "Bearer <token>" that SoundCloud requires.
type SoundCloudTransport struct {
	TokenSource oauth2.TokenSource
	Base        http.RoundTripper
}

func (t *SoundCloudTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.TokenSource.Token()
	if err != nil {
		log.Printf("[SoundCloudTransport] Failed to get token: %v", err)
		return nil, err
	}

	// Set the header directly as "OAuth" instead of relying on Bearer rewrite
	req.Header.Set("Authorization", "OAuth "+token.AccessToken)

	return t.Base.RoundTrip(req)
}

func NewSoundCloudClient() *http.Client {
	ctx := context.Background()
	config, err := configs.LoadSoundCloudConfig()
	if err != nil {
		log.Fatalf("couldn't load soundcloud config: %v", err)
		return nil
	}
	soundcloudConfig := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     "https://secure.soundcloud.com/oauth/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	}

	// Get a TokenSource that handles caching and auto-refresh
	tokenSource := soundcloudConfig.TokenSource(ctx)

	client := &http.Client{
		Transport: &SoundCloudTransport{
			TokenSource: tokenSource,
			Base:        http.DefaultTransport,
		},
	}

	return client
}
