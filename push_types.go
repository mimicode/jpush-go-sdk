package goserversdk

// Platform constants
const (
	PlatformAll      = "all"
	PlatformAndroid  = "android"
	PlatformIOS      = "ios"
	PlatformWinPhone = "winphone"
)

// Platform 推送平台
type Platform interface {
	GetPlatforms() interface{}
}

// AllPlatform 所有平台
type AllPlatform struct{}

func (p AllPlatform) GetPlatforms() interface{} {
	return "all"
}

// SpecificPlatforms 指定平台
type SpecificPlatforms struct {
	Platforms []string `json:"platforms"`
}

func (p SpecificPlatforms) GetPlatforms() interface{} {
	return p.Platforms
}

// NewAllPlatform 创建所有平台
func NewAllPlatform() Platform {
	return AllPlatform{}
}

// NewSpecificPlatforms 创建指定平台
func NewSpecificPlatforms(platforms ...string) Platform {
	return SpecificPlatforms{Platforms: platforms}
}

// Audience 推送目标
type Audience struct {
	All              *string   `json:"all,omitempty"`              // 广播
	Tag              []string  `json:"tag,omitempty"`              // 标签OR
	TagAnd           []string  `json:"tag_and,omitempty"`          // 标签AND
	TagNot           []string  `json:"tag_not,omitempty"`          // 标签NOT
	Alias            []string  `json:"alias,omitempty"`            // 别名
	RegistrationID   []string  `json:"registration_id,omitempty"`  // 注册ID
	Segment          []string  `json:"segment,omitempty"`          // 用户分群
	ABTest           []string  `json:"abtest,omitempty"`           // A/B测试
	LiveActivityID   *string   `json:"live_activity_id,omitempty"` // 实时活动ID
}

// NewBroadcastAudience 创建广播推送目标
func NewBroadcastAudience() *Audience {
	all := "all"
	return &Audience{All: &all}
}

// NewTagAudience 创建标签推送目标
func NewTagAudience(tags ...string) *Audience {
	return &Audience{Tag: tags}
}

// NewTagAndAudience 创建标签AND推送目标
func NewTagAndAudience(tags ...string) *Audience {
	return &Audience{TagAnd: tags}
}

// NewTagNotAudience 创建标签NOT推送目标
func NewTagNotAudience(tags ...string) *Audience {
	return &Audience{TagNot: tags}
}

// NewAliasAudience 创建别名推送目标
func NewAliasAudience(aliases ...string) *Audience {
	return &Audience{Alias: aliases}
}

// NewRegistrationIDAudience 创建注册ID推送目标
func NewRegistrationIDAudience(regIDs ...string) *Audience {
	return &Audience{RegistrationID: regIDs}
}

// NewLiveActivityAudience 创建实时活动推送目标
func NewLiveActivityAudience(liveActivityID string) *Audience {
	return &Audience{LiveActivityID: &liveActivityID}
}

// Intent 意图
type Intent struct {
	URL string `json:"url"`
}

// AndroidNotification Android通知
type AndroidNotification struct {
	Alert       string                 `json:"alert"`                  // 通知内容
	Title       *string                `json:"title,omitempty"`        // 通知标题
	BuilderID   *int                   `json:"builder_id,omitempty"`   // 通知栏样式ID
	ChannelID   *string                `json:"channel_id,omitempty"`   // 通知渠道ID
	Category    *string                `json:"category,omitempty"`     // 通知分类
	Priority    *int                   `json:"priority,omitempty"`     // 优先级
	Style       *int                   `json:"style,omitempty"`        // 通知样式
	AlertType   *int                   `json:"alert_type,omitempty"`   // 提醒类型
	BigText     *string                `json:"big_text,omitempty"`     // 大文本
	Inbox       map[string]interface{} `json:"inbox,omitempty"`        // 收件箱样式
	BigPicPath  *string                `json:"big_pic_path,omitempty"` // 大图路径
	LargeIcon   *string                `json:"large_icon,omitempty"`   // 大图标
	Intent      *Intent                `json:"intent,omitempty"`       // 意图
	Extras      map[string]interface{} `json:"extras,omitempty"`       // 附加字段
}

// IOSNotification iOS通知
type IOSNotification struct {
	Alert            interface{}            `json:"alert"`                       // 通知内容
	Sound            *string                `json:"sound,omitempty"`             // 声音
	Badge            interface{}            `json:"badge,omitempty"`             // 角标
	ContentAvailable *bool                  `json:"content-available,omitempty"` // 静默推送
	MutableContent   *bool                  `json:"mutable-content,omitempty"`   // 可变内容
	Category         *string                `json:"category,omitempty"`          // 分类
	ThreadID         *string                `json:"thread-id,omitempty"`         // 线程ID
	Extras           map[string]interface{} `json:"extras,omitempty"`            // 附加字段
}

// HMOSNotification 鸿蒙通知
type HMOSNotification struct {
	Alert        string                 `json:"alert"`                   // 通知内容
	Title        *string                `json:"title,omitempty"`         // 通知标题
	Intent       *Intent                `json:"intent,omitempty"`        // 意图
	BadgeAddNum  *int                   `json:"badge_add_num,omitempty"` // 角标增加数
	BadgeSetNum  *int                   `json:"badge_set_num,omitempty"` // 角标设置数
	Extras       map[string]interface{} `json:"extras,omitempty"`        // 附加字段
	Category     *string                `json:"category,omitempty"`      // 分类
	TestMessage  *bool                  `json:"test_message,omitempty"`  // 测试消息
	ReceiptID    *string                `json:"receipt_id,omitempty"`    // 回执ID
	LargeIcon    *string                `json:"large_icon,omitempty"`    // 大图标
	Style        *int                   `json:"style,omitempty"`         // 样式
	PushType     *int                   `json:"push_type,omitempty"`     // 推送类型
}

// QuickAppNotification 快应用通知
type QuickAppNotification struct {
	Alert  string                 `json:"alert"`            // 通知内容
	Title  *string                `json:"title,omitempty"`  // 通知标题
	Page   *string                `json:"page,omitempty"`   // 页面路径
	Extras map[string]interface{} `json:"extras,omitempty"` // 附加字段
}

// VOIPNotification VOIP通知
type VOIPNotification map[string]interface{}

// Notification 通知
type Notification struct {
	Alert    string                `json:"alert,omitempty"`    // 通用通知内容
	Android  *AndroidNotification  `json:"android,omitempty"`  // Android通知
	IOS      *IOSNotification      `json:"ios,omitempty"`      // iOS通知
	HMOS     *HMOSNotification     `json:"hmos,omitempty"`     // 鸿蒙通知
	QuickApp *QuickAppNotification `json:"quickapp,omitempty"` // 快应用通知
	VOIP     VOIPNotification      `json:"voip,omitempty"`     // VOIP通知
}

// Message 自定义消息
type Message struct {
	MsgContent  string                 `json:"msg_content"`            // 消息内容
	Title       *string                `json:"title,omitempty"`        // 消息标题
	ContentType *string                `json:"content_type,omitempty"` // 内容类型
	Extras      map[string]interface{} `json:"extras,omitempty"`       // 附加字段
}

// SMSMessage 短信补充
type SMSMessage struct {
	TempID       int                    `json:"temp_id"`                // 短信模板ID
	TempPara     map[string]interface{} `json:"temp_para,omitempty"`    // 短信模板参数
	DelayTime    *int                   `json:"delay_time,omitempty"`   // 延迟时间
	ActiveFilter *bool                  `json:"active_filter,omitempty"` // 活跃过滤
}

// Options 推送选项
type Options struct {
	TimeToLive      *int    `json:"time_to_live,omitempty"`      // 离线保留时长
	APNSProduction  *bool   `json:"apns_production,omitempty"`   // APNs生产环境
	APNSCollapseID  *string `json:"apns_collapse_id,omitempty"`  // APNs折叠ID
	BigPushDuration *int    `json:"big_push_duration,omitempty"` // 定速推送时长
}

// Callback 回调
type Callback struct {
	URL    string                 `json:"url"`              // 回调URL
	Params map[string]interface{} `json:"params,omitempty"` // 回调参数
	Type   *int                   `json:"type,omitempty"`   // 回调类型
}

// PushRequest 推送请求
type PushRequest struct {
	Platform     interface{}  `json:"platform"`               // 推送平台
	Audience     *Audience    `json:"audience"`               // 推送目标
	Notification *Notification `json:"notification,omitempty"` // 通知
	Message      *Message     `json:"message,omitempty"`      // 自定义消息
	SMSMessage   *SMSMessage  `json:"sms_message,omitempty"`  // 短信补充
	Options      *Options     `json:"options,omitempty"`      // 推送选项
	Callback     *Callback    `json:"callback,omitempty"`     // 回调
	CID          *string      `json:"cid,omitempty"`          // 防重复标识
}

// PushResponse 推送响应
type PushResponse struct {
	SendNo string `json:"sendno"` // 推送序号
	MsgID  string `json:"msg_id"` // 消息ID
}

// NewPushRequest 创建推送请求
func NewPushRequest() *PushRequest {
	return &PushRequest{}
}

// SetPlatform 设置推送平台
func (r *PushRequest) SetPlatform(platform Platform) *PushRequest {
	r.Platform = platform.GetPlatforms()
	return r
}

// SetAudience 设置推送目标
func (r *PushRequest) SetAudience(audience *Audience) *PushRequest {
	r.Audience = audience
	return r
}

// SetNotification 设置通知
func (r *PushRequest) SetNotification(notification *Notification) *PushRequest {
	r.Notification = notification
	return r
}

// SetMessage 设置自定义消息
func (r *PushRequest) SetMessage(message *Message) *PushRequest {
	r.Message = message
	return r
}

// SetSMSMessage 设置短信补充
func (r *PushRequest) SetSMSMessage(smsMessage *SMSMessage) *PushRequest {
	r.SMSMessage = smsMessage
	return r
}

// SetOptions 设置推送选项
func (r *PushRequest) SetOptions(options *Options) *PushRequest {
	r.Options = options
	return r
}

// SetCallback 设置回调
func (r *PushRequest) SetCallback(callback *Callback) *PushRequest {
	r.Callback = callback
	return r
}

// SetCID 设置防重复标识
func (r *PushRequest) SetCID(cid string) *PushRequest {
	r.CID = &cid
	return r
}