package agent

import (
	"context"
	"fmt"
	"sync"

	"gomanus/internal/llm"
	"gomanus/internal/schema"
	"gomanus/pkg/logger"
)

// ChatAgent 专门用于聊天对话的代理，不使用工具
type ChatAgent struct {
	*BaseAgent
	SystemPrompt string
	mu           sync.Mutex
}

// NewChatAgent 创建新的聊天代理
func NewChatAgent(name string, llm *llm.LLM) *ChatAgent {
	baseAgent := &BaseAgent{
		Name:        name,
		Description: "聊天代理 - 专门用于日常对话和问答",
		State:       StateIdle,
		LLM:         llm,
		Memory:      schema.NewMemory(),
		MaxSteps:    1, // 聊天只需要一步
		CurrentStep: 0,
	}

	systemPrompt := `你是一个友好的AI助手，专门用于聊天对话。你的任务是：

1. 进行自然、友好的对话
2. 回答用户的问题
3. 提供信息和建议
4. 保持轻松愉快的交流氛围

重要提示：
- 你在聊天模式下，不需要使用任何工具
- 直接用文字回答用户的问题
- 保持对话的连贯性和友好性
- 如果用户需要执行具体任务，建议他们切换到任务模式`

	return &ChatAgent{
		BaseAgent:    baseAgent,
		SystemPrompt: systemPrompt,
	}
}

// Run 重写Run方法，专门用于聊天
func (a *ChatAgent) Run(ctx context.Context, request string) (string, error) {
	logger.Info("聊天代理开始运行...")

	// 检查上下文是否已取消
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// 清空记忆，确保每次聊天都是独立的
	a.Memory = schema.NewMemory()

	// 添加系统提示
	a.AddMessage(schema.NewSystemMessage(a.SystemPrompt))

	// 添加用户输入
	a.AddMessage(schema.NewUserMessage(request))

	// 向LLM发送请求，不使用工具
	messages := a.Memory.GetMessages()
	response, err := a.LLM.AskWithOptions(ctx, messages, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("聊天请求失败: %w", err)
	}

	logger.Debug("聊天代理完成，返回响应: %s", response.Content)
	return response.Content, nil
}

// Step 重写Step方法，聊天代理不需要循环执行
func (a *ChatAgent) Step(ctx context.Context) (string, error) {
	// 聊天代理不应该被当作普通代理使用
	// 它只用于Run方法
	return "", fmt.Errorf("聊天代理不支持Step操作，请使用Run方法")
}

// AddMessage 添加消息到记忆中
func (a *ChatAgent) AddMessage(message schema.Message) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Memory.AddMessage(message)
}

// SetState 设置代理状态
func (a *ChatAgent) SetState(state AgentState) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.State = state
}
