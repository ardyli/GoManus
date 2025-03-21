# GoManus 项目

GoManus 是一个基于 Go 语言的 AI 代理系统，它可以帮助用户与 AI 进行交互，实现自动化和智能化的任务处理。

## 功能

- 与 LLM 进行交互
- 与工具进行交互（支持百度百科、Google、知乎、维基百科等搜索工具）
- 文件保存功能
- 浏览器使用功能
- 任务终止功能
- 多工具调用与规划功能
- 记忆管理功能

## 技术栈

- Go 1.21
- Gin v1.9.1
- Viper v1.18.2
- GORM v1.25.7
- SQLite v1.14.17
- Sonic v1.9.1

## 项目结构

- `config/`: 配置文件目录
  - `config.toml`: 配置文件
- `internal/`: 核心实现
  - `agent/`: AI 代理实现
    - `base.go`: 基础功能
    - `manus.go`: 主逻辑
    - `planning.go`: 规划功能
    - `react.go`: 反应功能
    - `toolcall.go`: 工具调用
  - `config/`: 配置管理
    - `config.go`: 配置加载
  - `llm/`: LLM 交互
    - `llm.go`: LLM 接口
  - `middleware/`: 中间件
    - `refresh.go`: 自动刷新
  - `schema/`: 数据结构
    - `agent.go`: 代理相关
    - `message.go`: 消息结构
    - `toolcall.go`: 工具调用结构
  - `tool/`: 工具实现
    - `baidu_baike_search.go`: 百度百科搜索
    - `base.go`: 工具基础
    - `browser_use.go`: 浏览器使用
    - `collection.go`: 工具集合
    - `file_saver.go`: 文件保存
    - `google_search.go`: Google 搜索
    - `planning.go`: 工具规划
    - `terminate.go`: 任务终止
    - `wikipedia_search.go`: 维基百科搜索
    - `zhihu_search.go`: 知乎搜索
- `main.go`: 项目入口文件
- `go.mod` 和 `go.sum`: Go 模块依赖管理文件

## 使用说明
 下载版：
  下载windows版本：
   命令窗口执行：
   进入目录:
   cd GoManus
   命令窗口执行：
   ./GoManus.exe


### 源码版：
1. 确保已安装 Go 1.21 或更高版本
2. 克隆项目到本地：
   ```bash
   git clone https://gitee.com/therebody/GoManus.git
   ```
3. 进入项目目录并安装依赖：
   ```bash
   go mod tidy
   ```
4. 运行项目：
   ```bash
   go run main.go
   ```


## 配置说明

1. 修改 `config/config.toml` 文件配置系统参数
2. 配置 LLM API 密钥
3. 配置工具相关参数

## 贡献指南

欢迎提交 PR 或 Issue 来改进项目。主要贡献方向包括：

- 新工具集成
- UI 功能改进
- 性能优化
- 文档完善
- 测试用例编写

贡献流程：

1. Fork 本仓库
2. 新建功能分支（Feat_xxx）或修复分支（Fix_xxx）
3. 提交代码变更
4. 新建 Pull Request
5. 等待代码审查与合并

## 许可证

本项目采用 BSD3 许可证，详情请查看 LICENSE 文件。
