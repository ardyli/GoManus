package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gomanus/internal/agent"
	"gomanus/internal/config"
	"gomanus/internal/llm"
	"gomanus/internal/tool"
	"gomanus/pkg/logger"

	"github.com/pterm/pterm"
)

func main() {
	// 设置日志级别
	logger.SetLevel(logger.LevelInfo)
	// 显示欢迎信息
	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).WithTextStyle(pterm.NewStyle(pterm.FgBlack)).Println("GoManus AI 助手")
	// 从配置文件加载配置
	pterm.Info.Println("⚙️  正在加载配置...")
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		logger.Fatal("加载配置失败: %v", err)
	}
	pterm.Success.Printf("✅ 配置加载成功: 使用模型 %s\n", cfg.LLM.Model)
	logger.Info("配置加载成功: 使用模型 %s", cfg.LLM.Model)

	// 获取工具配置
	pterm.Debug.Println("  📋 获取工具配置...")
	toolsCfg, err := config.GetToolsConfig()
	if err != nil {
		logger.Fatal("获取工具配置失败: %v", err)
	}
	pterm.Success.Println("✅ 工具配置获取成功")

	// 创建LLM实例
	pterm.Info.Println("🧠 正在初始化语言模型...")
	llmInstance, err := llm.NewLLM("") // 使用默认配置
	if err != nil {
		logger.Fatal("初始化语言模型失败: %v", err)
	}
	pterm.Success.Printf("✅ 语言模型初始化成功: %s\n", llmInstance.Model)
	logger.Info("语言模型初始化成功:\n %s \n %s \n %s \n %s \n %s", llmInstance.Model, llmInstance.BaseURL, llmInstance.MaxTokens, llmInstance.Temperature, llmInstance.APIKey)

	// 创建工具集合
	pterm.Info.Println("🔧 正在初始化工具集合...")
	tools := tool.NewToolCollection()

	// 根据配置添加工具
	pterm.Info.Println("📦 开始加载工具模块...")

	// 添加Terminate工具
	if toolsCfg.Terminate {
		pterm.Debug.Println("  ⚡ 加载 Terminate 工具")
		terminateTool := tool.NewTerminate()
		if err := tools.AddTool(terminateTool); err != nil {
			logger.Fatal("添加Terminate工具失败: %v", err)
		}
	}

	// 添加GoogleSearch工具
	if toolsCfg.GoogleSearch {
		pterm.Debug.Println("  🔍 加载 Google Search 工具")
		googleSearchTool := tool.NewGoogleSearch()
		if err := tools.AddTool(googleSearchTool); err != nil {
			logger.Fatal("添加GoogleSearch工具失败: %v", err)
		}
	}

	// 添加ZhihuSearch工具
	if toolsCfg.ZhihuSearch {
		pterm.Debug.Println("  📚 加载 Zhihu Search 工具")
		zhihuSearchTool := tool.NewZhihuSearch()
		if err := tools.AddTool(zhihuSearchTool); err != nil {
			logger.Fatal("添加ZhihuSearch工具失败: %v", err)
		}
	}

	// 添加BaiduBaikeSearch工具
	if toolsCfg.BaiduBaikeSearch {
		pterm.Debug.Println("  📖 加载 Baidu Baike Search 工具")
		baiduBaikeSearchTool := tool.NewBaiduBaikeSearch()
		if err := tools.AddTool(baiduBaikeSearchTool); err != nil {
			logger.Fatal("添加BaiduBaikeSearch工具失败: %v", err)
		}
	}

	// 添加WikipediaSearch工具
	if toolsCfg.WikipediaSearch {
		pterm.Debug.Println("  🌐 加载 Wikipedia Search 工具")
		wikipediaSearchTool := tool.NewWikipediaSearch()
		if err := tools.AddTool(wikipediaSearchTool); err != nil {
			logger.Fatal("添加WikipediaSearch工具失败: %v", err)
		}
	}

	// 添加BrowserUseTool工具
	if toolsCfg.BrowserUse {
		pterm.Debug.Println("  🌍 加载 Browser Use 工具")
		browserUseTool := tool.NewBrowserUseTool()
		if err := tools.AddTool(browserUseTool); err != nil {
			logger.Fatal("添加BrowserUseTool工具失败: %v", err)
		}
	}

	// 添加FileOperator工具
	if toolsCfg.FileOperator {
		pterm.Debug.Println("  📁 加载 File Operator 工具")
		fileOperatorTool := tool.NewFileOperator()
		if err := tools.AddTool(fileOperatorTool); err != nil {
			logger.Fatal("添加FileOperator工具失败: %v", err)
		}
	}

	pterm.Success.Println("✅ 工具模块加载完成")

	// 创建Manus代理
	pterm.Info.Println("🤖 正在创建 Manus 代理...")
	manusAgent := agent.NewManus("Manus", llmInstance, tools)
	pterm.Success.Println("✅ Manus 代理创建成功")

	// 创建聊天代理
	pterm.Info.Println("💬 正在创建聊天代理...")
	chatAgent := agent.NewChatAgent("ChatAgent", llmInstance)
	pterm.Success.Println("✅ 聊天代理创建成功")

	// 创建分类器代理
	pterm.Info.Println("🧠 正在创建输入分类器...")
	classifierAgent := agent.NewClassifierAgent("Classifier", llmInstance)
	pterm.Success.Println("✅ 输入分类器创建成功")

	// 根据配置创建规划代理
	var planningAgent *agent.PlanningAgent
	if toolsCfg.Planning {
		pterm.Info.Println("📋 正在创建规划代理...")
		planningAgent = agent.NewPlanningAgent("PlanningAgent", llmInstance, tools)
		pterm.Success.Println("✅ 规划代理创建成功")

		// 将Manus代理添加为规划代理的执行器
		planningAgent.AddExecutor("default", manusAgent.ToolCallAgent)
	}

	pterm.Success.Println("🎉 所有代理已准备就绪，开始交互式会话！")
	pterm.Println()

	pterm.Info.Println("欢迎使用GoManus！输入 'exit' 退出程序")
	pterm.Info.Println("🧠 智能分类功能已启用，系统会自动判断您的输入类型：")
	pterm.Info.Println("   💬 聊天模式：日常对话、问答交流")
	pterm.Info.Println("   ⚡ 任务模式：执行具体操作和任务")
	if toolsCfg.Planning {
		pterm.Info.Println("   📋 计划模式：制定复杂的多步骤计划")
	} else {
		pterm.Warning.Println("   📋 计划模式：未启用（需要在配置中开启）")
	}
	pterm.Println()

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动信号处理协程
	go func() {
		<-sigChan
		pterm.Warning.Println("\n⚠️  收到中断信号，正在取消当前任务...")
		cancel() // 取消当前执行的任务
		pterm.Success.Println("✅ 任务已取消")
		os.Exit(0)
	}()

	for {
		// 使用PTerm的交互式输入提示，添加panic恢复机制
		var input string
		func() {
			defer func() {
				if r := recover(); r != nil {
					pterm.Error.Println("⚠️  输入组件出现错误，请重新输入")
					input = ""
				}
			}()
			input, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("").WithTextStyle(pterm.NewStyle(pterm.FgCyan)).Show("🤖 请输入您的问题")
		}()
		input = strings.TrimSpace(input)
		if input == "exit" {
			pterm.Success.Println("👋 再见！感谢使用GoManus！")
			break
		}

		// 检查空输入
		if input == "" {
			pterm.Warning.Println("⚠️  请输入有效的问题或指令")
			continue
		}

		// 处理用户输入
		logger.Info("收到用户输入: %s", input)
		logger.Info("开始处理用户输入...")

		var response string
		var err error

		// 为每个请求创建新的可取消上下文
		requestCtx, requestCancel := context.WithCancel(ctx)
		defer requestCancel()

		// 使用分类器判断输入类型
		pterm.Info.Println("🔍 正在分析输入类型...")
		inputType, classifyErr := classifierAgent.ClassifyInput(requestCtx, input)
		if classifyErr != nil {
			logger.Error("输入分类失败: %v", classifyErr)
			pterm.Warning.Printf("⚠️  输入分类失败，使用默认模式: %v\n", classifyErr)
			inputType = agent.InputTypeTask // 默认为任务模式
		}

		// 显示分类结果
		switch inputType {
		case agent.InputTypeChat:
			pterm.Info.Println("💬 识别为：聊天模式")
		case agent.InputTypeTask:
			pterm.Info.Println("⚡ 识别为：任务模式")
		case agent.InputTypePlan:
			pterm.Info.Println("📋 识别为：计划模式")
		}

		// 根据输入类型选择处理方式
		switch inputType {
		case agent.InputTypePlan:
			// 计划模式
			if toolsCfg.Planning {
				logger.Info("使用规划代理处理计划请求: %s", input)
				pterm.Info.Println("📋 启用规划模式 (按 Ctrl+C 可取消)")
				spinner, _ := pterm.DefaultSpinner.Start("🧠 正在制定计划...")
				response, err = planningAgent.Run(requestCtx, input)
				spinner.Stop()
			} else {
				pterm.Warning.Println("⚠️  规划模式未启用，将使用任务模式处理")
				logger.Info("规划模式未启用，使用任务模式处理请求: %s", input)
				spinner, _ := pterm.DefaultSpinner.Start("⚡ 正在执行任务... (按 Ctrl+C 可取消)")
				response, err = manusAgent.Run(requestCtx, input)
				spinner.Stop()
			}
		case agent.InputTypeTask:
			// 任务模式
			logger.Info("使用任务模式处理请求: %s", input)
			spinner, _ := pterm.DefaultSpinner.Start("⚡ 正在执行任务... (按 Ctrl+C 可取消)")
			response, err = manusAgent.Run(requestCtx, input)
			spinner.Stop()
		case agent.InputTypeChat:
			// 聊天模式
			logger.Info("使用聊天模式处理请求: %s", input)
			spinner, _ := pterm.DefaultSpinner.Start("💬 正在聊天中... (按 Ctrl+C 可取消)")
			response, err = chatAgent.Run(requestCtx, input)
			spinner.Stop()
		default:
			// 默认使用任务模式
			logger.Info("使用默认任务模式处理请求: %s", input)
			spinner, _ := pterm.DefaultSpinner.Start("🤔 正在思考中... (按 Ctrl+C 可取消)")
			response, err = manusAgent.Run(requestCtx, input)
			spinner.Stop()
		}

		if err != nil {
			// 检查是否是上下文取消错误
			if err == context.Canceled {
				pterm.Warning.Println("⚠️  任务已被用户取消")
				continue
			}
			pterm.Error.Printf("❌ 处理消息时出错: %v\n", err)
			continue
		}

		// 输出响应
		logger.Info("处理完成，返回响应: %s", response)
		pterm.DefaultBox.WithTitle("🤖 GoManus 回复").WithTitleTopCenter().WithBoxStyle(pterm.NewStyle(pterm.FgCyan)).Println(response)
		pterm.Println()
	}
}
