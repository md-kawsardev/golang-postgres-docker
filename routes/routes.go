package routes

import (
	"github.com/gin-gonic/gin"

	"golang-postgres-docker/controllers"
	"golang-postgres-docker/middleware"
)

func SetupRouter(
	authController *controllers.AuthController,
	noteController *controllers.NoteController,
	jwtSecret string,
) *gin.Engine {

	router := gin.Default()

	api := router.Group("/api")

	// Public routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protected.POST("/notes", noteController.CreateNote)
		protected.GET("/notes", noteController.GetMyNotes)
		protected.PUT("/notes/:id", noteController.UpdateNote)
		protected.DELETE("/notes/:id", noteController.DeleteNote)

	}

	return router
}
