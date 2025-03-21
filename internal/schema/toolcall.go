package schema

// ToolCallFunction 表示工具调用的函数
type ToolCallFunction struct {
	Name      string `json:"name"`      // 函数名称
	Arguments string `json:"arguments"` // 函数参数，JSON格式
}

// ToolCall 表示LLM生成的工具调用
type ToolCall struct {
	ID       string           `json:"id"`       // 工具调用ID
	Type     string           `json:"type"`     // 工具调用类型，通常为"function"
	Function ToolCallFunction `json:"function"` // 工具调用函数
}

// LLMResponse 表示LLM的响应
type LLMResponse struct {
	Content   string     `json:"content"`    // 响应内容
	ToolCalls []ToolCall `json:"tool_calls"` // 工具调用
}
