package tool

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// BrowserUseTool 是一个用于与网页浏览器交互的工具
type BrowserUseTool struct {
	*BaseTool
	parameters map[string]interface{}
	mu         sync.Mutex
	sessions   map[string]*BrowserSession
}

// BrowserSession 表示一个浏览器会话
type BrowserSession struct {
	ID      string
	URL     string
	Content string
}

// NewBrowserUseTool 创建新的浏览器使用工具
func NewBrowserUseTool() *BrowserUseTool {
	description := `与网页浏览器交互，执行各种操作如导航、获取HTML内容和执行JavaScript。支持的操作包括：
- 'navigate': 导航到特定URL
- 'get_html': 获取页面HTML内容
- 'execute_js': 执行JavaScript代码
- 'new_tab': 打开新标签页
- 'close_tab': 关闭当前标签页`

	baseTool := NewBaseTool("browser_use", description)
	
	// 定义参数
	parameters := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type": "string",
				"enum": []string{
					"navigate",
					"get_html",
					"execute_js",
					"new_tab",
					"close_tab",
				},
				"description": "要执行的浏览器操作",
			},
			"url": map[string]interface{}{
				"type":        "string",
				"description": "'navigate'或'new_tab'操作的URL",
			},
			"script": map[string]interface{}{
				"type":        "string",
				"description": "'execute_js'操作的JavaScript代码",
			},
			"tab_id": map[string]interface{}{
				"type":        "string",
				"description": "操作的标签页ID",
			},
		},
		"required": []string{"action"},
	}
	
	return &BrowserUseTool{
		BaseTool:   baseTool,
		parameters: parameters,
		sessions:   make(map[string]*BrowserSession),
	}
}

// Parameters 返回工具参数定义
func (b *BrowserUseTool) Parameters() map[string]interface{} {
	return b.parameters
}

// Execute 执行工具
func (b *BrowserUseTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// 获取操作参数
	action, ok := params["action"].(string)
	if !ok || action == "" {
		return nil, fmt.Errorf("无效的操作参数")
	}
	
	// 获取URL参数
	url, _ := params["url"].(string)
	
	// 获取脚本参数
	script, _ := params["script"].(string)
	
	// 获取标签页ID参数
	tabID, _ := params["tab_id"].(string)
	if tabID == "" {
		tabID = "default"
	}
	
	// 执行相应的操作
	b.mu.Lock()
	defer b.mu.Unlock()
	
	// 确保会话存在
	if _, exists := b.sessions[tabID]; !exists && action != "new_tab" {
		b.sessions[tabID] = &BrowserSession{
			ID:  tabID,
			URL: "",
		}
	}
	
	switch action {
	case "navigate":
		return b.navigate(tabID, url)
	case "get_html":
		return b.getHTML(tabID)
	case "execute_js":
		return b.executeJS(tabID, script)
	case "new_tab":
		return b.newTab(url)
	case "close_tab":
		return b.closeTab(tabID)
	default:
		return nil, fmt.Errorf("不支持的操作: %s", action)
	}
}

// navigate 导航到指定URL
func (b *BrowserUseTool) navigate(tabID, url string) (interface{}, error) {
	if url == "" {
		return nil, fmt.Errorf("URL不能为空")
	}
	
	// 发送HTTP请求获取页面内容
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("导航到 %s 失败: %v", url, err)
	}
	defer resp.Body.Close()
	
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取页面内容失败: %v", err)
	}
	
	// 更新会话信息
	b.sessions[tabID] = &BrowserSession{
		ID:      tabID,
		URL:     url,
		Content: string(body),
	}
	
	return fmt.Sprintf("已导航到 %s", url), nil
}

// getHTML 获取页面HTML内容
func (b *BrowserUseTool) getHTML(tabID string) (interface{}, error) {
	session, exists := b.sessions[tabID]
	if !exists || session.Content == "" {
		return nil, fmt.Errorf("标签页 %s 不存在或未加载内容", tabID)
	}
	
	// 截取前2000个字符
	content := session.Content
	if len(content) > 2000 {
		content = content[:2000] + "..."
	}
	
	return content, nil
}

// executeJS 执行JavaScript代码
func (b *BrowserUseTool) executeJS(tabID, script string) (interface{}, error) {
	if script == "" {
		return nil, fmt.Errorf("脚本不能为空")
	}
	
	session, exists := b.sessions[tabID]
	if !exists {
		return nil, fmt.Errorf("标签页 %s 不存在", tabID)
	}
	
	// 模拟JavaScript执行
	// 在实际实现中，这里应该使用WebDriver或其他浏览器自动化工具
	return fmt.Sprintf("在 %s 上执行脚本: %s", session.URL, script), nil
}

// newTab 打开新标签页
func (b *BrowserUseTool) newTab(url string) (interface{}, error) {
	if url == "" {
		return nil, fmt.Errorf("URL不能为空")
	}
	
	// 生成新的标签页ID
	tabID := fmt.Sprintf("tab_%d", len(b.sessions)+1)
	
	// 导航到URL
	_, err := b.navigate(tabID, url)
	if err != nil {
		return nil, err
	}
	
	return fmt.Sprintf("已打开新标签页 %s 并导航到 %s", tabID, url), nil
}

// closeTab 关闭标签页
func (b *BrowserUseTool) closeTab(tabID string) (interface{}, error) {
	if _, exists := b.sessions[tabID]; !exists {
		return nil, fmt.Errorf("标签页 %s 不存在", tabID)
	}
	
	// 删除会话
	delete(b.sessions, tabID)
	
	return fmt.Sprintf("已关闭标签页 %s", tabID), nil
}

// GetToolDefinition 返回工具定义
func (b *BrowserUseTool) GetToolDefinition() map[string]interface{} {
	return map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        b.Name(),
			"description": b.Description(),
			"parameters":  b.Parameters(),
		},
	}
}
