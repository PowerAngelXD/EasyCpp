package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"easycpp/backend/internal/dto"
	"easycpp/backend/internal/middleware"
	"easycpp/backend/internal/repository"
	"easycpp/backend/internal/service"
	"easycpp/backend/pkg/response"
)

// CommentHandler 提供评论相关 HTTP 接口。
type CommentHandler struct {
	commentService service.CommentService // 评论业务逻辑。
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post id")
		return
	}
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid user context")
		return
	}
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	comment, err := h.commentService.CreateComment(c.Request.Context(), postID, userID, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "post not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to create comment")
		return
	}
	response.Success(c, http.StatusCreated, comment)
}

func (h *CommentHandler) ListComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post id")
		return
	}
	page := parsePage(c.Query("page"), 1)
	pageSize := parsePage(c.Query("pageSize"), 20)

	comments, err := h.commentService.ListComments(c.Request.Context(), postID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to list comments")
		return
	}
	response.Success(c, http.StatusOK, comments)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid comment id")
		return
	}
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid user context")
		return
	}
	if err := h.commentService.DeleteComment(c.Request.Context(), commentID, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "comment not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to delete comment")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"deleted": true, "commentId": commentID})
}

func parsePage(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return fallback
	}
	return value
}
