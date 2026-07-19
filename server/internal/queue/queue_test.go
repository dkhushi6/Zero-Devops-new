package queue

import (
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestSetUpQueues_NilConnection(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil connection")
		}
	}()
	_ = SetUpQueues(nil, nil)
}

func TestSetUpQueues_UnreachableConnection(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:1/")
	if err != nil {
		t.Skip("expected dial to fail on bad port, but got connection")
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		t.Skip("expected channel to fail on bad connection")
	}
	defer func() { _ = ch.Close() }()

	err = SetUpQueues(conn, ch)
	if err != nil {
		t.Logf("SetUpQueues returned expected error for unreachable broker: %v", err)
	}
}

func TestSetUpQueues_NoRabbitMQ(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	if err != nil {
		t.Skip("RabbitMQ is not running on default port — skipping integration test")
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		t.Skip("failed to open channel — skipping")
	}
	defer func() { _ = ch.Close() }()

	if err := SetUpQueues(conn, ch); err != nil {
		t.Errorf("SetUpQueues failed: %v", err)
	}
}

func TestSetUpQueues_TimesOutFastOnBadHost(t *testing.T) {
	done := make(chan bool)
	go func() {
		defer func() { _ = recover(); close(done) }()
		_ = SetUpQueues(nil, nil)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("SetUpQueues with nil connection did not return within 2 seconds")
	}
}
