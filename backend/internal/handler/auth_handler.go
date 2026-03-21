package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"easycpp/backend/internal/dto"
	"easycpp/backend/internal/middleware"
	"easycpp/backend/internal/repository"
	"easycpp/backend/internal/service"
	"easycpp/backend/pkg/response"
)

// AuthHandler 提供认证与会话相关 HTTP 接口。
type AuthHandler struct {
	authService service.AuthService // 认证与会话业务逻辑。
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			response.Error(c, http.StatusConflict, "email already exists")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to register")
		return
	}

	response.Success(c, http.StatusCreated, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		if errors.Is(err, repository.ErrUnauthorized) {
			response.Error(c, http.StatusUnauthorized, "invalid email or password")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to login")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, ok := middleware.GetSessionID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid session context")
		return
	}
	if err := h.authService.Logout(c.Request.Context(), sessionID); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to logout")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"loggedOut": true})
}

func (h *AuthHandler) ListSessions(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid user context")
		return
	}
	sessions, err := h.authService.ListSessions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to list sessions")
		return
	}
	response.Success(c, http.StatusOK, sessions)
}

func (h *AuthHandler) RevokeSession(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid user context")
		return
	}
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		response.Error(c, http.StatusBadRequest, "session id is required")
		return
	}
	if err := h.authService.RevokeSession(c.Request.Context(), userID, sessionID); err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			response.Error(c, http.StatusForbidden, "cannot revoke another user's session")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to revoke session")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"revoked": true, "sessionId": sessionID})
}
