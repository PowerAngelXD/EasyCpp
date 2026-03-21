package service

import (
	"context"
	"errors"
	"time"

	"easycpp/backend/internal/dto"
	"easycpp/backend/internal/model"
	"easycpp/backend/internal/repository"
)

type PostService interface {
	CreatePost(ctx context.Context, authorID uint64, req dto.CreatePostRequest) (*model.Post, error)
	GetPostByID(ctx context.Context, id uint64) (*model.Post, error)
	ListPosts(ctx context.Context, page, pageSize int) ([]model.Post, error)
	DeletePost(ctx context.Context, id, authorID uint64) error
}

// postService 是贴文领域服务实现。
type postService struct {
	postRepo repository.PostRepository // 贴文数据读写。
	userRepo repository.UserRepository // 作者身份校验。
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) PostService {
	return &postService{postRepo: postRepo, userRepo: userRepo}
}

func (s *postService) CreatePost(ctx context.Context, authorID uint64, req dto.CreatePostRequest) (*model.Post, error) {
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, repository.ErrUnauthorized
		}
		return nil, err
	}

	now := time.Now()
	post := &model.Post{
		AuthorID:     authorID,
		AuthorName:   author.Username,
		Title:        req.Title,
		Summary:      req.Summary,
		Content:      req.Content,
		Language:     req.Language,
		Difficulty:   req.Difficulty,
		Tags:         req.Tags,
		IsPublished:  true,
		PublishedAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
		LastEditedAt: now,
	}
	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *postService) GetPostByID(ctx context.Context, id uint64) (*model.Post, error) {
	return s.postRepo.GetByID(ctx, id)
}

func (s *postService) ListPosts(ctx context.Context, page, pageSize int) ([]model.Post, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.postRepo.List(ctx, pageSize, offset)
}

func (s *postService) DeletePost(ctx context.Context, id, authorID uint64) error {
	return s.postRepo.DeleteByID(ctx, id, authorID)
}
