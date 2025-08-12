package goserversdk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushService_Push(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sendno": "test-sendno", "msg_id": "123456789"}`))
	}))
	defer server.Close()

	client, err := NewTestClient()
	assert.NoError(t, err)
	
	client.baseURLs["push"] = server.URL

	all := "all"
	request := &PushRequest{
		Platform: NewAllPlatform().GetPlatforms(),
		Audience: &Audience{
			All: &all,
		},
		Notification: &Notification{
			Alert: "Test notification",
		},
	}

	result, err := client.Push.Push(request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-sendno", result.SendNo)
	assert.Equal(t, "123456789", result.MsgID)
}

func TestPushService_Push_InvalidRequest(t *testing.T) {
	client, err := NewTestClient()
	assert.NoError(t, err)

	// 测试空请求
	_, err = client.Push.Push(nil)
	assert.Error(t, err)

	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}
}

func TestPushService_Push_ValidationError(t *testing.T) {
	client, err := NewTestClient()
	assert.NoError(t, err)

	// 测试无效的推送请求（缺少platform）
	request := &PushRequest{
		Notification: &Notification{
			Alert: "Test notification",
		},
	}

	_, err = client.Push.Push(request)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidParams, jpushErr.Code)
	}
}

func TestPushService_Push_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"code": 1003, "message": "Invalid request"}}`))
	}))
	defer server.Close()

	client, err := NewTestClient()
	assert.NoError(t, err)
	
	client.baseURLs["push"] = server.URL

	all2 := "all"
	request := &PushRequest{
		Platform: NewAllPlatform().GetPlatforms(),
		Audience: &Audience{
			All: &all2,
		},
		Notification: &Notification{
			Alert: "Test notification",
		},
	}

	_, err = client.Push.Push(request)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCode(1003), jpushErr.Code)
	}
}

func TestPushService_Push_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client, err := NewTestClient()
	assert.NoError(t, err)
	
	client.baseURLs["push"] = server.URL

	all3 := "all"
	request := &PushRequest{
		Platform: NewAllPlatform().GetPlatforms(),
		Audience: &Audience{
			All: &all3,
		},
		Notification: &Notification{
			Alert: "Test notification",
		},
	}

	_, err = client.Push.Push(request)
	assert.Error(t, err)
	
	if jpushErr, ok := err.(*JPushError); ok {
		assert.Equal(t, ErrorCodeInvalidJSON, jpushErr.Code)
	}
}

func TestValidatePushRequest(t *testing.T) {
	all := "all"
	tests := []struct {
		name        string
		request     *PushRequest
		wantErr     bool
		expectedErr ErrorCode
	}{
		{
			name: "valid request with all audience",
			request: &PushRequest{
				Platform: NewAllPlatform().GetPlatforms(),
				Audience: &Audience{
					All: &all,
				},
				Notification: &Notification{
					Alert: "Test",
				},
			},
			wantErr: false,
		},
		{
			name: "valid request with tag audience",
			request: &PushRequest{
				Platform: NewAllPlatform().GetPlatforms(),
				Audience: &Audience{
					Tag: []string{"tag1", "tag2"},
				},
				Notification: &Notification{
					Alert: "Test",
				},
			},
			wantErr: false,
		},
		{
			name: "valid request with alias audience",
			request: &PushRequest{
				Platform: NewAllPlatform().GetPlatforms(),
				Audience: &Audience{
					Alias: []string{"alias1", "alias2"},
				},
				Notification: &Notification{
					Alert: "Test",
				},
			},
			wantErr: false,
		},
		{
			name: "valid request with registration_id audience",
			request: &PushRequest{
				Platform: NewAllPlatform().GetPlatforms(),
				Audience: &Audience{
					RegistrationID: []string{"reg1", "reg2"},
				},
				Notification: &Notification{
					Alert: "Test",
				},
			},
			wantErr: false,
		},
		{
			name: "valid request with message",
			request: &PushRequest{
				Platform: NewAllPlatform().GetPlatforms(),
				Audience: &Audience{
					All: &all,
				},
				Message: &Message{
					MsgContent: "Test message",
				},
			},
			wantErr: false,
		},
		{
			name:        "nil request",
			request:     nil,
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
		{
			name: "nil audience",
			request: &PushRequest{
				Platform: NewAllPlatform().GetPlatforms(),
				Notification: &Notification{
					Alert: "Test",
				},
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
		{
			name: "empty audience",
			request: &PushRequest{
				Platform:     NewAllPlatform().GetPlatforms(),
				Audience:     &Audience{},
				Notification: &Notification{Alert: "Test"},
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
		{
			name: "no notification or message",
			request: &PushRequest{
				Platform: NewAllPlatform().GetPlatforms(),
				Audience: &Audience{
					All: &all,
				},
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
		{
			name: "empty notification alert",
			request: &PushRequest{
				Audience: &Audience{
					All: &all,
				},
				Notification: &Notification{},
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
		{
			name: "empty message content",
			request: &PushRequest{
				Audience: &Audience{
					All: &all,
				},
				Message: &Message{},
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePushRequest(tt.request)
			
			if tt.wantErr {
				assert.Error(t, err)
				if jpushErr, ok := err.(*JPushError); ok {
					assert.Equal(t, tt.expectedErr, jpushErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAudience(t *testing.T) {
	all := "all"
	tests := []struct {
		name        string
		audience    *Audience
		wantErr     bool
		expectedErr ErrorCode
	}{
		{
			name: "valid all audience",
			audience: &Audience{
				All: &all,
			},
			wantErr: false,
		},
		{
			name: "valid tag audience",
			audience: &Audience{
				Tag: []string{"tag1"},
			},
			wantErr: false,
		},
		{
			name: "valid alias audience",
			audience: &Audience{
				Alias: []string{"alias1"},
			},
			wantErr: false,
		},
		{
			name: "valid registration_id audience",
			audience: &Audience{
				RegistrationID: []string{"reg1"},
			},
			wantErr: false,
		},
		{
			name:        "nil audience",
			audience:    nil,
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
		{
			name:        "empty audience",
			audience:    &Audience{},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
		{
			name: "empty tag array",
			audience: &Audience{
				Tag: []string{},
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
		{
			name: "empty alias array",
			audience: &Audience{
				Alias: []string{},
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
		{
			name: "empty registration_id array",
			audience: &Audience{
				RegistrationID: []string{},
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAudience(tt.audience)
			
			if tt.wantErr {
				assert.Error(t, err)
				if jpushErr, ok := err.(*JPushError); ok {
					assert.Equal(t, tt.expectedErr, jpushErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateNotification(t *testing.T) {
	tests := []struct {
		name         string
		notification *Notification
		wantErr      bool
		expectedErr  ErrorCode
	}{
		{
			name: "valid notification",
			notification: &Notification{
				Alert: "Test alert",
			},
			wantErr: false,
		},
		{
			name:         "nil notification",
			notification: nil,
			wantErr:      true,
			expectedErr:  ErrorCodeInvalidParams,
		},
		{
			name: "empty alert",
			notification: &Notification{
				Alert: "",
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNotification(tt.notification)
			
			if tt.wantErr {
				assert.Error(t, err)
				if jpushErr, ok := err.(*JPushError); ok {
					assert.Equal(t, tt.expectedErr, jpushErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     *Message
		wantErr     bool
		expectedErr ErrorCode
	}{
		{
			name: "valid message",
			message: &Message{
				MsgContent: "Test message",
			},
			wantErr: false,
		},
		{
			name:        "nil message",
			message:     nil,
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
		{
			name: "empty message content",
			message: &Message{
				MsgContent: "",
			},
			wantErr:     true,
			expectedErr: ErrorCodeInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMessage(tt.message)
			
			if tt.wantErr {
				assert.Error(t, err)
				if jpushErr, ok := err.(*JPushError); ok {
					assert.Equal(t, tt.expectedErr, jpushErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// 辅助函数，用于测试
func validatePushRequest(req *PushRequest) error {
	if req == nil {
		return NewJPushError(ErrorCodeInvalidParams, "push request cannot be nil")
	}
	
	if err := validateAudience(req.Audience); err != nil {
		return err
	}
	
	if req.Notification == nil && req.Message == nil {
		return NewJPushError(ErrorCodeInvalidParams, "either notification or message must be provided")
	}
	
	if req.Notification != nil {
		if err := validateNotification(req.Notification); err != nil {
			return err
		}
	}
	
	if req.Message != nil {
		if err := validateMessage(req.Message); err != nil {
			return err
		}
	}
	
	return nil
}

func validateAudience(audience *Audience) error {
	if audience == nil {
		return NewJPushError(ErrorCodeInvalidParams, "audience cannot be nil")
	}
	
	if audience.All == nil && len(audience.Tag) == 0 && len(audience.Alias) == 0 && len(audience.RegistrationID) == 0 {
		return NewJPushError(ErrorCodeInvalidParams, "audience must specify at least one target")
	}
	
	return nil
}

func validateNotification(notification *Notification) error {
	if notification == nil {
		return NewJPushError(ErrorCodeInvalidParams, "notification cannot be nil")
	}
	
	if notification.Alert == "" {
		return NewJPushError(ErrorCodeInvalidParams, "notification alert cannot be empty")
	}
	
	return nil
}

func validateMessage(message *Message) error {
	if message == nil {
		return NewJPushError(ErrorCodeInvalidParams, "message cannot be nil")
	}
	
	if message.MsgContent == "" {
		return NewJPushError(ErrorCodeInvalidParams, "message content cannot be empty")
	}
	
	return nil
}
