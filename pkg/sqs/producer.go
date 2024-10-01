package sqs

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSProducer struct {
	client   *sqs.Client
	queueUrl string
}

func NewSQSProducer(client *sqs.Client, queueUrl string) *SQSProducer {
	return &SQSProducer{
		client:   client,
		queueUrl: queueUrl,
	}
}

func (p *SQSProducer) EnqueueMessage(ctx context.Context, messageType string, messageBody interface{}) error {
	message := map[string]interface{}{
		"type": messageType,
		"body": messageBody,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueUrl),
		MessageBody: aws.String(string(messageBytes)),
	})

	if err != nil {
		log.Printf("Failed to send SQS message: %v", err)
		return err
	}

	log.Printf("Successfully enqueued message of type: %s", messageType)
	return nil
}
