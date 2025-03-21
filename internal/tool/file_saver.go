package tool

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gomanus/internal/llm"
	"gomanus/internal/schema"
	"gomanus/pkg/logger"
)

// FileOperator 是一个用于文件操作的工具，支持读取和保存
type FileOperator struct {
	*BaseTool
	parameters map[string]interface{}
}

// NewFileOperator 创建新的文件操作工具
func NewFileOperator() *FileOperator {
	description := "对文件进行读取和保存操作。可以读取txt、md、pdf、png、jpg等格式文件，也可以将内容保存到指定路径的本地文件。"
	baseTool := NewBaseTool("file_operator", description)
	
	// 定义参数
	parameters := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"operation": map[string]interface{}{
				"type":        "string",
				"description": "(必填) 操作类型。'read'表示读取文件，'write'表示写入文件。",
				"enum":        []string{"read", "write"},
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "(写入操作必填) 要保存到文件的内容。",
			},
			"file_path": map[string]interface{}{
				"type":        "string",
				"description": "(必填) 文件的路径，包括文件名和扩展名。",
			},
			"mode": map[string]interface{}{
				"type":        "string",
				"description": "(写入操作可选) 文件打开模式。默认为'w'表示写入。使用'a'表示追加。",
				"enum":        []string{"w", "a"},
				"default":     "w",
			},
			"encoding": map[string]interface{}{
				"type":        "string",
				"description": "(读取操作可选) 文件编码格式，用于文本文件。默认为'utf-8'。",
				"default":     "utf-8",
			},
			"max_size": map[string]interface{}{
				"type":        "integer",
				"description": "(读取操作可选) 读取文件的最大大小(字节)，用于限制大文件读取。默认为10MB。",
				"default":     10485760, // 10MB
			},
		},
		"required": []string{"operation", "file_path"},
	}
	
	return &FileOperator{
		BaseTool:   baseTool,
		parameters: parameters,
	}
}

// Parameters 返回工具参数定义
func (f *FileOperator) Parameters() map[string]interface{} {
	return f.parameters
}

// Execute 执行工具
func (f *FileOperator) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// 获取操作类型
	operation, ok := params["operation"].(string)
	if !ok || (operation != "read" && operation != "write") {
		return nil, fmt.Errorf("无效的操作类型参数，必须是'read'或'write'")
	}
	
	// 获取文件路径参数
	filePath, ok := params["file_path"].(string)
	if !ok || filePath == "" {
		return nil, fmt.Errorf("无效的文件路径参数")
	}
	
	// 根据操作类型执行不同的操作
	if operation == "read" {
		return f.readFile(ctx, filePath, params)
	} else {
		return f.writeFile(filePath, params)
	}
}

// readFile 读取文件内容
func (f *FileOperator) readFile(ctx context.Context, filePath string, params map[string]interface{}) (interface{}, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("文件不存在: %s", filePath)
	}
	
	// 获取最大读取大小
	maxSize := int64(10485760) // 默认10MB
	if maxSizeParam, ok := params["max_size"].(float64); ok {
		maxSize = int64(maxSizeParam)
	}
	
	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}
	
	// 检查文件大小
	if fileInfo.Size() > maxSize {
		return nil, fmt.Errorf("文件过大，超过最大读取限制 %d 字节", maxSize)
	}
	
	// 根据文件扩展名决定读取方式
	ext := strings.ToLower(filepath.Ext(filePath))
	
	// 文本文件处理 (txt, md)
	if ext == ".txt" || ext == ".md" {
		// 读取文本文件
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("读取文件失败: %v", err)
		}
		return string(content), nil
	}
	
	// 图像文件处理 (png, jpg等)
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		// 读取图像文件
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("打开文件失败: %v", err)
		}
		defer file.Close()
		
		// 读取文件内容
		content := make([]byte, fileInfo.Size())
		_, err = io.ReadFull(file, content)
		if err != nil {
			return nil, fmt.Errorf("读取文件内容失败: %v", err)
		}
		
		// 使用视觉模型处理图像
		description, err := f.processImageWithVisionModel(ctx, content, ext)
		if err != nil {
			logger.Warn("使用视觉模型处理图像失败: %v，返回基本信息", err)
			return fmt.Sprintf("成功读取%s文件，大小: %d 字节", ext[1:], fileInfo.Size()), nil
		}
		
		return description, nil
	}
	
	// PDF文件处理
	if ext == ".pdf" {
		// 读取PDF文件
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("打开文件失败: %v", err)
		}
		defer file.Close()
		
		// 读取文件内容
		content := make([]byte, fileInfo.Size())
		_, err = io.ReadFull(file, content)
		if err != nil {
			return nil, fmt.Errorf("读取文件内容失败: %v", err)
		}
		
		// 返回文件类型和大小信息
		return fmt.Sprintf("成功读取PDF文件，大小: %d 字节", fileInfo.Size()), nil
	}
	
	// 其他未明确支持的文件类型
	return nil, fmt.Errorf("不支持的文件类型: %s", ext)
}

// processImageWithVisionModel 使用视觉模型处理图像
func (f *FileOperator) processImageWithVisionModel(ctx context.Context, imageData []byte, fileExt string) (string, error) {
	logger.Info("使用视觉模型处理图像...")
	
	// 获取视觉模型配置
	visionLLM, err := llm.NewLLM("vision")
	if err != nil {
		return "", fmt.Errorf("创建视觉模型实例失败: %v", err)
	}
	
	// 将图像编码为base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	
	// 创建包含图像的消息
	imageMessage := schema.Message{
		Role: "user",
		Content: "",
	}
	
	// 创建图像内容
	contentParts := []map[string]interface{}{
		{
			"type": "text",
			"text": "这是一张图片，请详细描述图片中的内容。",
		},
		{
			"type": "image_url",
			"image_url": map[string]interface{}{
				"url": fmt.Sprintf("data:image/%s;base64,%s", strings.TrimPrefix(fileExt, "."), base64Image),
			},
		},
	}
	
	// 将内容部分序列化为JSON
	contentJSON, err := json.Marshal(contentParts)
	if err != nil {
		return "", fmt.Errorf("序列化图像内容失败: %v", err)
	}
	
	// 将JSON字符串设置为消息内容
	imageMessage.Content = string(contentJSON)
	
	// 创建系统消息
	systemMessage := schema.NewSystemMessage("你是一个图像分析助手，请详细描述图片中的内容，包括主体、背景、颜色、动作等细节。")
	
	// 发送请求到视觉模型
	response, err := visionLLM.AskWithOptions(ctx, []schema.Message{imageMessage}, []schema.Message{systemMessage}, nil, nil)
	if err != nil {
		return "", fmt.Errorf("发送请求到视觉模型失败: %v", err)
	}
	
	// 返回视觉模型的描述
	return response.Content, nil
}

// writeFile 写入文件内容
func (f *FileOperator) writeFile(filePath string, params map[string]interface{}) (interface{}, error) {
	// 获取内容参数
	content, ok := params["content"].(string)
	if !ok {
		return nil, fmt.Errorf("无效的内容参数")
	}
	
	// 获取模式参数
	mode := "w"
	if modeParam, ok := params["mode"].(string); ok && (modeParam == "w" || modeParam == "a") {
		mode = modeParam
	}
	
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建目录失败: %v", err)
		}
	}
	
	// 打开文件
	var flag int
	if mode == "a" {
		flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	} else {
		flag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}
	
	file, err := os.OpenFile(filePath, flag, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()
	
	// 写入内容
	if _, err := file.WriteString(content); err != nil {
		return nil, fmt.Errorf("写入文件失败: %v", err)
	}
	
	return fmt.Sprintf("内容已成功保存到 %s", filePath), nil
}

// GetToolDefinition 返回工具定义
func (f *FileOperator) GetToolDefinition() map[string]interface{} {
	return map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        f.Name(),
			"description": f.Description(),
			"parameters":  f.Parameters(),
		},
	}
}
