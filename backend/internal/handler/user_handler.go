package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"easycpp/backend/internal/repository"
	"easycpp/backend/internal/service"
	"easycpp/backend/pkg/response"
)

// UserHandler 提供用户相关 HTTP 接口。
type UserHandler struct {
	userService service.UserService // 用户查询类业务逻辑。
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "user not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get user")
		return
	}

	response.Success(c, http.StatusOK, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page := parsePage(c.Query("page"), 1)
	pageSize := parsePage(c.Query("pageSize"), 20)

	users, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to list users")
		return
	}

	response.Success(c, http.StatusOK, users)
}
