package auth

import (
	"encoding/json"
	"net/http"
  "context"
  "errors"

	"golang.org/x/oauth2"
	// "gitlab.com/playIt02/api/config"
	// "gitlab.com/playIt02/api/models"
	"github.com/google/uuid"
)

// OAuthUser represents user data fetched from Google/Spotify
type OAuthUser struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

// Fetch user info from Google
func FetchGoogleUser(accessToken string) (*GoogleUser, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Fetch user info from Spotify API
func FetchSpotifyUser(accessToken string) (*User, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func RefreshToken(userID uuid.UUID) (*oauth2.Token, error) {
	ctx := context.Background()

	// Get user from database
	user, err := models.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if user.RefreshToken == "" {
		return nil, errors.New("user has no refresh token")
	}

	// Create token from stored refresh token
	token := &oauth2.Token{
		RefreshToken: user.RefreshToken,
	}

	// Automatically refresh using Google's client
  var oauthConfig *oauth2.Config
  var url string
  if user.Role == "user"{
    oauthConfig = config.GoogleOAuthConfig
    url = "https://www.googleapis.com/oauth2/v2/userinfo"
  }else {
    oauthConfig = config.SpotifyOAuthConfig
    url = "https://api.spotify.com/v1/me"
  }

	client := oauthConfig.Client(ctx, token)

	// Call a simple API request to trigger token refresh
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if token was refreshed
	newToken, err := oauthConfig.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, err
	}

	// Store updated token in DB
	err = models.StoreOAuthTokens(userID, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}
