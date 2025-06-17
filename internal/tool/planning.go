package tool

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// PlanningTool 是一个用于创建和管理任务计划的工具
type PlanningTool struct {
	*BaseTool
	parameters map[string]interface{}
	plans      map[string]map[string]interface{} // 存储计划的数据
	activePlan string                            // 当前活动计划的ID
}

// NewPlanningTool 创建新的规划工具
func NewPlanningTool() *PlanningTool {
	description := "一个规划工具，允许代理创建和管理解决复杂任务的计划。该工具提供创建计划、更新计划步骤和跟踪进度的功能。"
	baseTool := NewBaseTool("planning", description)

	// 定义参数
	parameters := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "要执行的命令。可用命令：create, update, list, get, set_active, mark_step, delete",
				"enum":        []string{"create", "update", "list", "get", "set_active", "mark_step", "delete"},
			},
			"plan_id": map[string]interface{}{
				"type":        "string",
				"description": "计划的唯一标识符。对于create, update, set_active和delete命令是必需的。对于get和mark_step是可选的（如果未指定，则使用活动计划）。",
			},
			"title": map[string]interface{}{
				"type":        "string",
				"description": "计划的标题。对于create命令是必需的，对于update命令是可选的。",
			},
			"steps": map[string]interface{}{
				"type":        "array",
				"description": "计划步骤列表。对于create命令是必需的，对于update命令是可选的。",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"step_index": map[string]interface{}{
				"type":        "integer",
				"description": "要更新的步骤索引（从0开始）。对于mark_step命令是必需的。",
			},
			"step_status": map[string]interface{}{
				"type":        "string",
				"description": "为步骤设置的状态。与mark_step命令一起使用。",
				"enum":        []string{"not_started", "in_progress", "completed", "blocked"},
			},
			"step_notes": map[string]interface{}{
				"type":        "string",
				"description": "步骤的附加说明。对于mark_step命令是可选的。",
			},
		},
		"required": []string{"command"},
	}

	return &PlanningTool{
		BaseTool:   baseTool,
		parameters: parameters,
		plans:      make(map[string]map[string]interface{}),
	}
}

// Parameters 返回工具参数定义
func (p *PlanningTool) Parameters() map[string]interface{} {
	return p.parameters
}

// Execute 执行工具
func (p *PlanningTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// 获取命令参数
	command, ok := params["command"].(string)
	if !ok || command == "" {
		return nil, fmt.Errorf("无效的命令参数")
	}

	// 根据命令执行相应的操作
	switch command {
	case "create":
		return p.createPlan(params)
	case "update":
		return p.updatePlan(params)
	case "list":
		return p.listPlans()
	case "get":
		return p.getPlan(params)
	case "set_active":
		return p.setActivePlan(params)
	case "mark_step":
		return p.markStep(params)
	case "delete":
		return p.deletePlan(params)
	default:
		return nil, fmt.Errorf("未知命令: %s", command)
	}
}

// createPlan 创建新计划
func (p *PlanningTool) createPlan(params map[string]interface{}) (interface{}, error) {
	// 获取计划ID
	planID, ok := params["plan_id"].(string)
	if !ok || planID == "" {
		// 如果未提供计划ID，生成一个基于纳秒时间戳的唯一ID
		planID = fmt.Sprintf("plan_%d", time.Now().UnixNano())
	}

	// 检查计划ID是否已存在，如果存在则生成新的ID
	for {
		if _, exists := p.plans[planID]; !exists {
			break
		}
		// 如果ID已存在，生成新的ID
		planID = fmt.Sprintf("plan_%d", time.Now().UnixNano())
	}

	// 获取标题
	title, ok := params["title"].(string)
	if !ok || title == "" {
		return nil, fmt.Errorf("创建计划需要标题")
	}

	// 获取步骤
	stepsParam, ok := params["steps"].([]interface{})
	if !ok || len(stepsParam) == 0 {
		return nil, fmt.Errorf("创建计划需要步骤列表")
	}

	// 转换步骤为字符串数组
	steps := make([]string, len(stepsParam))
	for i, step := range stepsParam {
		stepStr, ok := step.(string)
		if !ok {
			return nil, fmt.Errorf("步骤必须是字符串")
		}
		steps[i] = stepStr
	}

	// 创建计划
	plan := map[string]interface{}{
		"plan_id":       planID,
		"title":         title,
		"steps":         steps,
		"step_statuses": make([]string, len(steps)),
		"step_notes":    make([]string, len(steps)),
	}

	// 初始化所有步骤状态为"not_started"
	statuses := plan["step_statuses"].([]string)
	for i := range statuses {
		statuses[i] = "not_started"
	}

	// 初始化所有步骤说明为空字符串
	notes := plan["step_notes"].([]string)
	for i := range notes {
		notes[i] = ""
	}

	// 保存计划
	p.plans[planID] = plan
	p.activePlan = planID // 设置为活动计划

	return p.formatPlan(plan), nil
}

// updatePlan 更新现有计划
func (p *PlanningTool) updatePlan(params map[string]interface{}) (interface{}, error) {
	// 获取计划ID
	planID, ok := params["plan_id"].(string)
	if !ok || planID == "" {
		return nil, fmt.Errorf("更新计划需要计划ID")
	}

	// 检查计划是否存在
	plan, exists := p.plans[planID]
	if !exists {
		return nil, fmt.Errorf("计划ID '%s' 不存在", planID)
	}

	// 更新标题（如果提供）
	if title, ok := params["title"].(string); ok && title != "" {
		plan["title"] = title
	}

	// 更新步骤（如果提供）
	if stepsParam, ok := params["steps"].([]interface{}); ok && len(stepsParam) > 0 {
		// 转换步骤为字符串数组
		steps := make([]string, len(stepsParam))
		for i, step := range stepsParam {
			stepStr, ok := step.(string)
			if !ok {
				return nil, fmt.Errorf("步骤必须是字符串")
			}
			steps[i] = stepStr
		}

		// 保存旧步骤和状态
		oldSteps := plan["steps"].([]string)
		oldStatuses := plan["step_statuses"].([]string)
		oldNotes := plan["step_notes"].([]string)

		// 创建新状态和说明数组
		newStatuses := make([]string, len(steps))
		newNotes := make([]string, len(steps))

		// 对于相同位置的相同步骤，保留原状态和说明
		for i, step := range steps {
			if i < len(oldSteps) && step == oldSteps[i] {
				newStatuses[i] = oldStatuses[i]
				newNotes[i] = oldNotes[i]
			} else {
				newStatuses[i] = "not_started"
				newNotes[i] = ""
			}
		}

		// 更新计划
		plan["steps"] = steps
		plan["step_statuses"] = newStatuses
		plan["step_notes"] = newNotes
	}

	return p.formatPlan(plan), nil
}

// listPlans 列出所有计划
func (p *PlanningTool) listPlans() (interface{}, error) {
	if len(p.plans) == 0 {
		return "没有可用的计划", nil
	}

	result := "可用计划:\n\n"
	for id, plan := range p.plans {
		title := plan["title"].(string)
		steps := plan["steps"].([]string)
		activeMarker := ""
		if id == p.activePlan {
			activeMarker = " (当前活动)"
		}
		result += fmt.Sprintf("- %s: %s%s (%d 步骤)\n", id, title, activeMarker, len(steps))
	}

	return result, nil
}

// getPlan 获取计划详情
func (p *PlanningTool) getPlan(params map[string]interface{}) (interface{}, error) {
	// 获取计划ID（如果未提供，使用活动计划）
	planID, _ := params["plan_id"].(string)
	if planID == "" {
		if p.activePlan == "" {
			return nil, fmt.Errorf("没有活动计划，请提供计划ID")
		}
		planID = p.activePlan
	}

	// 检查计划是否存在
	plan, exists := p.plans[planID]
	if !exists {
		return nil, fmt.Errorf("计划ID '%s' 不存在", planID)
	}

	return p.formatPlan(plan), nil
}

// setActivePlan 设置活动计划
func (p *PlanningTool) setActivePlan(params map[string]interface{}) (interface{}, error) {
	// 获取计划ID
	planID, ok := params["plan_id"].(string)
	if !ok || planID == "" {
		return nil, fmt.Errorf("设置活动计划需要计划ID")
	}

	// 检查计划是否存在
	if _, exists := p.plans[planID]; !exists {
		return nil, fmt.Errorf("计划ID '%s' 不存在", planID)
	}

	// 设置活动计划
	p.activePlan = planID

	return fmt.Sprintf("已将计划 '%s' 设置为活动计划", planID), nil
}

// markStep 标记步骤状态
func (p *PlanningTool) markStep(params map[string]interface{}) (interface{}, error) {
	// 获取计划ID（如果未提供，使用活动计划）
	planID, _ := params["plan_id"].(string)
	if planID == "" {
		if p.activePlan == "" {
			return nil, fmt.Errorf("没有活动计划，请提供计划ID")
		}
		planID = p.activePlan
	}

	// 检查计划是否存在
	plan, exists := p.plans[planID]
	if !exists {
		return nil, fmt.Errorf("计划ID '%s' 不存在", planID)
	}

	// 获取步骤索引
	stepIndexFloat, ok := params["step_index"].(float64)
	if !ok {
		return nil, fmt.Errorf("标记步骤需要步骤索引")
	}
	stepIndex := int(stepIndexFloat)

	// 检查步骤索引是否有效
	steps := plan["steps"].([]string)
	if stepIndex < 0 || stepIndex >= len(steps) {
		return nil, fmt.Errorf("步骤索引 %d 超出范围 (0-%d)", stepIndex, len(steps)-1)
	}

	// 获取步骤状态
	stepStatus, ok := params["step_status"].(string)
	if !ok || stepStatus == "" {
		return nil, fmt.Errorf("标记步骤需要步骤状态")
	}

	// 验证步骤状态
	validStatuses := []string{"not_started", "in_progress", "completed", "blocked"}
	isValidStatus := false
	for _, status := range validStatuses {
		if stepStatus == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return nil, fmt.Errorf("无效的步骤状态: %s", stepStatus)
	}

	// 更新步骤状态
	statuses := plan["step_statuses"].([]string)
	statuses[stepIndex] = stepStatus

	// 更新步骤说明（如果提供）
	if stepNotes, ok := params["step_notes"].(string); ok {
		notes := plan["step_notes"].([]string)
		notes[stepIndex] = stepNotes
	}

	return p.formatPlan(plan), nil
}

// deletePlan 删除计划
func (p *PlanningTool) deletePlan(params map[string]interface{}) (interface{}, error) {
	// 获取计划ID
	planID, ok := params["plan_id"].(string)
	if !ok || planID == "" {
		return nil, fmt.Errorf("删除计划需要计划ID")
	}

	// 检查计划是否存在
	if _, exists := p.plans[planID]; !exists {
		return nil, fmt.Errorf("计划ID '%s' 不存在", planID)
	}

	// 删除计划
	delete(p.plans, planID)

	// 如果删除的是活动计划，清除活动计划
	if planID == p.activePlan {
		p.activePlan = ""
	}

	return fmt.Sprintf("已删除计划 '%s'", planID), nil
}

// formatPlan 格式化计划为可读文本
func (p *PlanningTool) formatPlan(plan map[string]interface{}) string {
	planID := plan["plan_id"].(string)
	title := plan["title"].(string)
	steps := plan["steps"].([]string)
	statuses := plan["step_statuses"].([]string)
	notes := plan["step_notes"].([]string)

	// 计算完成步骤数
	completed := 0
	for _, status := range statuses {
		if status == "completed" {
			completed++
		}
	}

	// 计算进度百分比
	progress := 0.0
	if len(steps) > 0 {
		progress = float64(completed) / float64(len(steps)) * 100
	}

	// 统计各状态步骤数
	statusCounts := map[string]int{
		"completed":   0,
		"in_progress": 0,
		"blocked":     0,
		"not_started": 0,
	}
	for _, status := range statuses {
		statusCounts[status]++
	}

	// 格式化计划文本
	result := fmt.Sprintf("计划: %s (ID: %s)\n", title, planID)
	result += strings.Repeat("=", len(result)) + "\n\n"

	result += fmt.Sprintf("进度: %d/%d 步骤已完成 (%.1f%%)\n", completed, len(steps), progress)
	result += fmt.Sprintf("状态: %d 已完成, %d 进行中, %d 已阻塞, %d 未开始\n\n",
		statusCounts["completed"], statusCounts["in_progress"],
		statusCounts["blocked"], statusCounts["not_started"])

	result += "步骤:\n"
	for i, step := range steps {
		var statusMark string
		switch statuses[i] {
		case "completed":
			statusMark = "[✓]"
		case "in_progress":
			statusMark = "[→]"
		case "blocked":
			statusMark = "[!]"
		default: // not_started
			statusMark = "[ ]"
		}

		result += fmt.Sprintf("%d. %s %s\n", i+1, statusMark, step)
		if notes[i] != "" {
			result += fmt.Sprintf("   备注: %s\n", notes[i])
		}
	}

	return result
}

// GetToolDefinition 返回工具定义
func (p *PlanningTool) GetToolDefinition() map[string]interface{} {
	return map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        p.Name(),
			"description": p.Description(),
			"parameters":  p.Parameters(),
		},
	}
}
