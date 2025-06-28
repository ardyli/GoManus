# GoManus系统详细设计方案

## 1. 项目概述

### 1.1 项目背景
GoManus是一个基于Go语言开发的智能AI代理系统，旨在为用户提供强大的LLM交互能力、工具调用功能和桌面自动化操作。系统采用现代化的桌面UI设计，支持多模态输入输出，具备完整的会话管理和知识库管理功能。

### 1.2 设计目标
- 构建高性能、可扩展的AI代理系统架构
- 实现用户友好的桌面GUI界面
- 提供完整的会话管理和数据持久化能力
- 支持多种工具集成和桌面自动化操作
- 确保系统安全性和跨平台兼容性

### 1.3 系统架构

GoManus采用简化的分层架构设计，去除了内部API层，采用直接函数调用的方式，包含表示层、业务逻辑层、数据访问层和基础设施层。系统支持多种LLM提供商，集成桌面自动化功能，提供丰富的工具调用能力。架构特点：

- **直接调用架构**: UI层直接调用业务逻辑层，避免内部API的复杂性
- **本地化存储**: 使用SQLite数据库，确保数据安全和离线可用
- **外部服务集成**: 仅保留必要的外部API调用（LLM、搜索、语音等服务）
- **模块化设计**: 保持良好的模块边界，便于维护和扩展

### 1.4 技术选型
- **后端框架**: Go 
- **数据库**: SQLite + GORM
- **UI框架**: Gio UI
- **配置管理**: Viper + TOML
- **桌面自动化**: RobotGo
- **向量数据库**: SQLite + 向量扩展
- **文件监控**: fsnotify

## 2. 系统架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                    GoManus系统架构                          │
├─────────────────────────────────────────────────────────────┤
│  表示层 (Presentation Layer)                               │
│  ┌─────────────┐  ┌─────────────┐                          │
│  │  桌面GUI    │  │  系统托盘    │                          │
│  │ (Cogent Core)│  │  (Tray)     │                          │
│  └─────────────┘  └─────────────┘                          │
│                      │                                      │
│                   直接调用                                   │
│                      ↓                                      │
├─────────────────────────────────────────────────────────────┤
│  业务逻辑层 (Business Logic Layer)                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  会话管理    │  │  工具调用    │  │  知识库管理  │        │
│  │  Manager     │  │  Engine     │  │  KnowledgeDB │        │
│  │             │  │ ┌─────────┐ │  │             │        │
│  │             │  │ │桌面自动化│ │  │             │        │
│  │             │  │ │Automation│ │  │             │        │
│  │             │  │ └─────────┘ │  │             │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
│  ┌─────────────┐                    ┌─────────────┐        │
│  │  LLM交互     │                    │  多模态处理  │        │
│  │  Client      │                    │  Processor   │        │
│  └─────────────┘                    └─────────────┘        │
├─────────────────────────────────────────────────────────────┤
│  数据访问层 (Data Access Layer)                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  会话数据    │  │  配置数据    │  │  知识库数据  │        │
│  │  Repository  │  │  Repository  │  │  Repository  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
├─────────────────────────────────────────────────────────────┤
│  基础设施层 (Infrastructure Layer)                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  SQLite     │  │  文件系统    │  │  外部API     │        │
│  │  Database   │  │  FileSystem  │  │  HTTP Client │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```


### 2.2 核心模块划分

#### 2.2.1 会话管理模块 (Session Manager)
- **职责**: 管理用户会话的创建、保存、加载、删除
- **核心组件**:
  - SessionService: 会话业务逻辑
  - SessionRepository: 会话数据访问
  - MessageRepository: 消息数据访问
  - WorkspaceManager: 工作目录管理

#### 2.2.2 LLM交互模块 (LLM Client)
- **职责**: 处理与大语言模型的交互
- **核心组件**:
  - LLMClient: LLM API客户端
  - PromptBuilder: 提示词构建器
  - ResponseParser: 响应解析器
  - TokenManager: Token管理器

#### 2.2.3 工具调用模块 (Tool Engine)
- **职责**: 管理和执行各种工具调用，包括桌面自动化操作
- **核心组件**:
  - ToolRegistry: 工具注册表
  - ToolExecutor: 工具执行器
  - MCPManager: MCP服务管理器
  - ToolConfigManager: 工具配置管理器
  - DesktopAutomation: 桌面自动化工具
    - MouseController: 鼠标控制器
    - KeyboardController: 键盘控制器
    - WindowManager: 窗口管理器
    - ScreenCapture: 屏幕截图器

#### 2.2.4 知识库管理模块 (Knowledge Manager)
- **职责**: 管理本地知识库的索引和搜索
- **核心组件**:
  - IndexBuilder: 索引构建器
  - VectorStore: 向量存储
  - FileMonitor: 文件监控器
  - SearchEngine: 搜索引擎

#### 2.2.5 多模态处理模块 (Multimodal Processor)
- **职责**: 处理语音、图像、文本等多模态输入输出
- **核心组件**:
  - VoiceRecorder: 语音录制器
  - TTSEngine: 文本转语音引擎
  - MarkdownRenderer: Markdown渲染器
  - FileHandler: 文件处理器

### 2.3 模块交互流程

#### 2.3.1 用户消息处理流程
```
用户输入 → 多模态处理 → 会话管理 → LLM交互 → 工具调用 → 响应生成 → 界面显示
    ↓           ↓           ↓          ↓         ↓          ↓          ↓
  语音/文本   格式转换    消息保存   函数调用   工具执行   结果整合   UI更新
```

#### 2.3.2 知识库检索流程
```
用户查询 → 查询解析 → 向量搜索 → 结果排序 → 内容提取 → 结果返回
    ↓         ↓         ↓         ↓         ↓         ↓
  关键词    向量化    相似度计算  相关性排序  文档片段   格式化输出
```

#### 2.3.3 桌面自动化流程
```
操作指令 → 指令解析 → 权限验证 → 窗口最小化 → 操作执行 → 结果反馈 → 窗口恢复
    ↓         ↓         ↓         ↓         ↓         ↓         ↓
  AI生成    参数提取   安全检查   GoManus隐藏  RobotGo执行  状态监控  界面恢复
```

## 3. 数据库设计

### 3.1 数据库表结构

#### 3.1.1 会话相关表

**sessions表 - 会话基本信息**
```sql
CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL COMMENT '会话标题',
    theme VARCHAR(50) DEFAULT 'general' COMMENT '会话主题类型',
    created_at DATETIME NOT NULL COMMENT '创建时间',
    updated_at DATETIME NOT NULL COMMENT '最后更新时间',
    message_count INTEGER DEFAULT 0 COMMENT '消息数量',
    workspace_path VARCHAR(500) COMMENT '工作目录路径',
    is_active BOOLEAN DEFAULT FALSE COMMENT '是否为当前活跃会话',
    metadata TEXT COMMENT '会话元数据JSON'
);

CREATE INDEX idx_sessions_updated_at ON sessions(updated_at DESC);
CREATE INDEX idx_sessions_theme ON sessions(theme);
```

**messages表 - 消息详细内容**
```sql
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL COMMENT '会话ID',
    role VARCHAR(20) NOT NULL COMMENT '角色(user/assistant/system)',
    content TEXT NOT NULL COMMENT '消息内容',
    content_type VARCHAR(20) DEFAULT 'text' COMMENT '内容类型(text/image/audio)',
    timestamp DATETIME NOT NULL COMMENT '消息时间戳',
    token_count INTEGER DEFAULT 0 COMMENT 'Token数量',
    tool_calls TEXT COMMENT '工具调用信息JSON',
    metadata TEXT COMMENT '消息元数据JSON',
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

CREATE INDEX idx_messages_session_id ON messages(session_id);
CREATE INDEX idx_messages_timestamp ON messages(timestamp);
CREATE INDEX idx_messages_role ON messages(role);
```

**session_tools表 - 会话工具配置**
```sql
CREATE TABLE session_tools (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL COMMENT '会话ID',
    tool_name VARCHAR(100) NOT NULL COMMENT '工具名称',
    is_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    config_params TEXT COMMENT '工具配置参数JSON',
    created_at DATETIME NOT NULL COMMENT '创建时间',
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_session_tools_unique ON session_tools(session_id, tool_name);
```

**session_mcp_services表 - 会话MCP服务配置**
```sql
CREATE TABLE session_mcp_services (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL COMMENT '会话ID',
    service_name VARCHAR(100) NOT NULL COMMENT 'MCP服务名称',
    priority INTEGER DEFAULT 1 COMMENT '优先级(1-高,2-中,3-低)',
    config_info TEXT COMMENT '服务配置信息JSON',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    created_at DATETIME NOT NULL COMMENT '创建时间',
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

CREATE INDEX idx_session_mcp_session_id ON session_mcp_services(session_id);
```

#### 3.1.2 知识库相关表

**knowledge_directories表 - 知识库目录配置**
```sql
CREATE TABLE knowledge_directories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path VARCHAR(500) NOT NULL UNIQUE COMMENT '目录路径',
    alias VARCHAR(255) COMMENT '目录别名',
    file_types TEXT COMMENT '包含的文件类型JSON数组',
    exclude_patterns TEXT COMMENT '排除模式JSON数组',
    monitor_depth INTEGER DEFAULT 5 COMMENT '监控深度',
    is_monitoring BOOLEAN DEFAULT TRUE COMMENT '是否启用监控',
    index_status VARCHAR(20) DEFAULT 'pending' COMMENT '索引状态(pending/indexing/completed/error)',
    file_count INTEGER DEFAULT 0 COMMENT '文件数量',
    index_size INTEGER DEFAULT 0 COMMENT '索引大小(字节)',
    last_scan_at DATETIME COMMENT '最后扫描时间',
    created_at DATETIME NOT NULL COMMENT '创建时间',
    updated_at DATETIME NOT NULL COMMENT '更新时间'
);

CREATE INDEX idx_knowledge_directories_status ON knowledge_directories(index_status);
```

**knowledge_files表 - 知识库文件信息**
```sql
CREATE TABLE knowledge_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    directory_id INTEGER NOT NULL COMMENT '目录ID',
    file_path VARCHAR(1000) NOT NULL COMMENT '文件路径',
    file_name VARCHAR(255) NOT NULL COMMENT '文件名',
    file_type VARCHAR(20) NOT NULL COMMENT '文件类型',
    file_size INTEGER NOT NULL COMMENT '文件大小(字节)',
    content_hash VARCHAR(64) COMMENT '内容哈希值',
    chunk_count INTEGER DEFAULT 0 COMMENT '分块数量',
    index_status VARCHAR(20) DEFAULT 'pending' COMMENT '索引状态',
    last_modified DATETIME NOT NULL COMMENT '文件最后修改时间',
    indexed_at DATETIME COMMENT '索引时间',
    created_at DATETIME NOT NULL COMMENT '创建时间',
    FOREIGN KEY (directory_id) REFERENCES knowledge_directories(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_knowledge_files_path ON knowledge_files(file_path);
CREATE INDEX idx_knowledge_files_directory ON knowledge_files(directory_id);
CREATE INDEX idx_knowledge_files_type ON knowledge_files(file_type);
```

**knowledge_chunks表 - 知识库文档分块**
```sql
CREATE TABLE knowledge_chunks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_id INTEGER NOT NULL COMMENT '文件ID',
    chunk_index INTEGER NOT NULL COMMENT '分块索引',
    content TEXT NOT NULL COMMENT '分块内容',
    content_length INTEGER NOT NULL COMMENT '内容长度',
    start_position INTEGER COMMENT '在原文件中的起始位置',
    end_position INTEGER COMMENT '在原文件中的结束位置',
    embedding_vector BLOB COMMENT '向量嵌入数据',
    metadata TEXT COMMENT '分块元数据JSON',
    created_at DATETIME NOT NULL COMMENT '创建时间',
    FOREIGN KEY (file_id) REFERENCES knowledge_files(id) ON DELETE CASCADE
);

CREATE INDEX idx_knowledge_chunks_file ON knowledge_chunks(file_id);
CREATE INDEX idx_knowledge_chunks_index ON knowledge_chunks(chunk_index);
```

#### 3.1.3 系统配置相关表

**system_configs表 - 系统配置**
```sql
CREATE TABLE system_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    config_key VARCHAR(100) NOT NULL UNIQUE COMMENT '配置键',
    config_value TEXT COMMENT '配置值',
    config_type VARCHAR(20) DEFAULT 'string' COMMENT '配置类型(string/int/bool/json)',
    description TEXT COMMENT '配置描述',
    is_encrypted BOOLEAN DEFAULT FALSE COMMENT '是否加密存储',
    created_at DATETIME NOT NULL COMMENT '创建时间',
    updated_at DATETIME NOT NULL COMMENT '更新时间'
);

CREATE UNIQUE INDEX idx_system_configs_key ON system_configs(config_key);
```

**automation_logs表 - 桌面自动化操作日志**
```sql
CREATE TABLE automation_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER COMMENT '关联会话ID',
    operation_type VARCHAR(50) NOT NULL COMMENT '操作类型',
    target_info TEXT COMMENT '操作目标信息JSON',
    parameters TEXT COMMENT '操作参数JSON',
    execution_status VARCHAR(20) NOT NULL COMMENT '执行状态(success/failed/timeout)',
    error_message TEXT COMMENT '错误信息',
    execution_time INTEGER COMMENT '执行耗时(毫秒)',
    created_at DATETIME NOT NULL COMMENT '创建时间',
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE SET NULL
);

CREATE INDEX idx_automation_logs_session ON automation_logs(session_id);
CREATE INDEX idx_automation_logs_type ON automation_logs(operation_type);
CREATE INDEX idx_automation_logs_status ON automation_logs(execution_status);
```

### 3.2 数据库索引优化

#### 3.2.1 查询性能优化索引
```sql
-- 会话列表查询优化
CREATE INDEX idx_sessions_list_query ON sessions(updated_at DESC, is_active);

-- 消息历史查询优化
CREATE INDEX idx_messages_history ON messages(session_id, timestamp DESC);

-- 全文搜索优化
CREATE VIRTUAL TABLE messages_fts USING fts5(
    content,
    content='messages',
    content_rowid='id'
);

-- 知识库搜索优化
CREATE INDEX idx_knowledge_chunks_search ON knowledge_chunks(file_id, chunk_index);
```

#### 3.2.2 数据完整性约束
```sql
-- 会话状态约束
CREATE TRIGGER ensure_single_active_session
AFTER UPDATE OF is_active ON sessions
WHEN NEW.is_active = TRUE
BEGIN
    UPDATE sessions SET is_active = FALSE WHERE id != NEW.id AND is_active = TRUE;
END;

-- 消息计数自动更新
CREATE TRIGGER update_message_count_insert
AFTER INSERT ON messages
BEGIN
    UPDATE sessions SET 
        message_count = message_count + 1,
        updated_at = datetime('now')
    WHERE id = NEW.session_id;
END;

CREATE TRIGGER update_message_count_delete
AFTER DELETE ON messages
BEGIN
    UPDATE sessions SET 
        message_count = message_count - 1,
        updated_at = datetime('now')
    WHERE id = OLD.session_id;
END;
```

## 4. 外部API接口定义

### 4.1 LLM服务接口

**OpenAI API接口**
```http
POST https://api.openai.com/v1/chat/completions
Authorization: Bearer {API_KEY}
Content-Type: application/json

{
    "model": "gpt-4",
    "messages": [
        {
            "role": "user",
            "content": "用户消息内容"
        }
    ],
    "temperature": 0.7,
    "max_tokens": 2000,
    "tools": [
        {
            "type": "function",
            "function": {
                "name": "search_tool",
                "description": "搜索工具",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "query": {"type": "string"}
                    }
                }
            }
        }
    ]
}
```

**Anthropic Claude API接口**
```http
POST https://api.anthropic.com/v1/messages
X-API-Key: {API_KEY}
Content-Type: application/json
Anthropic-Version: 2023-06-01

{
    "model": "claude-3-sonnet-20240229",
    "max_tokens": 2000,
    "messages": [
        {
            "role": "user",
            "content": "用户消息内容"
        }
    ],
    "tools": [
        {
            "name": "search_tool",
            "description": "搜索工具",
            "input_schema": {
                "type": "object",
                "properties": {
                    "query": {"type": "string"}
                }
            }
        }
    ]
}
```

#### 4.2.2 搜索服务接口

**Google Search API**
```http
GET https://www.googleapis.com/customsearch/v1
?key={API_KEY}
&cx={SEARCH_ENGINE_ID}
&q={搜索关键词}
&num=10
&start=1
&lr=lang_zh-CN
&safe=active
```

**百度搜索API**
```http
GET https://api.baidu.com/json/sms/service/search
?access_token={ACCESS_TOKEN}
&query={搜索关键词}
&pn=0
&rn=10
&ie=utf-8
&oe=utf-8
```

#### 4.2.3 语音服务接口

**语音转文本API**
```http
POST https://api.openai.com/v1/audio/transcriptions
Authorization: Bearer {API_KEY}
Content-Type: multipart/form-data

file: {音频文件}
model: whisper-1
language: zh
response_format: json
temperature: 0
```

**文本转语音API**
```http
POST https://api.openai.com/v1/audio/speech
Authorization: Bearer {API_KEY}
Content-Type: application/json

{
    "model": "tts-1",
    "input": "要转换的文本内容",
    "voice": "alloy",
    "response_format": "mp3",
    "speed": 1.0
}
```

#### 4.2.4 向量化服务接口

**OpenAI Embeddings API**
```http
POST https://api.openai.com/v1/embeddings
Authorization: Bearer {API_KEY}
Content-Type: application/json

{
    "model": "text-embedding-3-small",
    "input": "要向量化的文本内容",
    "encoding_format": "float"
}
```


                        "function": "search_tool",
                        "arguments": {"query": "搜索关键词"}
                    }
                ],
                "created_at": "2024-01-01T12:01:00Z"
            }
        ],
        "total": 100,
        "has_more": true
    }
}
```

### 4.3 错误代码说明

#### 4.3.1 通用错误代码
- **200**: 请求成功
- **400**: 请求参数错误
- **401**: 未授权访问
- **403**: 权限不足
- **404**: 资源不存在
- **429**: 请求频率超限
- **500**: 服务器内部错误
- **503**: 服务不可用

#### 4.3.2 业务错误代码
- **1001**: 会话不存在
- **1002**: 会话已删除
- **1003**: 消息发送失败
- **1004**: LLM API调用失败
- **1005**: 工具执行失败
- **1006**: 知识库索引失败
- **1007**: 桌面自动化操作失败
- **1008**: 配置文件格式错误
- **1009**: 文件操作权限不足
- **1010**: 数据库连接失败

## 5. 功能设计

### 5.1 会话管理功能设计

#### 5.1.1 会话创建流程
1. **用户触发创建**: 点击"新建会话"按钮
2. **主题选择**: 显示主题选择对话框，包含预设主题模板
3. **工具配置**: 根据选择的主题显示推荐工具，用户可自定义选择
4. **MCP服务配置**: 配置相关的MCP服务和参数
5. **会话初始化**: 创建数据库记录，生成工作目录，保存配置
6. **界面更新**: 切换到新会话，更新UI状态

#### 5.1.2 会话数据结构
```go
type Session struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Title        string    `gorm:"size:255;not null" json:"title"`
    Theme        string    `gorm:"size:50;default:general" json:"theme"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    MessageCount int       `gorm:"default:0" json:"message_count"`
    WorkspacePath string   `gorm:"size:500" json:"workspace_path"`
    IsActive     bool      `gorm:"default:false" json:"is_active"`
    Metadata     string    `gorm:"type:text" json:"metadata"`
    
    Messages     []Message     `gorm:"foreignKey:SessionID" json:"messages,omitempty"`
    Tools        []SessionTool `gorm:"foreignKey:SessionID" json:"tools,omitempty"`
    MCPServices  []SessionMCPService `gorm:"foreignKey:SessionID" json:"mcp_services,omitempty"`
}

type Message struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    SessionID   uint      `gorm:"not null;index" json:"session_id"`
    Role        string    `gorm:"size:20;not null" json:"role"`
    Content     string    `gorm:"type:text;not null" json:"content"`
    ContentType string    `gorm:"size:20;default:text" json:"content_type"`
    Timestamp   time.Time `gorm:"not null;index" json:"timestamp"`
    TokenCount  int       `gorm:"default:0" json:"token_count"`
    ToolCalls   string    `gorm:"type:text" json:"tool_calls"`
    Metadata    string    `gorm:"type:text" json:"metadata"`
}

type SessionTool struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    SessionID    uint      `gorm:"not null;index" json:"session_id"`
    ToolName     string    `gorm:"size:100;not null" json:"tool_name"`
    IsEnabled    bool      `gorm:"default:true" json:"is_enabled"`
    ConfigParams string    `gorm:"type:text" json:"config_params"`
    CreatedAt    time.Time `json:"created_at"`
}
```

#### 5.1.3 工作目录管理
```go
type WorkspaceManager struct {
    basePath string
    db       *gorm.DB
}

func (wm *WorkspaceManager) CreateSessionWorkspace(sessionID uint) (string, error) {
    // 创建会话专用工作目录
    sessionPath := filepath.Join(wm.basePath, "sessions", fmt.Sprintf("%d", sessionID))
    
    // 创建子目录结构
    dirs := []string{
        "files",        // 用户文件
        "outputs",      // 工具输出
        "screenshots",  // 屏幕截图
        "temp",         // 临时文件
    }
    
    for _, dir := range dirs {
        if err := os.MkdirAll(filepath.Join(sessionPath, dir), 0755); err != nil {
            return "", err
        }
    }
    
    return sessionPath, nil
}

func (wm *WorkspaceManager) OpenWorkspaceInExplorer(sessionID uint) error {
    session := &Session{}
    if err := wm.db.First(session, sessionID).Error; err != nil {
        return err
    }
    
    return wm.openInSystemExplorer(session.WorkspacePath)
}

func (wm *WorkspaceManager) openInSystemExplorer(path string) error {
    var cmd *exec.Cmd
    
    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("explorer", path)
    case "darwin":
        cmd = exec.Command("open", path)
    case "linux":
        cmd = exec.Command("xdg-open", path)
    default:
        return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
    }
    
    return cmd.Start()
}
```

### 5.2 LLM交互功能设计

#### 5.2.1 LLM客户端架构
```go
type LLMClient interface {
    Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error)
    GetModels() ([]Model, error)
    ValidateConfig() error
}

type ChatRequest struct {
    SessionID    uint      `json:"session_id"`
    Messages     []Message `json:"messages"`
    Model        string    `json:"model"`
    Temperature  float64   `json:"temperature"`
    MaxTokens    int       `json:"max_tokens"`
    Tools        []Tool    `json:"tools,omitempty"`
    ToolChoice   string    `json:"tool_choice,omitempty"`
}

type ChatResponse struct {
    Message   Message `json:"message"`
    Usage     Usage   `json:"usage"`
    Model     string  `json:"model"`
    CreatedAt time.Time `json:"created_at"`
}

type OpenAIClient struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
    timeout    time.Duration
}

func (c *OpenAIClient) Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
    // 构建OpenAI API请求
    apiRequest := &openai.ChatCompletionRequest{
        Model:       request.Model,
        Messages:    convertMessages(request.Messages),
        Temperature: &request.Temperature,
        MaxTokens:   request.MaxTokens,
        Tools:       convertTools(request.Tools),
    }
    
    // 发送请求
    resp, err := c.sendRequest(ctx, apiRequest)
    if err != nil {
        return nil, err
    }
    
    // 解析响应
    return c.parseResponse(resp)
}
```

#### 5.2.2 提示词管理
```go
type PromptBuilder struct {
    systemPrompts map[string]string
    templates     map[string]*template.Template
}

func (pb *PromptBuilder) BuildSystemPrompt(sessionTheme string, tools []string) string {
    basePrompt := pb.systemPrompts["base"]
    themePrompt := pb.systemPrompts[sessionTheme]
    toolsPrompt := pb.buildToolsPrompt(tools)
    
    return fmt.Sprintf("%s\n\n%s\n\n%s", basePrompt, themePrompt, toolsPrompt)
}

func (pb *PromptBuilder) buildToolsPrompt(tools []string) string {
    if len(tools) == 0 {
        return "No tools are available for this session."
    }
    
    var builder strings.Builder
    builder.WriteString("Available tools for this session:\n")
    
    for _, tool := range tools {
        if desc, exists := pb.getToolDescription(tool); exists {
            builder.WriteString(fmt.Sprintf("- %s: %s\n", tool, desc))
        }
    }
    
    return builder.String()
}
```

### 5.3 工具调用功能设计

#### 5.3.1 工具注册系统
```go
type ToolRegistry struct {
    tools    map[string]Tool
    sessions map[uint][]string // 会话ID -> 启用的工具列表
    mu       sync.RWMutex
}

type Tool interface {
    Name() string
    Description() string
    Functions() []ToolFunction
    Execute(ctx context.Context, function string, params map[string]interface{}) (interface{}, error)
    Validate(params map[string]interface{}) error
}

type ToolFunction struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Parameters  map[string]interface{} `json:"parameters"`
}

func (tr *ToolRegistry) RegisterTool(tool Tool) error {
    tr.mu.Lock()
    defer tr.mu.Unlock()
    
    if _, exists := tr.tools[tool.Name()]; exists {
        return fmt.Errorf("tool %s already registered", tool.Name())
    }
    
    tr.tools[tool.Name()] = tool
    return nil
 }
```

#### 4.5 桌面自动化操作

```go
// 桌面自动化操作管理器
type DesktopAutomation struct {
    windowManager *WindowManager
    mouseHandler  *MouseHandler
    keyboardHandler *KeyboardHandler
    screenCapture *ScreenCapture
    fileManager   *FileManager
    operationLog  *OperationLogger
}

// 窗口操作
func (da *DesktopAutomation) MinimizeGoManus() error {
    // 最小化GoManus窗口
    return da.windowManager.MinimizeWindow("GoManus")
}

func (da *DesktopAutomation) RestoreGoManus() error {
    // 恢复GoManus窗口
    return da.windowManager.RestoreWindow("GoManus")
}

// 鼠标操作
func (da *DesktopAutomation) PerformMouseAction(action MouseAction) error {
    da.operationLog.LogOperation("mouse", action)
    
    switch action.Type {
    case "click":
        return da.mouseHandler.Click(action.X, action.Y, action.Button)
    case "double_click":
        return da.mouseHandler.DoubleClick(action.X, action.Y)
    case "drag":
        return da.mouseHandler.Drag(action.StartX, action.StartY, action.EndX, action.EndY)
    case "scroll":
        return da.mouseHandler.Scroll(action.X, action.Y, action.Direction, action.Amount)
    default:
        return fmt.Errorf("unsupported mouse action: %s", action.Type)
    }
}

// 键盘操作
func (da *DesktopAutomation) PerformKeyboardAction(action KeyboardAction) error {
    da.operationLog.LogOperation("keyboard", action)
    
    switch action.Type {
    case "type":
        return da.keyboardHandler.TypeText(action.Text)
    case "key_press":
        return da.keyboardHandler.PressKey(action.Key)
    case "key_combination":
        return da.keyboardHandler.PressKeyCombination(action.Keys)
    case "special_key":
        return da.keyboardHandler.PressSpecialKey(action.SpecialKey)
    default:
        return fmt.Errorf("unsupported keyboard action: %s", action.Type)
    }
}

// 屏幕截图
func (da *DesktopAutomation) CaptureScreen(region *ScreenRegion) (*Screenshot, error) {
    screenshot, err := da.screenCapture.Capture(region)
    if err != nil {
        return nil, err
    }
    
    da.operationLog.LogOperation("screenshot", map[string]interface{}{
        "region": region,
        "size":   screenshot.Size,
    })
    
    return screenshot, nil
}

// 文件操作
func (da *DesktopAutomation) PerformFileOperation(operation FileOperation) error {
    da.operationLog.LogOperation("file", operation)
    
    switch operation.Type {
    case "copy":
        return da.fileManager.CopyFile(operation.Source, operation.Destination)
    case "move":
        return da.fileManager.MoveFile(operation.Source, operation.Destination)
    case "delete":
        return da.fileManager.DeleteFile(operation.Path)
    case "create":
        return da.fileManager.CreateFile(operation.Path, operation.Content)
    case "read":
        content, err := da.fileManager.ReadFile(operation.Path)
        operation.Result = content
        return err
    default:
        return fmt.Errorf("unsupported file operation: %s", operation.Type)
    }
}
```

#### 4.6 动画系统

```go
// 动画系统
type AnimationSystem struct {
    animations map[string]*Animation
    renderer   *AnimationRenderer
    timeline   *AnimationTimeline
}

// 动画效果
func (as *AnimationSystem) PlayAnimation(name string, target interface{}, duration time.Duration) error {
    animation := &Animation{
        Name:     name,
        Target:   target,
        Duration: duration,
        StartTime: time.Now(),
        Status:   "playing",
    }
    
    as.animations[name] = animation
    
    switch name {
    case "popup":
        return as.playPopupAnimation(animation)
    case "scale":
        return as.playScaleAnimation(animation)
    case "bubble":
        return as.playBubbleAnimation(animation)
    case "transition":
        return as.playTransitionAnimation(animation)
    case "feedback":
        return as.playFeedbackAnimation(animation)
    default:
        return fmt.Errorf("unknown animation: %s", name)
    }
}

func (as *AnimationSystem) playPopupAnimation(animation *Animation) error {
    // 弹出动画实现
    keyframes := []Keyframe{
        {Time: 0.0, Scale: 0.0, Opacity: 0.0},
        {Time: 0.3, Scale: 1.1, Opacity: 0.8},
        {Time: 1.0, Scale: 1.0, Opacity: 1.0},
    }
    
    return as.renderer.RenderKeyframes(animation, keyframes)
}
```

## 5. UI设计

### 5.1 主窗口设计

#### 5.1.1 窗口属性
- **尺寸**: 默认 800x600 像素，最小 600x400 像素
- **背景色**: `rgba(248, 249, 250, 0.95)` (半透明白色)
- **边框**: 无边框设计，圆角半径 12px
- **阴影**: `box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12)`
- **毛玻璃效果**: `backdrop-filter: blur(20px)`

#### 5.1.2 布局结构
```
┌─────────────────────────────────────┐
│ 标题栏 (高度: 40px)                    │
├─────────────────────────────────────┤
│                                     │
│ 聊天对话区域                          │
│ (flex: 1, 自适应高度)                 │
│                                     │
├─────────────────────────────────────┤
│ 输入区域 (高度: 120px)                │
├─────────────────────────────────────┤
│ 状态栏 (高度: 24px)                   │
└─────────────────────────────────────┘
```

### 5.2 颜色主题

#### 5.2.1 浅色主题
- **主色调**: `#007AFF` (系统蓝)
- **背景色**: `#FFFFFF`
- **次要背景**: `#F8F9FA`
- **文字颜色**: `#1D1D1F`
- **次要文字**: `#6E6E73`
- **边框颜色**: `#E5E5E7`
- **阴影颜色**: `rgba(0, 0, 0, 0.1)`

#### 5.2.2 深色主题
- **主色调**: `#0A84FF` (深色模式蓝)
- **背景色**: `#1C1C1E`
- **次要背景**: `#2C2C2E`
- **文字颜色**: `#FFFFFF`
- **次要文字**: `#8E8E93`
- **边框颜色**: `#38383A`
- **阴影颜色**: `rgba(0, 0, 0, 0.3)`

### 5.3 组件设计

#### 5.3.1 消息气泡
```css
/* 用户消息 */
.user-message {
    background: linear-gradient(135deg, #007AFF 0%, #5856D6 100%);
    color: #FFFFFF;
    border-radius: 18px 18px 4px 18px;
    padding: 12px 16px;
    margin: 8px 0;
    max-width: 70%;
    align-self: flex-end;
    box-shadow: 0 2px 8px rgba(0, 122, 255, 0.3);
}

/* AI回复 */
.ai-message {
    background: #F2F2F7;
    color: #1D1D1F;
    border-radius: 18px 18px 18px 4px;
    padding: 12px 16px;
    margin: 8px 0;
    max-width: 70%;
    align-self: flex-start;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}
```

#### 5.3.2 输入框设计
```css
.input-container {
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(0, 0, 0, 0.1);
    border-radius: 24px;
    padding: 12px 20px;
    margin: 16px;
    transition: all 0.3s ease;
}

.input-container:focus-within {
    border-color: #007AFF;
    box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.1);
}
```

## 6. 技术选型说明

### 6.1 技术栈

#### 6.1.1 核心框架
- **Go语言**: v1.21+

#### 6.1.2 数据库
- **SQLite**: v3.45.0+
  - 轻量级嵌入式数据库，存储配置、会话数据
  - 无需服务器，文件型数据库
- **GORM**: v1.25.0
  - ORM框架，简化数据库操作
  - 支持自动迁移和关联查询

#### 6.1.3 向量数据库
- **sqlite-vec**: v0.1.0+
  - SQLite向量扩展，支持向量存储和相似性搜索
  - 与SQLite完美集成，统一数据存储
  - 支持高维向量索引和快速检索

### 6.2 前端技术栈

#### 6.2.1 桌面UI框架
- **Cogent Core**:v0.3.11
  - Go原生UI框架
  - 跨平台支持
  - GPU加速渲染
  - 声明式UI设计


### 6.3 中间件和工具

#### 6.3.1 日志系统
- **Logrus**: v1.9.0
  - 结构化日志
  - 多级别日志
  - 多输出格式

#### 6.3.2 配置管理
- **Viper**: v1.16.0
  - 多格式配置文件支持
  - 环境变量集成
  - 配置热重载

#### 6.3.3 任务队列
- **Asynq**: v0.24.0
  - Redis-based任务队列
  - 任务重试机制
  - 任务调度支持

### 6.4 外部服务集成

#### 6.4.1 LLM服务
- **OpenAI兼容API**: GPT-4, GPT-3.5-turbo
- **本地模型**: Ollama集成

#### 6.4.2 语音服务
- **OpenAI Whisper**: 语音转文字
- **Azure Speech**: 文字转语音
- **本地语音**: 离线语音处理

#### 6.4.3 搜索服务
- **Google Search API**: 网络搜索
- **Bing Search API**: 备用搜索
- **DuckDuckGo**: 隐私搜索

### 6.5 开发工具

#### 6.5.1 构建工具
- **Go Modules**: 依赖管理
- **Make**: 构建脚本,分别提供Linux、MacOS、Windows的构建脚本

#### 6.5.2 代码质量
- **golangci-lint**: 代码检查
- **gofmt**: 代码格式化
- **go test**: 单元测试


## 7. 部署架构

### 7.1 本地化部署架构

除大模型API外，所有组件均为本地部署，确保数据安全和离线可用性。

```
┌─────────────────────────────────────────────────────────────┐
│                    GoManus 桌面应用                        │
├─────────────────────────────────────────────────────────────┤
│ 用户界面层                                                  │
│ • Gio UI 界面 (原生桌面界面)                                │
│ • 主题管理 (浅色/深色主题)                                  │
│ • 多语言支持 (中文/英文)                                    │
├─────────────────────────────────────────────────────────────┤
│ 业务逻辑层                                                  │
│ • 会话管理 (多会话支持)                                     │
│ • 消息处理 (多模态消息)                                     │
│ • 工具调用 (内置工具集成)                                   │
│ • 桌面自动化 (鼠标键盘控制)                                 │
│ • 知识库管理 (文档索引和搜索)                               │
├─────────────────────────────────────────────────────────────┤
│ 数据存储层 (本地)                                           │
│ • SQLite 数据库 (关系数据存储)                              │
│ • sqlite-vec 向量存储 (语义搜索)                            │
│ • 文件系统存储 (附件和缓存)                                 │
│ • 本地配置文件 (用户设置)                                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      外部API服务                           │
├─────────────────────────────────────────────────────────────┤
│ 大模型API (远程调用)                                        │
│ • OpenAI 兼容API (GPT-4, Claude等)                         │
│ • Ollama 本地API (可选本地大模型)                           │
│ • 自定义API端点 (企业内部模型)                              │
├─────────────────────────────────────────────────────────────┤
│ 搜索API (可选)                                              │
│ • 网络搜索API (Google, Bing等)                             │
│ • 知识库API (Wikipedia等)                                   │
└─────────────────────────────────────────────────────────────┘
```

### 7.2 部署特点

#### 本地化优势
- **数据安全**: 所有用户数据和会话记录存储在本地
- **离线可用**: 除LLM调用外，所有功能均可离线使用
- **快速响应**: 本地数据库和文件系统，响应速度快
- **隐私保护**: 无需上传敏感文件到云端
- **简化部署**: 单一可执行文件，无需复杂配置

#### 外部依赖
- **大模型API**: 唯一的外部依赖，支持多种API提供商
- **网络搜索**: 可选功能，用于获取实时信息
- **在线更新**: 软件更新和工具库更新

## 8. 安全设计

### 8.1 数据安全
- **API密钥加密**: 不需要
- **本地数据加密**: 不需要
- **传输加密**: HTTPS/TLS 1.3
- **访问控制**: 基于会话的权限管理

### 8.2 隐私保护
- **本地优先**: 敏感数据本地存储
- **数据最小化**: 只收集必要数据
- **数据控制**: 本地数据完全可控



---

*本文档版本: v1.0*  
*最后更新: 2024年12月*  
*文档状态: 详细设计阶段*

func (tr *ToolRegistry) GetSessionTools(sessionID uint) []Tool {
    tr.mu.RLock()
    defer tr.mu.RUnlock()
    
    toolNames, exists := tr.sessions[sessionID]
    if !exists {
        return nil
    }
    
    var tools []Tool
    for _, name := range toolNames {
        if tool, exists := tr.tools[name]; exists {
            tools = append(tools, tool)
        }
    }
    
    return tools
}
```

#### 5.3.2 文件操作工具实现
```go
type FileOperationsTool struct {
    workspaceManager *WorkspaceManager
    allowedPaths     []string
}

func (f *FileOperationsTool) Name() string {
    return "file_operations"
}

func (f *FileOperationsTool) Functions() []ToolFunction {
    return []ToolFunction{
        {
            Name:        "read_file",
            Description: "读取文件内容",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "file_path": map[string]interface{}{
                        "type":        "string",
                        "description": "文件路径",
                    },
                    "encoding": map[string]interface{}{
                        "type":        "string",
                        "description": "文件编码",
                        "default":     "utf-8",
                    },
                },
                "required": []string{"file_path"},
            },
        },
        {
            Name:        "write_file",
            Description: "写入文件内容",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "file_path": map[string]interface{}{
                        "type":        "string",
                        "description": "文件路径",
                    },
                    "content": map[string]interface{}{
                        "type":        "string",
                        "description": "文件内容",
                    },
                    "encoding": map[string]interface{}{
                        "type":        "string",
                        "description": "文件编码",
                        "default":     "utf-8",
                    },
                },
                "required": []string{"file_path", "content"},
            },
        },
    }
}

func (f *FileOperationsTool) Execute(ctx context.Context, function string, params map[string]interface{}) (interface{}, error) {
    switch function {
    case "read_file":
        return f.readFile(params)
    case "write_file":
        return f.writeFile(params)
    default:
        return nil, fmt.Errorf("unknown function: %s", function)
    }
}

func (f *FileOperationsTool) readFile(params map[string]interface{}) (interface{}, error) {
    filePath, ok := params["file_path"].(string)
    if !ok {
        return nil, fmt.Errorf("file_path parameter is required")
    }
    
    // 安全检查：验证文件路径
    if !f.isPathAllowed(filePath) {
        return nil, fmt.Errorf("access to path %s is not allowed", filePath)
    }
    
    encoding, _ := params["encoding"].(string)
    if encoding == "" {
        encoding = "utf-8"
    }
    
    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %w", err)
    }
    
    stat, _ := os.Stat(filePath)
    
    return map[string]interface{}{
        "content":       string(content),
        "file_size":     len(content),
        "encoding":      encoding,
        "last_modified": stat.ModTime(),
    }, nil
}
```

### 5.4 知识库管理功能设计

#### 5.4.1 索引构建器
```go
type IndexBuilder struct {
    vectorStore   VectorStore
    textProcessor TextProcessor
    fileMonitor   *FileMonitor
    db           *gorm.DB
    concurrency  int
}

func (ib *IndexBuilder) BuildIndex(directoryID uint) error {
    directory := &KnowledgeDirectory{}
    if err := ib.db.First(directory, directoryID).Error; err != nil {
        return err
    }
    
    // 更新状态为索引中
    directory.IndexStatus = "indexing"
    ib.db.Save(directory)
    
    // 扫描文件
    files, err := ib.scanDirectory(directory)
    if err != nil {
        directory.IndexStatus = "error"
        ib.db.Save(directory)
        return err
    }
    
    // 并发处理文件
    return ib.processFilesParallel(directory, files)
}

func (ib *IndexBuilder) processFilesParallel(directory *KnowledgeDirectory, files []string) error {
    semaphore := make(chan struct{}, ib.concurrency)
    var wg sync.WaitGroup
    var mu sync.Mutex
    var errors []error
    
    for _, filePath := range files {
        wg.Add(1)
        go func(path string) {
            defer wg.Done()
            semaphore <- struct{}{} // 获取信号量
            defer func() { <-semaphore }()
            
            if err := ib.processFile(directory, path); err != nil {
                mu.Lock()
                errors = append(errors, err)
                mu.Unlock()
            }
        }(filePath)
    }
    
    wg.Wait()
    
    if len(errors) > 0 {
        directory.IndexStatus = "error"
        ib.db.Save(directory)
        return fmt.Errorf("failed to process %d files", len(errors))
    }
    
    directory.IndexStatus = "completed"
    directory.LastScanAt = time.Now()
    ib.db.Save(directory)
    
    return nil
}

func (ib *IndexBuilder) processFile(directory *KnowledgeDirectory, filePath string) error {
    // 检查文件是否已存在
    existingFile := &KnowledgeFile{}
    err := ib.db.Where("file_path = ?", filePath).First(existingFile).Error
    
    stat, err := os.Stat(filePath)
    if err != nil {
        return err
    }
    
    // 计算文件哈希
    hash, err := ib.calculateFileHash(filePath)
    if err != nil {
        return err
    }
    
    // 如果文件未变化，跳过处理
    if err == nil && existingFile.ContentHash == hash {
        return nil
    }
    
    // 提取文件内容
    content, err := ib.extractFileContent(filePath)
    if err != nil {
        return err
    }
    
    // 创建或更新文件记录
    file := &KnowledgeFile{
        DirectoryID:  directory.ID,
        FilePath:     filePath,
        FileName:     filepath.Base(filePath),
        FileType:     filepath.Ext(filePath),
        FileSize:     int(stat.Size()),
        ContentHash:  hash,
        LastModified: stat.ModTime(),
        IndexStatus:  "indexing",
    }
    
    if err == nil {
        file.ID = existingFile.ID
        ib.db.Save(file)
    } else {
        ib.db.Create(file)
    }
    
    // 分块处理
    chunks := ib.textProcessor.SplitIntoChunks(content, 1000)
    
    // 删除旧的分块
    ib.db.Where("file_id = ?", file.ID).Delete(&KnowledgeChunk{})
    
    // 处理新分块
    for i, chunk := range chunks {
        embedding, err := ib.vectorStore.GenerateEmbedding(chunk)
        if err != nil {
            continue
        }
        
        chunkRecord := &KnowledgeChunk{
            FileID:        file.ID,
            ChunkIndex:    i,
            Content:       chunk,
            ContentLength: len(chunk),
            StartPosition: i * 1000,
            EndPosition:   min((i+1)*1000, len(content)),
            EmbeddingVector: embedding,
        }
        
        ib.db.Create(chunkRecord)
    }
    
    file.ChunkCount = len(chunks)
    file.IndexStatus = "completed"
    file.IndexedAt = time.Now()
    ib.db.Save(file)
    
    return nil
}