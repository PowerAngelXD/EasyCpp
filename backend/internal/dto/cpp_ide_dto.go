package dto

// RunCPPRequest 表示在线 IDE 的 C++ 运行请求。
type RunCPPRequest struct {
	Code        string `json:"code" binding:"required,min=1,max=50000"` // C++ 源码。
	Stdin       string `json:"stdin" binding:"max=10000"`               // 标准输入。
	TimeLimitMs int    `json:"timeLimitMs"`                             // 运行超时（毫秒）。
}

// CPPCompileResult 表示编译阶段结果。
type CPPCompileResult struct {
	Succeeded bool   `json:"succeeded"` // 是否编译成功。
	ExitCode  int    `json:"exitCode"`  // 编译进程退出码。
	Stderr    string `json:"stderr"`    // 编译错误信息。
}

// CPPRunResult 表示运行阶段结果。
type CPPRunResult struct {
	Succeeded bool   `json:"succeeded"` // 是否运行成功（退出码为 0）。
	ExitCode  int    `json:"exitCode"`  // 运行进程退出码。
	Stdout    string `json:"stdout"`    // 程序标准输出。
	Stderr    string `json:"stderr"`    // 程序标准错误输出。
	TimedOut  bool   `json:"timedOut"`  // 是否超时。
}

// RunCPPResponse 表示一次完整编译执行结果。
type RunCPPResponse struct {
	Language   string           `json:"language"`   // 当前固定为 cpp。
	Compile    CPPCompileResult `json:"compile"`    // 编译阶段结果。
	Run        *CPPRunResult    `json:"run"`        // 运行阶段结果，编译失败时为 nil。
	DurationMs int64            `json:"durationMs"` // 总耗时（毫秒）。
}
