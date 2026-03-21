package dto

// CreateCommentRequest 定义创建评论请求体。
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"` // 评论正文。
}
