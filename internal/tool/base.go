package tool

import (
	"context"
	"fmt"
)

// Tool 定义工具的基本接口
type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

// BaseTool 提供工具的基础实现
type BaseTool struct {
	name        string
	description string
}

// NewBaseTool 创建新的基础工具
func NewBaseTool(name, desc string) *BaseTool {
	return &BaseTool{
		name:        name,
		description: desc,
	}
}

// Name 返回工具名称
func (t *BaseTool) Name() string {
	return t.name
}

// Description 返回工具描述
func (t *BaseTool) Description() string {
	return t.description
}

// ToolResult 表示工具执行的结果
type ToolResult struct {
	Output    interface{}
	Error     error
	SystemMsg string
}

// NewToolResult 创建新的工具结果
func NewToolResult(output interface{}, err error, sysMsg string) *ToolResult {
	return &ToolResult{
		Output:    output,
		Error:     err,
		SystemMsg: sysMsg,
	}
}

// IsSuccess 检查工具执行是否成功
func (r *ToolResult) IsSuccess() bool {
	return r.Error == nil
}

// String 返回工具结果的字符串表示
func (r *ToolResult) String() string {
	if r.Error != nil {
		return "Error: " + r.Error.Error()
	}
	if str, ok := r.Output.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", r.Output)
}

// Combine 合并两个工具结果
func (r *ToolResult) Combine(other *ToolResult) *ToolResult {
	if r == nil {
		return other
	}
	if other == nil {
		return r
	}

	var output interface{}
	if r.Output != nil && other.Output != nil {
		// 尝试合并输出
		if rStr, ok := r.Output.(string); ok {
			if oStr, ok := other.Output.(string); ok {
				output = rStr + oStr
			} else {
				output = r.Output
			}
		} else {
			output = r.Output
		}
	} else {
		output = r.Output
		if output == nil {
			output = other.Output
		}
	}

	var err error
	if r.Error != nil {
		err = r.Error
	} else {
		err = other.Error
	}

	var sysMsg string
	if r.SystemMsg != "" && other.SystemMsg != "" {
		sysMsg = r.SystemMsg + "\n" + other.SystemMsg
	} else if r.SystemMsg != "" {
		sysMsg = r.SystemMsg
	} else {
		sysMsg = other.SystemMsg
	}

	return &ToolResult{
		Output:    output,
		Error:     err,
		SystemMsg: sysMsg,
	}
}
