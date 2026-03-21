package dto

// CreatePostRequest 创建贴文请求。
type CreatePostRequest struct {
	Title      string   `json:"title" binding:"required,min=5,max=120"`                             // 贴文标题。
	Summary    string   `json:"summary" binding:"max=280"`                                          // 贴文摘要内容。
	Content    string   `json:"content" binding:"required,min=10"`                                  // 贴文正文内容。
	Language   string   `json:"language" binding:"required,oneof=cpp c"`                            // 代码语言类型。
	Difficulty string   `json:"difficulty" binding:"required,oneof=beginner intermediate advanced"` // 内容难度等级。
	Tags       []string `json:"tags" binding:"max=8,dive,max=20"`                                   // 便于检索的标签集合。
}
