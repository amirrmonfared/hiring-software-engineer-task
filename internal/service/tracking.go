package service

import (
	"sweng-task/internal/model"
	"sync"
	"time"

	"go.uber.org/zap"
)

type TrackingService struct {
	mu     sync.Mutex
	events []model.TrackingEvent
	log    *zap.SugaredLogger
}

func NewTrackingService(log *zap.SugaredLogger) *TrackingService {
	return &TrackingService{
		events: make([]model.TrackingEvent, 0),
		log:    log,
	}
}

func (ts *TrackingService) RecordEvent(evt model.TrackingEvent) error {
	if evt.Timestamp.IsZero() {
		evt.Timestamp = time.Now()
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.events = append(ts.events, evt)
	ts.log.Infow("Tracking event recorded", "event_type", evt.EventType, "line_item_id", evt.LineItemID)

	return nil
}
