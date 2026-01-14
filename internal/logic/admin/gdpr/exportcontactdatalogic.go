package gdpr

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExportContactDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportContactDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportContactDataLogic {
	return &ExportContactDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ContactExport represents the full data export for a contact
type ContactExport struct {
	ExportedAt string                 `json:"exported_at"`
	ContactID  string                 `json:"contact_id"`
	Profile    map[string]interface{} `json:"profile"`
	Tags       []string               `json:"tags"`
	Activity   []map[string]interface{} `json:"activity,omitempty"`
	Consent    map[string]interface{} `json:"consent"`
}

func (l *ExportContactDataLogic) ExportContactData(req *types.GDPRExportRequest) (resp *types.GDPRExportResponse, err error) {
	// Get contact
	contact, err := l.svcCtx.DB.GetContact(l.ctx, req.ContactId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contact not found")
		}
		return nil, err
	}

	// Build export data
	export := ContactExport{
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		ContactID:  contact.ID,
		Profile: map[string]interface{}{
			"email":             contact.Email,
			"name":              contact.Name,
			"source":            contact.Source.String,
			"email_verified":    contact.EmailVerified == 1,
			"verified_at":       contact.VerifiedAt.String,
			"unsubscribed_at":   contact.UnsubscribedAt.String,
			"created_at":        contact.CreatedAt.String,
			"updated_at":        contact.UpdatedAt.String,
		},
		Consent: map[string]interface{}{
			"gdpr_consent":    contact.GdprConsent.Valid && contact.GdprConsent.Int64 == 1,
			"gdpr_consent_at": contact.GdprConsentAt.String,
			"is_subscribed":   contact.UnsubscribedAt.String == "",
		},
	}

	// Get tags
	tags, err := l.svcCtx.DB.GetContactTags(l.ctx, sql.NullString{String: contact.ID, Valid: true})
	if err == nil {
		export.Tags = make([]string, len(tags))
		for i, t := range tags {
			export.Tags[i] = t.Tag
		}
	}

	// TODO: Add email activity, campaign interactions, etc.

	// Write export to file
	exportDir := filepath.Join(".", "exports", "gdpr")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create export directory: %w", err)
	}

	exportID := uuid.New().String()
	filename := fmt.Sprintf("gdpr-export-%s-%s.json", contact.ID[:8], exportID[:8])
	filePath := filepath.Join(exportDir, filename)

	exportData, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to serialize export: %w", err)
	}

	if err := os.WriteFile(filePath, exportData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write export file: %w", err)
	}

	// Generate download URL (expires in 24 hours)
	expiresAt := time.Now().Add(24 * time.Hour)

	return &types.GDPRExportResponse{
		DownloadURL: fmt.Sprintf("/api/admin/gdpr/export/%s/download", exportID),
		ExpiresAt:   expiresAt.Format(time.RFC3339),
		Format:      "json",
	}, nil
}
