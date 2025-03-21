package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"gomanus/internal/llm"
	"gomanus/internal/schema"
	"gomanus/internal/tool"
	"gomanus/pkg/logger"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// PlanningAgent 是一个用于任务规划和执行的代理
type PlanningAgent struct {
	*ToolCallAgent
	PlanningTool   *tool.PlanningTool
	ExecutorAgents map[string]*ToolCallAgent
	ActivePlanID   string
	CurrentStep    int
	MaxSteps       int
}

// NewPlanningAgent 创建新的规划代理
func NewPlanningAgent(name string, llm *llm.LLM, tools *tool.ToolCollection) *PlanningAgent {
	toolCallAgent := NewToolCallAgent(name, llm, tools)
	toolCallAgent.Description = "规划代理 - 用于任务规划和执行"

	// 创建规划工具
	planningTool := tool.NewPlanningTool()
	
	// 将规划工具添加到工具集合中
	if err := tools.AddTool(planningTool); err != nil {
		logger.Error("添加规划工具失败: %v", err)
	}

	// 生成基于时间戳的计划ID
	activePlanID := fmt.Sprintf("plan_%d", time.Now().Unix())

	return &PlanningAgent{
		ToolCallAgent:  toolCallAgent,
		PlanningTool:   planningTool,
		ExecutorAgents: make(map[string]*ToolCallAgent),
		ActivePlanID:   activePlanID,
		CurrentStep:    -1, // -1表示尚未开始执行步骤
		MaxSteps:       100, // 默认最大步骤数
	}
}

// AddExecutor 添加执行器代理
func (a *PlanningAgent) AddExecutor(name string, agent *ToolCallAgent) {
	a.ExecutorAgents[name] = agent
}

// GetExecutor 获取适合当前步骤的执行器代理
func (a *PlanningAgent) GetExecutor(stepType string) *ToolCallAgent {
	// 如果提供了步骤类型并且存在对应的执行器，使用该执行器
	if stepType != "" {
		if executor, exists := a.ExecutorAgents[stepType]; exists {
			return executor
		}
	}

	// 否则使用默认执行器（自身）
	return a.ToolCallAgent
}

// Run 重写Run方法，实现规划和执行流程
func (a *PlanningAgent) Run(ctx context.Context, request string) (string, error) {
	logger.Info("规划代理开始运行...")

	// 重置代理状态
	a.SetState(StateIdle)
	a.ResetSteps()
	a.CurrentStep = -1

	// 创建初始计划
	if err := a.CreateInitialPlan(ctx, request); err != nil {
		return "", fmt.Errorf("创建初始计划失败: %w", err)
	}
	
	// 获取计划步骤数量并设置最大步骤数
	planText, err := a.GetPlanText(ctx)
	if err != nil {
		logger.Warn("获取计划步骤数量失败，使用默认值: %v", err)
	} else {
		// 解析计划文本，获取步骤数量
		steps := a.GetStepsCount(planText)
		if steps > 0 {
			logger.Info("设置最大步骤数为: %d", steps)
			a.MaxSteps = steps
		}
	}

	// 执行计划
	result, err := a.ExecutePlan(ctx)
	if err != nil {
		return "", fmt.Errorf("执行计划失败: %w", err)
	}

	return result, nil
}

// CreateInitialPlan 创建初始计划
func (a *PlanningAgent) CreateInitialPlan(ctx context.Context, request string) error {
	logger.Info("创建初始计划，ID: %s", a.ActivePlanID)

	// 创建系统消息
	systemMessage := schema.NewSystemMessage(
		"你是一个规划助手。你的任务是创建一个详细的计划，包含清晰的步骤来完成用户的请求。" +
		"每个步骤应该具体且可执行。如果步骤涉及特定类型的操作，请使用方括号标记，例如[SEARCH]表示搜索操作，" +
		"[CODE]表示编码操作。这将帮助系统选择合适的工具来执行该步骤。")

	// 创建用户消息
	userMessage := schema.NewUserMessage(
		fmt.Sprintf("为完成以下任务创建一个详细的计划：%s", request))

	// 添加消息到记忆中
	a.AddMessage(systemMessage)
	a.AddMessage(userMessage)

	// 调用LLM，要求使用规划工具
	response, err := a.LLM.AskTool(ctx, a.Memory.GetMessages(), nil, 
		[]map[string]interface{}{a.PlanningTool.GetToolDefinition()}, "required")
	if err != nil {
		return fmt.Errorf("调用LLM创建计划失败: %w", err)
	}

	// 将LLM响应添加到记忆中
	a.AddMessage(schema.Message{
		Role:      "assistant",
		Content:   response.Content,
		ToolCalls: response.ToolCalls,
	})

	// 处理工具调用
	if len(response.ToolCalls) > 0 {
		for _, toolCall := range response.ToolCalls {
			if toolCall.Function.Name == "planning" {
				// 解析参数
				var args map[string]interface{}
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
					logger.Error("解析规划工具参数失败: %v", err)
					continue
				}

				// 确保使用正确的计划ID
				args["plan_id"] = a.ActivePlanID
				args["command"] = "create" // 确保是创建命令

				// 执行规划工具
				_, err := a.PlanningTool.Execute(ctx, args)
				if err != nil {
					logger.Error("执行规划工具失败: %v", err)
					return err
				}

				logger.Info("成功创建初始计划")
				return nil
			}
		}
	}

	// 如果没有成功创建计划，创建默认计划
	logger.Warn("LLM未创建计划，创建默认计划")
	defaultArgs := map[string]interface{}{
		"command": "create",
		"plan_id": a.ActivePlanID,
		"title":   fmt.Sprintf("计划: %s", request),
		"steps":   []interface{}{"分析请求", "执行任务", "验证结果"},
	}

	_, err = a.PlanningTool.Execute(ctx, defaultArgs)
	if err != nil {
		logger.Error("创建默认计划失败: %v", err)
		return err
	}

	return nil
}

// ExecutePlan 执行计划
func (a *PlanningAgent) ExecutePlan(ctx context.Context) (string, error) {
	logger.Info("开始执行计划: %s", a.ActivePlanID)
	
	// 初始化步骤索引
	a.CurrentStep = 0
	
	// 记录结果
	var results []string

	for {
		// 获取当前步骤信息
		stepIndex, stepInfo, err := a.GetCurrentStepInfo(ctx)
		if err != nil {
			return "", fmt.Errorf("获取当前步骤信息失败: %w", err)
		}
		
		// 如果没有更多步骤，完成计划
		if stepIndex == -1 {
			finalResult, err := a.FinalizePlan(ctx)
			if err != nil {
				return "", fmt.Errorf("完成计划失败: %w", err)
			}
			results = append(results, finalResult)
			break
		}
		
		// 更新当前步骤索引
		a.CurrentStep = stepIndex
		
		// 获取步骤类型和执行器
		stepType := a.ExtractStepType(stepInfo)
		executor := a.GetExecutor(stepType)
		
		// 执行当前步骤
		stepResult, err := a.ExecuteStep(ctx, executor, stepInfo)
		if err != nil {
			return "", fmt.Errorf("执行步骤 %d 失败: %w", stepIndex, err)
		}
		
		// 记录步骤结果
		results = append(results, stepResult)
		
		// 标记步骤为已完成
		a.MarkStepCompleted(ctx)
		
		// 增加步骤索引
		a.CurrentStep++
	}
	
	// 返回所有结果
	return strings.Join(results, "\n\n"), nil
}

// GetCurrentStepInfo 获取当前步骤信息
func (a *PlanningAgent) GetCurrentStepInfo(ctx context.Context) (int, string, error) {
	// 获取计划信息
	args := map[string]interface{}{
		"command": "get",
		"plan_id": a.ActivePlanID,
	}

	result, err := a.PlanningTool.Execute(ctx, args)
	if err != nil {
		return -1, "", fmt.Errorf("获取计划失败: %w", err)
	}

	planText := fmt.Sprintf("%v", result)
	
	// 解析计划文本，查找未完成的步骤
	lines := strings.Split(planText, "\n")
	stepRegex := regexp.MustCompile(`^(\d+)\. (\[.\]) (.+)$`)
	
	for _, line := range lines {
		matches := stepRegex.FindStringSubmatch(line)
		if len(matches) == 4 {
			stepNum, _ := strconv.Atoi(matches[1])
			statusMark := matches[2]
			stepText := matches[3]
			
			// 查找未开始或进行中的步骤
			if statusMark == "[ ]" || statusMark == "[→]" {
				// 将步骤标记为进行中
				markArgs := map[string]interface{}{
					"command":     "mark_step",
					"plan_id":     a.ActivePlanID,
					"step_index":  stepNum - 1, // 步骤索引从0开始
					"step_status": "in_progress",
				}
				
				_, err := a.PlanningTool.Execute(ctx, markArgs)
				if err != nil {
					logger.Warn("标记步骤为进行中失败: %v", err)
				}
				
				return stepNum - 1, stepText, nil
			}
		}
	}
	
	// 没有找到未完成的步骤
	return -1, "", nil
}

// ExtractStepType 从步骤文本中提取步骤类型
func (a *PlanningAgent) ExtractStepType(stepText string) string {
	// 查找方括号中的类型标记，例如[SEARCH]
	typeRegex := regexp.MustCompile(`\[([A-Z_]+)\]`)
	matches := typeRegex.FindStringSubmatch(stepText)
	
	if len(matches) > 1 {
		return strings.ToLower(matches[1])
	}
	
	return ""
}

// ExecuteStep 执行单个步骤
func (a *PlanningAgent) ExecuteStep(ctx context.Context, executor *ToolCallAgent, stepText string) (string, error) {
	// 获取计划状态
	planStatus, err := a.GetPlanText(ctx)
	if err != nil {
		logger.Warn("获取计划状态失败: %v", err)
		planStatus = fmt.Sprintf("计划ID: %s", a.ActivePlanID)
	}
	
	// 创建步骤提示
	stepPrompt := fmt.Sprintf(`
当前计划状态:
%s

你的当前任务:
你正在执行步骤 %d: "%s"

请使用适当的工具执行此步骤。完成后，提供一个总结说明你完成了什么。
`, planStatus, a.CurrentStep+1, stepText)
	
	// 使用执行器代理执行步骤
	executor.ResetSteps()
	stepResult, err := executor.Run(ctx, stepPrompt)
	if err != nil {
		return "", fmt.Errorf("执行步骤失败: %w", err)
	}
	
	// 标记步骤为已完成
	a.MarkStepCompleted(ctx)
	
	return fmt.Sprintf("步骤 %d 完成: %s\n\n%s", a.CurrentStep+1, stepText, stepResult), nil
}

// MarkStepCompleted 标记当前步骤为已完成
func (a *PlanningAgent) MarkStepCompleted(ctx context.Context) {
	if a.CurrentStep < 0 {
		return
	}
	
	// 标记步骤为已完成
	args := map[string]interface{}{
		"command":     "mark_step",
		"plan_id":     a.ActivePlanID,
		"step_index":  a.CurrentStep,
		"step_status": "completed",
	}
	
	_, err := a.PlanningTool.Execute(ctx, args)
	if err != nil {
		logger.Warn("标记步骤为已完成失败: %v", err)
	}
}

// GetPlanText 获取计划文本
func (a *PlanningAgent) GetPlanText(ctx context.Context) (string, error) {
	args := map[string]interface{}{
		"command": "get",
		"plan_id": a.ActivePlanID,
	}
	
	result, err := a.PlanningTool.Execute(ctx, args)
	if err != nil {
		return "", err
	}
	
	return fmt.Sprintf("%v", result), nil
}

// FinalizePlan 完成计划
func (a *PlanningAgent) FinalizePlan(ctx context.Context) (string, error) {
	// 获取计划状态
	planText, err := a.GetPlanText(ctx)
	if err != nil {
		return "", fmt.Errorf("获取计划状态失败: %w", err)
	}
	
	// 创建总结提示
	summaryPrompt := fmt.Sprintf(`
计划已完成:
%s

请提供一个简短的总结，说明已完成的工作和结果。
`, planText)
	
	// 创建系统消息
	systemMessage := schema.NewSystemMessage(
		"你是一个总结助手。请简明扼要地总结已完成的计划和结果。")
	
	// 创建用户消息
	userMessage := schema.NewUserMessage(summaryPrompt)
	
	// 重置记忆
	a.Memory.Clear()
	a.AddMessage(systemMessage)
	a.AddMessage(userMessage)
	
	// 调用LLM生成总结
	response, err := a.LLM.AskTool(ctx, a.Memory.GetMessages(), nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("生成总结失败: %w", err)
	}
	
	// 将LLM响应添加到记忆中
	a.AddMessage(schema.Message{
		Role:    "assistant",
		Content: response.Content,
	})
	
	return fmt.Sprintf("计划完成！\n\n%s\n\n总结:\n%s", planText, response.Content), nil
}

// GetStepsCount 从计划文本中获取步骤数量
func (a *PlanningAgent) GetStepsCount(planText string) int {
	// 查找进度信息行，格式如: "进度: 0/5 步骤已完成 (0.0%)"
	progressLine := ""
	lines := strings.Split(planText, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "进度:") {
			progressLine = line
			break
		}
	}
	
	if progressLine == "" {
		logger.Warn("无法从计划文本中找到进度信息")
		return 0
	}
	
	// 解析步骤总数
	parts := strings.Split(progressLine, "/")
	if len(parts) < 2 {
		logger.Warn("进度信息格式不正确: %s", progressLine)
		return 0
	}
	
	totalPart := strings.Split(parts[1], " ")[0]
	total, err := strconv.Atoi(totalPart)
	if err != nil {
		logger.Warn("解析步骤总数失败: %v", err)
		return 0
	}
	
	return total
}

// Step 实现BaseAgent.Step接口
func (a *PlanningAgent) Step(ctx context.Context) (string, error) {
	// 检查是否取消
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// 获取当前步骤信息
	stepIndex, stepText, err := a.GetCurrentStepInfo(ctx)
	if err != nil {
		return "", fmt.Errorf("获取当前步骤信息失败: %w", err)
	}

	// 记录当前执行的步骤
	logger.Info("执行计划步骤: %s", stepText)

	// 将步骤标记为进行中
	if err := a.MarkStepStatus(ctx, stepIndex, "in_progress", ""); err != nil {
		logger.Warn("标记步骤为进行中失败: %v", err)
	}

	// 提取步骤类型
	stepType := a.ExtractStepType(stepText)
	
	// 获取适合当前步骤的执行器
	executor := a.GetExecutor(stepType)

	// 执行步骤
	result, err := a.ExecuteStep(ctx, executor, stepText)
	if err != nil {
		// 标记步骤为阻塞状态
		if markErr := a.MarkStepStatus(ctx, stepIndex, "blocked", fmt.Sprintf("错误: %v", err)); markErr != nil {
			logger.Warn("标记步骤为阻塞状态失败: %v", markErr)
		}
		return "", fmt.Errorf("执行步骤失败: %w", err)
	}

	// 标记步骤为已完成
	if err := a.MarkStepStatus(ctx, stepIndex, "completed", fmt.Sprintf("结果: %s", result)); err != nil {
		logger.Warn("标记步骤为已完成失败: %v", err)
	}

	return result, nil
}

// MarkStepStatus 标记步骤状态
func (a *PlanningAgent) MarkStepStatus(ctx context.Context, stepIndex int, status string, note string) error {
	args := map[string]interface{}{
		"command":     "mark_step",
		"plan_id":     a.ActivePlanID,
		"step_index":  stepIndex,
		"step_status": status,
	}
	
	if note != "" {
		args["step_note"] = note
	}
	
	_, err := a.PlanningTool.Execute(ctx, args)
	return err
}
