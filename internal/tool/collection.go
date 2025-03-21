package tool

import (
	"context"
	"fmt"
	"sync"
)

// ToolCollection 管理一组工具
type ToolCollection struct {
	tools map[string]Tool
	mu    sync.RWMutex
}

// NewToolCollection 创建新的工具集合
func NewToolCollection() *ToolCollection {
	return &ToolCollection{
		tools: make(map[string]Tool),
	}
}

// AddTool 添加工具到集合
func (tc *ToolCollection) AddTool(tool Tool) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	name := tool.Name()
	if name == "" {
		return fmt.Errorf("工具名称不能为空")
	}

	if _, exists := tc.tools[name]; exists {
		return fmt.Errorf("工具 %s 已存在", name)
	}

	tc.tools[name] = tool
	return nil
}

// GetTool 根据名称获取工具
func (tc *ToolCollection) GetTool(name string) (Tool, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	tool, exists := tc.tools[name]
	if !exists {
		return nil, fmt.Errorf("工具 %s 不存在", name)
	}

	return tool, nil
}

// RemoveTool 从集合中移除工具
func (tc *ToolCollection) RemoveTool(name string) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if _, exists := tc.tools[name]; !exists {
		return fmt.Errorf("工具 %s 不存在", name)
	}

	delete(tc.tools, name)
	return nil
}

// GetAllTools 获取所有工具
func (tc *ToolCollection) GetAllTools() []Tool {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	tools := make([]Tool, 0, len(tc.tools))
	for _, tool := range tc.tools {
		tools = append(tools, tool)
	}

	return tools
}

// GetToolDefinitions 获取所有工具的定义，用于LLM
func (tc *ToolCollection) GetToolDefinitions() []map[string]interface{} {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	definitions := make([]map[string]interface{}, 0, len(tc.tools))
	for _, tool := range tc.tools {
		// 基本定义
		functionDef := map[string]interface{}{
			"name":        tool.Name(),
			"description": tool.Description(),
		}
		
		// 检查工具是否实现了Parameters方法
		if paramProvider, ok := tool.(interface{ Parameters() map[string]interface{} }); ok {
			functionDef["parameters"] = paramProvider.Parameters()
		}
		
		// 检查工具是否实现了GetToolDefinition方法
		if defProvider, ok := tool.(interface{ GetToolDefinition() map[string]interface{} }); ok {
			// 如果工具提供了完整的定义，直接使用它
			definitions = append(definitions, defProvider.GetToolDefinition())
			continue
		}
		
		// 否则使用构建的定义
		definition := map[string]interface{}{
			"type":     "function",
			"function": functionDef,
		}
		definitions = append(definitions, definition)
	}

	return definitions
}

// ExecuteTool 执行指定的工具
func (tc *ToolCollection) ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	tool, err := tc.GetTool(name)
	if err != nil {
		return nil, err
	}
	return tool.Execute(ctx, params)
}
