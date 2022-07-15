package notification

import "github.com/ok93-01-18/event_reporter"

type EventReporterNotifier struct {
	eventReporter *event_reporter.EventReporter
	topic         string
}

func (g *EventReporterNotifier) Publish(msg string) {
	g.eventReporter.Publish(g.topic, msg)
}

func NewEventReporterNotifier(reporter *event_reporter.EventReporter) Notifier {
	return &EventReporterNotifier{
		eventReporter: reporter,
		topic:         Global,
	}
}
