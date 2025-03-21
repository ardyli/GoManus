package agent

import (
	"context"
	"fmt"

	"gomanus/internal/llm"
	"gomanus/internal/schema"
	"gomanus/pkg/logger"
)

// ReActAgent 实现了思考-行动循环的代理
type ReActAgent struct {
	*BaseAgent
}

// NewReActAgent 创建新的ReAct代理
func NewReActAgent(name string, llm *llm.LLM) *ReActAgent {
	baseAgent := NewBaseAgent(name, llm)
	baseAgent.Description = "ReAct代理 - 实现思考-行动循环"
	
	return &ReActAgent{
		BaseAgent: baseAgent,
	}
}

// Step 实现代理的单个步骤，包括思考和行动
func (a *ReActAgent) Step(ctx context.Context) (string, error) {
	// 思考
	logger.Info("代理正在思考...")
	shouldAct, err := a.Think(ctx)
	if err != nil {
		return "", fmt.Errorf("思考失败: %w", err)
	}
	
	// 如果不需要行动，直接返回
	if !shouldAct {
		return "思考完成，无需行动", nil
	}
	
	// 行动
	logger.Info("代理正在行动...")
	result, err := a.Act(ctx)
	if err != nil {
		return "", fmt.Errorf("行动失败: %w", err)
	}
	
	return result, nil
}

// Think 思考下一步行动
func (a *ReActAgent) Think(ctx context.Context) (bool, error) {
	// 这是一个抽象方法，需要被子类实现
	return false, fmt.Errorf("Think方法需要被子类实现")
}

// Act 执行行动
func (a *ReActAgent) Act(ctx context.Context) (string, error) {
	// 这是一个抽象方法，需要被子类实现
	return "", fmt.Errorf("Act方法需要被子类实现")
}

// IsFinished 检查代理是否完成任务
func (a *ReActAgent) IsFinished() bool {
	return a.GetState() == StateFinished
}

// Finish 标记代理已完成任务
func (a *ReActAgent) Finish(result string) {
	a.SetState(StateFinished)
	a.AddMessage(schema.NewAssistantMessage(result))
}
