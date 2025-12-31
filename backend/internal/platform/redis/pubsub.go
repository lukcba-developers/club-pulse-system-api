package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// Event types
const (
	EventBookingCancelled = "booking.cancelled"
	EventSlotAvailable    = "slot.available"
	EventMaintenanceStart = "maintenance.start"
	EventMaintenanceEnd   = "maintenance.end"
)

// Channel names
const (
	ChannelBookings    = "clubpulse:bookings"
	ChannelMaintenance = "clubpulse:maintenance"
)

// BookingEvent represents a booking-related event
type BookingEvent struct {
	Type       string    `json:"type"`
	FacilityID string    `json:"facility_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	UserID     string    `json:"user_id,omitempty"`
	Message    string    `json:"message,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// EventPublisher publishes events to Redis Pub/Sub
type EventPublisher struct {
	redis *RedisClient
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		redis: GetClient(),
	}
}

// PublishBookingCancelled notifies that a booking was cancelled (slot is now available)
func (p *EventPublisher) PublishBookingCancelled(ctx context.Context, facilityID, userID string, start, end time.Time) error {
	event := BookingEvent{
		Type:       EventBookingCancelled,
		FacilityID: facilityID,
		StartTime:  start,
		EndTime:    end,
		UserID:     userID,
		Message:    "A booking slot has become available",
		Timestamp:  time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.redis.Publish(ctx, ChannelBookings, string(data))
}

// PublishSlotAvailable notifies that a slot is available
func (p *EventPublisher) PublishSlotAvailable(ctx context.Context, facilityID string, start, end time.Time) error {
	event := BookingEvent{
		Type:       EventSlotAvailable,
		FacilityID: facilityID,
		StartTime:  start,
		EndTime:    end,
		Message:    "Slot is now available for booking",
		Timestamp:  time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.redis.Publish(ctx, ChannelBookings, string(data))
}

// EventSubscriber subscribes to Redis Pub/Sub channels
type EventSubscriber struct {
	redis *RedisClient
}

// NewEventSubscriber creates a new event subscriber
func NewEventSubscriber() *EventSubscriber {
	return &EventSubscriber{
		redis: GetClient(),
	}
}

// SubscribeToBookings subscribes to booking events and calls the handler for each event
func (s *EventSubscriber) SubscribeToBookings(ctx context.Context, handler func(event BookingEvent)) error {
	pubsub := s.redis.Subscribe(ctx, ChannelBookings)

	// Process messages in a goroutine
	go func() {
		ch := pubsub.Channel()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				s.processMessage(msg, handler)
			case <-ctx.Done():
				pubsub.Close()
				return
			}
		}
	}()

	return nil
}

func (s *EventSubscriber) processMessage(msg *goredis.Message, handler func(event BookingEvent)) {
	var event BookingEvent
	if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return
	}
	handler(event)
}
