package dto

// CreateUserRequest 创建用户请求。
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"` // 用户公开显示名。
	Email    string `json:"email" binding:"required,email"`           // 用于登录和通知的邮箱地址。
	Password string `json:"password" binding:"required,min=6,max=64"` // 用户提交的明文密码（仅用于注册阶段）。
	Bio      string `json:"bio" binding:"max=200"`                    // 用户个人简介。
}
