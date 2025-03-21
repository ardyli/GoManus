package schema

import (
	"time"
)

// Message 表示代理与LLM之间交换的消息
type Message struct {
	Role        string     `json:"role"`         // 消息角色: system, user, assistant, tool
	Content     string     `json:"content"`      // 消息内容
	Name        string     `json:"name,omitempty"`        // 工具名称，仅用于tool角色
	ToolCallID  string     `json:"tool_call_id,omitempty"` // 工具调用ID，仅用于tool角色
	ToolCalls   []ToolCall `json:"tool_calls,omitempty"`   // 工具调用，仅用于assistant角色
	Timestamp   time.Time  `json:"timestamp"`    // 消息时间戳
}

// NewUserMessage 创建新的用户消息
func NewUserMessage(content string) Message {
	return Message{
		Role:      "user",
		Content:   content,
		Timestamp: time.Now(),
	}
}

// NewAssistantMessage 创建新的助手消息
func NewAssistantMessage(content string) Message {
	return Message{
		Role:      "assistant",
		Content:   content,
		Timestamp: time.Now(),
	}
}

// NewSystemMessage 创建新的系统消息
func NewSystemMessage(content string) Message {
	return Message{
		Role:      "system",
		Content:   content,
		Timestamp: time.Now(),
	}
}

// NewToolMessage 创建新的工具消息
func NewToolMessage(content, toolCallID, name string) Message {
	return Message{
		Role:       "tool",
		Content:    content,
		ToolCallID: toolCallID,
		Name:       name,
		Timestamp:  time.Now(),
	}
}

// Memory 表示代理的记忆，存储消息历史
type Memory struct {
	Messages []Message
}

// NewMemory 创建新的记忆
func NewMemory() *Memory {
	return &Memory{
		Messages: make([]Message, 0),
	}
}

// AddMessage 向记忆中添加消息
func (m *Memory) AddMessage(msg Message) {
	m.Messages = append(m.Messages, msg)
}

// GetMessages 获取所有消息
func (m *Memory) GetMessages() []Message {
	return m.Messages
}

// GetLastNMessages 获取最后N条消息
func (m *Memory) GetLastNMessages(n int) []Message {
	if len(m.Messages) <= n {
		return m.Messages
	}
	return m.Messages[len(m.Messages)-n:]
}

// Clear 清空记忆
func (m *Memory) Clear() {
	m.Messages = make([]Message, 0)
}
