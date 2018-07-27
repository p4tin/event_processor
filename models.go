package event_processor

type Plugin struct {
	MessageProc
	Collector
	Type      string
	Provider  string
	QueueName string
}
