package tool

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// TerminalExecutor 是一个用于执行终端命令的工具
type TerminalExecutor struct {
	*BaseTool
	parameters map[string]interface{}
}

// NewTerminalExecutor 创建新的终端命令执行工具
func NewTerminalExecutor() *TerminalExecutor {
	description := "执行操作系统终端命令并返回执行结果。支持Windows的cmd和PowerShell，以及Linux的bash终端。可以执行系统命令、文件操作、程序启动等操作。"
	baseTool := NewBaseTool("terminal_executor", description)
	
	// 定义参数
	parameters := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "(必填) 要执行的终端命令。",
			},
			"shell_type": map[string]interface{}{
				"type":        "string",
				"description": "(可选) 指定使用的shell类型。Windows支持'cmd'和'powershell'，Linux支持'bash'。默认为系统默认shell。",
				"enum":        []string{"cmd", "powershell", "bash", "auto"},
				"default":     "auto",
			},
			"working_directory": map[string]interface{}{
				"type":        "string",
				"description": "(可选) 命令执行的工作目录。默认为当前目录。",
			},
			"timeout": map[string]interface{}{
				"type":        "integer",
				"description": "(可选) 命令执行超时时间(秒)。默认为30秒。",
				"default":     30,
			},
			"capture_stderr": map[string]interface{}{
				"type":        "boolean",
				"description": "(可选) 是否捕获标准错误输出。默认为true。",
				"default":     true,
			},
		},
		"required": []string{"command"},
	}
	
	return &TerminalExecutor{
		BaseTool:   baseTool,
		parameters: parameters,
	}
}

// Parameters 返回工具参数定义
func (t *TerminalExecutor) Parameters() map[string]interface{} {
	return t.parameters
}

// Execute 执行工具
func (t *TerminalExecutor) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// 获取命令参数
	command, ok := params["command"].(string)
	if !ok || command == "" {
		return nil, fmt.Errorf("无效的命令参数")
	}
	
	// 获取shell类型参数
	shellType := "auto"
	if shellTypeParam, ok := params["shell_type"].(string); ok {
		shellType = shellTypeParam
	}
	
	// 获取工作目录参数
	workingDir := ""
	if workingDirParam, ok := params["working_directory"].(string); ok {
		workingDir = workingDirParam
	}
	
	// 获取超时参数
	timeout := 30
	if timeoutParam, ok := params["timeout"]; ok {
		if timeoutFloat, ok := timeoutParam.(float64); ok {
			timeout = int(timeoutFloat)
		}
	}
	
	// 获取是否捕获stderr参数
	captureStderr := true
	if captureStderrParam, ok := params["capture_stderr"]; ok {
		if captureStderrBool, ok := captureStderrParam.(bool); ok {
			captureStderr = captureStderrBool
		}
	}
	
	// 执行命令
	result, err := t.executeCommand(ctx, command, shellType, workingDir, timeout, captureStderr)
	if err != nil {
		return nil, fmt.Errorf("命令执行失败: %v", err)
	}
	
	return result, nil
}

// executeCommand 执行实际的命令
func (t *TerminalExecutor) executeCommand(ctx context.Context, command, shellType, workingDir string, timeout int, captureStderr bool) (map[string]interface{}, error) {
	// 创建超时上下文
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	
	// 根据操作系统和shell类型构建命令
	var cmd *exec.Cmd
	var actualShell string
	
	switch runtime.GOOS {
	case "windows":
		actualShell = t.getWindowsShell(shellType)
		switch actualShell {
		case "powershell":
			cmd = exec.CommandContext(ctxWithTimeout, "powershell", "-Command", command)
		case "cmd":
			cmd = exec.CommandContext(ctxWithTimeout, "cmd", "/C", command)
		default:
			// 默认使用cmd
			cmd = exec.CommandContext(ctxWithTimeout, "cmd", "/C", command)
			actualShell = "cmd"
		}
	case "linux", "darwin":
		actualShell = t.getUnixShell(shellType)
		switch actualShell {
		case "bash":
			cmd = exec.CommandContext(ctxWithTimeout, "bash", "-c", command)
		default:
			// 默认使用sh
			cmd = exec.CommandContext(ctxWithTimeout, "sh", "-c", command)
			actualShell = "sh"
		}
	default:
		return nil, fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
	
	// 设置工作目录
	if workingDir != "" {
		if _, err := os.Stat(workingDir); os.IsNotExist(err) {
			return nil, fmt.Errorf("工作目录不存在: %s", workingDir)
		}
		cmd.Dir = workingDir
	}
	
	// 执行命令并捕获输出
	stdout, err := cmd.Output()
	var stderr []byte
	var exitCode int
	
	if err != nil {
		// 检查是否是ExitError（命令执行了但返回非零退出码）
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
			if captureStderr {
				stderr = exitError.Stderr
			}
		} else {
			// 其他错误（如命令不存在、超时等）
			return nil, fmt.Errorf("命令执行错误: %v", err)
		}
	}
	
	// 构建返回结果
	result := map[string]interface{}{
		"command":     command,
		"shell":       actualShell,
		"os":          runtime.GOOS,
		"exit_code":   exitCode,
		"stdout":      string(stdout),
		"success":     exitCode == 0,
		"working_dir": cmd.Dir,
	}
	
	if captureStderr && len(stderr) > 0 {
		result["stderr"] = string(stderr)
	}
	
	return result, nil
}

// getWindowsShell 获取Windows下的shell类型
func (t *TerminalExecutor) getWindowsShell(shellType string) string {
	switch strings.ToLower(shellType) {
	case "powershell", "pwsh":
		return "powershell"
	case "cmd", "command":
		return "cmd"
	case "auto":
		// 检查PowerShell是否可用
		if _, err := exec.LookPath("powershell"); err == nil {
			return "powershell"
		}
		return "cmd"
	default:
		return "cmd"
	}
}

// getUnixShell 获取Unix系统下的shell类型
func (t *TerminalExecutor) getUnixShell(shellType string) string {
	switch strings.ToLower(shellType) {
	case "bash":
		return "bash"
	case "auto":
		// 检查bash是否可用
		if _, err := exec.LookPath("bash"); err == nil {
			return "bash"
		}
		return "sh"
	default:
		return "sh"
	}
}

// GetToolDefinition 返回工具定义
func (t *TerminalExecutor) GetToolDefinition() map[string]interface{} {
	return map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        t.Name(),
			"description": t.Description(),
			"parameters":  t.Parameters(),
		},
	}
}