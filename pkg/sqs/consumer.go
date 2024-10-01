package sqs

import (
	"context"
	"encoding/json"
	"gcstatus/internal/usecases"
	"gcstatus/pkg/sqs/messages"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSConsumer struct {
	client          *sqs.Client
	queueUrl        string
	purchaseHandler *messages.PurchaseMessageHandler
}

func NewSQSConsumer(
	client *sqs.Client,
	queueUrl string,
	userService *usecases.UserService,
	transactionService *usecases.TransactionService,
	notificationService *usecases.NotificationService,
) *SQSConsumer {
	purchaseHandler := messages.NewPurchaseMessageHandler(
		userService,
		transactionService,
		notificationService,
	)

	return &SQSConsumer{
		client:          client,
		queueUrl:        queueUrl,
		purchaseHandler: purchaseHandler,
	}
}

func (c *SQSConsumer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping SQS consumer...")
			return
		default:
			c.pollMessages(ctx)
		}
	}
}

func (c *SQSConsumer) pollMessages(ctx context.Context) {
	output, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueUrl),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     10,
	})

	if err != nil {
		log.Printf("Failed to receive messages: %v", err)
		return
	}

	for _, message := range output.Messages {
		go c.processMessage(ctx, message)
	}
}

func (c *SQSConsumer) processMessage(ctx context.Context, message types.Message) {
	log.Printf("Processing message: %s", *message.Body)

	var messageType struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal([]byte(*message.Body), &messageType); err != nil {
		log.Printf("Error unmarshalling message type: %v", err)
		return
	}

	switch messageType.Type {
	case "PurchaseTitle":
		c.purchaseHandler.HandlePurchaseMessage(ctx, message)
	default:
		log.Printf("Unknown message type: %s", messageType.Type)
	}

	if err := c.deleteMessage(ctx, message); err != nil {
		log.Printf("Failed to delete message: %+v", err)
	}
}

func (c *SQSConsumer) deleteMessage(ctx context.Context, message types.Message) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueUrl),
		ReceiptHandle: message.ReceiptHandle,
	})

	return err
}
