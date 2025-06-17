package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"gomanus/internal/llm"
	"gomanus/internal/schema"
	"gomanus/pkg/logger"
)

// AgentState 表示代理的当前状态
type AgentState int

const (
	StateIdle AgentState = iota
	StateRunning
	StateFinished
	StateError
)

// String 返回代理状态的字符串表示
func (s AgentState) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateRunning:
		return "running"
	case StateFinished:
		return "finished"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}

// BaseAgent 提供代理的基础功能
type BaseAgent struct {
	Name        string
	Description string
	State       AgentState
	LLM         *llm.LLM
	Memory      *schema.Memory
	MaxSteps    int
	CurrentStep int
	mu          sync.Mutex
}

// NewBaseAgent 创建新的基础代理
func NewBaseAgent(name string, llmInstance *llm.LLM) *BaseAgent {
	return &BaseAgent{
		Name:        name,
		Description: "基础代理",
		LLM:         llmInstance,
		Memory:      schema.NewMemory(),
		MaxSteps:    300,
		CurrentStep: 0,
		State:       StateIdle,
	}
}

// Stepper 定义步骤执行接口
type Stepper interface {
	Step(ctx context.Context) (string, error)
}

// Run 运行代理的主循环
func (a *BaseAgent) Run(ctx context.Context, request string) (string, error) {
	return a.RunWithStepper(ctx, request, a)
}

// RunWithStepper 使用指定的步骤执行器运行代理
func (a *BaseAgent) RunWithStepper(ctx context.Context, request string, stepper Stepper) (string, error) {
	// 检查代理状态
	a.mu.Lock()
	if a.State != StateIdle {
		a.mu.Unlock()
		return "", fmt.Errorf("无法从状态 %s 运行代理", a.State)
	}

	// 重置步骤计数并设置状态为运行中
	a.State = StateRunning
	a.CurrentStep = 0
	a.mu.Unlock()

	// 确保在函数返回时将状态重置为空闲
	defer func() {
		a.mu.Lock()
		a.State = StateIdle
		a.mu.Unlock()
	}()

	// 如果有请求，添加到记忆中
	if request != "" {
		logger.Info("添加用户请求到记忆: %s", request)
		a.AddMessage(schema.NewUserMessage(request))

		// 立即向AI咨询，生成初始步骤
		logger.Info("向AI咨询初始步骤...")
		initialStep, err := stepper.Step(ctx)
		if err != nil {
			a.mu.Lock()
			a.State = StateError
			a.mu.Unlock()
			logger.Error("初始步骤生成失败: %v", err)
			return "", fmt.Errorf("初始步骤生成失败: %w", err)
		}

		// 记录初始步骤结果
		a.mu.Lock()
		a.CurrentStep++
		stepNum := a.CurrentStep
		a.mu.Unlock()

		stepResult := fmt.Sprintf("步骤 %d: %s", stepNum, initialStep)
		logger.Info(stepResult)

		// 如果子类没有实现Step方法，这里可能已经得到了最终结果
		// 检查是否需要继续执行更多步骤
		if a.GetState() == StateFinished {
			return initialStep, nil
		}
	}

	// 执行步骤直到达到最大步骤数或代理状态变为已完成
	var results []string
	for a.CurrentStep < a.MaxSteps && a.GetState() != StateFinished {
		a.mu.Lock()
		a.CurrentStep++
		stepNum := a.CurrentStep
		a.mu.Unlock()

		logger.Info("执行步骤 %d/%d", stepNum, a.MaxSteps)

		// 执行单个步骤
		result, err := stepper.Step(ctx)
		if err != nil {
			a.mu.Lock()
			a.State = StateError
			a.mu.Unlock()
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
	if a.CurrentStep >= a.MaxSteps {
		maxStepsMsg := fmt.Sprintf("终止: 达到最大步骤数 (%d)", a.MaxSteps)
		logger.Warn(maxStepsMsg)
		results = append(results, maxStepsMsg)
	}

	// 返回所有步骤的结果
	if len(results) == 0 {
		return "未执行任何步骤", nil
	}
	return strings.Join(results, "\n"), nil
}

// Step 执行代理的单个步骤，需要被子类实现
func (a *BaseAgent) Step(ctx context.Context) (string, error) {
	// 这是一个基础实现，实际应用中应该被子类重写
	logger.Warn("BaseAgent.Step被调用，但未被子类重写")

	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		// 继续执行
	}

	// 检查是否有足够的消息可处理
	messages := a.Memory.GetMessages()
	if len(messages) == 0 {
		return "", fmt.Errorf("没有消息可处理")
	}

	// 检查是否已达到最大步骤数
	if a.CurrentStep >= a.MaxSteps {
		return "", fmt.Errorf("已达到最大步骤数 (%d)", a.MaxSteps)
	}

	// 在实际应用中，子类应该实现具体的逻辑，例如：
	// 1. 调用LLM进行思考
	// 2. 执行工具调用
	// 3. 处理工具调用结果
	// 4. 生成下一步行动

	// 子类实现应该包含以下步骤：
	// 1. 获取当前记忆中的所有消息
	// 2. 构建提示，包括系统提示、历史消息等
	// 3. 调用LLM获取响应
	// 4. 解析响应，提取工具调用或其他行动
	// 5. 执行行动并处理结果
	// 6. 将结果添加到记忆中
	// 7. 返回步骤执行的结果描述

	return "基础步骤实现 - 请在子类中重写此方法", nil
}

// AddMessage 向代理的记忆中添加一条消息
func (a *BaseAgent) AddMessage(msg schema.Message) {
	a.Memory.AddMessage(msg)
}

// GetState 获取代理的当前状态
func (a *BaseAgent) GetState() AgentState {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.State
}

// SetState 设置代理的状态
func (a *BaseAgent) SetState(state AgentState) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.State = state
}

// GetName 获取代理的名称
func (a *BaseAgent) GetName() string {
	return a.Name
}

// isStuck 检查代理是否陷入循环
func (a *BaseAgent) isStuck() bool {
	messages := a.Memory.GetMessages()
	if len(messages) < 4 {
		return false
	}

	// 获取最后一条助手消息
	var lastAssistantMsg schema.Message
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "assistant" {
			lastAssistantMsg = messages[i]
			break
		}
	}

	// 如果没有找到助手消息或内容为空，则不认为陷入循环
	if lastAssistantMsg.Role == "" || lastAssistantMsg.Content == "" {
		return false
	}

	// 统计相同内容的助手消息数量
	duplicateCount := 0
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "assistant" && messages[i].Content == lastAssistantMsg.Content {
			duplicateCount++
		}
	}

	// 设置阈值为2，如果有2个或以上相同内容的助手消息，则认为陷入循环
	duplicateThreshold := 2
	return duplicateCount >= duplicateThreshold
}

// handleStuckState 处理代理陷入循环的情况
func (a *BaseAgent) handleStuckState() {
	stuckPrompt := "检测到重复响应，任务可能已完成或遇到问题，正在终止执行。"
	logger.Warn("代理检测到循环状态，强制终止: %s", stuckPrompt)
	a.AddMessage(schema.NewSystemMessage(stuckPrompt))
	// 强制设置状态为完成，避免无限循环
	a.SetState(StateFinished)
}

// GetMaxSteps 获取代理的最大步骤数
func (a *BaseAgent) GetMaxSteps() int {
	return a.MaxSteps
}

// SetMaxSteps 设置代理的最大步骤数
func (a *BaseAgent) SetMaxSteps(maxSteps int) {
	if maxSteps > 0 {
		a.MaxSteps = maxSteps
	}
}

// GetCurrentStep 获取代理的当前步骤数
func (a *BaseAgent) GetCurrentStep() int {
	return a.CurrentStep
}

// ResetSteps 重置代理的步骤计数
func (a *BaseAgent) ResetSteps() {
	a.CurrentStep = 0
}
