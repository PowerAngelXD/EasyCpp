package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"easycpp/backend/internal/model"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	ListByPostID(ctx context.Context, postID uint64, limit, offset int) ([]model.Comment, error)
	DeleteByID(ctx context.Context, id, authorID uint64) error
}

// PGCommentRepository 提供基于 PostgreSQL 的评论数据访问实现。
type PGCommentRepository struct {
	db *pgxpool.Pool // PostgreSQL 连接池。
}

func NewPGCommentRepository(db *pgxpool.Pool) *PGCommentRepository {
	return &PGCommentRepository{db: db}
}

func (r *PGCommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	const query = `
		INSERT INTO comments (post_id, author_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, comment.PostID, comment.AuthorID, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PGCommentRepository) ListByPostID(ctx context.Context, postID uint64, limit, offset int) ([]model.Comment, error) {
	const query = `
		SELECT id, post_id, author_id, content, created_at, updated_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]model.Comment, 0, limit)
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

func (r *PGCommentRepository) DeleteByID(ctx context.Context, id, authorID uint64) error {
	const query = `DELETE FROM comments WHERE id = $1 AND author_id = $2`
	result, err := r.db.Exec(ctx, query, id, authorID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PGCommentRepository) ExistsPost(ctx context.Context, postID uint64) error {
	const query = `SELECT id FROM posts WHERE id = $1`
	var id uint64
	err := r.db.QueryRow(ctx, query, postID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
