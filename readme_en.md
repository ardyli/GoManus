# GoManus Project
<p align="center">
  <img src="logo/logo.jpg" width="200"/>
</p>

## Project Introduction

### GoManus: The Revolution of AI Agent Systems, Making Automation So Simple!
Imagine an AI agent system that not only understands every command you give but also acts like a real assistant, helping you complete various complex tasks. This is not a plot from science fiction, but GoManus—a Go-based AI agent system that is making all of this a reality!

### What is GoManus?
GoManus, a name that sounds full of power, is actually a fully open-source AI agent system. It not only helps users interact with AI but also enables automated and intelligent task processing. This means that whether it's data analysis, file management, or web searching, GoManus can handle it all for you.

### Why Choose GoManus?
The reasons to choose GoManus are countless! First, it's completely open source, which means you can develop on top of this project and create an AI assistant that's entirely your own. Second, GoManus doesn't require complex environment deployment. Forget about those headache-inducing conda and python environment packages—GoManus only requires you to download a GoManus.exe file, then run it via CMD, and it can run on your computer!

### How to Use GoManus?
Using GoManus is simpler than eating cake. You just need to download the GoManus.exe file and run it. Next, configure the config/config.toml file, and your AI assistant is ready. Yes, it's that simple! GoManus's powerful features are incredibly impressive. It supports mainstream LLM models on the market, can interact with AI, and achieve automated and intelligent task processing. Additionally, it can interact with various tools, including search tools like Baidu Baike, Google, Zhihu, Wikipedia, and more. File saving, browser usage, task termination, multi-tool calling and planning, memory management—GoManus has all these features.

## Project Repository

- Project Repository (China): https://gitee.com/therebody/go-manus

- Project Repository: https://github.com/ardyli/GoManus

## Features

- Interact with LLM, supporting mainstream LLM models on the market
- Interact with AI to achieve automated and intelligent task processing
- Interact with tools (supporting search tools like Baidu Baike, Google, Zhihu, Wikipedia, etc.)
- File saving functionality
- Browser usage functionality
- Task termination functionality
- Multi-tool calling and planning functionality
- Memory management functionality
- Command terminal operations, capable of operating windows\linux\mac through shell commands

## Introduction
**Startup Interface:**
<p align="center">
  <img src="images/boot.png" width="100%"/>
</p>

## Technology Stack

- Go 1.24
- Viper v1.18.2
- GORM v1.25.7
- SQLite v1.14.17
- Sonic v1.9.1
- PTerm   https://pterm.sh/

## Project Structure

- `config/`: Configuration file directory
  - `config.toml`: Configuration file
- `internal/`: Core implementation
  - `agent/`: AI agent implementation
    - `base.go`: Basic functionality
    - `manus.go`: Main logic
    - `planning.go`: Planning functionality
    - `react.go`: Reaction functionality
    - `toolcall.go`: Tool calling
  - `config/`: Configuration management
    - `config.go`: Configuration loading
  - `llm/`: LLM interaction
    - `llm.go`: LLM interface
  - `middleware/`: Middleware
    - `refresh.go`: Auto refresh
  - `schema/`: Data structures
    - `agent.go`: Agent related
    - `message.go`: Message structure
    - `toolcall.go`: Tool call structure
  - `tool/`: Tool implementation
    - `baidu_baike_search.go`: Baidu Baike search
    - `base.go`: Tool base
    - `browser_use.go`: Browser usage
    - `collection.go`: Tool collection
    - `file_saver.go`: File saving
    - `google_search.go`: Google search
    - `planning.go`: Tool planning
    - `terminate.go`: Task termination
    - `wikipedia_search.go`: Wikipedia search
    - `zhihu_search.go`: Zhihu search
- `main.go`: Project entry file
- `go.mod` and `go.sum`: Go module dependency management files

## Usage Instructions
Download version:
  Download Windows version:
   Execute in command window:
   Enter directory:
   cd GoManus
   Execute in command window:
   ./GoManus.exe
  Download Linux version:
  chmod +x GoManus
  Execute in command window:
  ./GoManus

### Source code version:
1. Ensure Go 1.21 or higher is installed
2. Clone the project locally:
   ```bash
   git clone https://gitee.com/therebody/GoManus.git
   ```
3. Enter the project directory and install dependencies:
   ```bash
   go mod tidy
   ```
4. Run the project:
   ```bash
   go run main.go
   ```

## Configuration Instructions

1. Modify the `config/config.toml` file to configure system parameters
2. Configure LLM API keys
3. Configure tool-related parameters

## Contribution Guide

Welcome to submit PRs or Issues to improve the project. Main contribution directions include:

- New tool integration
- UI functionality improvements
- Performance optimization
- Documentation improvement
- Test case writing

Contribution process:

1. Fork this repository
2. Create a feature branch (Feat_xxx) or fix branch (Fix_xxx)
3. Submit code changes
4. Create a Pull Request
5. Wait for code review and merge

## License

This project uses the BSD3 license. For details, please see the LICENSE file.