package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	// "gitlab.com/playIt02/api/db/gen"
)

// JWTSecret should be loaded once to avoid multiple env lookups
var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWT creates a JWT token for a user
func GenerateJWT(user db.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID.String(), // Ensure it's a string
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}
