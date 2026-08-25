package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang-postgres-docker/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (c *AuthController) Register(ctx *gin.Context) {
	var request RegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := c.authService.Register(
		ctx.Request.Context(),
		request.Name,
		request.Email,
		request.Password,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"data":    user,
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (c *AuthController) Login(ctx *gin.Context) {
	var request LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	result, err := c.authService.Login(
		ctx.Request.Context(),
		request.Email,
		request.Password,
	)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   result.Token,
	})
}

func (c *AuthController) ProtectedTest(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "user id not found in context",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "protected route accessed successfully",
		"user_id": userID,
	})
}
