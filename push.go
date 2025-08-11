package goserversdk

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// PushService 推送服务
type PushService struct {
	client *Client
}

// NewPushService 创建推送服务
func NewPushService(client *Client) *PushService {
	return &PushService{client: client}
}

// Push 创建推送
// 向某单个设备或者某设备列表推送一条通知、或者消息
func (s *PushService) Push(req *PushRequest) (*PushResponse, error) {
	s.client.logger.Info("开始创建推送", zap.Any("request", req))

	// 验证必填参数
	if err := s.validatePushRequest(req); err != nil {
		s.client.logger.Error("推送请求参数验证失败", zap.Error(err))
		return nil, err
	}

	resp, err := s.client.makePushRequest(http.MethodPost, "/v3/push", req)
	if err != nil {
		s.client.logger.Error("推送请求失败", zap.Error(err))
		return nil, err
	}

	// 记录频率限制信息
	limit, remaining, reset := s.client.GetRateLimitInfo(resp)
	s.client.logger.Info("推送频率限制信息",
		zap.Int("limit", limit),
		zap.Int("remaining", remaining),
		zap.Int("reset", reset))

	var pushResp PushResponse
	if err := s.parseResponse(resp.Body, &pushResp); err != nil {
		s.client.logger.Error("解析推送响应失败", zap.Error(err))
		return nil, err
	}

	s.client.logger.Info("推送创建成功",
		zap.String("sendno", pushResp.SendNo),
		zap.String("msg_id", pushResp.MsgID))

	return &pushResp, nil
}

// validatePushRequest 验证推送请求参数
func (s *PushService) validatePushRequest(req *PushRequest) error {
	if req == nil {
		return NewJPushError(ErrorCodeInvalidParams, "推送请求不能为空")
	}

	if req.Platform == nil {
		return NewJPushError(ErrorCodeInvalidPlatform, "推送平台不能为空")
	}

	if req.Audience == nil {
		return NewJPushError(ErrorCodeInvalidAudience, "推送目标不能为空")
	}

	// 验证推送目标
	if err := s.validateAudience(req.Audience); err != nil {
		return err
	}

	// 验证推送内容
	if req.Notification == nil && req.Message == nil {
		return NewJPushError(ErrorCodeInvalidParams, "通知和消息至少需要有一个")
	}

	// 验证通知内容
	if req.Notification != nil {
		if err := s.validateNotification(req.Notification); err != nil {
			return err
		}
	}

	// 验证消息内容
	if req.Message != nil {
		if err := s.validateMessage(req.Message); err != nil {
			return err
		}
	}

	return nil
}

// validateAudience 验证推送目标
func (s *PushService) validateAudience(audience *Audience) error {
	hasTarget := false

	if audience.All != nil {
		hasTarget = true
	}

	if len(audience.Tag) > 0 {
		hasTarget = true
		if len(audience.Tag) > 20 {
			return NewJPushError(ErrorCodeInvalidAudience, "标签数量不能超过20个")
		}
	}

	if len(audience.TagAnd) > 0 {
		hasTarget = true
		if len(audience.TagAnd) > 20 {
			return NewJPushError(ErrorCodeInvalidAudience, "标签AND数量不能超过20个")
		}
	}

	if len(audience.TagNot) > 0 {
		hasTarget = true
		if len(audience.TagNot) > 20 {
			return NewJPushError(ErrorCodeInvalidAudience, "标签NOT数量不能超过20个")
		}
	}

	if len(audience.Alias) > 0 {
		hasTarget = true
		if len(audience.Alias) > 1000 {
			return NewJPushError(ErrorCodeInvalidAudience, "别名数量不能超过1000个")
		}
	}

	if len(audience.RegistrationID) > 0 {
		hasTarget = true
		if len(audience.RegistrationID) > 1000 {
			return NewJPushError(ErrorCodeInvalidAudience, "注册ID数量不能超过1000个")
		}
	}

	if len(audience.Segment) > 0 {
		hasTarget = true
		if len(audience.Segment) > 1 {
			return NewJPushError(ErrorCodeInvalidAudience, "用户分群只能指定一个")
		}
	}

	if len(audience.ABTest) > 0 {
		hasTarget = true
		if len(audience.ABTest) > 1 {
			return NewJPushError(ErrorCodeInvalidAudience, "A/B测试只能指定一个")
		}
	}

	if audience.LiveActivityID != nil {
		hasTarget = true
		// 实时活动不能与其他目标组合使用
		if audience.All != nil || len(audience.Tag) > 0 || len(audience.TagAnd) > 0 ||
			len(audience.TagNot) > 0 || len(audience.Alias) > 0 ||
			len(audience.RegistrationID) > 0 || len(audience.Segment) > 0 ||
			len(audience.ABTest) > 0 {
			return NewJPushError(ErrorCodeInvalidAudience, "实时活动不能与其他推送目标组合使用")
		}
	}

	if !hasTarget {
		return NewJPushError(ErrorCodeInvalidAudience, "必须指定至少一个推送目标")
	}

	return nil
}

// validateNotification 验证通知内容
func (s *PushService) validateNotification(notification *Notification) error {
	if notification.Alert == "" && notification.Android == nil &&
		notification.IOS == nil && notification.HMOS == nil &&
		notification.QuickApp == nil {
		return NewJPushError(ErrorCodeInvalidNotification, "通知内容不能为空")
	}

	return nil
}

// validateMessage 验证消息内容
func (s *PushService) validateMessage(message *Message) error {
	if message.MsgContent == "" {
		return NewJPushError(ErrorCodeInvalidMessage, "消息内容不能为空")
	}

	return nil
}

// parseResponse 解析响应
func (s *PushService) parseResponse(body map[string]interface{}, result interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return NewJPushError(ErrorCodeInvalidJSON, "响应序列化失败")
	}

	if err := json.Unmarshal(jsonData, result); err != nil {
		return NewJPushError(ErrorCodeInvalidJSON, "响应解析失败")
	}

	return nil
}