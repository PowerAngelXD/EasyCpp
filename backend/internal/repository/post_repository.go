package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"

	"easycpp/backend/internal/model"
)

type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	GetByID(ctx context.Context, id uint64) (*model.Post, error)
	List(ctx context.Context, limit, offset int) ([]model.Post, error)
	DeleteByID(ctx context.Context, id, authorID uint64) error
}

// PGPostRepository 提供基于 PostgreSQL 的贴文数据访问实现。
type PGPostRepository struct {
	db *pgxpool.Pool // PostgreSQL 连接池。
}

func NewPGPostRepository(db *pgxpool.Pool) *PGPostRepository {
	return &PGPostRepository{db: db}
}

func (r *PGPostRepository) Create(ctx context.Context, post *model.Post) error {
	const query = `
		INSERT INTO posts (
			author_id, author_name, title, summary, content, language, difficulty, tags,
			view_count, like_count, comment_count, is_pinned, is_published, published_at, last_edited_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14, $15
		)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query,
		post.AuthorID,
		post.AuthorName,
		post.Title,
		post.Summary,
		post.Content,
		post.Language,
		post.Difficulty,
		pq.Array(post.Tags),
		post.ViewCount,
		post.LikeCount,
		post.CommentCount,
		post.IsPinned,
		post.IsPublished,
		post.PublishedAt,
		post.LastEditedAt,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	return err
}

func (r *PGPostRepository) GetByID(ctx context.Context, id uint64) (*model.Post, error) {
	const query = `
		SELECT
			id, author_id, author_name, title, summary, content, language, difficulty,
			tags, view_count, like_count, comment_count, is_pinned, is_published,
			published_at, last_edited_at, created_at, updated_at
		FROM posts WHERE id = $1`
	var post model.Post
	err := r.db.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.AuthorID,
		&post.AuthorName,
		&post.Title,
		&post.Summary,
		&post.Content,
		&post.Language,
		&post.Difficulty,
		pq.Array(&post.Tags),
		&post.ViewCount,
		&post.LikeCount,
		&post.CommentCount,
		&post.IsPinned,
		&post.IsPublished,
		&post.PublishedAt,
		&post.LastEditedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (r *PGPostRepository) List(ctx context.Context, limit, offset int) ([]model.Post, error) {
	const query = `
		SELECT
			id, author_id, author_name, title, summary, content, language, difficulty,
			tags, view_count, like_count, comment_count, is_pinned, is_published,
			published_at, last_edited_at, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]model.Post, 0, limit)
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.AuthorName,
			&post.Title,
			&post.Summary,
			&post.Content,
			&post.Language,
			&post.Difficulty,
			pq.Array(&post.Tags),
			&post.ViewCount,
			&post.LikeCount,
			&post.CommentCount,
			&post.IsPinned,
			&post.IsPublished,
			&post.PublishedAt,
			&post.LastEditedAt,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

func (r *PGPostRepository) DeleteByID(ctx context.Context, id, authorID uint64) error {
	const query = `DELETE FROM posts WHERE id = $1 AND author_id = $2`
	result, err := r.db.Exec(ctx, query, id, authorID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
