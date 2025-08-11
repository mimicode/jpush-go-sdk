package goserversdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// AdvancedService 高级功能服务
type AdvancedService struct {
	client *Client
}

// CIDType CID类型
type CIDType string

const (
	CIDTypePush     CIDType = "push"
	CIDTypeSchedule CIDType = "schedule"
)

// CIDRequest 获取CID请求参数
type CIDRequest struct {
	Count int     `json:"count,omitempty"` // CID数量，VIP应用范围[1,1000]，非VIP应用范围[1,10]
	Type  CIDType `json:"type,omitempty"`  // CID类型，默认为push
}

// CIDResponse 获取CID响应
type CIDResponse struct {
	CIDList []string `json:"cidlist"` // CID列表
}

// QuotaData 厂商配额数据
type QuotaData struct {
	XiaomiQuota *VendorQuota `json:"xiaomi_quota,omitempty"` // 小米配额
	OppoQuota   *VendorQuota `json:"oppo_quota,omitempty"`   // OPPO配额
	VivoQuota   *VivoQuota   `json:"vivo_quota,omitempty"`   // vivo配额
}

// VendorQuota 厂商配额信息
type VendorQuota struct {
	Operation *QuotaInfo `json:"operation,omitempty"` // 运营消息配额
}

// VivoQuota vivo配额信息
type VivoQuota struct {
	System    *QuotaInfo `json:"system,omitempty"`    // 系统消息配额
	Operation *QuotaInfo `json:"operation,omitempty"` // 运营消息配额
}

// QuotaInfo 配额信息
type QuotaInfo struct {
	Total int `json:"total"` // 可用总额度，开通不限量时返回-1
	Used  int `json:"used"`  // 已使用额度，开通不限量时返回-1
}

// QuotaResponse 厂商配额查询响应
type QuotaResponse struct {
	Code    int        `json:"code"`    // 返回码，0表示成功
	Message string     `json:"message"` // 返回消息
	Data    *QuotaData `json:"data"`    // 配额数据
}

// FileAudience 文件推送目标
type FileAudience struct {
	File *FileTarget `json:"file,omitempty"` // 文件目标
}

// FileTarget 文件目标
type FileTarget struct {
	FileID string `json:"file_id"` // 文件唯一标识
}

// FilePushRequest 文件推送请求
type FilePushRequest struct {
	Platform     interface{}   `json:"platform"`              // 推送平台
	Audience     *FileAudience `json:"audience"`              // 推送目标（仅支持file）
	Notification *Notification `json:"notification,omitempty"` // 通知内容
	Message      *Message      `json:"message,omitempty"`      // 自定义消息
	SMSMessage   *SMSMessage   `json:"sms_message,omitempty"`  // 短信补充
	Options      *Options      `json:"options,omitempty"`      // 推送选项
	Callback     *Callback     `json:"callback,omitempty"`     // 回调参数
}

// GetCID 获取推送唯一标识符
// count: CID数量，VIP应用范围[1,1000]，非VIP应用范围[1,10]
// cidType: CID类型，默认为push
func (s *AdvancedService) GetCID(count int, cidType CIDType) (*CIDResponse, error) {
	if count <= 0 {
		return nil, NewJPushError(ErrorCodeInvalidParams, "count must be greater than 0")
	}

	// 构建查询参数
	params := make(map[string]string)
	params["count"] = strconv.Itoa(count)
	if cidType != "" {
		params["type"] = string(cidType)
	}

	// 构建URL
	url := "/v3/push/cid"
	if len(params) > 0 {
		var paramPairs []string
		for k, v := range params {
			paramPairs = append(paramPairs, fmt.Sprintf("%s=%s", k, v))
		}
		url += "?" + strings.Join(paramPairs, "&")
	}

	resp, err := s.client.makePushRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var cidResp CIDResponse
	bodyBytes, _ := json.Marshal(resp.Body)
	if err := json.Unmarshal(bodyBytes, &cidResp); err != nil {
		return nil, NewJPushError(ErrorCodeInvalidJSON, fmt.Sprintf("failed to parse CID response: %v", err))
	}

	return &cidResp, nil
}

// ValidatePush 推送校验API，验证推送调用是否能够成功，不向用户发送任何消息
func (s *AdvancedService) ValidatePush(req *PushRequest) (*PushResponse, error) {
	if req == nil {
		return nil, NewJPushError(ErrorCodeInvalidParams, "push request cannot be nil")
	}

	// 验证推送请求参数
	if err := s.validatePushRequest(req); err != nil {
		return nil, err
	}

	resp, err := s.client.makePushRequest(http.MethodPost, "/v3/push/validate", req)
	if err != nil {
		return nil, err
	}

	var pushResp PushResponse
	bodyBytes, _ := json.Marshal(resp.Body)
	if err := json.Unmarshal(bodyBytes, &pushResp); err != nil {
		return nil, NewJPushError(ErrorCodeInvalidJSON, fmt.Sprintf("failed to parse push response: %v", err))
	}

	return &pushResp, nil
}

// CancelPush 推送撤销API，撤销指定的推送消息
// msgID: 推送消息ID
func (s *AdvancedService) CancelPush(msgID string) error {
	if msgID == "" {
		return NewJPushError(ErrorCodeInvalidParams, "message ID cannot be empty")
	}

	url := fmt.Sprintf("/v3/push/%s", msgID)
	_, err := s.client.makePushRequest(http.MethodDelete, url, nil)
	return err
}

// GetVendorQuota 查询厂商配额信息
func (s *AdvancedService) GetVendorQuota() (*QuotaResponse, error) {
	resp, err := s.client.makePushRequest(http.MethodGet, "/v3/push/quota", nil)
	if err != nil {
		return nil, err
	}

	var quotaResp QuotaResponse
	bodyBytes, _ := json.Marshal(resp.Body)
	if err := json.Unmarshal(bodyBytes, &quotaResp); err != nil {
		return nil, NewJPushError(ErrorCodeInvalidJSON, fmt.Sprintf("failed to parse quota response: %v", err))
	}

	return &quotaResp, nil
}

// PushByFile 文件推送API，通过文件ID进行推送
func (s *AdvancedService) PushByFile(req *FilePushRequest) (*PushResponse, error) {
	if req == nil {
		return nil, NewJPushError(ErrorCodeInvalidParams, "file push request cannot be nil")
	}

	// 验证文件推送请求参数
	if err := s.validateFilePushRequest(req); err != nil {
		return nil, err
	}

	resp, err := s.client.makePushRequest(http.MethodPost, "/v3/push/file", req)
	if err != nil {
		return nil, err
	}

	var pushResp PushResponse
	bodyBytes, _ := json.Marshal(resp.Body)
	if err := json.Unmarshal(bodyBytes, &pushResp); err != nil {
		return nil, NewJPushError(ErrorCodeInvalidJSON, fmt.Sprintf("failed to parse file push response: %v", err))
	}

	return &pushResp, nil
}

// validatePushRequest 验证推送请求参数
func (s *AdvancedService) validatePushRequest(req *PushRequest) error {
	if req.Platform == nil {
		return NewJPushError(ErrorCodeInvalidParams, "platform is required")
	}

	if req.Audience == nil {
		return NewJPushError(ErrorCodeInvalidParams, "audience is required")
	}

	if req.Notification == nil && req.Message == nil {
		return NewJPushError(ErrorCodeInvalidParams, "at least one of notification or message is required")
	}

	return nil
}

// validateFilePushRequest 验证文件推送请求参数
func (s *AdvancedService) validateFilePushRequest(req *FilePushRequest) error {
	if req.Platform == nil {
		return NewJPushError(ErrorCodeInvalidParams, "platform is required")
	}

	if req.Audience == nil || req.Audience.File == nil {
		return NewJPushError(ErrorCodeInvalidParams, "file audience is required")
	}

	if req.Audience.File.FileID == "" {
		return NewJPushError(ErrorCodeInvalidParams, "file_id is required")
	}

	if req.Notification == nil && req.Message == nil {
		return NewJPushError(ErrorCodeInvalidParams, "at least one of notification or message is required")
	}

	return nil
}

// SetFileAudience 设置文件推送目标
func (req *FilePushRequest) SetFileAudience(fileID string) *FilePushRequest {
	req.Audience = &FileAudience{
		File: &FileTarget{
			FileID: fileID,
		},
	}
	return req
}

// SetPlatform 设置推送平台
func (req *FilePushRequest) SetPlatform(platform interface{}) *FilePushRequest {
	req.Platform = platform
	return req
}

// SetNotification 设置通知内容
func (req *FilePushRequest) SetNotification(notification *Notification) *FilePushRequest {
	req.Notification = notification
	return req
}

// SetMessage 设置自定义消息
func (req *FilePushRequest) SetMessage(message *Message) *FilePushRequest {
	req.Message = message
	return req
}

// SetSMSMessage 设置短信补充
func (req *FilePushRequest) SetSMSMessage(smsMessage *SMSMessage) *FilePushRequest {
	req.SMSMessage = smsMessage
	return req
}

// SetOptions 设置推送选项
func (req *FilePushRequest) SetOptions(options *Options) *FilePushRequest {
	req.Options = options
	return req
}

// SetCallback 设置回调参数
func (req *FilePushRequest) SetCallback(callback *Callback) *FilePushRequest {
	req.Callback = callback
	return req
}

// NewFilePushRequest 创建新的文件推送请求
func NewFilePushRequest() *FilePushRequest {
	return &FilePushRequest{}
}