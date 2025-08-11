# 极光推送 Go SDK

极光推送服务端 Go SDK，支持推送、统计、高级功能等完整的极光推送服务。

## 安装

```bash
go get github.com/mimicode/jpush-go-sdk
```

## 快速开始

### 1. 初始化客户端

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/mimicode/jpush-go-sdk"
)

func main() {
    // 创建客户端
    client, err := goserversdk.NewClient("your-app-key", "your-master-secret")
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建推送服务
    pushService := goserversdk.NewPushService(client)
}
```

### 2. 基本推送

```go
// 创建推送请求
pushReq := &goserversdk.PushRequest{
    Platform: goserversdk.PlatformAll,
    Audience: &goserversdk.Audience{
        All: true,
    },
    Notification: &goserversdk.Notification{
        Alert: "Hello, JPush!",
        Android: &goserversdk.AndroidNotification{
            Alert: "Android 通知",
            Title: "测试标题",
        },
        IOS: &goserversdk.IOSNotification{
            Alert: "iOS 通知",
            Badge: 1,
            Sound: "default",
        },
    },
}

// 发送推送
resp, err := pushService.Push(pushReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("推送成功，消息ID: %d\n", resp.MsgID)
```

### 3. 指定用户推送

```go
// 按注册ID推送
pushReq := &goserversdk.PushRequest{
    Platform: goserversdk.PlatformAll,
    Audience: &goserversdk.Audience{
        RegistrationID: []string{"registration_id_1", "registration_id_2"},
    },
    Notification: &goserversdk.Notification{
        Alert: "定向推送消息",
    },
}

resp, err := pushService.Push(pushReq)
```

### 4. 按标签推送

```go
// 按标签推送
pushReq := &goserversdk.PushRequest{
    Platform: goserversdk.PlatformAll,
    Audience: &goserversdk.Audience{
        Tag: []string{"tag1", "tag2"},
    },
    Notification: &goserversdk.Notification{
        Alert: "标签推送消息",
    },
}

resp, err := pushService.Push(pushReq)
```

### 5. 自定义消息

```go
pushReq := &goserversdk.PushRequest{
    Platform: goserversdk.PlatformAll,
    Audience: &goserversdk.Audience{
        All: true,
    },
    Message: &goserversdk.Message{
        MsgContent: "自定义消息内容",
        Title:      "消息标题",
        ContentType: "text",
        Extras: map[string]interface{}{
            "key1": "value1",
            "key2": "value2",
        },
    },
}

resp, err := pushService.Push(pushReq)
```

## 高级功能

### 1. 获取CID

```go
advancedService := goserversdk.NewAdvancedService(client)

// 获取单个CID
cid, err := advancedService.GetCID()
if err != nil {
    log.Fatal(err)
}

// 获取多个CID
cids, err := advancedService.GetCIDs(10)
if err != nil {
    log.Fatal(err)
}
```

### 2. 推送验证

```go
// 验证推送请求
isValid, err := advancedService.ValidatePush(pushReq)
if err != nil {
    log.Fatal(err)
}

if isValid {
    fmt.Println("推送请求有效")
}
```

### 3. 取消推送

```go
// 取消推送
err := advancedService.CancelPush("msg_id")
if err != nil {
    log.Fatal(err)
}
```

### 4. 获取厂商配额

```go
quota, err := advancedService.GetVendorQuota()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("华为配额: %d\n", quota.Huawei)
fmt.Printf("小米配额: %d\n", quota.Xiaomi)
```

## 统计功能

### 1. 获取送达统计

```go
reportService := goserversdk.NewReportService(client)

// 获取送达统计详情
msgIDs := []string{"msg_id_1", "msg_id_2"}
details, err := reportService.GetReceivedDetail(msgIDs)
if err != nil {
    log.Fatal(err)
}

for _, detail := range details {
    fmt.Printf("消息ID: %s, 目标数: %d, 送达数: %d\n", 
        detail.MsgID, detail.Target, detail.Received)
}
```

### 2. 查询消息状态

```go
// 查询消息送达状态（VIP功能）
statusReq := &goserversdk.MessageStatusRequest{
    MsgID: 123456789,
    RegistrationIDs: []string{"reg_id_1", "reg_id_2"},
}

status, err := reportService.GetMessageStatus(statusReq)
if err != nil {
    log.Fatal(err)
}

for regID, info := range status {
    fmt.Printf("设备 %s 状态: %+v\n", regID, info)
}
```

## 错误处理

SDK 提供了详细的错误信息：

```go
resp, err := pushService.Push(pushReq)
if err != nil {
    if jpushErr, ok := err.(*goserversdk.JPushError); ok {
        fmt.Printf("错误码: %d\n", jpushErr.Code)
        fmt.Printf("错误信息: %s\n", jpushErr.Message)
        
        // 根据错误码处理
        switch jpushErr.Code {
        case goserversdk.ErrorCodeInvalidAuth:
            fmt.Println("认证失败，请检查AppKey和MasterSecret")
        case goserversdk.ErrorCodeRateLimitExceeded:
            fmt.Println("请求频率超限，请稍后重试")
        case goserversdk.ErrorCodeInvalidParams:
            fmt.Println("参数错误，请检查请求参数")
        }
    }
    return
}
```

## 配置选项

### 1. 自定义HTTP客户端

```go
import (
    "net/http"
    "time"
)

// 创建自定义HTTP客户端
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}

client, err := goserversdk.NewClientWithHTTPClient("app-key", "master-secret", httpClient)
```

### 2. 设置日志

```go
import "go.uber.org/zap"

// 创建自定义logger
logger, _ := zap.NewProduction()

client, err := goserversdk.NewClientWithLogger("app-key", "master-secret", logger)
```

## API 参考

### 错误码

| 错误码 | 说明 |
|--------|------|
| `ErrorCodeInvalidParams` | 参数错误 |
| `ErrorCodeInvalidAuth` | 认证失败 |
| `ErrorCodeInvalidAppKey` | AppKey无效 |
| `ErrorCodeRateLimitExceeded` | 频率限制 |
| `ErrorCodeInternalError` | 内部错误 |
| `ErrorCodeTimeout` | 请求超时 |

### 平台常量

- `PlatformAll`: 所有平台
- `PlatformAndroid`: Android平台
- `PlatformIOS`: iOS平台
- `PlatformWinPhone`: Windows Phone平台

## 相关链接

- [极光推送官方文档](https://docs.jiguang.cn/jpush/server/push/server_overview)
- [REST API v3 推送](https://docs.jiguang.cn/jpush/server/push/rest_api_v3_push)
- [高级功能](https://docs.jiguang.cn/jpush/server/push/rest_api_v3_push_advanced)
- [设备管理](https://docs.jiguang.cn/jpush/server/push/rest_api_v3_device)

## 许可证

MIT License
 