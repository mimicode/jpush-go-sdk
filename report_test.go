package goserversdk

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestReportService_GetReceivedDetail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"msg_id": "123", "target": 100, "received": 95}]`))
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
	
	client.baseURLs["report"] = server.URL

	msgIDs := []string{"123", "456"}
	result, err := client.Report.GetReceivedDetail(msgIDs)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestReportService_GetReceivedDetail_InvalidParams(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)

	// 测试空的msgIDs
	_, err = client.Report.GetReceivedDetail([]string{})
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}

	// 测试nil msgIDs
	_, err = client.Report.GetReceivedDetail(nil)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}
}

func TestReportService_GetReceivedDetail_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"code": 1003, "message": "Invalid params"}}`))
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
	
	client.baseURLs["report"] = server.URL

	msgIDs := []string{"123", "456"}
	_, err = client.Report.GetReceivedDetail(msgIDs)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCode(1003), jpushErr.Code)
	}
}

func TestReportService_GetReceived(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"time": "2023-01-01", "android": {"target": 100, "received": 95}}]`))
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
	
	client.baseURLs["report"] = server.URL

	msgIDs := []string{"123", "456"}
	result, err := client.Report.GetReceived(msgIDs)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestReportService_GetReceived_InvalidParams(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)

	// 测试空的msgIDs
	_, err = client.Report.GetReceived([]string{})
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}
}

func TestReportService_GetMessageStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "sent", "msg_id": "123"}`))
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
	
	client.baseURLs["report"] = server.URL

	request := &MessageStatusRequest{
		MsgID:           int64(123),
		RegistrationIDs: []string{"test-reg-id"},
		Date:            "2023-01-01",
	}

	result, err := client.Report.GetMessageStatus(request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestReportService_GetMessageStatus_InvalidParams(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)

	// 测试nil请求
	_, err = client.Report.GetMessageStatus(nil)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}

	// 测试无效的msgID
	request := &MessageStatusRequest{
		MsgID:           int64(0),
		RegistrationIDs: []string{"test-reg-id"},
		Date:            "2023-01-01",
	}

	_, err = client.Report.GetMessageStatus(request)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}
}

func TestReportService_GetMessageDetail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"msg_id": "123", "platform": "android", "audience": "all"}`))
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
	
	client.baseURLs["report"] = server.URL

	result, err := client.Report.GetMessageDetail([]string{"123"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestReportService_GetMessageDetail_InvalidParams(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)

	// 测试空的msgID
	_, err = client.Report.GetMessageDetail([]string{})
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}
}

func TestReportService_InvalidJSON(t *testing.T) {
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
	
	client.baseURLs["report"] = server.URL

	// 测试GetReceivedDetail的无效JSON响应
	msgIDs := []string{"123", "456"}
	_, err = client.Report.GetReceivedDetail(msgIDs)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidJSON, jpushErr.Code)
	}

	// 测试GetReceived的无效JSON响应
	_, err = client.Report.GetReceived(msgIDs)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidJSON, jpushErr.Code)
	}

	// 测试GetMessageStatus的无效JSON响应
	request := &MessageStatusRequest{
		MsgID:           int64(123),
		RegistrationIDs: []string{"test-reg-id"},
		Date:            "2023-01-01",
	}

	_, err = client.Report.GetMessageStatus(request)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidJSON, jpushErr.Code)
	}

	// 测试GetMessageDetail的无效JSON响应
	_, err = client.Report.GetMessageDetail([]string{"123"})
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidJSON, jpushErr.Code)
	}
}

func TestReportService_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
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
	
	client.baseURLs["report"] = server.URL

	msgIDs := []string{"123", "456"}
	result, err := client.Report.GetReceivedDetail(msgIDs)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestReportService_NetworkTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       "test-key",
		MasterSecret: "test-secret",
		Logger:       logger,
		Timeout:      1 * time.Millisecond,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)
	
	client.baseURLs["report"] = server.URL

	msgIDs := []string{"123", "456"}
	_, err = client.Report.GetReceivedDetail(msgIDs)
	assert.Error(t, err)
}
