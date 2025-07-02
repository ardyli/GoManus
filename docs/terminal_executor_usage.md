# 终端执行工具使用说明

## 概述

终端执行工具（Terminal Executor）是GoManus系统中的一个重要工具，它允许AI代理执行操作系统级别的终端命令，并获取命令的执行结果。该工具支持跨平台操作，能够在Windows、Linux和macOS系统上正常工作。

## 功能特性

- **跨平台支持**：支持Windows的cmd和PowerShell，以及Linux/macOS的bash和sh
- **智能Shell选择**：可以自动选择最适合的shell，也可以手动指定
- **工作目录设置**：支持指定命令执行的工作目录
- **超时控制**：可以设置命令执行的超时时间，防止长时间阻塞
- **错误捕获**：能够捕获标准输出和标准错误输出
- **详细结果**：返回包含退出码、输出内容、执行环境等详细信息

## 参数说明

### 必填参数

- `command` (string): 要执行的终端命令

### 可选参数

- `shell_type` (string): 指定使用的shell类型
  - `"auto"` (默认): 自动选择最适合的shell
  - `"cmd"`: Windows命令提示符
  - `"powershell"`: Windows PowerShell
  - `"bash"`: Unix/Linux Bash shell

- `working_directory` (string): 命令执行的工作目录
  - 默认为当前目录
  - 必须是存在的目录路径

- `timeout` (integer): 命令执行超时时间（秒）
  - 默认为30秒
  - 超时后命令会被强制终止

- `capture_stderr` (boolean): 是否捕获标准错误输出
  - 默认为true
  - 设置为false可以忽略错误输出

## 返回结果

工具执行后返回一个包含以下字段的JSON对象：

```json
{
  "command": "执行的命令",
  "shell": "实际使用的shell类型",
  "os": "操作系统类型",
  "exit_code": 0,
  "stdout": "标准输出内容",
  "stderr": "标准错误输出内容（如果有）",
  "success": true,
  "working_dir": "工作目录路径"
}
```

## 使用示例

### 示例1：基本命令执行

```json
{
  "command": "echo Hello World"
}
```

### 示例2：指定PowerShell执行

```json
{
  "command": "Get-Process | Select-Object -First 5",
  "shell_type": "powershell"
}
```

### 示例3：在指定目录执行命令

```json
{
  "command": "ls -la",
  "working_directory": "/home/user/documents",
  "shell_type": "bash"
}
```

### 示例4：设置超时时间

```json
{
  "command": "ping google.com -c 4",
  "timeout": 10,
  "capture_stderr": true
}
```

## 平台差异

### Windows系统

- 默认优先使用PowerShell（如果可用），否则使用cmd
- 支持所有Windows内置命令和PowerShell cmdlet
- 路径分隔符使用反斜杠（\）

**常用命令示例：**
```bash
# 列出目录内容
dir
Get-ChildItem  # PowerShell

# 查看当前路径
cd
Get-Location   # PowerShell

# 创建目录
mkdir newfolder
New-Item -ItemType Directory -Name "newfolder"  # PowerShell
```

### Linux/macOS系统

- 默认优先使用bash，否则使用sh
- 支持所有Unix/Linux标准命令
- 路径分隔符使用正斜杠（/）

**常用命令示例：**
```bash
# 列出目录内容
ls -la

# 查看当前路径
pwd

# 创建目录
mkdir newfolder

# 查看系统信息
uname -a
```

## 安全注意事项

1. **命令验证**：在执行命令前，应该验证命令的安全性
2. **权限控制**：避免执行需要管理员权限的危险命令
3. **路径安全**：确保工作目录和文件路径的安全性
4. **超时设置**：为长时间运行的命令设置合理的超时时间
5. **错误处理**：妥善处理命令执行失败的情况

## 常见问题

### Q: 命令执行失败怎么办？
A: 检查返回结果中的`exit_code`和`stderr`字段，根据错误信息进行调试。

### Q: 如何执行需要交互的命令？
A: 该工具不支持交互式命令，请使用非交互式参数或脚本。

### Q: 可以执行多行命令吗？
A: 可以使用shell的命令连接符（如`&&`、`;`、`|`）来执行多个命令。

### Q: 如何处理包含特殊字符的命令？
A: 确保命令字符串正确转义，特别是引号和反斜杠。

## 配置说明

在`config/config.toml`文件中启用终端执行工具：

```toml
[tools]
terminal_executor = true  # 启用终端命令执行工具
```

## 集成说明

终端执行工具已经集成到GoManus的工具系统中，可以通过以下方式使用：

1. 确保在配置文件中启用了该工具
2. 在代理对话中请求执行终端命令
3. AI代理会自动调用该工具执行命令并返回结果

该工具为AI代理提供了强大的系统操作能力，使其能够执行文件操作、系统管理、程序运行等各种任务。