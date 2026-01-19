package email

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// SESConfig holds AWS SES configuration
type SESConfig struct {
	Region      string
	AccessKey   string
	SecretKey   string
	FromAddress string
	FromName    string
	ReplyTo     string
}

// SendEmailViaSES sends an email using the AWS SES API directly
// This is the preferred method when AWS credentials are configured
func SendEmailViaSES(ctx context.Context, sesConfig *SESConfig, to, subject, htmlBody string) error {
	if sesConfig.Region == "" {
		sesConfig.Region = "us-east-1"
	}

	var cfg aws.Config
	var err error

	if sesConfig.AccessKey != "" && sesConfig.SecretKey != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(sesConfig.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				sesConfig.AccessKey,
				sesConfig.SecretKey,
				"",
			)),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(sesConfig.Region))
	}

	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := ses.NewFromConfig(cfg)

	// Build the from address
	from := sesConfig.FromAddress
	if sesConfig.FromName != "" {
		from = fmt.Sprintf("%s <%s>", sesConfig.FromName, sesConfig.FromAddress)
	}

	// Build reply-to addresses
	var replyToAddresses []string
	if sesConfig.ReplyTo != "" {
		replyToAddresses = []string{sesConfig.ReplyTo}
	}

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source:           aws.String(from),
		ReplyToAddresses: replyToAddresses,
	}

	_, err = client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("SES SendEmail failed: %w", err)
	}

	return nil
}
