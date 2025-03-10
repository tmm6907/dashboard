package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tmm6907/dashboard/routes"
	"github.com/tmm6907/dashboard/worker"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func initDB() (*sqlx.DB, error) {
	dbName := "mashboard.db"
	buildFile := "build.sql"
	buildCommands := []string{
		"PRAGMA load_extension = 1;",
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_keys=1;",
	}
	buildCommand := strings.Join(buildCommands, " ")
	db, err := sqlx.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("Unable to open mashboard connection: %w", err)
	}

	if _, err := db.Exec(buildCommand); err != nil {
		return nil, fmt.Errorf("Unable to run build commands: %w", err)
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
	server := fiber.New()
	db, err := initDB()
	if err != nil {
		log.Error(err)
		panic(err)
	}
	defer db.Close()

	err = godotenv.Load()
	if err != nil {
		panic(err)
	}

	var GoogleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	routesHandler := routes.NewHandler(db, GoogleOAuthConfig)
	workerHandler := worker.NewHandler(db)

	apiRoutes := server.Group("/api")
	apiRoutes.Get("/feeds", routesHandler.GetAllFeeds)
	apiRoutes.Post("/feeds", routesHandler.CreateFeed)
	apiRoutes.Get("/feeds/items", routesHandler.GetFeedItems)

	authRoutes := server.Group("/auth")
	authRoutes.Get("/login", routesHandler.LoginHandler)
	authRoutes.Get("/callback", routesHandler.CallbackHandler())

	server.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(
			map[string]any{
				"endpoints": []string{"/feeds"},
			},
		)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	host := fmt.Sprintf(":%s", port)
	log.Debug(host)
	// interval := 15 * time.Minute
	go workerHandler.StartRSSFetcher(nil)
	server.Listen(host)
}
