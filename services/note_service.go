package services

import (
	"context"
	"fmt"
	"strings"

	"golang-postgres-docker/models"
	"golang-postgres-docker/repository"
)

type NoteService struct {
	noteRepo *repository.NoteRepository
}

func NewNoteService(noteRepo *repository.NoteRepository) *NoteService {
	return &NoteService{
		noteRepo: noteRepo,
	}
}

func (s *NoteService) CreateNote(
	ctx context.Context,
	userID int64,
	title string,
	content string,
) (*models.Note, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	if content == "" {
		return nil, fmt.Errorf("content is required")
	}

	note := &models.Note{
		Title:   title,
		Content: content,

		// IMPORTANT:
		// userID comes from the authenticated JWT,
		// not from the request body.
		UserID: userID,
	}

	if err := s.noteRepo.CreateNote(ctx, note); err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	return note, nil
}

func (s *NoteService) GetMyNotes(
	ctx context.Context,
	userID int64,
) ([]models.Note, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}

	notes, err := s.noteRepo.FindNotesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}

	return notes, nil
}

func (s *NoteService) UpdateNote(
	ctx context.Context,
	noteID int64,
	userID int64,
	title string,
	content string,
) (*models.Note, error) {
	if noteID <= 0 {
		return nil, fmt.Errorf("invalid note id")
	}

	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	if content == "" {
		return nil, fmt.Errorf("content is required")
	}

	note, err := s.noteRepo.UpdateNote(
		ctx,
		noteID,
		userID,
		title,
		content,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	return note, nil
}

func (s *NoteService) DeleteNote(
	ctx context.Context,
	noteID int64,
	userID int64,
) error {
	if noteID <= 0 {
		return fmt.Errorf("invalid note id")
	}

	if userID <= 0 {
		return fmt.Errorf("invalid user id")
	}

	if err := s.noteRepo.DeleteNote(ctx, noteID, userID); err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	return nil
}
