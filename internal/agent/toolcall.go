package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"gomanus/internal/llm"
	"gomanus/internal/schema"
	"gomanus/internal/tool"
	"gomanus/pkg/logger"
)

// ToolCallAgent 实现了工具调用的代理
type ToolCallAgent struct {
	*ReActAgent
	Tools *tool.ToolCollection
}

// NewToolCallAgent 创建新的工具调用代理
func NewToolCallAgent(name string, llm *llm.LLM, tools *tool.ToolCollection) *ToolCallAgent {
	reactAgent := NewReActAgent(name, llm)
	reactAgent.Description = "工具调用代理 - 能够使用工具执行任务"
	
	return &ToolCallAgent{
		ReActAgent: reactAgent,
		Tools:      tools,
	}
}

// Think 思考下一步行动，解析LLM响应中的工具调用
func (a *ToolCallAgent) Think(ctx context.Context) (bool, error) {
	// 获取所有消息
	messages := a.Memory.GetMessages()
	if len(messages) == 0 {
		return false, fmt.Errorf("没有消息可处理")
	}
	
	// 向LLM发送请求
	logger.Info("向LLM发送请求...")
	response, err := a.LLM.AskTool(ctx, messages, nil, a.Tools.GetToolDefinitions(), "auto")
	if err != nil {
		return false, fmt.Errorf("发送消息到LLM失败: %w", err)
	}
	
	// 将LLM响应添加到记忆中
	a.AddMessage(schema.Message{
		Role:      "assistant",
		Content:   response.Content,
		ToolCalls: response.ToolCalls,
	})
	
	// 检查是否有工具调用
	if len(response.ToolCalls) == 0 {
		logger.Info("LLM响应中没有工具调用")
		return false, nil
	}
	
	// 处理工具调用
	logger.Info("发现 %d 个工具调用", len(response.ToolCalls))
	for _, tc := range response.ToolCalls {
		a.AddMessage(schema.Message{
			Role:       "tool",
			ToolCallID: tc.ID,
			Content:    "",
		})
	}
	
	return true, nil
}

// Act 执行工具调用
func (a *ToolCallAgent) Act(ctx context.Context) (string, error) {
	// 获取最近的消息
	messages := a.Memory.GetMessages()
	if len(messages) < 2 {
		return "", fmt.Errorf("没有足够的消息来执行工具调用")
	}
	
	// 获取最近的LLM响应
	var llmResponse schema.Message
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "assistant" && len(messages[i].ToolCalls) > 0 {
			llmResponse = messages[i]
			break
		}
	}
	
	if len(llmResponse.ToolCalls) == 0 {
		return "", fmt.Errorf("没有找到工具调用")
	}
	
	// 执行每个工具调用
	var results []string
	for _, tc := range llmResponse.ToolCalls {
		logger.Info("执行工具调用: %s", tc.Function.Name)
		
		// 查找工具
		tool, err := a.Tools.GetTool(tc.Function.Name)
		if err != nil {
			errMsg := fmt.Sprintf("找不到工具 %s: %v", tc.Function.Name, err)
			logger.Error(errMsg)
			
			// 添加错误消息
			a.AddMessage(schema.Message{
				Role:       "tool",
				ToolCallID: tc.ID,
				Content:    errMsg,
			})
			
			results = append(results, errMsg)
			continue
		}
		
		// 解析参数
		var params map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &params); err != nil {
			errMsg := fmt.Sprintf("解析工具参数失败: %v", err)
			logger.Error(errMsg)
			
			// 添加错误消息
			a.AddMessage(schema.Message{
				Role:       "tool",
				ToolCallID: tc.ID,
				Content:    errMsg,
			})
			
			results = append(results, errMsg)
			continue
		}
		
		// 执行工具
		result, err := tool.Execute(ctx, params)
		if err != nil {
			errMsg := fmt.Sprintf("执行工具失败: %v", err)
			logger.Error(errMsg)
			
			// 添加错误消息
			a.AddMessage(schema.Message{
				Role:       "tool",
				ToolCallID: tc.ID,
				Content:    errMsg,
			})
			
			results = append(results, errMsg)
			continue
		}
		
		// 添加工具结果
		resultStr := fmt.Sprintf("%v", result)
		a.AddMessage(schema.Message{
			Role:       "tool",
			ToolCallID: tc.ID,
			Content:    resultStr,
		})
		
		results = append(results, fmt.Sprintf("工具 %s 执行结果: %s", tc.Function.Name, resultStr))
	}
	
	return fmt.Sprintf("执行了 %d 个工具调用:\n%s", len(results), results), nil
}
