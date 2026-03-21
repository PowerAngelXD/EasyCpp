package service

import (
	"context"

	"easycpp/backend/internal/model"
	"easycpp/backend/internal/repository"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uint64) (*model.User, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]model.User, error)
}

// userService 是用户领域服务实现。
type userService struct {
	userRepo repository.UserRepository // 用户数据读写。
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(ctx context.Context, id uint64) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]model.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.userRepo.List(ctx, pageSize, offset)
}
