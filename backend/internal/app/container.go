package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"easycpp/backend/internal/config"
	"easycpp/backend/internal/handler"
	"easycpp/backend/internal/platform/postgres"
	"easycpp/backend/internal/platform/rediscache"
	"easycpp/backend/internal/repository"
	"easycpp/backend/internal/service"
)

type Container struct {
	Config         config.Config           // 应用运行配置。
	DB             *pgxpool.Pool           // PostgreSQL 连接池实例。
	Redis          *redis.Client           // Redis 客户端实例。
	AuthService    service.AuthService     // 认证与会话业务能力。
	UserHandler    *handler.UserHandler    // 用户相关 HTTP 请求处理器。
	PostHandler    *handler.PostHandler    // 贴文相关 HTTP 请求处理器。
	AuthHandler    *handler.AuthHandler    // 注册、登录与会话管理处理器。
	CommentHandler *handler.CommentHandler // 评论相关 HTTP 请求处理器。
	CPPIdeHandler  *handler.CPPIdeHandler  // C++ 在线编译执行处理器。
}

func NewContainer(ctx context.Context) (*Container, error) {
	cfg := config.Load()

	db, err := postgres.NewPool(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}

	redisClient, err := rediscache.NewClient(ctx, cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		db.Close()
		return nil, err
	}

	if err := repository.EnsureSchema(ctx, db); err != nil {
		_ = redisClient.Close()
		db.Close()
		return nil, err
	}

	userRepo := repository.NewPGUserRepository(db)
	postRepo := repository.NewPGPostRepository(db)
	commentRepo := repository.NewPGCommentRepository(db)
	sessionRepo := repository.NewRedisSessionRepository(redisClient)

	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo, userRepo)
	authService := service.NewAuthService(cfg, userRepo, sessionRepo)
	commentService := service.NewCommentService(commentRepo, postRepo)
	cppIdeService := service.NewCPPIdeService()

	return &Container{
		Config:         cfg,
		DB:             db,
		Redis:          redisClient,
		AuthService:    authService,
		UserHandler:    handler.NewUserHandler(userService),
		PostHandler:    handler.NewPostHandler(postService),
		AuthHandler:    handler.NewAuthHandler(authService),
		CommentHandler: handler.NewCommentHandler(commentService),
		CPPIdeHandler:  handler.NewCPPIdeHandler(cppIdeService),
	}, nil
}

func (c *Container) Close() {
	if c.Redis != nil {
		_ = c.Redis.Close()
	}
	if c.DB != nil {
		c.DB.Close()
	}
}
