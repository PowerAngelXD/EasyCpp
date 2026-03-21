package response

import "github.com/gin-gonic/gin"

// Envelope 定义统一响应结构。
type Envelope struct {
	Code    int    `json:"code"`           // 业务响应状态码。
	Message string `json:"message"`        // 响应说明信息。
	Data    any    `json:"data,omitempty"` // 响应载荷。
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, Envelope{
		Code:    status,
		Message: "ok",
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Envelope{
		Code:    status,
		Message: message,
	})
}
