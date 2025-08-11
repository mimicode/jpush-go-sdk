package goserversdk

import "fmt"

// ErrorCode 定义JPush API错误码
type ErrorCode int

const (
	// 通用错误码
	ErrorCodeSuccess           ErrorCode = 0    // 成功
	ErrorCodeInvalidParams     ErrorCode = 1000 // 参数错误
	ErrorCodeMissingAuth       ErrorCode = 1001 // 缺少认证信息
	ErrorCodeInvalidAuth       ErrorCode = 1002 // 认证信息错误
	ErrorCodeInvalidAppKey     ErrorCode = 1003 // AppKey错误
	ErrorCodeInvalidJSON       ErrorCode = 1004 // JSON格式错误
	ErrorCodeTimeout           ErrorCode = 1005 // 请求超时
	ErrorCodeInternalError     ErrorCode = 1006 // 内部错误
	ErrorCodeRateLimitExceeded ErrorCode = 2002 // 频率限制
	ErrorCodeAppKeyBlacklisted ErrorCode = 2003 // AppKey被加入黑名单
	ErrorCodeBroadcastLimit    ErrorCode = 2008 // 广播推送频率限制

	// 推送相关错误码
	ErrorCodeInvalidPlatform    ErrorCode = 3001 // 无效的平台
	ErrorCodeInvalidAudience    ErrorCode = 3002 // 无效的推送目标
	ErrorCodeInvalidNotification ErrorCode = 3003 // 无效的通知内容
	ErrorCodeInvalidMessage     ErrorCode = 3004 // 无效的消息内容
	ErrorCodeInvalidOptions     ErrorCode = 3005 // 无效的推送选项

	// 设备相关错误码
	ErrorCodeInvalidRegistrationID ErrorCode = 7001 // 无效的注册ID
	ErrorCodeInvalidTag            ErrorCode = 7002 // 无效的标签
	ErrorCodeInvalidAlias          ErrorCode = 7003 // 无效的别名
	ErrorCodeTagLimitExceeded      ErrorCode = 7004 // 标签数量超限
	ErrorCodeAliasLimitExceeded    ErrorCode = 7015 // 别名绑定设备数量超限
	ErrorCodeIllegalRegistrationID ErrorCode = 7013 // 非法的注册ID
	ErrorCodeTagOperationFailed    ErrorCode = 7016 // 标签操作失败
)

// JPushError JPush API错误
type JPushError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *JPushError) Error() string {
	return fmt.Sprintf("JPush API Error [%d]: %s", e.Code, e.Message)
}

// NewJPushError 创建JPush错误
func NewJPushError(code ErrorCode, message string) *JPushError {
	return &JPushError{
		Code:    code,
		Message: message,
	}
}

// IsJPushError 判断是否为JPush错误
func IsJPushError(err error) bool {
	_, ok := err.(*JPushError)
	return ok
}

// GetErrorCode 获取错误码
func GetErrorCode(err error) ErrorCode {
	if jpushErr, ok := err.(*JPushError); ok {
		return jpushErr.Code
	}
	return ErrorCodeInternalError
}

// GetErrorCodeFromHTTPStatus 根据HTTP状态码获取错误码
func GetErrorCodeFromHTTPStatus(httpStatus int) ErrorCode {
	switch httpStatus {
	case 200:
		return ErrorCodeSuccess
	case 400:
		return ErrorCodeInvalidParams
	case 401:
		return ErrorCodeInvalidAuth
	case 403:
		return ErrorCodeAppKeyBlacklisted
	case 429:
		return ErrorCodeRateLimitExceeded
	default:
		return ErrorCodeInternalError
	}
}