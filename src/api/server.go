package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tmm6907/dashboard/routes"
	"github.com/tmm6907/dashboard/worker"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false // File does not exist
		}
	}
	if info.IsDir() {
		log.Debug("is a dir2ectory")
		return false
	}
	return !info.IsDir() // Returns false if it's a directory
}

func initDB() (*sqlx.DB, error) {
	dbName := "mashboard.sqlite"
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
	err := godotenv.Load()
	if err != nil {
		log.Error(err)
	}

	server := fiber.New()
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	allowedOrigins := "http://localhost:3030, http://localhost:8080, http://localhost/, http://localhost:4173, http://50.116.53.73:4173, http://50.116.53.73:3030"

	server.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	server.Use(func(c *fiber.Ctx) error {
		if c.Method() == "OPTIONS" {
			c.Set("Access-Control-Allow-Origin", allowedOrigins)
			c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			return c.SendStatus(204)
		}
		return c.Next()
	})

	client_id := os.Getenv("GOOGLE_CLIENT_ID")
	if client_id == "" {
		log.Fatal("Missing google client id")
	}

	client_secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if client_secret == "" {
		log.Fatal("Missing google client secret")
	}

	var GoogleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "https://masboard.app:8080/auth/callback",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	routesHandler := routes.NewHandler(db, GoogleOAuthConfig)
	workerHandler := worker.NewHandler(db)

	apiRoutes := server.Group("/api")
	feedRoutes := apiRoutes.Group("/feeds")
	feedItemRoutes := feedRoutes.Group("/items")
	feedRoutes.Get("/", routesHandler.CheckAuthHandler(), routesHandler.GetFeeds)
	feedRoutes.Post("/", routesHandler.CreateFeed)
	feedRoutes.Post("/find", routesHandler.CheckAuthHandler(), routesHandler.SearchForFeed)
	feedRoutes.Post("/follow", routesHandler.CheckAuthHandler(), routesHandler.FollowFeed)
	feedItemRoutes.Get("/", routesHandler.CheckAuthHandler(), routesHandler.GetFeedItems)
	feedItemRoutes.Get("/saved", routesHandler.CheckAuthHandler(), routesHandler.GetSavedFeedItems)
	feedItemRoutes.Get("/:id", routesHandler.CheckAuthHandler(), routesHandler.GetFeedItem)
	feedItemRoutes.Post("/:id/bookmark", routesHandler.CheckAuthHandler(), routesHandler.SaveFeedItem)

	authRoutes := server.Group("/auth")
	authRoutes.Get("/login", routesHandler.LoginHandler)
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
		port = "8080"
	}
	host := fmt.Sprintf(":%s", port)
	log.Debug(host)
	go workerHandler.StartRSSFetcher(nil)
	server.Listen(host)
}
