-- ============================================
-- GoManus系统数据库初始化脚本
-- 数据库类型：SQLite
-- 创建时间：2024年
-- 说明：本脚本用于初始化GoManus系统的SQLite数据库结构
-- ============================================

-- 启用外键约束
PRAGMA foreign_keys = ON;

-- 设置SQLite优化参数
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = 10000;
PRAGMA temp_store = memory;

-- ============================================
-- 会话管理相关表
-- ============================================

-- 会话表
CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 会话ID，主键
    title TEXT NOT NULL, -- 会话标题
    description TEXT, -- 会话描述
    session_type TEXT DEFAULT 'chat', -- 会话类型：chat=聊天，task=任务，automation=自动化
    status TEXT DEFAULT 'active', -- 会话状态：active=活跃，archived=归档，deleted=已删除
    workspace_path TEXT, -- 工作目录路径
    model_config TEXT, -- LLM模型配置，JSON格式
    tool_config TEXT, -- 工具配置，JSON格式
    context_length INTEGER DEFAULT 4096, -- 上下文长度限制
    message_count INTEGER DEFAULT 0, -- 消息总数
    total_tokens INTEGER DEFAULT 0, -- 总token消耗
    last_message_at DATETIME, -- 最后消息时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

-- 会话表索引
CREATE INDEX idx_session_type ON sessions(session_type);
CREATE INDEX idx_status ON sessions(status);
CREATE INDEX idx_last_message ON sessions(last_message_at);
CREATE INDEX idx_created_at ON sessions(created_at);

-- 消息表
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 消息ID，主键
    session_id INTEGER NOT NULL, -- 会话ID，外键关联sessions表
    parent_id INTEGER, -- 父消息ID，用于消息树结构
    role TEXT NOT NULL, -- 消息角色：user=用户，assistant=助手，system=系统
    content TEXT NOT NULL, -- 消息内容
    content_type TEXT DEFAULT 'text', -- 内容类型：text=文本，image=图片，file=文件，code=代码
    metadata TEXT, -- 消息元数据，JSON格式
    attachments TEXT, -- 附件信息，JSON格式
    token_count INTEGER DEFAULT 0, -- 消息token数量
    model_name TEXT, -- 使用的模型名称
    finish_reason TEXT, -- 完成原因：stop=正常结束，length=长度限制，error=错误
    is_deleted INTEGER DEFAULT 0, -- 是否已删除，1=已删除，0=正常
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES messages(id) ON DELETE SET NULL
);

-- 消息表索引
CREATE INDEX idx_messages_session ON messages(session_id);
CREATE INDEX idx_messages_parent ON messages(parent_id);
CREATE INDEX idx_messages_role ON messages(role);
CREATE INDEX idx_messages_content_type ON messages(content_type);
CREATE INDEX idx_messages_is_deleted ON messages(is_deleted);
CREATE INDEX idx_messages_created_at ON messages(created_at);

-- ============================================
-- 配置管理相关表
-- ============================================

-- 系统配置表
CREATE TABLE system_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 配置ID，主键
    config_key TEXT NOT NULL UNIQUE, -- 配置键名
    config_value TEXT, -- 配置值
    config_type TEXT DEFAULT 'string', -- 配置类型：string=字符串，json=JSON，number=数字，boolean=布尔
    category TEXT DEFAULT 'general', -- 配置分类：general=通用，llm=语言模型，ui=界面，security=安全
    description TEXT, -- 配置描述
    is_encrypted INTEGER DEFAULT 0, -- 是否加密存储，1=加密，0=明文
    is_readonly INTEGER DEFAULT 0, -- 是否只读，1=只读，0=可修改
    validation_rule TEXT, -- 验证规则
    default_value TEXT, -- 默认值
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

-- 系统配置表索引
CREATE INDEX idx_system_configs_category ON system_configs(category);
CREATE INDEX idx_system_configs_type ON system_configs(config_type);

-- 用户配置表
CREATE TABLE user_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 配置ID，主键
    user_id INTEGER NOT NULL, -- 用户ID，外键关联users表
    config_key TEXT NOT NULL, -- 配置键名
    config_value TEXT, -- 配置值
    config_type TEXT DEFAULT 'string', -- 配置类型：string=字符串，json=JSON，number=数字，boolean=布尔
    category TEXT DEFAULT 'general', -- 配置分类
    description TEXT, -- 配置描述
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    UNIQUE(user_id, config_key),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 用户配置表索引
CREATE INDEX idx_user_configs_user ON user_configs(user_id);
CREATE INDEX idx_user_configs_category ON user_configs(category);

-- ============================================
-- 知识库管理相关表
-- ============================================

-- 知识库目录表
CREATE TABLE knowledge_directories (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 目录ID，主键
    user_id INTEGER NOT NULL, -- 用户ID，外键关联users表
    name TEXT NOT NULL, -- 目录名称
    path TEXT NOT NULL, -- 目录路径
    description TEXT, -- 目录描述
    is_active INTEGER DEFAULT 1, -- 是否启用，1=启用，0=禁用
    auto_scan INTEGER DEFAULT 1, -- 是否自动扫描，1=自动，0=手动
    scan_interval INTEGER DEFAULT 3600, -- 扫描间隔，单位秒
    file_patterns TEXT, -- 文件匹配模式，JSON数组
    exclude_patterns TEXT, -- 排除模式，JSON数组
    max_file_size INTEGER DEFAULT 10485760, -- 最大文件大小，单位字节，默认10MB
    index_status TEXT DEFAULT 'pending', -- 索引状态：pending=待处理，indexing=索引中，completed=完成，error=错误
    last_scan_at DATETIME, -- 最后扫描时间
    file_count INTEGER DEFAULT 0, -- 文件总数
    indexed_count INTEGER DEFAULT 0, -- 已索引文件数
    total_size INTEGER DEFAULT 0, -- 总文件大小，单位字节
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 知识库目录表索引
CREATE INDEX idx_knowledge_dirs_user ON knowledge_directories(user_id);
CREATE INDEX idx_knowledge_dirs_active ON knowledge_directories(is_active);
CREATE INDEX idx_knowledge_dirs_status ON knowledge_directories(index_status);
CREATE INDEX idx_knowledge_dirs_scan ON knowledge_directories(last_scan_at);

-- 知识库文件表
CREATE TABLE knowledge_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 文件ID，主键
    directory_id INTEGER NOT NULL, -- 目录ID，外键关联knowledge_directories表
    file_path TEXT NOT NULL UNIQUE, -- 文件完整路径
    file_name TEXT NOT NULL, -- 文件名
    file_type TEXT, -- 文件类型扩展名
    file_size INTEGER DEFAULT 0, -- 文件大小，单位字节
    content_hash TEXT, -- 文件内容哈希值，用于检测变更
    mime_type TEXT, -- MIME类型
    encoding TEXT DEFAULT 'utf-8', -- 文件编码
    language TEXT, -- 编程语言或文档语言
    chunk_count INTEGER DEFAULT 0, -- 分块数量
    index_status TEXT DEFAULT 'pending', -- 索引状态：pending=待处理，indexing=索引中，completed=完成，error=错误
    error_message TEXT, -- 错误信息
    last_modified DATETIME, -- 文件最后修改时间
    indexed_at DATETIME, -- 索引完成时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    FOREIGN KEY (directory_id) REFERENCES knowledge_directories(id) ON DELETE CASCADE
);

-- 知识库文件表索引
CREATE INDEX idx_knowledge_files_dir ON knowledge_files(directory_id);
CREATE INDEX idx_knowledge_files_type ON knowledge_files(file_type);
CREATE INDEX idx_knowledge_files_status ON knowledge_files(index_status);
CREATE INDEX idx_knowledge_files_modified ON knowledge_files(last_modified);
CREATE INDEX idx_knowledge_files_hash ON knowledge_files(content_hash);

-- 知识库文件分块表
CREATE TABLE knowledge_chunks (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 分块ID，主键
    file_id INTEGER NOT NULL, -- 文件ID，外键关联knowledge_files表
    chunk_index INTEGER NOT NULL, -- 分块索引，从0开始
    content TEXT NOT NULL, -- 分块内容
    content_length INTEGER DEFAULT 0, -- 内容长度，字符数
    start_position INTEGER DEFAULT 0, -- 在原文件中的起始位置
    end_position INTEGER DEFAULT 0, -- 在原文件中的结束位置
    embedding_vector TEXT, -- 向量嵌入，JSON格式存储
    embedding_model TEXT, -- 嵌入模型名称
    metadata TEXT, -- 分块元数据，JSON格式
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    UNIQUE(file_id, chunk_index),
    FOREIGN KEY (file_id) REFERENCES knowledge_files(id) ON DELETE CASCADE
);

-- 知识库分块表索引
CREATE INDEX idx_knowledge_chunks_file ON knowledge_chunks(file_id);
CREATE INDEX idx_knowledge_chunks_index ON knowledge_chunks(chunk_index);
CREATE INDEX idx_knowledge_chunks_length ON knowledge_chunks(content_length);

-- 会话与知识库关系表
CREATE TABLE session_knowledge_directories (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 关系ID，主键
    session_id INTEGER NOT NULL, -- 会话ID，外键关联sessions表
    directory_id INTEGER NOT NULL, -- 知识库目录ID，外键关联knowledge_directories表
    is_enabled INTEGER DEFAULT 1, -- 是否启用，1=启用，0=禁用
    priority INTEGER DEFAULT 0, -- 优先级，数值越大优先级越高
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    UNIQUE(session_id, directory_id),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (directory_id) REFERENCES knowledge_directories(id) ON DELETE CASCADE
);

-- 会话与知识库关系表索引
CREATE INDEX idx_session_knowledge_session ON session_knowledge_directories(session_id);
CREATE INDEX idx_session_knowledge_directory ON session_knowledge_directories(directory_id);
CREATE INDEX idx_session_knowledge_enabled ON session_knowledge_directories(is_enabled);
CREATE INDEX idx_session_knowledge_priority ON session_knowledge_directories(priority);

-- 会话与工具关系表
CREATE TABLE session_tools (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 关系ID，主键
    session_id INTEGER NOT NULL, -- 会话ID，外键关联sessions表
    tool_id INTEGER NOT NULL, -- 工具ID，外键关联tools表
    is_enabled INTEGER DEFAULT 1, -- 是否启用，1=启用，0=禁用
    priority INTEGER DEFAULT 0, -- 优先级，数值越大优先级越高
    config_override TEXT, -- 配置覆盖，JSON格式，覆盖用户默认工具配置
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    UNIQUE(session_id, tool_id),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE
);

-- 会话与工具关系表索引
CREATE INDEX idx_session_tools_session ON session_tools(session_id);
CREATE INDEX idx_session_tools_tool ON session_tools(tool_id);
CREATE INDEX idx_session_tools_enabled ON session_tools(is_enabled);
CREATE INDEX idx_session_tools_priority ON session_tools(priority);

-- ============================================
-- 工具管理相关表
-- ============================================

-- 工具定义表
CREATE TABLE tools (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 工具ID，主键
    name TEXT NOT NULL UNIQUE, -- 工具名称，唯一标识
    display_name TEXT, -- 显示名称
    description TEXT, -- 工具描述
    category TEXT DEFAULT 'general', -- 工具分类：search=搜索，file=文件，automation=自动化，llm=语言模型
    tool_type TEXT DEFAULT 'builtin', -- 工具类型：builtin=内置，plugin=插件，mcp=MCP服务
    version TEXT DEFAULT '1.0.0', -- 工具版本
    config_schema TEXT, -- 配置模式，JSON Schema格式
    default_config TEXT, -- 默认配置，JSON格式
    capabilities TEXT, -- 工具能力列表，JSON数组
    requirements TEXT, -- 依赖要求，JSON格式
    is_active INTEGER DEFAULT 1, -- 是否启用，1=启用，0=禁用
    is_system INTEGER DEFAULT 0, -- 是否系统工具，1=系统，0=用户
    author TEXT, -- 工具作者
    documentation_url TEXT, -- 文档链接
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

-- 工具定义表索引
CREATE INDEX idx_tools_category ON tools(category);
CREATE INDEX idx_tools_type ON tools(tool_type);
CREATE INDEX idx_tools_active ON tools(is_active);

-- 用户工具配置表
CREATE TABLE user_tool_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 配置ID，主键
    user_id INTEGER NOT NULL, -- 用户ID，外键关联users表
    tool_id INTEGER NOT NULL, -- 工具ID，外键关联tools表
    config_data TEXT, -- 工具配置数据，JSON格式
    is_enabled INTEGER DEFAULT 1, -- 是否启用，1=启用，0=禁用
    priority INTEGER DEFAULT 0, -- 优先级，数值越大优先级越高
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    UNIQUE(user_id, tool_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE
);

-- 用户工具配置表索引
CREATE INDEX idx_user_tool_configs_user ON user_tool_configs(user_id);
CREATE INDEX idx_user_tool_configs_tool ON user_tool_configs(tool_id);
CREATE INDEX idx_user_tool_configs_enabled ON user_tool_configs(is_enabled);
CREATE INDEX idx_user_tool_configs_priority ON user_tool_configs(priority);

-- MCP服务表
CREATE TABLE mcp_services (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- MCP服务ID，主键
    name TEXT NOT NULL UNIQUE, -- MCP服务名称，唯一标识
    display_name TEXT, -- 显示名称
    description TEXT, -- 服务描述
    service_type TEXT DEFAULT 'external', -- 服务类型：external=外部服务，builtin=内置服务
    endpoint_url TEXT, -- 服务端点URL
    api_version TEXT DEFAULT '1.0', -- API版本
    authentication_type TEXT DEFAULT 'none', -- 认证类型：none=无认证，apikey=API密钥，oauth=OAuth
    config_schema TEXT, -- 配置模式，JSON Schema格式
    default_config TEXT, -- 默认配置，JSON格式
    capabilities TEXT, -- 服务能力列表，JSON数组
    is_active INTEGER DEFAULT 1, -- 是否启用，1=启用，0=禁用
    health_check_url TEXT, -- 健康检查URL
    last_health_check DATETIME, -- 最后健康检查时间
    health_status TEXT DEFAULT 'unknown', -- 健康状态：healthy=健康，unhealthy=不健康，unknown=未知
    timeout_seconds INTEGER DEFAULT 30, -- 超时时间，单位秒
    retry_count INTEGER DEFAULT 3, -- 重试次数
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

-- MCP服务表索引
CREATE INDEX idx_mcp_services_type ON mcp_services(service_type);
CREATE INDEX idx_mcp_services_active ON mcp_services(is_active);
CREATE INDEX idx_mcp_services_health ON mcp_services(health_status);
CREATE INDEX idx_mcp_services_check ON mcp_services(last_health_check);

-- 用户MCP服务配置表
CREATE TABLE user_mcp_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 配置ID，主键
    user_id INTEGER NOT NULL, -- 用户ID，外键关联users表
    mcp_service_id INTEGER NOT NULL, -- MCP服务ID，外键关联mcp_services表
    config_data TEXT, -- 服务配置数据，JSON格式
    is_enabled INTEGER DEFAULT 1, -- 是否启用，1=启用，0=禁用
    priority INTEGER DEFAULT 0, -- 优先级，数值越大优先级越高
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    UNIQUE(user_id, mcp_service_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (mcp_service_id) REFERENCES mcp_services(id) ON DELETE CASCADE
);

-- 用户MCP服务配置表索引
CREATE INDEX idx_user_mcp_configs_user ON user_mcp_configs(user_id);
CREATE INDEX idx_user_mcp_configs_service ON user_mcp_configs(mcp_service_id);
CREATE INDEX idx_user_mcp_configs_enabled ON user_mcp_configs(is_enabled);
CREATE INDEX idx_user_mcp_configs_priority ON user_mcp_configs(priority);

-- 会话与MCP服务关系表
CREATE TABLE session_mcp_services (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 关系ID，主键
    session_id INTEGER NOT NULL, -- 会话ID，外键关联sessions表
    mcp_service_id INTEGER NOT NULL, -- MCP服务ID，外键关联mcp_services表
    is_enabled INTEGER DEFAULT 1, -- 是否启用，1=启用，0=禁用
    priority INTEGER DEFAULT 0, -- 优先级，数值越大优先级越高
    config_override TEXT, -- 配置覆盖，JSON格式，覆盖用户默认MCP配置
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    UNIQUE(session_id, mcp_service_id),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (mcp_service_id) REFERENCES mcp_services(id) ON DELETE CASCADE
);

-- 会话与MCP服务关系表索引
CREATE INDEX idx_session_mcp_session ON session_mcp_services(session_id);
CREATE INDEX idx_session_mcp_service ON session_mcp_services(mcp_service_id);
CREATE INDEX idx_session_mcp_enabled ON session_mcp_services(is_enabled);
CREATE INDEX idx_session_mcp_priority ON session_mcp_services(priority);

-- ============================================
-- 任务管理相关表
-- ============================================

-- 任务表
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 任务ID，主键
    session_id INTEGER, -- 会话ID，外键关联sessions表
    user_id INTEGER NOT NULL, -- 用户ID，外键关联users表
    task_type TEXT NOT NULL, -- 任务类型：automation=自动化，indexing=索引，search=搜索，llm=语言模型
    title TEXT NOT NULL, -- 任务标题
    description TEXT, -- 任务描述
    input_data TEXT, -- 输入数据，JSON格式
    output_data TEXT, -- 输出数据，JSON格式
    status TEXT DEFAULT 'pending', -- 任务状态：pending=待处理，running=运行中，completed=完成，failed=失败，cancelled=取消
    progress INTEGER DEFAULT 0, -- 进度百分比，0-100
    error_message TEXT, -- 错误信息
    started_at DATETIME, -- 开始时间
    completed_at DATETIME, -- 完成时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 任务表索引
CREATE INDEX idx_tasks_session ON tasks(session_id);
CREATE INDEX idx_tasks_user ON tasks(user_id);
CREATE INDEX idx_tasks_type ON tasks(task_type);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_created ON tasks(created_at);

-- ============================================
-- 操作日志相关表
-- ============================================

-- 操作日志表
CREATE TABLE operation_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 日志ID，主键
    user_id INTEGER, -- 用户ID，外键关联users表
    session_id INTEGER, -- 会话ID，外键关联sessions表
    task_id INTEGER, -- 任务ID，外键关联tasks表
    operation_type TEXT NOT NULL, -- 操作类型：mouse=鼠标，keyboard=键盘，file=文件，window=窗口，api=接口调用
    operation_name TEXT NOT NULL, -- 操作名称
    target TEXT, -- 操作目标
    parameters TEXT, -- 操作参数，JSON格式
    result TEXT, -- 操作结果，JSON格式
    status TEXT DEFAULT 'success', -- 操作状态：success=成功，failed=失败，timeout=超时
    duration_ms INTEGER DEFAULT 0, -- 操作耗时，单位毫秒
    error_message TEXT, -- 错误信息
    ip_address TEXT, -- IP地址
    user_agent TEXT, -- 用户代理
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE SET NULL,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE SET NULL
);

-- 操作日志表索引
CREATE INDEX idx_operation_logs_user ON operation_logs(user_id);
CREATE INDEX idx_operation_logs_session ON operation_logs(session_id);
CREATE INDEX idx_operation_logs_task ON operation_logs(task_id);
CREATE INDEX idx_operation_logs_type ON operation_logs(operation_type);
CREATE INDEX idx_operation_logs_status ON operation_logs(status);
CREATE INDEX idx_operation_logs_created ON operation_logs(created_at);

-- ============================================
-- 文件附件相关表
-- ============================================

-- 文件附件表
CREATE TABLE file_attachments (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 附件ID，主键
    message_id INTEGER, -- 消息ID，外键关联messages表
    user_id INTEGER NOT NULL, -- 用户ID，外键关联users表
    original_name TEXT NOT NULL, -- 原始文件名
    stored_name TEXT NOT NULL, -- 存储文件名
    file_path TEXT NOT NULL, -- 文件存储路径
    file_size INTEGER DEFAULT 0, -- 文件大小，单位字节
    mime_type TEXT, -- MIME类型
    file_hash TEXT, -- 文件哈希值
    thumbnail_path TEXT, -- 缩略图路径
    metadata TEXT, -- 文件元数据，JSON格式
    is_temporary INTEGER DEFAULT 0, -- 是否临时文件，1=临时，0=永久
    expires_at DATETIME, -- 过期时间，仅临时文件
    download_count INTEGER DEFAULT 0, -- 下载次数
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 文件附件表索引
CREATE INDEX idx_file_attachments_message ON file_attachments(message_id);
CREATE INDEX idx_file_attachments_user ON file_attachments(user_id);
CREATE INDEX idx_file_attachments_hash ON file_attachments(file_hash);
CREATE INDEX idx_file_attachments_temporary ON file_attachments(is_temporary);
CREATE INDEX idx_file_attachments_expires ON file_attachments(expires_at);

-- ============================================
-- 系统统计相关表
-- ============================================

-- 系统统计表
CREATE TABLE system_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 统计ID，主键
    stat_date DATE NOT NULL, -- 统计日期
    stat_type TEXT NOT NULL, -- 统计类型：daily=日统计，hourly=小时统计
    metric_name TEXT NOT NULL, -- 指标名称
    metric_value REAL DEFAULT 0.0, -- 指标值
    additional_data TEXT, -- 附加数据，JSON格式
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    UNIQUE(stat_date, stat_type, metric_name)
);

-- 系统统计表索引
CREATE INDEX idx_system_stats_date ON system_stats(stat_date);
CREATE INDEX idx_system_stats_type ON system_stats(stat_type);
CREATE INDEX idx_system_stats_metric ON system_stats(metric_name);

-- ============================================
-- 初始化数据
-- ============================================

-- 插入系统配置初始数据
INSERT INTO system_configs (config_key, config_value, config_type, category, description, default_value) VALUES
('app_name', 'GoManus', 'string', 'general', '应用程序名称', 'GoManus'),
('app_version', '1.0.0', 'string', 'general', '应用程序版本', '1.0.0'),
('max_sessions', '100', 'number', 'general', '最大会话数量', '100'),
('default_model', 'gpt-4', 'string', 'llm', '默认LLM模型', 'gpt-4'),
('max_tokens', '4096', 'number', 'llm', '最大token数量', '4096'),
('temperature', '0.7', 'number', 'llm', '模型温度参数', '0.7'),
('ui_theme', 'light', 'string', 'ui', 'UI主题', 'light'),
('auto_save', 'true', 'boolean', 'general', '自动保存', 'true'),
('log_level', 'info', 'string', 'general', '日志级别', 'info'),
('max_file_size', '10485760', 'number', 'general', '最大文件大小（字节）', '10485760');

-- 插入工具初始数据
INSERT INTO tools (name, display_name, description, category, tool_type, capabilities, is_system) VALUES
('file_search', '文件搜索', '在指定目录中搜索文件', 'file', 'builtin', '["search", "filter"]', 1),
('text_editor', '文本编辑器', '编辑文本文件', 'file', 'builtin', '["read", "write", "edit"]', 1),
('web_search', '网络搜索', '在互联网上搜索信息', 'search', 'builtin', '["search", "crawl"]', 1),
('mouse_control', '鼠标控制', '控制鼠标操作', 'automation', 'builtin', '["click", "move", "drag"]', 1),
('keyboard_control', '键盘控制', '控制键盘输入', 'automation', 'builtin', '["type", "hotkey"]', 1),
('window_control', '窗口控制', '控制应用程序窗口', 'automation', 'builtin', '["focus", "resize", "move"]', 1),
('llm_chat', 'LLM对话', '与语言模型进行对话', 'llm', 'builtin', '["chat", "completion"]', 1);

-- 插入MCP服务初始数据
INSERT INTO mcp_services (name, display_name, description, service_type, capabilities) VALUES
('filesystem_mcp', '文件系统MCP', '提供文件系统操作的MCP服务', 'builtin', '["file_read", "file_write", "directory_list"]'),
('database_mcp', '数据库MCP', '提供数据库操作的MCP服务', 'builtin', '["query", "insert", "update", "delete"]'),
('api_gateway_mcp', 'API网关MCP', '提供外部API调用的MCP服务', 'external', '["http_request", "webhook", "api_proxy"]'),
('notification_mcp', '通知MCP', '提供消息通知的MCP服务', 'external', '["email", "sms", "push_notification"]'),
('scheduler_mcp', '调度器MCP', '提供任务调度的MCP服务', 'builtin', '["cron", "timer", "delayed_task"]');

-- ============================================
-- 创建复合索引
-- ============================================

-- 消息相关复合索引
CREATE INDEX idx_messages_session_created ON messages(session_id, created_at);
CREATE INDEX idx_knowledge_chunks_file_index ON knowledge_chunks(file_id, chunk_index);
CREATE INDEX idx_operation_logs_user_created ON operation_logs(user_id, created_at);
CREATE INDEX idx_tasks_user_status ON tasks(user_id, status);
CREATE INDEX idx_file_attachments_user_created ON file_attachments(user_id, created_at);

-- 会话关系表复合索引
CREATE INDEX idx_session_knowledge_session_enabled ON session_knowledge_directories(session_id, is_enabled);
CREATE INDEX idx_session_tools_session_enabled ON session_tools(session_id, is_enabled);
CREATE INDEX idx_session_mcp_session_enabled ON session_mcp_services(session_id, is_enabled);
CREATE INDEX idx_session_knowledge_enabled_priority ON session_knowledge_directories(is_enabled, priority);
CREATE INDEX idx_session_tools_enabled_priority ON session_tools(is_enabled, priority);
CREATE INDEX idx_session_mcp_enabled_priority ON session_mcp_services(is_enabled, priority);

-- ============================================
-- 创建视图
-- ============================================

-- 会话摘要视图
CREATE VIEW v_session_summary AS
SELECT 
    s.id,
    s.title,
    s.session_type,
    s.status,
    s.message_count,
    s.total_tokens,
    s.last_message_at,
    s.created_at,
    COUNT(m.id) as actual_message_count,
    MAX(m.created_at) as latest_message_time
FROM sessions s
LEFT JOIN messages m ON s.id = m.session_id AND m.is_deleted = 0
GROUP BY s.id, s.title, s.session_type, s.status, s.message_count, s.total_tokens, s.last_message_at, s.created_at;

-- 知识库目录统计视图
CREATE VIEW v_knowledge_directory_stats AS
SELECT 
    kd.id,
    kd.name,
    kd.path,
    kd.index_status,
    kd.file_count,
    kd.indexed_count,
    kd.total_size,
    COUNT(kf.id) as actual_file_count,
    SUM(kf.file_size) as actual_total_size,
    COUNT(CASE WHEN kf.index_status = 'completed' THEN 1 END) as actual_indexed_count
FROM knowledge_directories kd
LEFT JOIN knowledge_files kf ON kd.id = kf.directory_id
GROUP BY kd.id, kd.name, kd.path, kd.index_status, kd.file_count, kd.indexed_count, kd.total_size;

-- ============================================
-- 数据库初始化完成
-- ============================================

-- 输出完成信息
SELECT 'GoManus SQLite数据库初始化完成！' as message;
SELECT 'sqlite-vec向量扩展需要单独安装和配置' as note;
SELECT '数据库文件建议存储路径：./data/gomanus.db' as suggestion;