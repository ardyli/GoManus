package tool

import (
	"context"
	"fmt"
)

// Terminate 是一个用于终止代理执行的工具
type Terminate struct {
	*BaseTool
	parameters map[string]interface{}
}

// NewTerminate 创建新的终止工具
func NewTerminate() *Terminate {
	description := "当请求满足或助手无法继续任务时终止交互"
	baseTool := NewBaseTool("terminate", description)
	
	// 定义参数
	parameters := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"status": map[string]interface{}{
				"type":        "string",
				"description": "交互的完成状态",
				"enum":        []string{"success", "failure"},
			},
		},
		"required": []string{"status"},
	}
	
	return &Terminate{
		BaseTool:   baseTool,
		parameters: parameters,
	}
}

// Parameters 返回工具参数定义
func (t *Terminate) Parameters() map[string]interface{} {
	return t.parameters
}

// Execute 执行工具
func (t *Terminate) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// 获取状态参数
	status, ok := params["status"].(string)
	if !ok {
		return nil, fmt.Errorf("无效的状态参数")
	}
	
	// 验证状态值
	if status != "success" && status != "failure" {
		return nil, fmt.Errorf("状态必须是 'success' 或 'failure'")
	}
	
	// 返回完成消息
	return fmt.Sprintf("交互已完成，状态: %s", status), nil
}

// GetToolDefinition 返回工具定义
func (t *Terminate) GetToolDefinition() map[string]interface{} {
	return map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        t.Name(),
			"description": t.Description(),
			"parameters":  t.Parameters(),
		},
	}
}
