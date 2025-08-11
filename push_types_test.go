package goserversdk

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func TestPlatformConstants(t *testing.T) {
	assert.Equal(t, "all", PlatformAll)
	assert.Equal(t, "android", PlatformAndroid)
	assert.Equal(t, "ios", PlatformIOS)
	assert.Equal(t, "winphone", PlatformWinPhone)
}

func TestPushRequest_JSON(t *testing.T) {
	all := "all"
	req := &PushRequest{
		Platform: PlatformAll,
		Audience: &Audience{
			All: &all,
		},
		Notification: &Notification{
			Alert: "Test notification",
		},
		Message: &Message{
			MsgContent: "Test message",
			Title:      stringPtr("Test title"),
		},
		Options: &Options{
			TimeToLive:     intPtr(3600),
			APNSProduction: boolPtr(false),
		},
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded PushRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, req.Platform, decoded.Platform)
	assert.Equal(t, req.Audience.All, decoded.Audience.All)
	assert.Equal(t, req.Notification.Alert, decoded.Notification.Alert)
	assert.Equal(t, req.Message.MsgContent, decoded.Message.MsgContent)
	assert.Equal(t, req.Options.TimeToLive, decoded.Options.TimeToLive)
}

func TestAudience_JSON(t *testing.T) {
	all := "all"
	tests := []struct {
		name     string
		audience *Audience
	}{
		{
			name: "all audience",
			audience: &Audience{
				All: &all,
			},
		},
		{
			name: "registration id audience",
			audience: &Audience{
				RegistrationID: []string{"reg1", "reg2", "reg3"},
			},
		},
		{
			name: "tag audience",
			audience: &Audience{
				Tag: []string{"tag1", "tag2"},
			},
		},
		{
			name: "tag and audience",
			audience: &Audience{
				TagAnd: []string{"tag1", "tag2"},
			},
		},
		{
			name: "tag not audience",
			audience: &Audience{
				TagNot: []string{"tag1", "tag2"},
			},
		},
		{
			name: "alias audience",
			audience: &Audience{
				Alias: []string{"alias1", "alias2"},
			},
		},
		{
			name: "segment audience",
			audience: &Audience{
				Segment: []string{"seg1", "seg2"},
			},
		},
		{
			name: "ab test audience",
			audience: &Audience{
				ABTest: []string{"test_id_123"},
			},
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.audience)
			assert.NoError(t, err)
			assert.NotEmpty(t, data)

			var decoded Audience
			err = json.Unmarshal(data, &decoded)
			assert.NoError(t, err)
			
			// Compare specific fields based on test case
			if tt.audience.All != nil {
				assert.Equal(t, tt.audience.All, decoded.All)
			}
			if tt.audience.RegistrationID != nil {
				assert.Equal(t, tt.audience.RegistrationID, decoded.RegistrationID)
			}
			if tt.audience.Tag != nil {
				assert.Equal(t, tt.audience.Tag, decoded.Tag)
			}
			if tt.audience.TagAnd != nil {
				assert.Equal(t, tt.audience.TagAnd, decoded.TagAnd)
			}
			if tt.audience.TagNot != nil {
				assert.Equal(t, tt.audience.TagNot, decoded.TagNot)
			}
			if tt.audience.Alias != nil {
				assert.Equal(t, tt.audience.Alias, decoded.Alias)
			}
			if tt.audience.Segment != nil {
				assert.Equal(t, tt.audience.Segment, decoded.Segment)
			}
		})
	}
}

func TestNotification_JSON(t *testing.T) {
	notification := &Notification{
		Alert: "Global alert",
		Android: &AndroidNotification{
			Alert:     "Android alert",
			Title:     stringPtr("Android title"),
			BuilderID: intPtr(1),
			Priority:  intPtr(1),
			Category:  stringPtr("test"),
			Style:     intPtr(1),
			AlertType: intPtr(1),
			BigText:   stringPtr("Big text content"),
			Inbox: map[string]interface{}{
				"line1": "Inbox line 1",
				"line2": "Inbox line 2",
			},
			BigPicPath: stringPtr("http://example.com/pic.jpg"),
			Extras: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
		},
		IOS: &IOSNotification{
			Alert: map[string]interface{}{
				"title": "iOS title",
				"body":  "iOS body",
			},
			Sound:            stringPtr("default"),
			Badge:            stringPtr("+1"),
			ContentAvailable: boolPtr(true),
			MutableContent:   boolPtr(true),
			Category:         stringPtr("test_category"),
			Extras: map[string]interface{}{
				"ios_key": "ios_value",
			},
		},

	}

	data, err := json.Marshal(notification)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded Notification
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, notification.Alert, decoded.Alert)
	assert.Equal(t, notification.Android.Alert, decoded.Android.Alert)
	assert.Equal(t, notification.IOS.Sound, decoded.IOS.Sound)
}

func TestMessage_JSON(t *testing.T) {
	message := &Message{
		MsgContent:      "Message content",
		Title:           stringPtr("Message title"),
		ContentType:     stringPtr("text"),
		Extras: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
			"key3": true,
		},
	}

	data, err := json.Marshal(message)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded Message
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, message.MsgContent, decoded.MsgContent)
	assert.Equal(t, message.Title, decoded.Title)
	assert.Equal(t, message.ContentType, decoded.ContentType)
	assert.Equal(t, message.Extras["key1"], decoded.Extras["key1"])
}

func TestOptions_JSON(t *testing.T) {
	options := &Options{
		TimeToLive:      intPtr(3600),
		APNSProduction:  boolPtr(false),
		APNSCollapseID:  stringPtr("collapse_id"),
		BigPushDuration: intPtr(60),
	}

	data, err := json.Marshal(options)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded Options
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, options.TimeToLive, decoded.TimeToLive)
	assert.Equal(t, options.APNSProduction, decoded.APNSProduction)
	assert.Equal(t, options.APNSCollapseID, decoded.APNSCollapseID)
	assert.Equal(t, options.BigPushDuration, decoded.BigPushDuration)
}

func TestPushResponse_JSON(t *testing.T) {
	response := &PushResponse{
		MsgID:  "123456789",
		SendNo: "test_sendno",
	}

	data, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded PushResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, response.MsgID, decoded.MsgID)
	assert.Equal(t, response.SendNo, decoded.SendNo)
}



func TestMessageStatusRequest_JSON(t *testing.T) {
	statusReq := &MessageStatusRequest{
		MsgID:           int64(123456789),
		RegistrationIDs: []string{"test_registration_id"},
		Date:            "2023-10-01",
	}

	data, err := json.Marshal(statusReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded MessageStatusRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, statusReq.MsgID, decoded.MsgID)
	assert.Equal(t, statusReq.RegistrationIDs, decoded.RegistrationIDs)
	assert.Equal(t, statusReq.Date, decoded.Date)
}

func TestComplexNotificationStructure(t *testing.T) {
	// Test complex nested structure
	notification := &Notification{
		Alert: "Complex notification",
		Android: &AndroidNotification{
			Alert:     "Android complex alert",
			Title:     stringPtr("Android complex title"),
			BuilderID: intPtr(2),
			Priority:  intPtr(2),
			Category:  stringPtr("complex"),
			Style:     intPtr(2),
			AlertType: intPtr(2),
			BigText:   stringPtr("This is a very long text for big text style notification"),
			Inbox: map[string]interface{}{
				"line1": "First inbox line",
				"line2": "Second inbox line",
				"line3": "Third inbox line",
			},
			BigPicPath: stringPtr("https://example.com/big_picture.jpg"),
			Extras: map[string]interface{}{
				"android_extra_1": "android_value_1",
				"android_extra_2": 456,
				"android_extra_3": true,
				"android_extra_4": map[string]interface{}{
					"nested_key": "nested_value",
				},
			},
		},
		IOS: &IOSNotification{
			Alert: map[string]interface{}{
				"title":    "iOS Complex Title",
				"body":     "iOS Complex Body",
				"subtitle": "iOS Subtitle",
			},
			Sound:            stringPtr("custom_sound.wav"),
			Badge:            "10",
			ContentAvailable: boolPtr(true),
			MutableContent:   boolPtr(true),
			Category:         stringPtr("complex_category"),
			ThreadID:         stringPtr("thread_123"),
			Extras: map[string]interface{}{
				"ios_extra_1": "ios_value_1",
				"ios_extra_2": 789,
				"ios_extra_3": false,
				"ios_extra_4": []interface{}{
					"item1", "item2", "item3",
				},
			},
		},
	}

	data, err := json.Marshal(notification)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded Notification
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	
	// Verify complex structure is preserved
	assert.Equal(t, notification.Alert, decoded.Alert)
	assert.Equal(t, notification.Android.BigText, decoded.Android.BigText)
	assert.Equal(t, notification.Android.Extras["android_extra_1"], decoded.Android.Extras["android_extra_1"])
	assert.Equal(t, notification.IOS.ThreadID, decoded.IOS.ThreadID)
	
	// Verify nested structures
	iosAlert, ok := decoded.IOS.Alert.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "iOS Complex Title", iosAlert["title"])
	assert.Equal(t, "iOS Subtitle", iosAlert["subtitle"])
}

func TestEmptyStructures(t *testing.T) {
	// Test empty structures
	emptyReq := &PushRequest{}
	data, err := json.Marshal(emptyReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded PushRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	emptyAudience := &Audience{}
	data, err = json.Marshal(emptyAudience)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decodedAudience Audience
	err = json.Unmarshal(data, &decodedAudience)
	assert.NoError(t, err)
}

func TestNilPointers(t *testing.T) {
	// Test structures with nil pointers
	req := &PushRequest{
		Platform:     PlatformAll,
		Audience:     nil,
		Notification: nil,
		Message:      nil,
		Options:      nil,
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded PushRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, req.Platform, decoded.Platform)
	assert.Nil(t, decoded.Audience)
	assert.Nil(t, decoded.Notification)
	assert.Nil(t, decoded.Message)
	assert.Nil(t, decoded.Options)
}