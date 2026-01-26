package email

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// DNSRecord represents a DNS record that needs to be configured
type DNSRecord struct {
	Type     string `json:"type"`     // CNAME, TXT, MX
	Name     string `json:"name"`     // Record name/host
	Value    string `json:"value"`    // Record value
	Priority int    `json:"priority"` // For MX records
	Purpose  string `json:"purpose"`  // dkim, verification, mail_from
}

// DomainIdentityResult contains the result of verifying a domain
type DomainIdentityResult struct {
	Domain             string       `json:"domain"`
	VerificationToken  string       `json:"verification_token"`
	VerificationStatus string       `json:"verification_status"`
	DKIMTokens         []string     `json:"dkim_tokens"`
	DKIMStatus         string       `json:"dkim_status"`
	DNSRecords         []*DNSRecord `json:"dns_records"`
}

// DomainIdentityStatus contains the verification status of a domain
type DomainIdentityStatus struct {
	Domain             string `json:"domain"`
	VerificationStatus string `json:"verification_status"`
	DKIMStatus         string `json:"dkim_status"`
	MailFromStatus     string `json:"mail_from_status"`
}

// getSESClient creates an SES client with the provided credentials
func getSESClient(ctx context.Context, region, accessKey, secretKey string) (*ses.Client, error) {
	if region == "" {
		region = "us-east-1"
	}

	var cfg aws.Config
	var err error

	if accessKey != "" && secretKey != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return ses.NewFromConfig(cfg), nil
}

// getSNSClient creates an SNS client with the provided credentials
func getSNSClient(ctx context.Context, region, accessKey, secretKey string) (*sns.Client, error) {
	if region == "" {
		region = "us-east-1"
	}

	var cfg aws.Config
	var err error

	if accessKey != "" && secretKey != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return sns.NewFromConfig(cfg), nil
}

// SetupBounceNotifications configures SES to send bounce/complaint notifications to our webhook
// This creates an SNS topic and subscribes our webhook endpoint, then configures SES to use it
// If notifications are already configured, this is a no-op
func SetupBounceNotifications(ctx context.Context, region, accessKey, secretKey, domain, webhookURL string) error {
	sesClient, err := getSESClient(ctx, region, accessKey, secretKey)
	if err != nil {
		return err
	}

	// Check if bounce notifications are already configured
	notifAttrs, err := sesClient.GetIdentityNotificationAttributes(ctx, &ses.GetIdentityNotificationAttributesInput{
		Identities: []string{domain},
	})
	if err == nil {
		if attrs, ok := notifAttrs.NotificationAttributes[domain]; ok {
			// If bounce topic is already set, don't modify anything
			if attrs.BounceTopic != nil && *attrs.BounceTopic != "" {
				return nil
			}
		}
	}

	snsClient, err := getSNSClient(ctx, region, accessKey, secretKey)
	if err != nil {
		return err
	}

	// Create SNS topic for this domain (idempotent - returns existing topic if exists)
	// Use a sanitized domain name for the topic name
	topicName := fmt.Sprintf("outlet-ses-%s", strings.ReplaceAll(domain, ".", "-"))
	createTopicResult, err := snsClient.CreateTopic(ctx, &sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil {
		return fmt.Errorf("failed to create SNS topic: %w", err)
	}
	topicArn := aws.ToString(createTopicResult.TopicArn)

	// Subscribe our webhook to the topic (idempotent)
	_, err = snsClient.Subscribe(ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicArn),
		Protocol: aws.String("https"),
		Endpoint: aws.String(webhookURL),
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe webhook to SNS topic: %w", err)
	}

	// Configure SES to send bounce notifications to the SNS topic
	_, err = sesClient.SetIdentityNotificationTopic(ctx, &ses.SetIdentityNotificationTopicInput{
		Identity:         aws.String(domain),
		NotificationType: types.NotificationTypeBounce,
		SnsTopic:         aws.String(topicArn),
	})
	if err != nil {
		return fmt.Errorf("failed to set bounce notification topic: %w", err)
	}

	// Configure SES to send complaint notifications to the SNS topic
	_, err = sesClient.SetIdentityNotificationTopic(ctx, &ses.SetIdentityNotificationTopicInput{
		Identity:         aws.String(domain),
		NotificationType: types.NotificationTypeComplaint,
		SnsTopic:         aws.String(topicArn),
	})
	if err != nil {
		return fmt.Errorf("failed to set complaint notification topic: %w", err)
	}

	// Configure SES to send delivery notifications to the SNS topic
	_, err = sesClient.SetIdentityNotificationTopic(ctx, &ses.SetIdentityNotificationTopicInput{
		Identity:         aws.String(domain),
		NotificationType: types.NotificationTypeDelivery,
		SnsTopic:         aws.String(topicArn),
	})
	if err != nil {
		// Non-fatal - delivery notifications are nice to have but not critical
	}

	// Disable email notifications for bounces/complaints (we're using SNS)
	_, err = sesClient.SetIdentityFeedbackForwardingEnabled(ctx, &ses.SetIdentityFeedbackForwardingEnabledInput{
		Identity:          aws.String(domain),
		ForwardingEnabled: false,
	})
	if err != nil {
		// Non-fatal - email forwarding can still work alongside SNS
	}

	return nil
}

// VerifyDomainIdentity initiates domain verification with AWS SES and returns DNS records
// This sets up:
// 1. Domain verification (TXT record)
// 2. DKIM signing (CNAME records)
// 3. Custom MAIL FROM domain for SPF alignment (MX + TXT records)
func VerifyDomainIdentity(ctx context.Context, region, accessKey, secretKey, domain string) (*DomainIdentityResult, error) {
	client, err := getSESClient(ctx, region, accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	// Verify the domain identity
	verifyResult, err := client.VerifyDomainIdentity(ctx, &ses.VerifyDomainIdentityInput{
		Domain: aws.String(domain),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify domain identity: %w", err)
	}

	result := &DomainIdentityResult{
		Domain:             domain,
		VerificationToken:  aws.ToString(verifyResult.VerificationToken),
		VerificationStatus: "pending",
		DNSRecords:         []*DNSRecord{},
	}

	// Add TXT record for domain verification
	result.DNSRecords = append(result.DNSRecords, &DNSRecord{
		Type:    "TXT",
		Name:    fmt.Sprintf("_amazonses.%s", domain),
		Value:   aws.ToString(verifyResult.VerificationToken),
		Purpose: "verification",
	})

	// Enable DKIM signing and get DKIM tokens
	dkimResult, err := client.VerifyDomainDkim(ctx, &ses.VerifyDomainDkimInput{
		Domain: aws.String(domain),
	})
	if err != nil {
		// Don't fail completely, DKIM can be added later
		result.DKIMStatus = "failed"
	} else {
		result.DKIMTokens = dkimResult.DkimTokens
		result.DKIMStatus = "pending"

		// Add CNAME records for DKIM
		for _, token := range dkimResult.DkimTokens {
			result.DNSRecords = append(result.DNSRecords, &DNSRecord{
				Type:    "CNAME",
				Name:    fmt.Sprintf("%s._domainkey.%s", token, domain),
				Value:   fmt.Sprintf("%s.dkim.amazonses.com", token),
				Purpose: "dkim",
			})
		}
	}

	// Set up custom MAIL FROM domain for SPF alignment
	// This ensures the envelope sender (Return-Path) uses your domain instead of amazonses.com
	// Benefits: SPF alignment with your domain, better DMARC, cleaner bounce handling
	mailFromSubdomain := "mail"
	mailFromDomain := fmt.Sprintf("%s.%s", mailFromSubdomain, domain)

	_, mailFromErr := client.SetIdentityMailFromDomain(ctx, &ses.SetIdentityMailFromDomainInput{
		Identity:            aws.String(domain),
		MailFromDomain:      aws.String(mailFromDomain),
		BehaviorOnMXFailure: types.BehaviorOnMXFailureUseDefaultValue,
	})
	if mailFromErr != nil {
		// Log but don't fail - MAIL FROM can be set up later
		// The domain will still work, just with amazonses.com as envelope sender
	} else {
		// Add DNS records for custom MAIL FROM
		result.DNSRecords = append(result.DNSRecords, &DNSRecord{
			Type:     "MX",
			Name:     mailFromDomain,
			Value:    fmt.Sprintf("feedback-smtp.%s.amazonses.com", region),
			Priority: 10,
			Purpose:  "mail_from",
		})
		result.DNSRecords = append(result.DNSRecords, &DNSRecord{
			Type:    "TXT",
			Name:    mailFromDomain,
			Value:   "v=spf1 include:amazonses.com ~all",
			Purpose: "mail_from",
		})
	}

	// Add DMARC record for email authentication policy
	// p=none allows monitoring without rejecting - recommended for initial setup
	result.DNSRecords = append(result.DNSRecords, &DNSRecord{
		Type:    "TXT",
		Name:    fmt.Sprintf("_dmarc.%s", domain),
		Value:   "v=DMARC1; p=none;",
		Purpose: "dmarc",
	})

	return result, nil
}

// GetDomainIdentityStatus checks the current verification status of a domain
func GetDomainIdentityStatus(ctx context.Context, region, accessKey, secretKey, domain string) (*DomainIdentityStatus, error) {
	client, err := getSESClient(ctx, region, accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	// Get verification attributes
	attrsResult, err := client.GetIdentityVerificationAttributes(ctx, &ses.GetIdentityVerificationAttributesInput{
		Identities: []string{domain},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get verification attributes: %w", err)
	}

	status := &DomainIdentityStatus{
		Domain:             domain,
		VerificationStatus: "not_started",
		DKIMStatus:         "not_started",
		MailFromStatus:     "not_started",
	}

	if attrs, ok := attrsResult.VerificationAttributes[domain]; ok {
		status.VerificationStatus = mapVerificationStatus(attrs.VerificationStatus)
	}

	// Get DKIM attributes
	dkimResult, err := client.GetIdentityDkimAttributes(ctx, &ses.GetIdentityDkimAttributesInput{
		Identities: []string{domain},
	})
	if err == nil {
		if dkim, ok := dkimResult.DkimAttributes[domain]; ok {
			status.DKIMStatus = mapDKIMStatus(dkim.DkimVerificationStatus)
		}
	}

	// Get Mail From attributes
	mailFromResult, err := client.GetIdentityMailFromDomainAttributes(ctx, &ses.GetIdentityMailFromDomainAttributesInput{
		Identities: []string{domain},
	})
	if err == nil {
		if mailFrom, ok := mailFromResult.MailFromDomainAttributes[domain]; ok {
			status.MailFromStatus = mapMailFromStatus(mailFrom.MailFromDomainStatus)
		}
	}

	return status, nil
}

// SetupMailFrom configures a custom MAIL FROM domain
func SetupMailFrom(ctx context.Context, region, accessKey, secretKey, domain, mailFromSubdomain string) ([]*DNSRecord, error) {
	client, err := getSESClient(ctx, region, accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	mailFromDomain := fmt.Sprintf("%s.%s", mailFromSubdomain, domain)

	_, err = client.SetIdentityMailFromDomain(ctx, &ses.SetIdentityMailFromDomainInput{
		Identity:         aws.String(domain),
		MailFromDomain:   aws.String(mailFromDomain),
		BehaviorOnMXFailure: types.BehaviorOnMXFailureUseDefaultValue,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set mail from domain: %w", err)
	}

	// Return the DNS records needed
	records := []*DNSRecord{
		{
			Type:     "MX",
			Name:     mailFromDomain,
			Value:    fmt.Sprintf("feedback-smtp.%s.amazonses.com", region),
			Priority: 10,
			Purpose:  "mail_from",
		},
		{
			Type:    "TXT",
			Name:    mailFromDomain,
			Value:   "v=spf1 include:amazonses.com ~all",
			Purpose: "mail_from",
		},
	}

	return records, nil
}

// DeleteDomainIdentity removes a domain identity from SES
func DeleteDomainIdentity(ctx context.Context, region, accessKey, secretKey, domain string) error {
	client, err := getSESClient(ctx, region, accessKey, secretKey)
	if err != nil {
		return err
	}

	_, err = client.DeleteIdentity(ctx, &ses.DeleteIdentityInput{
		Identity: aws.String(domain),
	})
	if err != nil {
		return fmt.Errorf("failed to delete domain identity: %w", err)
	}

	return nil
}

// ExtractDomainFromEmail extracts the domain from an email address
func ExtractDomainFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return strings.ToLower(parts[1])
}

// DNSRecordsToJSON converts DNS records to JSON string
func DNSRecordsToJSON(records []*DNSRecord) (string, error) {
	data, err := json.Marshal(records)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DNSRecordsFromJSON parses DNS records from JSON string
func DNSRecordsFromJSON(data string) ([]*DNSRecord, error) {
	var records []*DNSRecord
	if err := json.Unmarshal([]byte(data), &records); err != nil {
		return nil, err
	}
	return records, nil
}

// Helper functions to map AWS status types to our string statuses
func mapVerificationStatus(status types.VerificationStatus) string {
	switch status {
	case types.VerificationStatusPending:
		return "pending"
	case types.VerificationStatusSuccess:
		return "success"
	case types.VerificationStatusFailed:
		return "failed"
	case types.VerificationStatusTemporaryFailure:
		return "temporary_failure"
	case types.VerificationStatusNotStarted:
		return "not_started"
	default:
		return "pending"
	}
}

func mapDKIMStatus(status types.VerificationStatus) string {
	return mapVerificationStatus(status)
}

func mapMailFromStatus(status types.CustomMailFromStatus) string {
	switch status {
	case types.CustomMailFromStatusPending:
		return "pending"
	case types.CustomMailFromStatusSuccess:
		return "success"
	case types.CustomMailFromStatusFailed:
		return "failed"
	default:
		return "not_started"
	}
}
