package goserversdk

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Client JPush客户端
type Client struct {
	appKey       string
	masterSecret string
	logger       *zap.Logger
	httpClient   *http.Client
	baseURLs     map[string]string
	Push         *PushService
	Advanced     *AdvancedService
	Report       *ReportService
}

// Config 客户端配置
type Config struct {
	AppKey       string        // JPush应用的AppKey
	MasterSecret string        // JPush应用的MasterSecret
	Logger       *zap.Logger   // 日志记录器
	Timeout      time.Duration // HTTP请求超时时间，默认30秒
}

// NewClient 创建JPush客户端
func NewClient(config *Config) (*Client, error) {
	if config.AppKey == "" {
		return nil, NewJPushError(ErrorCodeInvalidAppKey, "AppKey不能为空")
	}
	if config.MasterSecret == "" {
		return nil, NewJPushError(ErrorCodeMissingAuth, "MasterSecret不能为空")
	}
	if config.Logger == nil {
		return nil, NewJPushError(ErrorCodeInvalidParams, "Logger不能为空")
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	client := &Client{
		appKey:       config.AppKey,
		masterSecret: config.MasterSecret,
		logger:       config.Logger,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURLs: map[string]string{
			"push":   "https://api.jpush.cn",
			"device": "https://device.jpush.cn",
			"report": "https://report.jpush.cn",
		},
		Push:     &PushService{},
		Advanced: &AdvancedService{},
		Report:   &ReportService{},
	}

	// 初始化服务
	client.Push.client = client
	client.Advanced.client = client
	client.Report.client = client

	return client, nil
}

// APIResponse API响应结构
type APIResponse struct {
	StatusCode int                    `json:"-"`
	Headers    map[string][]string    `json:"-"`
	Body       map[string]interface{} `json:"-"`
	Error      *JPushError            `json:"error,omitempty"`
}

// makeRequest 发送HTTP请求
func (c *Client) makeRequest(ctx context.Context, method, baseURL, path string, body interface{}) (*APIResponse, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			c.logger.Error("序列化请求体失败", zap.Error(err))
			return nil, NewJPushError(ErrorCodeInvalidJSON, "请求体序列化失败")
		}
		reqBody = bytes.NewBuffer(jsonData)
		c.logger.Debug("请求体", zap.String("body", string(jsonData)))
	}

	url := baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		c.logger.Error("创建HTTP请求失败", zap.Error(err))
		return nil, NewJPushError(ErrorCodeInternalError, "创建HTTP请求失败")
	}

	// 设置认证头
	auth := base64.StdEncoding.EncodeToString([]byte(c.appKey + ":" + c.masterSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	c.logger.Debug("发送HTTP请求",
		zap.String("method", method),
		zap.String("url", url),
		zap.Any("headers", req.Header))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("HTTP请求失败", zap.Error(err))
		return nil, NewJPushError(ErrorCodeTimeout, "HTTP请求失败")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("读取响应体失败", zap.Error(err))
		return nil, NewJPushError(ErrorCodeInternalError, "读取响应体失败")
	}

	c.logger.Debug("收到HTTP响应",
		zap.Int("status_code", resp.StatusCode),
		zap.Any("headers", resp.Header),
		zap.String("body", string(respBody)))

	apiResp := &APIResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}

	// 解析响应体
	if len(respBody) > 0 {
		var bodyMap map[string]interface{}
		if err := json.Unmarshal(respBody, &bodyMap); err != nil {
			c.logger.Error("解析响应体失败", zap.Error(err))
			return nil, NewJPushError(ErrorCodeInvalidJSON, "响应体解析失败")
		}
		apiResp.Body = bodyMap

		// 检查是否有错误
		if errorData, exists := bodyMap["error"]; exists {
			if errorMap, ok := errorData.(map[string]interface{}); ok {
				code := ErrorCodeInternalError
				message := "未知错误"

				if codeFloat, exists := errorMap["code"]; exists {
					if codeVal, ok := codeFloat.(float64); ok {
						code = ErrorCode(int(codeVal))
					}
				}

				if msgStr, exists := errorMap["message"]; exists {
					if msgVal, ok := msgStr.(string); ok {
						message = msgVal
					}
				}

				apiResp.Error = NewJPushError(code, message)
			}
		}
	}

	// 检查HTTP状态码
	if resp.StatusCode >= 400 {
		if apiResp.Error == nil {
			message := fmt.Sprintf("HTTP错误: %d", resp.StatusCode)
			switch resp.StatusCode {
			case 400:
				apiResp.Error = NewJPushError(ErrorCodeInvalidParams, message)
			case 401:
				apiResp.Error = NewJPushError(ErrorCodeInvalidAuth, message)
			case 403:
				apiResp.Error = NewJPushError(ErrorCodeAppKeyBlacklisted, message)
			case 429:
				apiResp.Error = NewJPushError(ErrorCodeRateLimitExceeded, message)
			default:
				apiResp.Error = NewJPushError(ErrorCodeInternalError, message)
			}
		}
		c.logger.Error("API请求失败", zap.Any("error", apiResp.Error))
		return apiResp, apiResp.Error
	}

	return apiResp, nil
}

// makeRequestWithoutContext 发送HTTP请求（不需要context）
func (c *Client) makeRequestWithoutContext(method, baseURL, path string, body interface{}) (*APIResponse, error) {
	return c.makeRequest(context.Background(), method, baseURL, path, body)
}

// makeReportRequest 发送Report API请求
func (c *Client) makeReportRequest(method, path string, body interface{}) (*APIResponse, error) {
	return c.makeRequestWithoutContext(method, c.baseURLs["report"], path, body)
}

// makePushRequest 发送Push API请求
func (c *Client) makePushRequest(method, path string, body interface{}) (*APIResponse, error) {
	return c.makeRequestWithoutContext(method, c.baseURLs["push"], path, body)
}

// makeDeviceRequest 发送Device API请求
func (c *Client) makeDeviceRequest(method, path string, body interface{}) (*APIResponse, error) {
	return c.makeRequestWithoutContext(method, c.baseURLs["device"], path, body)
}

// GetRateLimitInfo 获取频率限制信息
func (c *Client) GetRateLimitInfo(resp *APIResponse) (limit, remaining, reset int) {
	if resp == nil || resp.Headers == nil {
		return 0, 0, 0
	}

	if limitHeaders := resp.Headers["X-Rate-Limit-Limit"]; len(limitHeaders) > 0 {
		fmt.Sscanf(limitHeaders[0], "%d", &limit)
	}

	if remainingHeaders := resp.Headers["X-Rate-Limit-Remaining"]; len(remainingHeaders) > 0 {
		fmt.Sscanf(remainingHeaders[0], "%d", &remaining)
	}

	if resetHeaders := resp.Headers["X-Rate-Limit-Reset"]; len(resetHeaders) > 0 {
		fmt.Sscanf(resetHeaders[0], "%d", &reset)
	}

	return limit, remaining, reset
}