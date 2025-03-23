# GoManus Project

## Project Introduction

### GoManus: Revolutionizing AI Agent Systems, Making Automation So Simple!
Imagine an AI agent system that not only understands your every command but also acts like a real assistant to help you complete various complex tasks. This is not a scene from a science fiction novel but GoManus—a Go-based AI agent system that is turning this into reality!

### What is GoManus?
GoManus, a name that sounds powerful, is actually a fully open-source AI agent system. It not only helps users interact with AI but also enables automated and intelligent task processing. This means that whether it's data analysis, file management, or web searches, GoManus has got you covered.

### Why Choose GoManus?
The reasons to choose GoManus are countless! First, it is completely open-source, which means you can build your own AI assistant based on this project. Second, GoManus does not require a complex environment setup. Forget about those headache-inducing conda and Python environment packages; GoManus only requires you to download a single GoManus.exe file, and with CMD, it can run on your computer!

### How to Use GoManus?
Using GoManus is simpler than eating a piece of cake. You just need to download the GoManus.exe file and run it. Then, configure the `config/config.toml` file, and your AI assistant is ready. Yes, it's that simple! The powerful features of GoManus are incredible. It supports mainstream LLM models, enabling interaction with AI and automated intelligent task processing. Additionally, it can interact with various tools, including Baidu Baike, Google, Zhihu, Wikipedia, and more. File saving, browser usage, task termination, multi-tool invocation and planning, memory management—GoManus has it all.

## Project Links

- Project Link (China): https://gitee.com/therebody/go-manus

- Project Link: https://github.com/ardyli/GoManus

## Features

- Interact with LLMs, supporting mainstream LLM models
- Interact with AI for automated and intelligent task processing
- Interact with tools (supports Baidu Baike, Google, Zhihu, Wikipedia, etc.)
- File saving functionality
- Browser usage functionality
- Task termination functionality
- Multi-tool invocation and planning functionality
- Memory management functionality

## Tech Stack

- Go 1.21
- Gin v1.9.1
- Viper v1.18.2
- GORM v1.25.7
- SQLite v1.14.17
- Sonic v1.9.1

## Project Structure

- `config/`: Configuration file directory
  - `config.toml`: Configuration file
- `internal/`: Core implementation
  - `agent/`: AI agent implementation
    - `base.go`: Basic functionality
    - `manus.go`: Main logic
    - `planning.go`: Planning functionality
    - `react.go`: Reaction functionality
    - `toolcall.go`: Tool invocation
  - `config/`: Configuration management
    - `config.go`: Configuration loading
  - `llm/`: LLM interaction
    - `llm.go`: LLM interface
  - `middleware/`: Middleware
    - `refresh.go`: Auto-refresh
  - `schema/`: Data structures
    - `agent.go`: Agent-related
    - `message.go`: Message structure
    - `toolcall.go`: Tool invocation structure
  - `tool/`: Tool implementation
    - `baidu_baike_search.go`: Baidu Baike search
    - `base.go`: Tool basics
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

### Download Version:
1. Download the Windows version.
2. Execute in the command window:
   Navigate to the directory:
   ```bash
   cd GoManus
   ```
   Execute in the command window:
   ```bash
   ./GoManus.exe
   ```

### Source Code Version:
1. Ensure Go 1.21 or higher is installed.
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

1. Modify the `config/config.toml` file to configure system parameters.
2. Configure the LLM API key.
3. Configure tool-related parameters.

## Contribution Guidelines

We welcome PRs or Issues to improve the project. Major contribution directions include:

- New tool integration
- UI feature improvements
- Performance optimization
- Documentation enhancement
- Test case writing

Contribution Process:

1. Fork this repository.
2. Create a feature branch (Feat_xxx) or fix branch (Fix_xxx).
3. Commit code changes.
4. Create a Pull Request.
5. Wait for code review and merging.

## License

This project is licensed under the BSD3 License. For details, please refer to the LICENSE file.
