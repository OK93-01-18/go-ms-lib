package notification

import (
	"github.com/ok93-01-18/event_reporter"
)

const (
	Global = "global"
)

type Notifier interface {
	// Publish - add msg to buffer
	Publish(string)
}

type GlobalNotifier struct {
	eventReporter *event_reporter.EventReporter
	topic         string
}

func (g *GlobalNotifier) Publish(msg string) {
	g.eventReporter.Publish(g.topic, msg)
}
