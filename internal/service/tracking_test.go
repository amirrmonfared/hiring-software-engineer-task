package service_test

import (
	"sweng-task/internal/model"
	"sweng-task/internal/service"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestTrackingService_RecordEvent(t *testing.T) {
	logger, _ := zap.NewProduction()
	log := logger.Sugar()

	trackingSvc := service.NewTrackingService(log)

	type testCase struct {
		name      string
		event     model.TrackingEvent
		expectErr bool
	}

	now := time.Now()

	tests := []testCase{
		{
			name: "Basic Impression",
			event: model.TrackingEvent{
				EventType:  model.TrackingEventTypeImpression,
				LineItemID: "li_1234",
				Timestamp:  now,
				UserID:     "user_1",
				Metadata: map[string]string{
					"referrer": "https://example.com",
				},
			},
			expectErr: false,
		},
		{
			name: "Click without timestamp",
			event: model.TrackingEvent{
				EventType:  model.TrackingEventTypeClick,
				LineItemID: "li_5678",
				UserID:     "user_2",
			},
			expectErr: false,
		},
		{
			name: "Conversion with metadata",
			event: model.TrackingEvent{
				EventType:  model.TrackingEventTypeConversion,
				LineItemID: "li_7777",
				Timestamp:  now.Add(-time.Hour),
				Metadata: map[string]string{
					"conversion_value": "99.99",
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := trackingSvc.RecordEvent(tc.event)
			if (err != nil) != tc.expectErr {
				t.Fatalf("RecordEvent() error = %v, expectErr = %v", err, tc.expectErr)
			}
		})
	}
}
