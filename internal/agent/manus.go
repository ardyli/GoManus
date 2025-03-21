package agent

import (
	"context"
	"fmt"
	"gomanus/internal/llm"
	"gomanus/internal/schema"
	"gomanus/internal/tool"
	"gomanus/pkg/logger"
	"strings"
)

// Manus 是gomanus的主代理，继承自ToolCallAgent
type Manus struct {
	*ToolCallAgent
	SystemPrompt string
}

// NewManus 创建新的Manus代理
func NewManus(name string, llm *llm.LLM, tools *tool.ToolCollection) *Manus {
	toolCallAgent := NewToolCallAgent(name, llm, tools)
	toolCallAgent.Description = "GoManus代理 - gomanus的主代理"

	// 创建包含工具说明的系统提示
	systemPrompt := `你是GoManus，一个强大的AI助手，能够使用各种工具帮助用户完成任务。

你可以使用以下工具：

1. google_search - 执行Google搜索并返回相关链接列表，用于查找网络信息
   参数: 
   - query: 搜索查询（必填）
   - num_results: 返回结果数量（可选，默认10）

2. zhihu_search - 执行知乎搜索并返回相关问题和回答的链接列表
   参数:
   - query: 搜索查询（必填）
   - num_results: 返回结果数量（可选，默认10）
   - search_type: 搜索类型（可选，"general"综合、"question"问题、"article"文章，默认"general"）

3. baidu_baike_search - 执行百度百科搜索并返回相关词条的链接和摘要
   参数:
   - query: 搜索查询（必填）
   - num_results: 返回结果数量（可选，默认5）

4. wikipedia_search - 执行维基百科搜索并返回相关条目的链接和摘要
   参数:
   - query: 搜索查询（必填）
   - language: 搜索的语言版本（可选，"zh"中文、"en"英文，默认"zh"）
   - num_results: 返回结果数量（可选，默认5）

5. file_operator - 对文件进行读取和保存操作，支持txt、md、pdf、png、jpg等格式
   参数:
   - operation: 操作类型（必填，"read"表示读取，"write"表示写入）
   - file_path: 文件路径（必填）
   - content: 要保存的内容（写入操作必填）
   - mode: 文件打开模式（写入操作可选，"w"写入或"a"追加，默认"w"）
   - encoding: 文件编码格式（读取操作可选，默认"utf-8"）
   - max_size: 读取文件的最大大小（读取操作可选，默认10MB）
   
   注意：当读取png或jpg图像文件时，系统会自动调用视觉模型对图像内容进行分析，并返回详细的文本描述。

6. browser_use - 与网页浏览器交互，执行各种操作
   参数:
   - action: 要执行的浏览器操作（必填，可选值: "navigate", "get_html", "execute_js", "new_tab", "close_tab"）
   - url: 'navigate'或'new_tab'操作的URL
   - script: 'execute_js'操作的JavaScript代码
   - tab_id: 操作的标签页ID

7. terminate - 终止当前执行
   参数:
   - status: 终止状态（必填，可选值: "success", "failure"）
   - message: 终止消息（可选）

当用户请求需要使用这些工具的任务时，请主动调用适当的工具来完成任务。每个工具都有特定的用途，请根据用户的需求选择最合适的工具。`

	return &Manus{
		ToolCallAgent: toolCallAgent,
		SystemPrompt:  systemPrompt,
	}
}

// SetSystemPrompt 设置系统提示
func (a *Manus) SetSystemPrompt(prompt string) {
	a.SystemPrompt = prompt
}

// Run 重写Run方法，添加系统提示并确保与AI交互
func (a *Manus) Run(ctx context.Context, request string) (string, error) {
	logger.Info("Manus代理开始运行...")

	// 重置代理状态
	a.SetState(StateIdle)
	a.ResetSteps() // 这个方法是从BaseAgent继承的

	// 添加系统提示
	if a.SystemPrompt != "" {
		logger.Info("添加系统提示: %s", a.SystemPrompt)
		a.AddMessage(schema.NewSystemMessage(a.SystemPrompt))
	}

	// 添加用户消息到记忆中
	logger.Info("添加用户消息: %s", request)
	a.AddMessage(schema.NewUserMessage(request))

	// 设置状态为运行中
	a.SetState(StateRunning)
	defer a.SetState(StateIdle)

	// 执行步骤直到达到最大步骤数或代理状态变为已完成
	var results []string
	// 使用BaseAgent的方法访问和修改步骤计数
	for a.BaseAgent.GetCurrentStep() < a.BaseAgent.GetMaxSteps() && a.GetState() != StateFinished {
		a.BaseAgent.CurrentStep++ // 直接访问BaseAgent的字段
		stepNum := a.BaseAgent.CurrentStep

		logger.Info("执行步骤 %d/%d", stepNum, a.BaseAgent.MaxSteps)

		// 执行单个步骤，这里会调用Manus的Step方法，而不是基类的
		result, err := a.Step(ctx)
		if err != nil {
			a.SetState(StateError)
			logger.Error("步骤 %d 执行失败: %v", stepNum, err)
			return "", fmt.Errorf("步骤 %d 执行失败: %w", stepNum, err)
		}

		// 记录步骤结果
		stepResult := fmt.Sprintf("步骤 %d: %s", stepNum, result)
		logger.Info(stepResult)
		results = append(results, stepResult)

		// 检查是否陷入循环
		if a.isStuck() {
			a.handleStuckState()
		}

		// 添加上下文取消检查
		select {
		case <-ctx.Done():
			logger.Warn("代理执行被上下文取消")
			results = append(results, "执行被取消")
			return strings.Join(results, "\n"), ctx.Err()
		default:
			// 继续执行
		}
	}

	// 检查是否达到最大步骤数
	if a.BaseAgent.CurrentStep >= a.BaseAgent.MaxSteps {
		maxStepsMsg := fmt.Sprintf("终止: 达到最大步骤数 (%d)", a.BaseAgent.MaxSteps)
		logger.Warn(maxStepsMsg)
		results = append(results, maxStepsMsg)
	}

	// 返回所有步骤的结果
	if len(results) == 0 {
		// 如果没有执行任何步骤，尝试获取最后的助手消息
		messages := a.Memory.GetMessages()
		for i := len(messages) - 1; i >= 0; i-- {
			if messages[i].Role == "assistant" {
				return messages[i].Content, nil
			}
		}
		return "未执行任何步骤", nil
	}
	return strings.Join(results, "\n"), nil
}

// Think 重写Think方法，使用系统提示
func (a *Manus) Think(ctx context.Context) (bool, error) {
	// 获取所有消息
	messages := a.Memory.GetMessages()
	if len(messages) == 0 {
		return false, fmt.Errorf("没有消息可处理")
	}

	// 向LLM发送请求，包括系统提示
	logger.Info("Manus代理向LLM发送请求...")

	// 准备系统消息
	var systemMsgs []schema.Message
	if a.SystemPrompt != "" {
		systemMsgs = []schema.Message{schema.NewSystemMessage(a.SystemPrompt)}
	}

	response, err := a.LLM.AskTool(ctx, messages, systemMsgs, a.Tools.GetToolDefinitions(), "auto")
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

// ProcessMessage 处理用户消息
func (a *Manus) ProcessMessage(ctx context.Context, message string) (string, error) {
	logger.Info("处理用户消息: %s", message)

	// 重置代理状态
	a.SetState(StateIdle)

	// 运行代理
	return a.Run(ctx, message)
}

// Step 实现Manus代理的单个步骤，重写以确保正确处理AI交互
func (a *Manus) Step(ctx context.Context) (string, error) {
	logger.Info("Manus代理执行步骤...")

	// 检查上下文是否已取消
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// 获取当前消息
	messages := a.Memory.GetMessages()
	logger.Info("当前消息数量: %d", len(messages))
	for i, msg := range messages {
		logger.Info("消息 %d - 角色: %s, 内容: %s", i, msg.Role, msg.Content)
	}

	// 思考
	logger.Info("Manus代理正在思考...")
	shouldAct, err := a.Think(ctx)
	if err != nil {
		logger.Error("思考失败: %v", err)
		return "", fmt.Errorf("思考失败: %w", err)
	}

	// 如果不需要行动，直接返回最后的助手消息
	if !shouldAct {
		// 获取最后的助手消息
		messages = a.Memory.GetMessages() // 重新获取，因为Think可能添加了新消息
		logger.Info("思考后消息数量: %d", len(messages))

		for i := len(messages) - 1; i >= 0; i-- {
			if messages[i].Role == "assistant" {
				logger.Info("Manus代理思考完成，返回最后的助手消息: %s", messages[i].Content)
				return messages[i].Content, nil
			}
		}
		logger.Info("未找到助手消息，返回默认消息")
		return "思考完成，无需行动", nil
	}

	// 行动
	logger.Info("Manus代理正在行动...")
	result, err := a.Act(ctx)
	if err != nil {
		logger.Error("行动失败: %v", err)
		return "", fmt.Errorf("行动失败: %w", err)
	}

	logger.Info("Manus代理步骤完成: %s", result)
	return result, nil
}
