package goserversdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ReportService 统计服务
type ReportService struct {
	client *Client
}

// ReceivedDetailResponse 送达统计详情响应
type ReceivedDetailResponse struct {
	MsgID                   string `json:"msg_id"`                     // 消息ID
	JPushReceived           *int   `json:"jpush_received"`             // 极光通道用户送达数
	AndroidPNSSent          *int   `json:"android_pns_sent"`           // Android厂商用户推送到厂商服务器成功数
	AndroidPNSReceived      *int   `json:"android_pns_received"`       // Android厂商用户推送达到设备数
	IOSAPNSSent             *int   `json:"ios_apns_sent"`              // iOS通知推送到APNs成功
	IOSAPNSReceived         *int   `json:"ios_apns_received"`          // iOS通知送达到设备并成功展示
	IOSMsgReceived          *int   `json:"ios_msg_received"`           // iOS自定义消息送达数
	LiveActivitySent        *int   `json:"live_acivity_sent"`          // 实时活动消息推送到APNs成功的用户数量
	LiveActivityReceived    *int   `json:"live_acivity_received"`      // 实时活动消息送达成功的用户数量
	WPMPNSSent              *int   `json:"wp_mpns_sent"`               // WP推送数
	QuickAppJPushReceived   *int   `json:"quickapp_jpush_received"`    // 快应用推送走极光通道送达设备成功的用户数量
	QuickAppPNSSent         *int   `json:"quickapp_pns_sent"`          // 快应用推送走厂商通道请求成功的用户数量
	HMOSHMPNSReceived       *int   `json:"hmos_hmpns_received"`        // 鸿蒙通知送达到设备数
	HMOSHMPNSSent           *int   `json:"hmos_hmpns_sent"`            // 鸿蒙通知推送到厂商服务器成功数
	HMOSMsgReceived         *int   `json:"hmos_msg_received"`          // 鸿蒙自定义消息送达到设备数
	HMOSMsgSent             *int   `json:"hmos_msg_sent"`              // 鸿蒙自定义消息推送到厂商服务器成功数
}

// ReceivedResponse 送达统计响应（旧接口）
type ReceivedResponse struct {
	MsgID             string `json:"msg_id"`               // 消息ID
	AndroidReceived   *int   `json:"android_received"`     // Android送达
	IOSAPNSSent       *int   `json:"ios_apns_sent"`        // iOS通知推送到APNs成功
	IOSAPNSReceived   *int   `json:"ios_apns_received"`    // iOS通知送达到设备并成功展示
	IOSMsgReceived    *int   `json:"ios_msg_received"`     // iOS自定义消息送达数
	WPMPNSSent        *int   `json:"wp_mpns_sent"`         // WP推送数
}

// MessageStatusRequest 送达状态查询请求
type MessageStatusRequest struct {
	MsgID           int64    `json:"msg_id"`            // 消息ID
	RegistrationIDs []string `json:"registration_ids"`  // 设备注册ID列表，最多1000个
	Date            string   `json:"date,omitempty"`    // 查询日期，格式yyyy-mm-dd，默认当天
}

// MessageStatusResponse 送达状态查询响应
type MessageStatusResponse map[string]MessageStatus

// MessageStatus 消息状态
type MessageStatus struct {
	Status int `json:"status"` // 状态：0-送达，1-未送达，2-registration_id不属于该应用，3-不是该条message的推送目标，4-系统异常
}

// ChannelStats 通道统计数据
type ChannelStats struct {
	Target   int `json:"target"`   // 目标数量
	Sent     int `json:"sent"`     // 发送数量
	Received int `json:"received"` // 接收数量
	Display  int `json:"display"`  // 展示数量
	Click    int `json:"click"`    // 点击数量
}

// AndroidSubChannels Android子通道统计
type AndroidSubChannels struct {
	JGAndroid *ChannelStats `json:"jg_android,omitempty"` // 极光Android通道
	Huawei    *ChannelStats `json:"huawei,omitempty"`     // 华为通道
	Xiaomi    *ChannelStats `json:"xiaomi,omitempty"`     // 小米通道
	Oppo      *ChannelStats `json:"oppo,omitempty"`       // OPPO通道
	Vivo      *ChannelStats `json:"vivo,omitempty"`       // vivo通道
	Meizu     *ChannelStats `json:"meizu,omitempty"`      // 魅族通道
	FCM       *ChannelStats `json:"fcm,omitempty"`        // FCM通道
	Asus      *ChannelStats `json:"asus,omitempty"`       // 华硕通道
	Tuibida   *ChannelStats `json:"tuibida,omitempty"`    // 推必达通道
	Honor     *ChannelStats `json:"honor,omitempty"`      // 荣耀通道
	Nio       *ChannelStats `json:"nio,omitempty"`        // 蔚来通道
}

// IOSSubChannels iOS子通道统计
type IOSSubChannels struct {
	VOIP  *ChannelStats `json:"voip,omitempty"`   // VOIP通道
	APNS  *ChannelStats `json:"apns,omitempty"`   // APNS通道
	JGIOS *ChannelStats `json:"jg_ios,omitempty"` // 极光iOS通道
}

// HMOSSubChannels 鸿蒙子通道统计
type HMOSSubChannels struct {
	HMPNS  *ChannelStats `json:"hmpns,omitempty"`   // 鸿蒙推送服务
	JGHMOS *ChannelStats `json:"jg_hmos,omitempty"` // 极光鸿蒙通道
}

// QuickAppSubChannels 快应用子通道统计
type QuickAppSubChannels struct {
	QuickJG     *ChannelStats `json:"quick_jg,omitempty"`     // 快应用极光通道
	QuickHuawei *ChannelStats `json:"quick_huawei,omitempty"` // 快应用华为通道
	QuickXiaomi *ChannelStats `json:"quick_xiaomi,omitempty"` // 快应用小米通道
	QuickOppo   *ChannelStats `json:"quick_oppo,omitempty"`   // 快应用OPPO通道
}

// NotificationStats 通知统计
type NotificationStats struct {
	Target         int                  `json:"target"`           // 目标数量
	Sent           int                  `json:"sent"`             // 发送数量
	Received       int                  `json:"received"`         // 接收数量
	Display        int                  `json:"display"`          // 展示数量
	Click          int                  `json:"click"`            // 点击数量
	SubAndroid     *AndroidSubChannels  `json:"sub_android"`      // Android子通道
	SubIOS         *IOSSubChannels      `json:"sub_ios"`          // iOS子通道
	SubHMOS        *HMOSSubChannels     `json:"sub_hmos"`         // 鸿蒙子通道
	SubQuickApp    *QuickAppSubChannels `json:"sub_quickapp"`     // 快应用子通道
}

// MessageStats 自定义消息统计
type MessageStats struct {
	Target         int                  `json:"target"`           // 目标数量
	Sent           int                  `json:"sent"`             // 发送数量
	Received       int                  `json:"received"`         // 接收数量
	Display        int                  `json:"display"`          // 展示数量
	Click          int                  `json:"click"`            // 点击数量
	SubAndroid     *AndroidSubChannels  `json:"sub_android"`      // Android子通道
	SubIOS         *IOSSubChannels      `json:"sub_ios"`          // iOS子通道
	SubHMOS        *HMOSSubChannels     `json:"sub_hmos"`         // 鸿蒙子通道
	SubQuickApp    *QuickAppSubChannels `json:"sub_quickapp"`     // 快应用子通道
}

// InAppStats 应用内提醒统计
type InAppStats struct {
	Target         int                  `json:"target"`           // 目标数量
	Sent           int                  `json:"sent"`             // 发送数量
	Received       int                  `json:"received"`         // 接收数量
	Display        int                  `json:"display"`          // 展示数量
	Click          int                  `json:"click"`            // 点击数量
	SubAndroid     *AndroidSubChannels  `json:"sub_android"`      // Android子通道
	SubIOS         *IOSSubChannels      `json:"sub_ios"`          // iOS子通道
}

// MessageDetailStats 消息详细统计
type MessageDetailStats struct {
	Notification *NotificationStats `json:"notification,omitempty"` // 通知统计
	Message      *MessageStats      `json:"message,omitempty"`      // 自定义消息统计
	InApp        *InAppStats        `json:"inapp,omitempty"`        // 应用内提醒统计
}

// MessageDetailResponse 消息统计详情响应
type MessageDetailResponse struct {
	MsgID   string              `json:"msg_id"`  // 消息ID
	Details *MessageDetailStats `json:"details"` // 详细统计数据
}

// GetReceivedDetail 获取送达统计详情
// msgIDs: 消息ID列表，最多支持100个
func (s *ReportService) GetReceivedDetail(msgIDs []string) ([]ReceivedDetailResponse, error) {
	if len(msgIDs) == 0 {
		return nil, NewJPushError(ErrorCodeInvalidParams, "msg_ids cannot be empty")
	}

	if len(msgIDs) > 100 {
		return nil, NewJPushError(ErrorCodeInvalidParams, "msg_ids cannot exceed 100")
	}

	url := fmt.Sprintf("/v3/received/detail?msg_ids=%s", strings.Join(msgIDs, ","))
	
	// 使用report域名
	resp, err := s.client.makeReportRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var receivedResp []ReceivedDetailResponse
	bodyBytes, _ := json.Marshal(resp.Body)
	if err := json.Unmarshal(bodyBytes, &receivedResp); err != nil {
		return nil, NewJPushError(ErrorCodeInvalidJSON, fmt.Sprintf("failed to parse received detail response: %v", err))
	}

	return receivedResp, nil
}

// GetReceived 获取送达统计（旧接口）
// msgIDs: 消息ID列表，最多支持100个
func (s *ReportService) GetReceived(msgIDs []string) ([]ReceivedResponse, error) {
	if len(msgIDs) == 0 {
		return nil, NewJPushError(ErrorCodeInvalidParams, "msg_ids cannot be empty")
	}

	if len(msgIDs) > 100 {
		return nil, NewJPushError(ErrorCodeInvalidParams, "msg_ids cannot exceed 100")
	}

	url := fmt.Sprintf("/v3/received?msg_ids=%s", strings.Join(msgIDs, ","))
	
	// 使用report域名
	resp, err := s.client.makeReportRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var receivedResp []ReceivedResponse
	bodyBytes, _ := json.Marshal(resp.Body)
	if err := json.Unmarshal(bodyBytes, &receivedResp); err != nil {
		return nil, NewJPushError(ErrorCodeInvalidJSON, fmt.Sprintf("failed to parse received response: %v", err))
	}

	return receivedResp, nil
}

// GetMessageStatus 查询消息送达状态（VIP功能）
func (s *ReportService) GetMessageStatus(req *MessageStatusRequest) (MessageStatusResponse, error) {
	if req == nil {
		return nil, NewJPushError(ErrorCodeInvalidParams, "request cannot be nil")
	}

	if req.MsgID == 0 {
		return nil, NewJPushError(ErrorCodeInvalidParams, "msg_id is required")
	}

	if len(req.RegistrationIDs) == 0 {
		return nil, NewJPushError(ErrorCodeInvalidParams, "registration_ids cannot be empty")
	}

	if len(req.RegistrationIDs) > 1000 {
		return nil, NewJPushError(ErrorCodeInvalidParams, "registration_ids cannot exceed 1000")
	}

	// 使用report域名
	resp, err := s.client.makeReportRequest(http.MethodPost, "/v3/status/message", req)
	if err != nil {
		return nil, err
	}

	var statusResp MessageStatusResponse
	bodyBytes, _ := json.Marshal(resp.Body)
	if err := json.Unmarshal(bodyBytes, &statusResp); err != nil {
		return nil, NewJPushError(ErrorCodeInvalidJSON, fmt.Sprintf("failed to parse message status response: %v", err))
	}

	return statusResp, nil
}

// GetMessageDetail 获取消息统计详情（VIP功能）
// msgIDs: 消息ID列表，最多支持100个
func (s *ReportService) GetMessageDetail(msgIDs []string) ([]MessageDetailResponse, error) {
	if len(msgIDs) == 0 {
		return nil, NewJPushError(ErrorCodeInvalidParams, "msg_ids cannot be empty")
	}

	if len(msgIDs) > 100 {
		return nil, NewJPushError(ErrorCodeInvalidParams, "msg_ids cannot exceed 100")
	}

	url := fmt.Sprintf("/v3/messages/detail?msg_ids=%s", strings.Join(msgIDs, ","))
	
	// 使用report域名
	resp, err := s.client.makeReportRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var detailResp []MessageDetailResponse
	bodyBytes, _ := json.Marshal(resp.Body)
	if err := json.Unmarshal(bodyBytes, &detailResp); err != nil {
		return nil, NewJPushError(ErrorCodeInvalidJSON, fmt.Sprintf("failed to parse message detail response: %v", err))
	}

	return detailResp, nil
}