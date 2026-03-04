package spotify_client

import (
	"context"
	"log"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2/clientcredentials"

	"music-domain/configs"

	"github.com/zmb3/spotify/v2"
)

func NewSpotifyClient() (*spotify.Client, error) {
	ctx := context.Background()
	credentials, err := configs.LoadSpotifyConfig()
	if err != nil {
		log.Fatalf("couldn't load spotify config: %v", err)
		return nil, err
	}
	config := &clientcredentials.Config{
		ClientID:     credentials.SpotifyClientID,
		ClientSecret: credentials.SpotifyClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
		return nil, err
	}
	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	return client, nil
}
