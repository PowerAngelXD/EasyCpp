package service

import (
	"context"
	"errors"

	"easycpp/backend/internal/dto"
	"easycpp/backend/internal/model"
	"easycpp/backend/internal/repository"
)

type CommentService interface {
	CreateComment(ctx context.Context, postID, authorID uint64, req dto.CreateCommentRequest) (*model.Comment, error)
	ListComments(ctx context.Context, postID uint64, page, pageSize int) ([]model.Comment, error)
	DeleteComment(ctx context.Context, commentID, authorID uint64) error
}

// commentService 是评论领域服务实现。
type commentService struct {
	commentRepo *repository.PGCommentRepository // 评论数据读写。
	postRepo    repository.PostRepository       // 贴文存在性校验。
}

func NewCommentService(commentRepo *repository.PGCommentRepository, postRepo repository.PostRepository) CommentService {
	return &commentService{commentRepo: commentRepo, postRepo: postRepo}
}

func (s *commentService) CreateComment(ctx context.Context, postID, authorID uint64, req dto.CreateCommentRequest) (*model.Comment, error) {
	if _, err := s.postRepo.GetByID(ctx, postID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	comment := &model.Comment{
		PostID:   postID,
		AuthorID: authorID,
		Content:  req.Content,
	}
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *commentService) ListComments(ctx context.Context, postID uint64, page, pageSize int) ([]model.Comment, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.commentRepo.ListByPostID(ctx, postID, pageSize, offset)
}

func (s *commentService) DeleteComment(ctx context.Context, commentID, authorID uint64) error {
	return s.commentRepo.DeleteByID(ctx, commentID, authorID)
}
