package goserversdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJPushError(t *testing.T) {
	tests := []struct {
		name    string
		code    ErrorCode
		message string
	}{
		{
			name:    "invalid params error",
			code:    ErrorCodeInvalidParams,
			message: "Invalid parameters",
		},
		{
			name:    "auth error",
			code:    ErrorCodeInvalidAuth,
			message: "Authentication failed",
		},
		{
			name:    "rate limit error",
			code:    ErrorCodeRateLimitExceeded,
			message: "Rate limit exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewJPushError(tt.code, tt.message)
			
			assert.NotNil(t, err)
			assert.Equal(t, tt.code, err.Code)
			assert.Equal(t, tt.message, err.Message)
			assert.Contains(t, err.Error(), tt.message)
			assert.Contains(t, err.Error(), "JPush Error")
		})
	}
}

func TestJPushError_Error(t *testing.T) {
	err := NewJPushError(ErrorCodeInvalidParams, "Test error message")
	errorString := err.Error()
	
	assert.Contains(t, errorString, "JPush Error")
	assert.Contains(t, errorString, "1003")
	assert.Contains(t, errorString, "Test error message")
}

func TestGetErrorCode(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode ErrorCode
	}{
		{
			name:         "jpush error - invalid params",
			err:          NewJPushError(ErrorCodeInvalidParams, "test"),
			expectedCode: ErrorCodeInvalidParams,
		},
		{
			name:         "jpush error - invalid auth",
			err:          NewJPushError(ErrorCodeInvalidAuth, "test"),
			expectedCode: ErrorCodeInvalidAuth,
		},
		{
			name:         "jpush error - invalid app key",
			err:          NewJPushError(ErrorCodeInvalidAppKey, "test"),
			expectedCode: ErrorCodeInvalidAppKey,
		},
		{
			name:         "non-jpush error",
			err:          fmt.Errorf("regular error"),
			expectedCode: ErrorCodeInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := GetErrorCode(tt.err)
			assert.Equal(t, tt.expectedCode, code)
		})
	}
}

func TestGetErrorCodeFromHTTPStatus(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedCode   ErrorCode
	}{
		{
			name:         "bad request",
			statusCode:   http.StatusBadRequest,
			expectedCode: ErrorCodeInvalidParams,
		},
		{
			name:         "unauthorized",
			statusCode:   http.StatusUnauthorized,
			expectedCode: ErrorCodeInvalidAuth,
		},
		{
			name:         "forbidden",
			statusCode:   http.StatusForbidden,
			expectedCode: ErrorCodeAppKeyBlacklisted,
		},
		{
			name:         "too many requests",
			statusCode:   http.StatusTooManyRequests,
			expectedCode: ErrorCodeRateLimitExceeded,
		},
		{
			name:         "internal server error",
			statusCode:   http.StatusInternalServerError,
			expectedCode: ErrorCodeInternalError,
		},
		{
			name:         "bad gateway",
			statusCode:   http.StatusBadGateway,
			expectedCode: ErrorCodeInternalError,
		},
		{
			name:         "service unavailable",
			statusCode:   http.StatusServiceUnavailable,
			expectedCode: ErrorCodeInternalError,
		},
		{
			name:         "gateway timeout",
			statusCode:   http.StatusGatewayTimeout,
			expectedCode: ErrorCodeTimeout,
		},
		{
			name:         "unknown status",
			statusCode:   418, // I'm a teapot
			expectedCode: ErrorCodeInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := GetErrorCodeFromHTTPStatus(tt.statusCode)
			assert.Equal(t, tt.expectedCode, code)
		})
	}
}

func TestErrorConstants(t *testing.T) {
	// 测试所有错误码常量是否定义正确
	assert.Equal(t, 1001, ErrorCodeInvalidAuth)
	assert.Equal(t, 1002, ErrorCodeMissingAuth)
	assert.Equal(t, 1003, ErrorCodeInvalidParams)
	assert.Equal(t, 1004, ErrorCodeInvalidAppKey)
	assert.Equal(t, 1005, ErrorCodeAppKeyBlacklisted)
	assert.Equal(t, 2002, ErrorCodeRateLimitExceeded)
	assert.Equal(t, 1011, ErrorCodeInvalidPlatform)
	assert.Equal(t, 1012, ErrorCodeInvalidAudience)
	assert.Equal(t, 1013, ErrorCodeInvalidNotification)
	assert.Equal(t, 1014, ErrorCodeInvalidMessage)
	assert.Equal(t, 1015, ErrorCodeInvalidJSON)
	assert.Equal(t, 5000, ErrorCodeInternalError)
	assert.Equal(t, 5001, ErrorCodeTimeout)
}

func TestJPushError_IsType(t *testing.T) {
	// 测试错误类型判断
	var err error = NewJPushError(ErrorCodeInvalidParams, "test")
	
	// 测试是否为JPushError类型
	jpushErr, ok := err.(*JPushError)
	assert.True(t, ok)
	assert.NotNil(t, jpushErr)
	
	// 测试错误码匹配
	assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	assert.NotEqual(t, ErrorCodeInvalidAuth, jpushErr.Code)
}