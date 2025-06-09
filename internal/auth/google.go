package auth

import (
	"context"
	"fmt"

	// "fmt"
	// "log"
	"os"

	"github.com/gofiber/fiber/v2"
	// "gitlab.com/playIt02/api/config"
	// "gitlab.com/playIt02/api/models"
	// "gitlab.com/playIt02/api/utils"
	"golang.org/x/oauth2"
)

// JWT secret
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GoogleUser represents Google user info response
type GoogleUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GoogleAuthHandler handles Google OAuth sign-in
// @Summary Sign in with Google
// @Description Redirects user to Google OAuth
// @Tags auth
// @Produce json
// @Success 302 {string} string "Redirecting to Google OAuth"
// @Router /api/auth/google [get]
func GoogleAuthHandler(c *fiber.Ctx) error {
	url := config.GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	return c.Redirect(url)
}

// GoogleCallbackHandler handles Google OAuth callback
// @Summary Google OAuth Callback
// @Description Handles Google OAuth login response
// @Tags auth
// @Accept  json
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Success 200 {object} map[string]interface{} "Returns JWT token and user data"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Failed to authenticate"
// @Router /api/auth/google/callback [get]
func GoogleCallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Code not found"})
	}

	// Exchange code for token
	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to exchange token"})
	}

	// Fetch user info from Google
	userInfo, err := FetchGoogleUser(token.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user info"})
	}

	// Check if user exists in DB
	user, err := models.FindUserByEmail(userInfo.Email)
	if err != nil {
		// Create user if not found
		user, err = models.CreateUser(userInfo.Name, userInfo.Email, "", "user")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}
	}

	// Store tokens
	err = models.StoreOAuthTokens(user.ID, token.AccessToken, token.RefreshToken, token.Expiry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store tokens"})
	}

	// Generate JWT
	jwtToken, err := utils.GenerateJWT(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate JWT"})
	}

	// Redirect to frontend with JWT as query parameter
	frontendURL := fmt.Sprintf("%s/auth/callback?token=%s", config.BASE_URL, jwtToken)
	return c.Redirect(frontendURL, fiber.StatusTemporaryRedirect)
}
