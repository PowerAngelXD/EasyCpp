package model

import "time"

// User 表示社区用户信息。
type User struct {
	ID           uint64    `json:"id"`                  // 用户主键。
	Username     string    `json:"username"`            // 用户显示名称。
	Email        string    `json:"email"`               // 用户登录邮箱。
	PasswordHash string    `json:"-"`                   // 密码哈希值，不对外暴露。
	AvatarURL    string    `json:"avatarUrl,omitempty"` // 用户头像地址。
	Bio          string    `json:"bio,omitempty"`       // 用户简介。
	Role         string    `json:"role"`                // 用户角色（如 user/admin）。
	Status       string    `json:"status"`              // 用户状态（如 active/disabled）。
	CreatedAt    time.Time `json:"createdAt"`           // 创建时间。
	UpdatedAt    time.Time `json:"updatedAt"`           // 更新时间。
}
