package errors

import (
	"fmt"
	"time"
)

// ErrorType 定义错误类型
type ErrorType string

const (
	// NetworkError 网络错误
	NetworkError ErrorType = "network_error"

	// ConfigError 配置错误
	ConfigError ErrorType = "config_error"

	// APIError API调用错误
	APIError ErrorType = "api_error"

	// ParseError 解析错误
	ParseError ErrorType = "parse_error"

	// TimeoutError 超时错误
	TimeoutError ErrorType = "timeout_error"

	// InternalError 内部错误
	InternalError ErrorType = "internal_error"
)

// SpeedTestError 自定义错误结构
type SpeedTestError struct {
	Type    ErrorType  // 错误类型
	Code    string     // 错误代码
	Message string     // 错误消息
	Detail  string     // 详细信息
	Time    time.Time  // 错误发生时间
	Cause   error      // 原始错误
}

// Error 实现error接口
func (e *SpeedTestError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("[%s:%s] %s", e.Type, e.Code, e.Message)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap 返回原始错误
func (e *SpeedTestError) Unwrap() error {
	return e.Cause
}

// WithDetail 添加详细信息
func (e *SpeedTestError) WithDetail(detail string) *SpeedTestError {
	e.Detail = detail
	return e
}

// NewNetworkError 创建网络错误
func NewNetworkError(message string, cause error) *SpeedTestError {
	return &SpeedTestError{
		Type:    NetworkError,
		Message: message,
		Cause:   cause,
		Time:    time.Now(),
	}
}

// NewConfigError 创建配置错误
func NewConfigError(message string, cause error) *SpeedTestError {
	return &SpeedTestError{
		Type:    ConfigError,
		Message: message,
		Cause:   cause,
		Time:    time.Now(),
	}
}

// NewAPIError 创建API错误
func NewAPIError(code, message string, cause error) *SpeedTestError {
	return &SpeedTestError{
		Type:    APIError,
		Code:    code,
		Message: message,
		Cause:   cause,
		Time:    time.Now(),
	}
}

// NewParseError 创建解析错误
func NewParseError(message string, cause error) *SpeedTestError {
	return &SpeedTestError{
		Type:    ParseError,
		Message: message,
		Cause:   cause,
		Time:    time.Now(),
	}
}

// NewTimeoutError 创建超时错误
func NewTimeoutError(message string, cause error) *SpeedTestError {
	return &SpeedTestError{
		Type:    TimeoutError,
		Message: message,
		Cause:   cause,
		Time:    time.Now(),
	}
}

// NewInternalError 创建内部错误
func NewInternalError(message string, cause error) *SpeedTestError {
	return &SpeedTestError{
		Type:    InternalError,
		Message: message,
		Cause:   cause,
		Time:    time.Now(),
	}
}

// IsSpeedTestError 检查是否为SpeedTestError类型错误
func IsSpeedTestError(err error) bool {
	_, ok := err.(*SpeedTestError)
	return ok
}

// GetType 获取错误类型
func GetType(err error) ErrorType {
	if stErr, ok := err.(*SpeedTestError); ok {
		return stErr.Type
	}
	return InternalError
}