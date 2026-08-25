package repository

import (
	"context"
	"fmt"

	"golang-postgres-docker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteRepository struct {
	db *pgxpool.Pool
}

func NewNoteRepository(db *pgxpool.Pool) *NoteRepository {
	return &NoteRepository{
		db: db,
	}
}

func (r *NoteRepository) CreateNote(
	ctx context.Context,
	note *models.Note,
) error {
	query := `
		INSERT INTO notes (title, content, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		note.Title,
		note.Content,
		note.UserID,
	).Scan(
		&note.ID,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	return nil
}

func (r *NoteRepository) FindNotesByUserID(
	ctx context.Context,
	userID int64,
) ([]models.Note, error) {
	query := `
		SELECT id, title, content, user_id, created_at, updated_at
		FROM notes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notes: %w", err)
	}
	defer rows.Close()

	notes := make([]models.Note, 0)

	for rows.Next() {
		var note models.Note

		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.UserID,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading notes: %w", err)
	}

	return notes, nil
}

func (r *NoteRepository) UpdateNote(
	ctx context.Context,
	noteID int64,
	userID int64,
	title string,
	content string,
) (*models.Note, error) {
	query := `
		UPDATE notes
		SET title = $1,
		    content = $2,
		    updated_at = NOW()
		WHERE id = $3
		  AND user_id = $4
		RETURNING id, title, content, user_id, created_at, updated_at
	`

	note := &models.Note{}

	err := r.db.QueryRow(
		ctx,
		query,
		title,
		content,
		noteID,
		userID,
	).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.UserID,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	return note, nil
}

func (r *NoteRepository) DeleteNote(
	ctx context.Context,
	noteID int64,
	userID int64,
) error {
	query := `
		DELETE FROM notes
		WHERE id = $1
		  AND user_id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		noteID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}
