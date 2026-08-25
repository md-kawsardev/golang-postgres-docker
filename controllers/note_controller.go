package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"golang-postgres-docker/services"
)

type NoteController struct {
	noteService *services.NoteService
}

func NewNoteController(noteService *services.NoteService) *NoteController {
	return &NoteController{
		noteService: noteService,
	}
}

type createNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (c *NoteController) CreateNote(ctx *gin.Context) {
	var req createNoteRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user authentication required",
		})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	note, err := c.noteService.CreateNote(
		ctx.Request.Context(),
		userID,
		req.Title,
		req.Content,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "note created successfully",
		"data":    note,
	})
}

func (c *NoteController) GetMyNotes(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user authentication required",
		})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	notes, err := c.noteService.GetMyNotes(
		ctx.Request.Context(),
		userID,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "notes retrieved successfully",
		"data":    notes,
	})
}

func (c *NoteController) UpdateNote(ctx *gin.Context) {
	noteID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || noteID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid note id",
		})
		return
	}

	var req createNoteRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user authentication required",
		})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	note, err := c.noteService.UpdateNote(
		ctx.Request.Context(),
		noteID,
		userID,
		req.Title,
		req.Content,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "note updated successfully",
		"data":    note,
	})
}

func (c *NoteController) DeleteNote(ctx *gin.Context) {
	noteID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || noteID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid note id",
		})
		return
	}

	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user authentication required",
		})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	if err := c.noteService.DeleteNote(
		ctx.Request.Context(),
		noteID,
		userID,
	); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "note deleted successfully",
	})
}
