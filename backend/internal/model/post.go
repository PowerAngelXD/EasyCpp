package model

import "time"

// Post 表示社区贴文信息。
type Post struct {
	ID           uint64    `json:"id"`                     // 贴文主键。
	AuthorID     uint64    `json:"authorId"`               // 作者用户 ID。
	AuthorName   string    `json:"authorName"`             // 作者展示名快照。
	Title        string    `json:"title"`                  // 贴文标题。
	Summary      string    `json:"summary,omitempty"`      // 贴文摘要。
	Content      string    `json:"content"`                // 贴文正文。
	Language     string    `json:"language"`               // 代码语言类型。
	Difficulty   string    `json:"difficulty"`             // 内容难度等级。
	Tags         []string  `json:"tags,omitempty"`         // 贴文标签集合。
	ViewCount    int64     `json:"viewCount"`              // 浏览次数。
	LikeCount    int64     `json:"likeCount"`              // 点赞次数。
	CommentCount int64     `json:"commentCount"`           // 评论总数。
	IsPinned     bool      `json:"isPinned"`               // 是否置顶。
	IsPublished  bool      `json:"isPublished"`            // 是否已发布。
	PublishedAt  time.Time `json:"publishedAt,omitempty"`  // 发布时间。
	LastEditedAt time.Time `json:"lastEditedAt,omitempty"` // 最近编辑时间。
	CreatedAt    time.Time `json:"createdAt"`              // 创建时间。
	UpdatedAt    time.Time `json:"updatedAt"`              // 更新时间。
}
