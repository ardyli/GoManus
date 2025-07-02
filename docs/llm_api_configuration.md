# LLM API 配置指南

本文档说明如何配置GoManus系统以支持不同类型的LLM API接口。

## 支持的API类型

系统现在支持以下两种API类型：

1. **ollama** - Ollama本地部署的模型API
2. **openai** - OpenAI兼容的API接口

## 配置说明

### 基本配置结构

在 `config/config.toml` 文件中，每个LLM配置都包含以下字段：

```toml
[llm]
model = "模型名称"
base_url = "API基础URL"
api_key = "API密钥"
api_type = "API类型"  # "ollama" 或 "openai"
max_tokens = 最大令牌数
temperature = 温度参数
```

### Ollama配置示例

```toml
# 默认Ollama配置
[llm]
model = "qwen2.5:14b"
base_url = "http://localhost:11434/v1/"
api_key = "ollama"
api_type = "ollama"
max_tokens = 131072
temperature = 0.0

# 视觉模型配置
[llm_types.vision]
model = "llava:13b"
base_url = "http://localhost:11434/v1/"
api_key = "ollama"
api_type = "ollama"
```

### OpenAI兼容接口配置示例

```toml
# OpenAI官方API
[llm_types.openai_gpt4]
model = "gpt-4"
base_url = "https://api.openai.com/v1/"
api_key = "sk-your-openai-api-key"
api_type = "openai"
max_tokens = 4096
temperature = 0.7

# DeepSeek API
[llm_types.deepseek]
model = "deepseek-chat"
base_url = "https://api.deepseek.com/v1/"
api_key = "your-deepseek-api-key"
api_type = "openai"
max_tokens = 8192
temperature = 0.3

# 智谱AI API
[llm_types.zhipu]
model = "glm-4"
base_url = "https://open.bigmodel.cn/api/paas/v4/"
api_key = "your-zhipu-api-key"
api_type = "openai"
max_tokens = 4096
temperature = 0.5
```

## API类型差异说明

### Ollama API
- 通常部署在本地或内网
- 认证方式较为宽松，api_key可以设置为"ollama"或留空
- 工具调用格式与OpenAI基本兼容

### OpenAI兼容API
- 需要有效的API密钥进行认证
- 使用标准的Bearer Token认证
- 严格遵循OpenAI API规范

## 使用方法

### 在代码中使用

```go
// 使用默认配置
llm, err := llm.NewLLM("")
if err != nil {
    log.Fatal(err)
}

// 使用特定配置
llm, err := llm.NewLLM("openai_gpt4")
if err != nil {
    log.Fatal(err)
}

// 发送请求
response, err := llm.AskWithOptions(ctx, messages, systemMsgs, tools, &toolChoice)
```

### 配置切换

只需要修改配置文件中的相应字段，无需修改代码：

1. 修改 `api_type` 字段来切换API类型
2. 修改 `base_url` 和 `api_key` 来切换服务提供商
3. 修改 `model` 来切换具体模型

## 注意事项

1. **API密钥安全**：请妥善保管API密钥，不要将其提交到版本控制系统
2. **网络连接**：确保系统能够访问配置的API端点
3. **模型兼容性**：不同的模型可能支持不同的功能，请根据实际需求选择
4. **费用控制**：使用付费API时请注意token消耗和费用控制

## 故障排除

### 常见问题

1. **认证失败**
   - 检查API密钥是否正确
   - 确认API类型设置是否匹配

2. **连接超时**
   - 检查网络连接
   - 确认base_url是否正确

3. **模型不存在**
   - 确认模型名称是否正确
   - 检查API提供商是否支持该模型

### 调试方法

系统会输出详细的日志信息，包括：
- LLM实例创建信息
- 请求和响应详情
- 错误信息和状态码

查看日志可以帮助快速定位问题。