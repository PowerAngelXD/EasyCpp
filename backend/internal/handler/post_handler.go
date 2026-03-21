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

// PostHandler 提供贴文相关 HTTP 接口。
type PostHandler struct {
	postService service.PostService // 贴文业务逻辑。
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	authorID, ok := middleware.GetAuthUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid user context")
		return
	}

	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.postService.CreatePost(c.Request.Context(), authorID, req)
	if err != nil {
		if errors.Is(err, repository.ErrUnauthorized) {
			response.Error(c, http.StatusUnauthorized, "author not found or unauthorized")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to create post")
		return
	}

	response.Success(c, http.StatusCreated, post)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := h.postService.GetPostByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "post not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get post")
		return
	}

	response.Success(c, http.StatusOK, post)
}

func (h *PostHandler) ListPosts(c *gin.Context) {
	page := parsePage(c.Query("page"), 1)
	pageSize := parsePage(c.Query("pageSize"), 20)

	posts, err := h.postService.ListPosts(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to list posts")
		return
	}

	response.Success(c, http.StatusOK, posts)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post id")
		return
	}
	authorID, ok := middleware.GetAuthUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid user context")
		return
	}

	if err := h.postService.DeletePost(c.Request.Context(), postID, authorID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "post not found or not owned by user")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to delete post")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"deleted": true, "postId": postID})
}
