package model

import "time"

// Comment 表示贴文评论信息。
type Comment struct {
	ID        uint64    `json:"id"`        // 评论主键。
	PostID    uint64    `json:"postId"`    // 所属贴文 ID。
	AuthorID  uint64    `json:"authorId"`  // 评论作者用户 ID。
	Content   string    `json:"content"`   // 评论正文。
	CreatedAt time.Time `json:"createdAt"` // 创建时间。
	UpdatedAt time.Time `json:"updatedAt"` // 更新时间。
}
