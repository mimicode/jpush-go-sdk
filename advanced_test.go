package goserversdk

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAdvancedService_GetCID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"cidlist": ["cid1", "cid2", "cid3"]}`))
	}))
	defer server.Close()

	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)
	
	client.baseURLs["push"] = server.URL

	result, err := client.Advanced.GetCID(3, CIDTypePush)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result.CIDList))
}

func TestAdvancedService_GetCID_InvalidParams(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)

	// 测试无效的count参数
	_, err = client.Advanced.GetCID(0, CIDTypePush)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}

	// 测试count超过限制
	_, err = client.Advanced.GetCID(1001, CIDTypePush)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}
}

func TestAdvancedService_GetCID_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": {"code": 1000, "message": "Internal Server Error"}}`))
	}))
	defer server.Close()

	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)
	
	client.baseURLs["push"] = server.URL

	_, err = client.Advanced.GetCID(3, CIDTypePush)
	assert.Error(t, err)
}

func TestAdvancedService_ValidatePush(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sendno": "123", "msg_id": "456"}`))
	}))
	defer server.Close()

	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)
	
	client.baseURLs["push"] = server.URL

	allValue := "all"
	request := &PushRequest{
		Platform: []string{"android", "ios"},
		Audience: &Audience{All: &allValue},
		Notification: &Notification{
			Alert: "Test notification",
		},
	}

	result, err := client.Advanced.ValidatePush(request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "123", result.SendNo)
	assert.Equal(t, "456", result.MsgID)
}

func TestAdvancedService_ValidatePush_InvalidRequest(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)

	// 测试nil请求
	_, err = client.Advanced.ValidatePush(nil)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}

	// 测试无效的请求（缺少audience）
	request := &PushRequest{
		Platform: []string{"android", "ios"},
		Notification: &Notification{
			Alert: "Test notification",
		},
	}

	_, err = client.Advanced.ValidatePush(request)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}
}

func TestAdvancedService_CancelPush(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)
	
	client.baseURLs["push"] = server.URL

	err = client.Advanced.CancelPush("123456")
	assert.NoError(t, err)
}

func TestAdvancedService_CancelPush_InvalidMsgID(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)

	// 测试空的msgID
	err = client.Advanced.CancelPush("")
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}
}

func TestAdvancedService_GetVendorQuota(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 0, "message": "success", "data": {"xiaomi_quota": {"operation": {"total": 1000, "used": 100}}}}`))
	}))
	defer server.Close()

	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)
	
	client.baseURLs["push"] = server.URL

	result, err := client.Advanced.GetVendorQuota()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Code)
	assert.Equal(t, "success", result.Message)
}

func TestAdvancedService_PushByFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sendno": "123", "msg_id": "456"}`))
	}))
	defer server.Close()

	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)
	
	client.baseURLs["push"] = server.URL

	request := &FilePushRequest{
		Platform: []string{"android", "ios"},
		Audience: &FileAudience{
			File: &FileTarget{
				FileID: "test-file-id",
			},
		},
		Notification: &Notification{
			Alert: "Test notification",
		},
	}

	result, err := client.Advanced.PushByFile(request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "123", result.SendNo)
	assert.Equal(t, "456", result.MsgID)
}

func TestAdvancedService_PushByFile_InvalidRequest(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)

	// 测试nil请求
	_, err = client.Advanced.PushByFile(nil)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}

	// 测试无效的请求（缺少file_id）
	request := &FilePushRequest{
		Platform: []string{"android", "ios"},
		Audience: &FileAudience{
			File: &FileTarget{
				FileID: "",
			},
		},
		Notification: &Notification{
			Alert: "Test notification",
		},
	}

	_, err = client.Advanced.PushByFile(request)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}
}

func TestAdvancedService_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)
	
	client.baseURLs["push"] = server.URL

	// 测试GetCID的无效JSON响应
	_, err = client.Advanced.GetCID(3, CIDTypePush)
	assert.Error(t, err)

	// 测试ValidatePush的无效JSON响应
	allValue3 := "all"
	request := &PushRequest{
		Platform: []string{"android", "ios"},
		Audience: &Audience{All: &allValue3},
		Notification: &Notification{
			Alert: "Test notification",
		},
	}

	_, err = client.Advanced.ValidatePush(request)
	assert.Error(t, err)

	// 测试GetVendorQuota的无效JSON响应
	_, err = client.Advanced.GetVendorQuota()
	assert.Error(t, err)
}

func TestAdvancedService_NetworkTimeout(t *testing.T) {
	// 创建一个不响应的服务器来模拟网络超时
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 不写任何响应，让请求超时
	}))
	defer server.Close()

	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
		Timeout:      1 * time.Second, // 设置1秒超时
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)
	
	client.baseURLs["push"] = server.URL

	// 测试GetCID超时
	_, err = client.Advanced.GetCID(3, CIDTypePush)
	assert.Error(t, err)

	// 测试ValidatePush超时
	allValue := "all"
	request := &PushRequest{
		Platform: []string{"android", "ios"},
		Audience: &Audience{All: &allValue},
		Notification: &Notification{
			Alert: "Test notification",
		},
	}

	_, err = client.Advanced.ValidatePush(request)
	assert.Error(t, err)
}