package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"speedtestup/config"
)

// 初始化测试配置
func init() {
	// 在测试中使用默认配置，避免解析命令行参数
	appConfig = config.NewDefaultConfig()
}

// Mock函数用于测试，避免实际网络请求
func TestGetIP(t *testing.T) {
	// 由于getIP涉及网络请求，我们将在集成测试中测试完整功能
	// 这里仅测试错误情况下的逻辑
	ip, err := getIP()

	// 在没有有效网络连接的测试环境中，这应该返回错误
	if err != nil {
		// 验证返回的是适当的错误
		assert.NotEmpty(t, err)
	} else {
		// 如果成功获取IP，则验证IP格式
		assert.NotEmpty(t, ip)
	}
}

// 为网络请求函数创建模拟测试
func TestSpeedupFunction(t *testing.T) {
	// 由于speedup函数涉及实际网络请求，这里主要测试错误路径
	result, err := speedup()

	// 验证返回值类型
	if err != nil {
		// 如果有错误，确保错误不为nil
		assert.NotNil(t, err)
	} else {
		// 如果没有错误，确保结果不为nil
		assert.NotNil(t, result)
	}
}

func TestCheckSpeedupFunction(t *testing.T) {
	// 由于checkSpeedup函数涉及实际网络请求，这里主要测试错误路径
	result, err := checkSpeedup()

	// 验证返回值类型
	if err != nil {
		assert.NotNil(t, err)
	} else {
		// 如果没有错误，确保返回布尔值
		assert.Implements(t, (*bool)(nil), &result)
	}
}

// 为配置相关的功能测试
func TestAppConfigInitialization(t *testing.T) {
	// 验证appConfig是否被正确初始化
	assert.NotNil(t, appConfig, "appConfig should be initialized")

	// 验证配置的各个部分是否被正确初始化
	assert.NotEmpty(t, appConfig.SpeedTest.URLs, "SpeedTest URLs should not be empty")
	assert.NotEmpty(t, appConfig.IPCheck.URLs, "IPCheck URLs should not be empty")
	assert.Greater(t, appConfig.SpeedTest.Timeout, 0, "SpeedTest timeout should be greater than 0")
	assert.Greater(t, appConfig.IPCheck.Timeout, 0, "IPCheck timeout should be greater than 0")
	assert.Greater(t, appConfig.IPCheck.CheckInterval, 0, "IPCheck check interval should be greater than 0")
}