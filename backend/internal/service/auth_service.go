package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"easycpp/backend/internal/config"
	"easycpp/backend/internal/dto"
	"easycpp/backend/internal/model"
	"easycpp/backend/internal/repository"
	"easycpp/backend/internal/security"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req dto.LoginRequest, userAgent, remoteAddr string) (*dto.LoginResponse, error)
	Logout(ctx context.Context, sessionID string) error
	ListSessions(ctx context.Context, userID uint64) ([]dto.SessionResponse, error)
	RevokeSession(ctx context.Context, userID uint64, sessionID string) error
	ValidateSessionToken(ctx context.Context, token string) (uint64, string, error)
}

// authService 是认证与会话服务实现。
type authService struct {
	cfg         config.Config                // 认证相关配置参数。
	userRepo    repository.UserRepository    // 用户信息查询与创建。
	sessionRepo repository.SessionRepository // 会话存储与校验。
}

func NewAuthService(cfg config.Config, userRepo repository.UserRepository, sessionRepo repository.SessionRepository) AuthService {
	return &authService{cfg: cfg, userRepo: userRepo, sessionRepo: sessionRepo}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, repository.ErrConflict
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.PasswordBcryptCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Bio:          req.Bio,
		Role:         "user",
		Status:       "active",
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest, userAgent, remoteAddr string) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, repository.ErrUnauthorized
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, repository.ErrUnauthorized
	}

	sessionID := uuid.NewString()
	now := time.Now()
	session := repository.Session{
		SessionID:  sessionID,
		UserID:     user.ID,
		UserAgent:  userAgent,
		RemoteAddr: remoteAddr,
		IssuedAt:   now.Unix(),
		ExpiresAt:  now.Add(s.cfg.SessionTTL).Unix(),
	}
	if err := s.sessionRepo.Save(ctx, session, s.cfg.SessionTTL); err != nil {
		return nil, err
	}

	token, expiresAt, err := security.GenerateAccessToken(s.cfg.JWTSecret, user.ID, sessionID, s.cfg.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresAt,
		SessionID:   sessionID,
	}, nil
}

func (s *authService) Logout(ctx context.Context, sessionID string) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

func (s *authService) ListSessions(ctx context.Context, userID uint64) ([]dto.SessionResponse, error) {
	sessions, err := s.sessionRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.SessionResponse, 0, len(sessions))
	for _, session := range sessions {
		if session.ExpiresAt < time.Now().Unix() {
			continue
		}
		result = append(result, dto.SessionResponse{
			SessionID:  session.SessionID,
			UserID:     session.UserID,
			UserAgent:  session.UserAgent,
			RemoteAddr: session.RemoteAddr,
			IssuedAt:   session.IssuedAt,
			ExpiresAt:  session.ExpiresAt,
		})
	}
	return result, nil
}

func (s *authService) RevokeSession(ctx context.Context, userID uint64, sessionID string) error {
	session, err := s.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return err
	}
	if session.UserID != userID {
		return repository.ErrForbidden
	}
	return s.sessionRepo.Delete(ctx, sessionID)
}

func (s *authService) ValidateSessionToken(ctx context.Context, token string) (uint64, string, error) {
	claims, err := security.ParseAccessToken(s.cfg.JWTSecret, token)
	if err != nil {
		return 0, "", repository.ErrUnauthorized
	}

	session, err := s.sessionRepo.Get(ctx, claims.SessionID)
	if err != nil {
		return 0, "", repository.ErrUnauthorized
	}
	if session.UserID != claims.UserID || session.ExpiresAt < time.Now().Unix() {
		return 0, "", repository.ErrUnauthorized
	}
	return claims.UserID, claims.SessionID, nil
}
