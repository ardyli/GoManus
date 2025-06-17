package agent

import (
	"context"
	"fmt"
	"strings"

	"gomanus/internal/llm"
	"gomanus/internal/schema"
	"gomanus/pkg/logger"
)

// InputType 表示输入类型
type InputType string

const (
	InputTypeChat InputType = "chat"
	InputTypeTask InputType = "task"
	InputTypePlan InputType = "plan"
)

// ClassifierAgent 用于分类用户输入的代理
type ClassifierAgent struct {
	*BaseAgent
	SystemPrompt string
}

// NewClassifierAgent 创建新的分类代理
func NewClassifierAgent(name string, llm *llm.LLM) *ClassifierAgent {
	baseAgent := &BaseAgent{
		Name:        name,
		Description: "输入分类代理 - 判断用户输入是聊天、任务还是计划",
		State:       StateIdle,
		LLM:         llm,
		Memory:      schema.NewMemory(),
		MaxSteps:    1, // 分类只需要一步
		CurrentStep: 0,
	}

	systemPrompt := `你是一个智能输入分类器，需要判断用户输入属于以下哪种类型：

1. **chat（聊天）**：
   - 日常对话、问候、闲聊
   - 询问信息、知识问答
   - 情感表达、观点讨论
   - 不需要执行具体任务的交流
   - 例如："你好"、"今天天气怎么样？"、"什么是人工智能？"

2. **task（任务）**：
   - 需要执行具体操作的请求
   - 文件操作、搜索、计算等
   - 明确的行动指令
   - 例如："帮我搜索关于机器学习的资料"、"保存这个文件"、"计算一下这个数据"

3. **plan（计划）**：
   - 复杂的多步骤任务
   - 需要制定详细计划的项目
   - 包含"计划"、"规划"、"方案"等关键词
   - 例如："制定一个学习计划"、"规划项目开发流程"、"plan:制定营销策略"

请仔细分析用户输入，只返回以下三个词之一：chat、task、plan
不要返回任何其他内容，不要解释原因。`

	return &ClassifierAgent{
		BaseAgent:    baseAgent,
		SystemPrompt: systemPrompt,
	}
}

// ClassifyInput 分类用户输入
func (a *ClassifierAgent) ClassifyInput(ctx context.Context, input string) (InputType, error) {
	logger.Info("开始分类用户输入: %s", input)

	// 检查上下文是否已取消
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// 清空记忆，确保每次分类都是独立的
	a.Memory = schema.NewMemory()

	// 添加系统提示
	a.AddMessage(schema.NewSystemMessage(a.SystemPrompt))

	// 添加用户输入
	a.AddMessage(schema.NewUserMessage(input))

	// 向LLM发送请求
	messages := a.Memory.GetMessages()
	response, err := a.LLM.AskWithOptions(ctx, messages, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("LLM分类失败: %w", err)
	}

	// 解析响应
	classification := strings.TrimSpace(strings.ToLower(response.Content))
	logger.Info("LLM分类结果: %s", classification)

	// 验证分类结果
	switch classification {
	case "chat":
		return InputTypeChat, nil
	case "task":
		return InputTypeTask, nil
	case "plan":
		return InputTypePlan, nil
	default:
		// 如果LLM返回了无效结果，使用备用逻辑
		logger.Warn("LLM返回无效分类结果: %s，使用备用逻辑", classification)
		return a.fallbackClassify(input), nil
	}
}

// fallbackClassify 备用分类逻辑
func (a *ClassifierAgent) fallbackClassify(input string) InputType {
	inputLower := strings.ToLower(input)

	// 检查计划关键词
	planKeywords := []string{"plan:", "计划", "规划", "方案", "策略", "流程", "步骤"}
	for _, keyword := range planKeywords {
		if strings.Contains(inputLower, keyword) {
			return InputTypePlan
		}
	}

	// 检查任务关键词
	taskKeywords := []string{"帮我", "搜索", "查找", "保存", "下载", "计算", "执行", "处理", "分析"}
	for _, keyword := range taskKeywords {
		if strings.Contains(inputLower, keyword) {
			return InputTypeTask
		}
	}

	// 检查聊天关键词
	chatKeywords := []string{"你好", "hello", "hi", "什么是", "为什么", "怎么样", "如何"}
	for _, keyword := range chatKeywords {
		if strings.Contains(inputLower, keyword) {
			return InputTypeChat
		}
	}

	// 默认根据长度判断
	if len(input) < 20 {
		return InputTypeChat
	}
	return InputTypeTask
}

// Step 重写Step方法，分类器不需要循环执行
func (a *ClassifierAgent) Step(ctx context.Context) (string, error) {
	// 分类器不应该被当作普通代理使用
	// 它只用于ClassifyInput方法
	return "", fmt.Errorf("分类器不支持Step操作，请使用ClassifyInput方法")
}

// AddMessage 添加消息到记忆中
func (a *ClassifierAgent) AddMessage(message schema.Message) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Memory.AddMessage(message)
}

// SetState 设置代理状态
func (a *ClassifierAgent) SetState(state AgentState) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.State = state
}