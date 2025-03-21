package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

// LLMConfig 表示LLM的配置
type LLMConfig struct {
	Model       string  `mapstructure:"model"`
	BaseURL     string  `mapstructure:"base_url"`
	APIKey      string  `mapstructure:"api_key"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"`
}

// ToolsConfig 表示工具的配置
type ToolsConfig struct {
	Terminate       bool `mapstructure:"terminate"`
	GoogleSearch    bool `mapstructure:"google_search"`
	ZhihuSearch     bool `mapstructure:"zhihu_search"`
	BaiduBaikeSearch bool `mapstructure:"baidu_baike_search"`
	WikipediaSearch bool `mapstructure:"wikipedia_search"`
	BrowserUse      bool `mapstructure:"browser_use"`
	FileOperator    bool `mapstructure:"file_operator"`
	Planning        bool `mapstructure:"planning"`
}

// Config 表示应用程序的配置
type Config struct {
	LLM      LLMConfig            `mapstructure:"llm"`
	LLMTypes map[string]LLMConfig `mapstructure:"llm_types"`
	Tools    ToolsConfig          `mapstructure:"tools"`
}

var (
	config Config
	once   sync.Once
)

// LoadConfig 从指定路径加载配置
func LoadConfig(configPath string) (*Config, error) {
	var err error
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		
		// 如果提供了配置路径，则使用它
		if configPath != "" {
			viper.AddConfigPath(configPath)
		} else {
			// 默认配置路径
			viper.AddConfigPath("./config")
			viper.AddConfigPath("../config")
			viper.AddConfigPath("../../config")
			
			// 获取可执行文件所在目录
			execPath, execErr := os.Executable()
			if execErr == nil {
				execDir := filepath.Dir(execPath)
				viper.AddConfigPath(filepath.Join(execDir, "config"))
			}
		}
		
		// 读取配置文件
		if readErr := viper.ReadInConfig(); readErr != nil {
			err = fmt.Errorf("读取配置文件失败: %w", readErr)
			return
		}
		
		// 解析配置
		if unmarshalErr := viper.Unmarshal(&config); unmarshalErr != nil {
			err = fmt.Errorf("解析配置失败: %w", unmarshalErr)
			return
		}
	})
	
	if err != nil {
		return nil, err
	}
	
	return &config, nil
}

// GetLLMConfig 获取指定名称的LLM配置
func GetLLMConfig(name string) (*LLMConfig, error) {
	cfg, err := LoadConfig("")
	if err != nil {
		return nil, err
	}
	
	// 如果请求的是默认配置
	if name == "" || name == "default" {
		// 使用顶级LLM配置
		return &cfg.LLM, nil
	}
	
	// 查找特定名称的配置
	if llmConfig, exists := cfg.LLMTypes[name]; exists {
		return &llmConfig, nil
	}
	
	return nil, fmt.Errorf("未找到名为 %s 的LLM配置", name)
}

// GetToolsConfig 获取工具配置
func GetToolsConfig() (*ToolsConfig, error) {
	cfg, err := LoadConfig("")
	if err != nil {
		return nil, err
	}
	
	return &cfg.Tools, nil
}
