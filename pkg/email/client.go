package email

import (
	"context"
	"fmt"

	envConfig "gcstatus/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SendEmailFunc func(recipient, body, subject string) error

func Send(recipient string, body string, subject string) error {
	env := envConfig.LoadConfig()

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(env.AwsMailRegion))
	if err != nil {
		return fmt.Errorf("failed to load configuration, %v", err)
	}

	svc := ses.NewFromConfig(cfg)

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				recipient,
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(env.AwsMailFrom),
	}

	_, err = svc.SendEmail(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to send email, %v", err)
	}

	return nil
}
