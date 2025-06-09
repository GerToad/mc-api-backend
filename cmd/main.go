package main

import (
	// "context"
	// "fmt"
	"log"
	// "os"
	// "time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	// "gitlab.com/playIt02/api/config"
	// _ "gitlab.com/playIt02/api/docs"
	// "gitlab.com/playIt02/api/jobs"
	// "gitlab.com/playIt02/api/routes"
)

var dbpool *pgxpool.Pool

func main() {
  // Configuration
	config.SetVariables()
	config.InitDB()
	defer config.DB.Close()

  // Jobs in the background
  // Authentication rotate and Session management.
  go jobs.RefreshTokensJob() // 30 minutes window. Fetch new access tokens from Spotify and Google for users.
  go jobs.ManageSessions() // 10 seconds window. Only uses Spotify for fetching the curent playlist.
  go jobs.RenewSessionsJob()

  // Track playlist
  go jobs.TrackAndReorderPlaylistJob() // Keeps track of context state and updates the order in the databse.

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
    		c.Locals("origin", origin)
	    	// Define allowed origins
	    	allowedOrigins := map[string]bool{
          "http://localhost:3000":          true,
          config.BASE_URL:      true,
	    	}
        	if allowedOrigins[origin] {
            		c.Set("Access-Control-Allow-Origin", origin)
        	} else {
            		c.Set("Access-Control-Allow-Origin", "*")
        	}
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Set("Access-Control-Allow-Credentials", "true")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	// Sample route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	// Swagger route
	app.Get("/documentation/*", swagger.HandlerDefault)

  	// Routes
  	routes.SetupRoutes(app)

	// log.Println("📌 Registered Routes:")
	// for _, route := range app.Stack() {
	//   for _, r := range route {
	//     log.Printf("➡ %s %s", r.Method, r.Path)
	//   }
	// }

	log.Fatal(app.Listen(":8080"))
}

// @title PlayIt API
// @version 1.0
// @description API for PlayIt authentication and queue management
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
