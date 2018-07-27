package event_processor

import "github.com/aws/aws-sdk-go/service/sqs"

type MessageProc interface {
	Register() map[string]string
	Process(data interface{}) error
}

type Collector interface {
	ReceiveMessage() ([]*sqs.Message, error)
	Ack(messageId string)
	Collect()
}
