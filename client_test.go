package goserversdk

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	client, err := NewTestClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.logger)
	// 验证配置已正确加载
	assert.NotEmpty(t, client.appKey)
	assert.NotEmpty(t, client.masterSecret)
}

func TestNewClient_ValidationErrors(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		expectedErr ErrorCode
	}{
		{
			name: "valid credentials",
			config: &Config{
				AppKey:       "test-app-key",
				MasterSecret: "test-master-secret",
				Logger:       logger,
			},
			wantErr: false,
		},
		{
			name: "empty app key",
			config: &Config{
				AppKey:       "",
				MasterSecret: "test-master-secret",
				Logger:       logger,
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidAppKey,
		},
		{
			name: "empty master secret",
			config: &Config{
				AppKey:       "test-app-key",
				MasterSecret: "",
				Logger:       logger,
			},
			wantErr:     true,
			expectedErr: ErrorCodeMissingAuth,
		},
		{
			name: "nil logger",
			config: &Config{
				AppKey:       "test-app-key",
				MasterSecret: "test-master-secret",
				Logger:       nil,
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
				if jpushErr, ok := err.(*JPushError); ok {
					assert.Equal(t, tt.expectedErr, jpushErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tt.config.AppKey, client.appKey)
				assert.Equal(t, tt.config.MasterSecret, client.masterSecret)
				assert.NotNil(t, client.httpClient)
				assert.NotNil(t, client.logger)
			}
		})
	}
}

func TestNewClientWithTimeout(t *testing.T) {
	testConfig, err := LoadTestConfig()
	assert.NoError(t, err)
	
	logger, _ := zap.NewDevelopment()
	config := &Config{
		AppKey:       testConfig.AppKey,
		MasterSecret: testConfig.MasterSecret,
		Logger:       logger,
		Timeout:      10 * time.Second,
	}
	
	client, err := NewClient(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, 10*time.Second, client.httpClient.Timeout)
}

func TestClient_MakeRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "data": "test"}`))
	}))
	defer server.Close()

	client, err := NewTestClient()
	assert.NoError(t, err)

	resp, err := client.makeRequestWithoutContext("GET", server.URL, "/test", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Body)
}

func TestClient_MakeRequest_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"code": 1000, "message": "Invalid request"}}`))
	}))
	defer server.Close()

	client, err := NewTestClient()
	assert.NoError(t, err)

	_, err = client.makeRequestWithoutContext("GET", server.URL, "/test", nil)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCode(1000), jpushErr.Code)
	}
}

func TestClient_MakeRequest_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client, err := NewTestClient()
	assert.NoError(t, err)

	_, err = client.makeRequestWithoutContext("GET", server.URL, "/test", nil)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidJSON, jpushErr.Code)
	}
}

func TestClient_MakeRequest_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewTestClient()
	assert.NoError(t, err)
	
	// 设置超时时间
	client.httpClient.Timeout = 1 * time.Millisecond

	_, err = client.makeRequestWithoutContext("GET", server.URL, "/test", nil)
	assert.Error(t, err)
}