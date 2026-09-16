package notify

import "context"

type Message struct {
	Channel  string
	To       string
	Template string
	Payload  map[string]any
}

type Notifier interface {
	Send(ctx context.Context, msg Message) error
}

type Noop struct{}

func (Noop) Send(context.Context, Message) error { return nil }
