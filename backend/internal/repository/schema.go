package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EnsureSchema(ctx context.Context, db *pgxpool.Pool) error {
	const schemaSQL = `
CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	username VARCHAR(32) NOT NULL,
	email VARCHAR(255) NOT NULL UNIQUE,
	password_hash VARCHAR(255) NOT NULL,
	avatar_url TEXT NOT NULL DEFAULT '',
	bio TEXT NOT NULL DEFAULT '',
	role VARCHAR(32) NOT NULL DEFAULT 'user',
	status VARCHAR(32) NOT NULL DEFAULT 'active',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
	id BIGSERIAL PRIMARY KEY,
	author_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	author_name VARCHAR(64) NOT NULL,
	title VARCHAR(120) NOT NULL,
	summary VARCHAR(280) NOT NULL DEFAULT '',
	content TEXT NOT NULL,
	language VARCHAR(16) NOT NULL,
	difficulty VARCHAR(32) NOT NULL,
	tags TEXT[] NOT NULL DEFAULT '{}',
	view_count BIGINT NOT NULL DEFAULT 0,
	like_count BIGINT NOT NULL DEFAULT 0,
	comment_count BIGINT NOT NULL DEFAULT 0,
	is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
	is_published BOOLEAN NOT NULL DEFAULT TRUE,
	published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	last_edited_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments (
	id BIGSERIAL PRIMARY KEY,
	post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
	author_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_author_id ON comments(author_id);
`
	_, err := db.Exec(ctx, schemaSQL)
	return err
}
