package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"easycpp/backend/internal/middleware"
	"easycpp/backend/pkg/response"
)

func NewRouter(ctx context.Context) (*gin.Engine, *Container, error) {
	container, err := NewContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			response.Success(c, http.StatusOK, gin.H{"status": "ok"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", container.AuthHandler.Register)
			auth.POST("/login", container.AuthHandler.Login)
		}

		api.GET("/users", container.UserHandler.ListUsers)
		api.GET("/users/:id", container.UserHandler.GetUserByID)

		api.GET("/posts/:id", container.PostHandler.GetPostByID)
		api.GET("/posts", container.PostHandler.ListPosts)
		api.GET("/posts/:id/comments", container.CommentHandler.ListComments)
		api.POST("/ide/cpp/run", container.CPPIdeHandler.RunCPP)

		protected := api.Group("")
		protected.Use(middleware.RequireAuth(container.AuthService))
		{
			protected.POST("/auth/logout", container.AuthHandler.Logout)
			protected.GET("/auth/sessions", container.AuthHandler.ListSessions)
			protected.DELETE("/auth/sessions/:sessionId", container.AuthHandler.RevokeSession)

			protected.POST("/posts", container.PostHandler.CreatePost)
			protected.DELETE("/posts/:id", container.PostHandler.DeletePost)

			protected.POST("/posts/:id/comments", container.CommentHandler.CreateComment)
			protected.DELETE("/comments/:commentId", container.CommentHandler.DeleteComment)
		}
	}

	return r, container, nil
}
