package middleware

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWT Secret Key (Ensure it's loaded from env)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// JWTProtected Middleware for Authentication
func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// log.Println("🔹 Checking JWT Token for request:", c.Method(), c.OriginalURL())

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Println("🚨 Missing Authorization Header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization token"})
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Println("🚨 Invalid or expired token:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// Extract user info
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["email"] == nil {
			log.Println("🚨 Invalid token claims")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Store user email
		// log.Println("✅ Authenticated:", claims["email"].(string))
		c.Locals("userEmail", claims["email"].(string))
		c.Locals("user_id", claims["user_id"].(string))
		return c.Next()
	}
}

// LoggerMiddleware logs requests
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		log.Printf(
			"📌 %s %s | Status: %d | Duration: %s",
			c.Method(), c.OriginalURL(), c.Response().StatusCode(), duration,
		)
		return err
	}
}

// RecoverMiddleware handles panics and prevents crashes
func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("🔥 Panic recovered: %v", err)
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
			}
		}()
		return c.Next()
	}
}
