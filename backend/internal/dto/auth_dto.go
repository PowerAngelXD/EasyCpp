package dto

// RegisterRequest 定义注册请求体。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"` // 注册时的用户名。
	Email    string `json:"email" binding:"required,email"`           // 注册账号邮箱。
	Password string `json:"password" binding:"required,min=8,max=72"` // 注册明文密码。
	Bio      string `json:"bio" binding:"max=200"`                    // 可选的用户简介。
}

// LoginRequest 定义登录请求体。
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`           // 登录邮箱。
	Password string `json:"password" binding:"required,min=8,max=72"` // 登录明文密码。
}

// LoginResponse 定义登录成功响应体。
type LoginResponse struct {
	AccessToken string `json:"accessToken"` // JWT 访问令牌。
	TokenType   string `json:"tokenType"`   // 鉴权头类型，通常为 Bearer。
	ExpiresIn   int64  `json:"expiresIn"`   // 访问令牌过期时间戳（Unix 秒）。
	SessionID   string `json:"sessionId"`   // 服务端会话标识。
}

// SessionResponse 定义会话信息响应体。
type SessionResponse struct {
	SessionID  string `json:"sessionId"`            // 当前会话唯一标识。
	UserID     uint64 `json:"userId"`               // 会话所属用户 ID。
	UserAgent  string `json:"userAgent,omitempty"`  // 发起会话的客户端标识。
	RemoteAddr string `json:"remoteAddr,omitempty"` // 发起会话的来源 IP。
	IssuedAt   int64  `json:"issuedAt"`             // 会话签发时间戳（Unix 秒）。
	ExpiresAt  int64  `json:"expiresAt"`            // 会话到期时间戳（Unix 秒）。
}
