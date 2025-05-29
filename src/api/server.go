package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tmm6907/dashboard/routes"
	"github.com/tmm6907/dashboard/utils"
	"github.com/tmm6907/dashboard/worker"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func initDB() (*sqlx.DB, error) {
	dbName := "db/mashboard.sqlite"
	buildFile := "build.sql"
	pragmaCommands := []string{
		"PRAGMA load_extension = 1;",
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_keys=1;",
	}
	pragmaCommand := strings.Join(pragmaCommands, " ")

	db, err := sqlx.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("Unable to open mashboard connection: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Unable to ping mashboard connection: %w", err)
	}

	if _, err := db.Exec(pragmaCommand); err != nil {
		return nil, fmt.Errorf("Unable to run PRAGMA commands: %w", err)
	}
	sqlData, err := os.ReadFile(buildFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to open build file: %w", err)
	}
	_, err = db.Exec(string(sqlData))
	if err != nil {
		return nil, fmt.Errorf("Unable to build sql: %w", err)
	}

	return db, nil
}

func main() {
	allowedOrigins := []string{
		"https://mashboard.app",
		"https://50.116.53.73:4173",
		"https://50.116.53.73:3030",
	}
	server := fiber.New()
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	allowedOriginsStr := strings.Join(allowedOrigins, ", ")

	server.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOriginsStr,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	server.Use(func(c *fiber.Ctx) error {
		if c.Method() == "OPTIONS" {
			c.Set("Access-Control-Allow-Origin", allowedOriginsStr)
			c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			return c.SendStatus(204)
		}
		return c.Next()
	})

	client_id, err := utils.ReadDockerSecret("google_client_id")
	if err != nil {
		log.Fatalf("Error reading google_client_id secret: %v", err)
	}
	// Trim whitespace if necessary, as secrets can have newlines
	client_id = strings.TrimSpace(client_id)
	if client_id == "" {
		log.Fatal("Google client ID secret is empty")
	}

	// Read GOOGLE_CLIENT_SECRET from Docker secret
	client_secret, err := utils.ReadDockerSecret("google_client_secret")
	if err != nil {
		log.Fatalf("Error reading google_client_secret secret: %v", err)
	}
	// Trim whitespace if necessary
	client_secret = strings.TrimSpace(client_secret)
	if client_secret == "" {
		log.Fatal("Google client secret secret is empty")
	}

	var GoogleOAuthConfig = &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		RedirectURL:  "https://api.mashboard.app/auth/callback",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	routesHandler := routes.NewHandler(db, GoogleOAuthConfig)
	workerHandler := worker.NewHandler(db)

	apiRoutes := server.Group("/api")
	feedRoutes := apiRoutes.Group("/feeds")
	feedItemRoutes := feedRoutes.Group("/items")
	feedRoutes.Get("/", routesHandler.CheckAuthHandler(), routesHandler.GetFeeds)
	feedRoutes.Post("/", routesHandler.CheckAuthHandler(), routesHandler.GetFeeds)
	// feedRoutes.Post("/data", routesHandler.CheckAuthHandler(), routesHandler.CreateFeed)
	feedRoutes.Post("/search", routesHandler.CheckAuthHandler(), routesHandler.GetFeeds)
	feedRoutes.Post("/search/new", routesHandler.CheckAuthHandler(), routesHandler.SearchForNewFeedByURL)
	feedRoutes.Post("/follow", routesHandler.CheckAuthHandler(), routesHandler.FollowFeed)
	feedRoutes.Get("/followed", routesHandler.CheckAuthHandler(), routesHandler.GetFollowedFeeds)
	feedItemRoutes.Get("/", routesHandler.CheckAuthHandler(), routesHandler.GetFeedItems)
	feedRoutes.Get("/followed", routesHandler.CheckAuthHandler(), routesHandler.GetFollowedFeedItems)
	feedItemRoutes.Get("/saved", routesHandler.CheckAuthHandler(), routesHandler.GetSavedFeedItems)
	feedItemRoutes.Get("/:id", routesHandler.CheckAuthHandler(), routesHandler.GetFeedItem)
	feedItemRoutes.Post("/:id/bookmark", routesHandler.CheckAuthHandler(), routesHandler.SaveFeedItem)

	userRoutes := apiRoutes.Group("/user")
	userRoutes.Get("/", routesHandler.GetUser)

	authRoutes := server.Group("/auth")
	authRoutes.Get("/login", routesHandler.GoogleOauthLoginHandler)
	authRoutes.Get("/logout", routesHandler.Logout)
	authRoutes.Get("/callback", routesHandler.CallbackHandler())

	server.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(
			map[string]any{
				"endpoints": []string{
					"/auth/login",
					"/api/feeds",
					"/api/feeds/items",
				},
			},
		)
	})
	port := os.Getenv("PORT")
	if port == "" {
		log.Error("port not configured")
	}
	host := fmt.Sprintf(":%s", port)
	go workerHandler.StartRSSFetcher(nil)
	server.Listen(host)
}
