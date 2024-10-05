package sqs

import (
	"context"
	"fmt"

	envconfig "gcstatus/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSClientInterface interface {
	SendMessage(ctx context.Context, queueUrl string, messageBody string) error
	GetAWSClient() *sqs.Client
}

type SQSClient struct {
	client *sqs.Client
}

var GlobalSQSClient SQSClientInterface

func NewSQSClient() *SQSClient {
	env := envconfig.LoadConfig()
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(env.AwsBucketRegion))
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config: %+v", err))
	}

	return &SQSClient{
		client: sqs.NewFromConfig(awsCfg),
	}
}

func (s *SQSClient) SendMessage(ctx context.Context, queueUrl string, messageBody string) error {
	_, err := s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueUrl),
		MessageBody: aws.String(messageBody),
	})

	return err
}

func (s *SQSClient) GetAWSClient() *sqs.Client {
	return s.client
}

func GetAwsQueue() string {
	env := envconfig.LoadConfig()

	return env.AwsSqsUrl
}
