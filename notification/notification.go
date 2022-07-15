package notification

const (
	Global = "global"
)

type Notifier interface {
	// Publish - add msg to buffer
	Publish(string)
}
