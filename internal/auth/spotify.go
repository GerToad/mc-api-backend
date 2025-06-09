package auth

import (
	"context"
	"fmt"
  "bytes"

	"github.com/gofiber/fiber/v2"
	// "gitlab.com/playIt02/api/config"
	// "gitlab.com/playIt02/api/models"
	// "gitlab.com/playIt02/api/utils"
	"golang.org/x/oauth2"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"display_name"`
	Email string `json:"email"`
  Images []struct {
    URL string `json:"url"`
  } `json:"images"`
}

// SpotifyAuthHandler handles Spotify OAuth sign-in
// @Summary Sign in with Spotify
// @Description Redirects user to Spotify OAuth
// @Tags auth
// @Produce json
// @Success 302 {string} string "Redirecting to Spotify OAuth"
// @Router /api/auth/spotify [get]
func SpotifyAuthHandler(c *fiber.Ctx) error {
	url := config.SpotifyOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	fmt.Println("Redirecting to:", url)
	return c.Redirect(url)
}

// SpotifyCallbackHandler handles Spotify OAuth callback
// @Summary Spotify OAuth Callback
// @Description Handles Spotify OAuth login response
// @Tags auth
// @Accept  json
// @Produce json
// @Param code query string true "Authorization code from Spotify"
// @Success 200 {object} map[string]interface{} "Returns JWT token and user data"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Failed to authenticate"
// @Router /api/auth/spotify/callback [get]
func SpotifyCallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing authorization code"})
	}

	// Exchange code for token
	token, err := config.SpotifyOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to exchange token"})
	}

	// Fetch user info from Spotify
	userInfo, err := FetchSpotifyUser(token.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user info"})
	}

	// Check if user exists in DB, if not register as a manager
	user, err := models.FindUserByEmail(userInfo.Email)
	if err != nil {
		// Create user as a "manager"
		user, err = models.CreateUser(
			userInfo.Name,
			userInfo.Email,
			userInfo.Images[0].URL,
			"manager",
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}
	}

	// Store tokens
	err = models.StoreOAuthTokens(user.ID, token.AccessToken, token.RefreshToken, token.Expiry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store tokens"})
	}

	// ✅ Generate QR Code
	qrCode, err := utils.GenerateQRCode(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate QR code"})
	}

	// Return the PNG as image data
	c.Set("Content-Type", "image/png")
	return c.Status(fiber.StatusOK).SendStream(bytes.NewReader(qrCode))
}
