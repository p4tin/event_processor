package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"event_processor"
)

const TEN = 10

var maxMessages = int64(TEN)
var maxTimeout = int64(TEN)

type Connection struct {
	Sqs *sqs.SQS
}

type Collector struct {
	Connection
	ID       string
	QueueUrl string
	Delay    int
	LateAck  bool
	Plugin   event_processor.Plugin
}

func CreateCollector(plugin event_processor.Plugin) (Collector, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return Collector{}, err
	}

	collector := Collector{
		QueueUrl: plugin.QueueName,
		Delay:    0,
		Plugin:   plugin,
		LateAck:  false,
		Connection: Connection{
			Sqs: sqs.New(sess),
		},
	}

	return collector, nil
}

func (c Collector) Ack(messageId string) {
	delParams := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.QueueUrl),
		ReceiptHandle: aws.String(messageId),
	}
	c.Sqs.DeleteMessage(delParams)
}

func (c Collector) ReceiveMessage() ([]*sqs.Message, error) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.QueueUrl), // Required
		MaxNumberOfMessages: &maxMessages,
		WaitTimeSeconds:     &maxTimeout,
	}
	resp, err := c.Sqs.ReceiveMessage(params)
	if err != nil {
		return []*sqs.Message{}, errors.New(fmt.Sprintf("Could not get AWS messages, error: %s", err.Error()))
	}
	if !c.LateAck {
		if len(resp.Messages) > 0 {
			for _, msg := range resp.Messages {
				fmt.Println(msg)
				c.Ack(*msg.ReceiptHandle)
			}
		}
	}
	return resp.Messages, nil
}

func (c Collector) Collect() {
	for {
		messages, err := c.ReceiveMessage()
		if err != nil {
			fmt.Printf("AWS Collector %s: %s\n", c.ID, err.Error())
		}
		for message := range messages {
			c.Plugin.Process(message)
		}
	}
}
