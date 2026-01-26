package websocket

import "time"

// Message types
const (
	TypePing                     = "ping"
	TypePong                     = "pong"
	TypeDomainIdentityUpdate     = "domain_identity_update"
	TypeDomainIdentityCreated    = "domain_identity_created"
	TypeDomainIdentityVerified   = "domain_identity_verified"
	TypeSubscribe                = "subscribe"
	TypeUnsubscribe              = "unsubscribe"
	TypeBackupUpdate             = "backup_update"
)

// Message is the base WebSocket message structure
type Message struct {
	Type      string      `json:"type"`
	Channel   string      `json:"channel,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

// DomainIdentityUpdate is sent when a domain identity status changes
type DomainIdentityUpdate struct {
	ID                 string `json:"id"`
	OrgID              string `json:"org_id"`
	Domain             string `json:"domain"`
	VerificationStatus string `json:"verification_status"`
	DKIMStatus         string `json:"dkim_status"`
	MailFromStatus     string `json:"mail_from_status"`
	LastCheckedAt      string `json:"last_checked_at"`
}

// NewMessage creates a new message with timestamp
func NewMessage(msgType string, data interface{}) *Message {
	return &Message{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewDomainIdentityUpdate creates a domain identity update message
func NewDomainIdentityUpdate(id, orgID, domain, verificationStatus, dkimStatus, mailFromStatus, lastCheckedAt string) *Message {
	return NewMessage(TypeDomainIdentityUpdate, DomainIdentityUpdate{
		ID:                 id,
		OrgID:              orgID,
		Domain:             domain,
		VerificationStatus: verificationStatus,
		DKIMStatus:         dkimStatus,
		MailFromStatus:     mailFromStatus,
		LastCheckedAt:      lastCheckedAt,
	})
}

// BackupUpdate is sent when a backup status changes
type BackupUpdate struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Filename string `json:"filename"`
	FileSize int64  `json:"file_size"`
	Error    string `json:"error,omitempty"`
}

// NewBackupUpdate creates a backup update message
func NewBackupUpdate(id, status, filename string, fileSize int64, errorMsg string) *Message {
	return NewMessage(TypeBackupUpdate, BackupUpdate{
		ID:       id,
		Status:   status,
		Filename: filename,
		FileSize: fileSize,
		Error:    errorMsg,
	})
}
