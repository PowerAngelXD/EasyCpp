package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"easycpp/backend/internal/repository"
	"easycpp/backend/internal/service"
	"easycpp/backend/pkg/response"
)

const (
	contextUserIDKey    = "auth_user_id"
	contextSessionIDKey = "auth_session_id"
)

func RequireAuth(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearerToken(c.GetHeader("Authorization"))
		if token == "" {
			response.Error(c, http.StatusUnauthorized, "missing authorization token")
			c.Abort()
			return
		}

		userID, sessionID, err := authService.ValidateSessionToken(c.Request.Context(), token)
		if err != nil {
			if err == repository.ErrUnauthorized {
				response.Error(c, http.StatusUnauthorized, "invalid or expired token")
				c.Abort()
				return
			}
			response.Error(c, http.StatusInternalServerError, "failed to authorize")
			c.Abort()
			return
		}

		c.Set(contextUserIDKey, userID)
		c.Set(contextSessionIDKey, sessionID)
		c.Next()
	}
}

func GetAuthUserID(c *gin.Context) (uint64, bool) {
	value, ok := c.Get(contextUserIDKey)
	if !ok {
		return 0, false
	}
	userID, ok := value.(uint64)
	return userID, ok
}

func GetSessionID(c *gin.Context) (string, bool) {
	value, ok := c.Get(contextSessionIDKey)
	if !ok {
		return "", false
	}
	sessionID, ok := value.(string)
	return sessionID, ok
}

func bearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
