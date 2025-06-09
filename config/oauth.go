package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/spotify"
)

// GoogleOAuthConfig is a global configuration for Google OAuth
var GoogleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	Scopes:       []string{
    "https://www.googleapis.com/auth/userinfo.email", 
    "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

// SpotifyOAuthConfig is a global configuration for Spotify OAuth
var SpotifyOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
	ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("SPOTIFY_REDIRECT_URL"),
	Scopes:       []string{
		"user-read-email",
		"user-read-private",
		"user-modify-playback-state",
		"user-read-playback-state",
		"user-read-currently-playing",
  },
	Endpoint:     spotify.Endpoint,
}
