package main

import (
	"context"
	"log"

	"golang-postgres-docker/config"
	"golang-postgres-docker/controllers"
	"golang-postgres-docker/database"
	"golang-postgres-docker/repository"
	"golang-postgres-docker/routes"
	"golang-postgres-docker/services"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal("Config error")
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("database URL not configured")
	}

	ctx := context.Background()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL!")

	// Repositories
	userRepo := repository.NewUserRepository(db)
	noteRepo := repository.NewNoteRepository(db)

	// Services
	authService := services.NewAuthService(
		userRepo,
		cfg.JWTSecret,
	)

	noteService := services.NewNoteService(
		noteRepo,
	)

	// Controllers
	authController := controllers.NewAuthController(authService)

	noteController := controllers.NewNoteController(
		noteService,
	)

	// Router
	router := routes.SetupRouter(
		authController,
		noteController,
		cfg.JWTSecret,
	)

	log.Println("server running on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
