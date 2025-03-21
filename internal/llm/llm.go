package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gomanus/internal/config"
	"gomanus/internal/schema"
	"gomanus/pkg/logger"
)

// LLM 表示语言模型接口
type LLM struct {
	ConfigName  string
	Model       string
	BaseURL     string
	APIKey      string
	MaxTokens   int
	Temperature float64
	Client      *http.Client
}

// NewLLM 创建新的语言模型实例
func NewLLM(configName string) (*LLM, error) {
	// 从配置文件获取LLM配置
	cfg, err := config.GetLLMConfig(configName)
	if err != nil {
		return nil, fmt.Errorf("获取LLM配置失败: %w", err)
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	logger.Info("创建LLM实例: %s", configName)
	return &LLM{
		ConfigName:  configName,
		Model:       cfg.Model,
		BaseURL:     cfg.BaseURL,
		APIKey:      cfg.APIKey,
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
		Client:      client,
	}, nil
}

// AskTool 向语言模型发送消息并获取工具调用响应
func (l *LLM) AskTool(ctx context.Context, messages []schema.Message, systemMsgs []schema.Message, tools []map[string]interface{}, toolChoice string) (*schema.LLMResponse, error) {
	return l.AskWithOptions(ctx, messages, systemMsgs, tools, &toolChoice)
}

// AskWithOptions 向语言模型发送消息并获取响应，支持系统消息、工具和工具选择
func (l *LLM) AskWithOptions(
	ctx context.Context,
	messages []schema.Message,
	systemMsgs []schema.Message,
	tools []map[string]interface{},
	toolChoice *string,
) (*schema.LLMResponse, error) {
	// 准备请求体
	allMessages := make([]map[string]interface{}, 0)

	// 添加系统消息
	if systemMsgs != nil && len(systemMsgs) > 0 {
		for _, msg := range systemMsgs {
			logger.Info("添加系统消息: %s", msg.Content)
			allMessages = append(allMessages, map[string]interface{}{
				"role":    msg.Role,
				"content": msg.Content,
			})
		}
	}

	// 添加用户和助手消息
	for _, msg := range messages {
		msgMap := map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		}

		// 如果是工具消息，添加工具相关字段
		if msg.Role == "tool" {
			msgMap["tool_call_id"] = msg.ToolCallID
			if msg.Name != "" {
				msgMap["name"] = msg.Name
			}
		}

		// 如果是助手消息且有工具调用，添加工具调用
		if msg.Role == "assistant" && len(msg.ToolCalls) > 0 {
			toolCalls := make([]map[string]interface{}, 0)
			for _, tc := range msg.ToolCalls {
				toolCall := map[string]interface{}{
					"id":   tc.ID,
					"type": "function",
					"function": map[string]interface{}{
						"name":      tc.Function.Name,
						"arguments": tc.Function.Arguments,
					},
				}
				toolCalls = append(toolCalls, toolCall)
			}
			msgMap["tool_calls"] = toolCalls
		}

		allMessages = append(allMessages, msgMap)
	}

	// 准备请求体
	requestBody := map[string]interface{}{
		"model":       l.Model,
		"messages":    allMessages,
		"max_tokens":  l.MaxTokens,
		"temperature": l.Temperature,
	}

	// 添加工具和工具选择
	if tools != nil && len(tools) > 0 {
		requestBody["tools"] = tools

		if toolChoice != nil {
			if *toolChoice == "none" {
				requestBody["tool_choice"] = "none"
			} else if *toolChoice == "auto" {
				requestBody["tool_choice"] = "auto"
			} else if *toolChoice == "required" {
				requestBody["tool_choice"] = map[string]interface{}{
					"type": "function",
				}
			}
		}
	}

	// 序列化请求体
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 记录请求详情
	logger.Info("发送LLM请求到: %s", l.BaseURL+"chat/completions")
	logger.Info("请求模型: %s", l.Model)
	logger.Info("请求工具数量: %d", len(tools))
	logger.Info("请求消息数量: %d", len(allMessages))

	// 记录请求体（但不包含敏感信息）
	prettyRequest, _ := json.MarshalIndent(requestBody, "", "  ")
	logger.Info("LLM请求体: %s", string(prettyRequest))

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", l.BaseURL+"chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if l.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+l.APIKey)
	}

	// 发送请求
	logger.Info("发送LLM请求: %s", l.Model)
	start := time.Now()
	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 记录原始响应
	logger.Info("LLM响应状态码: %d", resp.StatusCode)
	logger.Info("LLM原始响应: %s", string(body))

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LLM请求失败: %s, 状态码: %d, 响应: %s", l.Model, resp.StatusCode, string(body))
	}

	// 解析响应
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取响应内容
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("无效的响应格式: %s", string(body))
	}

	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	content, _ := message["content"].(string)

	// 提取工具调用
	var toolCalls []schema.ToolCall
	if toolCallsRaw, ok := message["tool_calls"].([]interface{}); ok && len(toolCallsRaw) > 0 {
		for _, tcRaw := range toolCallsRaw {
			tc := tcRaw.(map[string]interface{})
			id := tc["id"].(string)
			tcType := tc["type"].(string)
			function := tc["function"].(map[string]interface{})
			name := function["name"].(string)
			args := function["arguments"].(string)

			toolCalls = append(toolCalls, schema.ToolCall{
				ID:   id,
				Type: tcType,
				Function: schema.ToolCallFunction{
					Name:      name,
					Arguments: args,
				},
			})

			// 记录工具调用详情
			logger.Info("检测到工具调用: %s, 参数: %s", name, args)
		}
	} else {
		logger.Info("LLM响应中没有工具调用")
	}

	// 计算耗时
	elapsed := time.Since(start)
	logger.Info("LLM响应耗时: %s", elapsed)

	// 返回响应
	return &schema.LLMResponse{
		Content:   content,
		ToolCalls: toolCalls,
	}, nil
}
