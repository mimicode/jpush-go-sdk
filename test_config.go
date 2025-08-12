package goserversdk

import (
	"bufio"
	"os"
	"strings"

	"go.uber.org/zap"
)

// TestConfig 测试配置结构
type TestConfig struct {
	AppKey       string
	MasterSecret string
}

// LoadTestConfig 从env.test文件加载测试配置
func LoadTestConfig() (*TestConfig, error) {
	file, err := os.Open("env.test")
	if err != nil {
		// 如果文件不存在，返回默认的测试配置
		return &TestConfig{
			AppKey:       "test-key",
			MasterSecret: "test-secret",
		}, nil
	}
	defer file.Close()

	config := &TestConfig{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过注释和空行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析键值对
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "APP_KEY":
			config.AppKey = value
		case "MASTER_SECRET":
			config.MasterSecret = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 如果配置为空，使用默认值
	if config.AppKey == "" {
		config.AppKey = "test-key"
	}
	if config.MasterSecret == "" {
		config.MasterSecret = "test-secret"
	}

	return config, nil
}

// NewTestClient 创建用于测试的客户端
func NewTestClient() (*Client, error) {
	testConfig, err := LoadTestConfig()
	if err != nil {
		return nil, err
	}

	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       testConfig.AppKey,
		MasterSecret: testConfig.MasterSecret,
		Logger:       logger,
	}

	return NewClient(config)
}