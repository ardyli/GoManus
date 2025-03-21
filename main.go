package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"gomanus/internal/agent"
	"gomanus/internal/config"
	"gomanus/internal/llm"
	"gomanus/internal/tool"
	"gomanus/pkg/logger"
)

func main() {
	// 设置日志级别
	logger.SetLevel(logger.LevelInfo)

	// 从配置文件加载配置
	logger.Info("正在加载配置...")
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		logger.Fatal("加载配置失败: %v", err)
	}
	logger.Info("配置加载成功: 使用模型 %s", cfg.LLM.Model)

	// 获取工具配置
	toolsCfg, err := config.GetToolsConfig()
	if err != nil {
		logger.Fatal("获取工具配置失败: %v", err)
	}

	// 创建LLM实例
	logger.Info("正在初始化语言模型...")
	llmInstance, err := llm.NewLLM("") // 使用默认配置
	if err != nil {
		logger.Fatal("初始化语言模型失败: %v", err)
	}
	logger.Info("语言模型初始化成功:\n %s \n %s \n %s \n %s \n %s", llmInstance.Model, llmInstance.BaseURL, llmInstance.MaxTokens, llmInstance.Temperature, llmInstance.APIKey)

	// 创建工具集合
	logger.Info("正在初始化工具集合...")
	tools := tool.NewToolCollection()

	// 根据配置添加工具
	// 添加Terminate工具
	if toolsCfg.Terminate {
		logger.Info("正在添加Terminate工具...")
		terminateTool := tool.NewTerminate()
		if err := tools.AddTool(terminateTool); err != nil {
			logger.Fatal("添加Terminate工具失败: %v", err)
		}
	}

	// 添加GoogleSearch工具
	if toolsCfg.GoogleSearch {
		logger.Info("正在添加GoogleSearch工具...")
		googleSearchTool := tool.NewGoogleSearch()
		if err := tools.AddTool(googleSearchTool); err != nil {
			logger.Fatal("添加GoogleSearch工具失败: %v", err)
		}
	}

	// 添加ZhihuSearch工具
	if toolsCfg.ZhihuSearch {
		logger.Info("正在添加ZhihuSearch工具...")
		zhihuSearchTool := tool.NewZhihuSearch()
		if err := tools.AddTool(zhihuSearchTool); err != nil {
			logger.Fatal("添加ZhihuSearch工具失败: %v", err)
		}
	}

	// 添加BaiduBaikeSearch工具
	if toolsCfg.BaiduBaikeSearch {
		logger.Info("正在添加BaiduBaikeSearch工具...")
		baiduBaikeSearchTool := tool.NewBaiduBaikeSearch()
		if err := tools.AddTool(baiduBaikeSearchTool); err != nil {
			logger.Fatal("添加BaiduBaikeSearch工具失败: %v", err)
		}
	}

	// 添加WikipediaSearch工具
	if toolsCfg.WikipediaSearch {
		logger.Info("正在添加WikipediaSearch工具...")
		wikipediaSearchTool := tool.NewWikipediaSearch()
		if err := tools.AddTool(wikipediaSearchTool); err != nil {
			logger.Fatal("添加WikipediaSearch工具失败: %v", err)
		}
	}

	// 添加BrowserUseTool工具
	if toolsCfg.BrowserUse {
		logger.Info("正在添加BrowserUseTool工具...")
		browserUseTool := tool.NewBrowserUseTool()
		if err := tools.AddTool(browserUseTool); err != nil {
			logger.Fatal("添加BrowserUseTool工具失败: %v", err)
		}
	}

	// 添加FileOperator工具
	if toolsCfg.FileOperator {
		logger.Info("正在添加FileOperator工具...")
		fileOperatorTool := tool.NewFileOperator()
		if err := tools.AddTool(fileOperatorTool); err != nil {
			logger.Fatal("添加FileOperator工具失败: %v", err)
		}
	}

	// 创建Manus代理
	logger.Info("正在创建Manus代理...")
	manusAgent := agent.NewManus("Manus", llmInstance, tools)

	// 根据配置创建规划代理
	var planningAgent *agent.PlanningAgent
	if toolsCfg.Planning {
		logger.Info("正在创建规划代理...")
		planningAgent = agent.NewPlanningAgent("PlanningAgent", llmInstance, tools)

		// 将Manus代理添加为规划代理的执行器
		planningAgent.AddExecutor("default", manusAgent.ToolCallAgent)
	}

	logger.Info("代理已准备就绪，开始交互式会话...")

	fmt.Println("欢迎使用GoManus！输入'exit'退出。")
	if toolsCfg.Planning {
		fmt.Println("提示: 输入'plan:你的请求'使用规划模式，或直接输入请求使用普通模式。")
	}

	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "exit" {
			fmt.Println("再见！")
			break
		}

		// 处理用户输入
		logger.Info("收到用户输入: %s", input)
		logger.Info("开始处理用户输入...")

		var response string
		var err error

		// 检查是否使用规划模式
		if strings.HasPrefix(input, "plan:") && toolsCfg.Planning {
			// 提取规划请求
			planRequest := strings.TrimPrefix(input, "plan:")
			planRequest = strings.TrimSpace(planRequest)

			logger.Info("使用规划模式处理请求: %s", planRequest)
			response, err = planningAgent.Run(ctx, planRequest)
		} else {
			// 使用普通模式或规划模式未启用
			if strings.HasPrefix(input, "plan:") && !toolsCfg.Planning {
				fmt.Println("规划模式未启用，将使用普通模式处理请求")
			}
			logger.Info("使用普通模式处理请求")
			response, err = manusAgent.Run(ctx, input)
		}

		if err != nil {
			logger.Error("处理消息失败: %v", err)
			fmt.Printf("处理消息时出错: %v\n", err)
			continue
		}

		// 输出响应
		logger.Info("处理完成，返回响应: %s", response)
		fmt.Println(response)
	}
}
