# Outlet.sh Completion Tasks

This document outlines all remaining work to complete the Outlet.sh platform. Tasks are organized by priority and category.

---

## Priority 1: Fix Broken Tests

These must be fixed first as they prevent `go test ./...` from passing.

### Task 1.1: Delete or Rewrite Email Service Tests

**Files:**
- `internal/services/email/sender_test.go`
- `internal/services/email/sequences_test.go`

**Problem:**
Tests use outdated `NewService` signature that accepted individual SMTP parameters:
```go
// OLD (broken)
NewService(smtpHost, smtpPort, smtpUser, smtpPass, fromAddress, fromName, replyTo)

// NEW (current)
NewService(store *db.Store)
```

**Solution Options:**

**Option A: Delete the tests** (quick fix)
```bash
rm internal/services/email/sender_test.go
rm internal/services/email/sequences_test.go
```

**Option B: Rewrite tests** (proper fix)

Create mock store and rewrite tests to work with database-backed configuration:

```go
// Example test setup
func TestEmailService(t *testing.T) {
    // Create in-memory SQLite for testing
    conn, _ := sql.Open("sqlite", ":memory:")
    store := db.NewStore(conn)

    // Run migrations
    migrations.Run(conn)

    // Seed test SMTP config in platform_settings
    store.UpsertPlatformSetting(ctx, db.UpsertPlatformSettingParams{
        Category:  "email",
        Key:       "smtp_host",
        ValueText: sql.NullString{String: "localhost", Valid: true},
    })
    // ... seed other settings

    service := email.NewService(store)
    // ... run tests
}
```

### Task 1.2: Fix Workers Test Nil Pointer

**File:** `internal/workers/email_worker_test.go`

**Problem:**
Line 182 passes `nil` for sender, causing panic:
```go
service := email.NewSequenceService(nil, nil) // panics
```

**Solution:**
Either delete these tests or fix `NewSequenceService` to handle nil sender:

**Option A: Fix the function** (`internal/services/email/sequences.go` line 28-33)
```go
func NewSequenceService(db *db.Store, sender *Service) *SequenceService {
    baseURL := "https://outlet.sh" // default
    if sender != nil {
        baseURL = sender.GetBaseURL()
    }
    return &SequenceService{
        db:      db,
        sender:  sender,
        baseURL: baseURL,
    }
}
```

**Option B: Delete the tests**
```bash
# Remove lines 180-220 from email_worker_test.go
```

---

## Priority 2: Implement SDK Contact Endpoints

These endpoints are exposed in the API but return nil. They're needed for the SDK and MCP integration.

### Task 2.1: GetContact

**File:** `internal/logic/sdk/contacts/getcontactlogic.go`

**Implementation:**
```go
func (l *GetContactLogic) GetContact(req *types.GetContactRequest) (*types.SDKContactInfo, error) {
    contact, err := l.svcCtx.DB.GetContactByID(l.ctx, req.ContactId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errorx.NewAPIError(404, "Contact not found")
        }
        return nil, err
    }

    // Get tags for contact
    tags, _ := l.svcCtx.DB.GetContactTags(l.ctx, contact.ID)

    return &types.SDKContactInfo{
        Id:            contact.ID,
        Email:         contact.Email,
        Name:          contact.Name.String,
        Status:        contact.Status,
        Tags:          tagsToStrings(tags),
        CustomFields:  jsonToMap(contact.CustomFields),
        CreatedAt:     contact.CreatedAt.Format(time.RFC3339),
        UpdatedAt:     contact.UpdatedAt.Format(time.RFC3339),
    }, nil
}
```

### Task 2.2: UpdateContact

**File:** `internal/logic/sdk/contacts/updatecontactlogic.go`

**Implementation:**
```go
func (l *UpdateContactLogic) UpdateContact(req *types.UpdateContactRequest) (*types.SDKContactInfo, error) {
    // Verify contact exists
    contact, err := l.svcCtx.DB.GetContactByID(l.ctx, req.ContactId)
    if err != nil {
        return nil, errorx.NewAPIError(404, "Contact not found")
    }

    // Update fields
    params := db.UpdateContactParams{
        ID: req.ContactId,
    }
    if req.Name != "" {
        params.Name = sql.NullString{String: req.Name, Valid: true}
    }
    if req.Status != "" {
        params.Status = req.Status
    }
    if req.CustomFields != nil {
        params.CustomFields = mapToJSON(req.CustomFields)
    }

    updated, err := l.svcCtx.DB.UpdateContact(l.ctx, params)
    if err != nil {
        return nil, err
    }

    return contactToSDKInfo(updated), nil
}
```

### Task 2.3: AddContactTags

**File:** `internal/logic/sdk/contacts/addcontacttagslogic.go`

**Implementation:**
```go
func (l *AddContactTagsLogic) AddContactTags(req *types.AddContactTagsRequest) (*types.SDKContactInfo, error) {
    // Verify contact exists
    contact, err := l.svcCtx.DB.GetContactByID(l.ctx, req.ContactId)
    if err != nil {
        return nil, errorx.NewAPIError(404, "Contact not found")
    }

    // Add each tag
    for _, tagName := range req.Tags {
        // Get or create tag
        tag, err := l.svcCtx.DB.GetOrCreateTag(l.ctx, db.GetOrCreateTagParams{
            OrgID: contact.OrgID,
            Name:  tagName,
        })
        if err != nil {
            continue
        }

        // Link tag to contact
        l.svcCtx.DB.AddContactTag(l.ctx, db.AddContactTagParams{
            ContactID: contact.ID,
            TagID:     tag.ID,
        })
    }

    // Emit event for rules engine
    l.svcCtx.Events.Publish(events.Event{
        Type:      events.ContactTagAdded,
        ContactID: contact.ID,
        Data:      map[string]interface{}{"tags": req.Tags},
    })

    return l.getContactInfo(contact.ID)
}
```

### Task 2.4: RemoveContactTags

**File:** `internal/logic/sdk/contacts/removecontacttagslogic.go`

**Implementation:**
```go
func (l *RemoveContactTagsLogic) RemoveContactTags(req *types.RemoveContactTagsRequest) (*types.SDKContactInfo, error) {
    contact, err := l.svcCtx.DB.GetContactByID(l.ctx, req.ContactId)
    if err != nil {
        return nil, errorx.NewAPIError(404, "Contact not found")
    }

    for _, tagName := range req.Tags {
        tag, err := l.svcCtx.DB.GetTagByName(l.ctx, db.GetTagByNameParams{
            OrgID: contact.OrgID,
            Name:  tagName,
        })
        if err != nil {
            continue
        }

        l.svcCtx.DB.RemoveContactTag(l.ctx, db.RemoveContactTagParams{
            ContactID: contact.ID,
            TagID:     tag.ID,
        })
    }

    return l.getContactInfo(contact.ID)
}
```

### Task 2.5: ListContactActivity

**File:** `internal/logic/sdk/contacts/listcontactactivitylogic.go`

**Implementation:**
```go
func (l *ListContactActivityLogic) ListContactActivity(req *types.ListContactActivityRequest) (*types.ContactActivityResponse, error) {
    // Get email events for this contact
    events, err := l.svcCtx.DB.GetContactEmailEvents(l.ctx, db.GetContactEmailEventsParams{
        ContactID: req.ContactId,
        Limit:     int32(req.Limit),
        Offset:    int32(req.Offset),
    })
    if err != nil {
        return nil, err
    }

    activities := make([]types.ContactActivity, len(events))
    for i, event := range events {
        activities[i] = types.ContactActivity{
            Type:      event.EventType,
            Timestamp: event.CreatedAt.Format(time.RFC3339),
            Details:   jsonToMap(event.Metadata),
        }
    }

    return &types.ContactActivityResponse{
        Activities: activities,
        Total:      len(activities),
    }, nil
}
```

### Task 2.6: GlobalUnsubscribe

**File:** `internal/logic/sdk/contacts/globalunsubscribelogic.go`

**Implementation:**
```go
func (l *GlobalUnsubscribeLogic) GlobalUnsubscribe(req *types.GlobalUnsubscribeRequest) (*types.GlobalUnsubscribeResponse, error) {
    // Add to global blocklist
    _, err := l.svcCtx.DB.AddToBlocklist(l.ctx, db.AddToBlocklistParams{
        Email:  req.Email,
        Reason: "global_unsubscribe",
        Source: "sdk",
    })
    if err != nil {
        return nil, err
    }

    // Update all contacts with this email to unsubscribed
    err = l.svcCtx.DB.UnsubscribeAllByEmail(l.ctx, req.Email)
    if err != nil {
        return nil, err
    }

    return &types.GlobalUnsubscribeResponse{
        Success: true,
        Message: "Email globally unsubscribed",
    }, nil
}
```

---

## Priority 3: Implement Import Jobs

CSV import functionality for bulk subscriber management.

### Task 3.1: Database Schema

First, add import_jobs table if not exists. Check `internal/db/migrations/` for existing schema.

```sql
-- Add to migrations if needed
CREATE TABLE IF NOT EXISTS import_jobs (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL REFERENCES organizations(id),
    list_id TEXT NOT NULL REFERENCES email_lists(id),
    status TEXT NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed, cancelled
    filename TEXT,
    total_rows INTEGER DEFAULT 0,
    processed_rows INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    error_count INTEGER DEFAULT 0,
    errors_json TEXT, -- JSON array of error details
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME
);
```

### Task 3.2: ListImportJobs

**File:** `internal/logic/admin/imports/listimportjobslogic.go`

```go
func (l *ListImportJobsLogic) ListImportJobs(req *types.ListImportJobsRequest) (*types.ListImportJobsResponse, error) {
    orgID := l.ctx.Value(middleware.OrgIDKey).(string)

    jobs, err := l.svcCtx.DB.ListImportJobs(l.ctx, db.ListImportJobsParams{
        OrgID:  orgID,
        Limit:  int32(req.PageSize),
        Offset: int32((req.Page - 1) * req.PageSize),
    })
    if err != nil {
        return nil, err
    }

    total, _ := l.svcCtx.DB.CountImportJobs(l.ctx, orgID)

    return &types.ListImportJobsResponse{
        Jobs:  jobsToTypes(jobs),
        Total: int(total),
        Page:  req.Page,
    }, nil
}
```

### Task 3.3: GetImportJob

**File:** `internal/logic/admin/imports/getimportjoblogic.go`

```go
func (l *GetImportJobLogic) GetImportJob(req *types.GetImportJobRequest) (*types.ImportJob, error) {
    job, err := l.svcCtx.DB.GetImportJob(l.ctx, req.Id)
    if err != nil {
        return nil, errorx.NewAPIError(404, "Import job not found")
    }

    return jobToType(job), nil
}
```

### Task 3.4: CancelImportJob

**File:** `internal/logic/admin/imports/cancelimportjoblogic.go`

```go
func (l *CancelImportJobLogic) CancelImportJob(req *types.CancelImportJobRequest) (*types.ImportJob, error) {
    job, err := l.svcCtx.DB.GetImportJob(l.ctx, req.Id)
    if err != nil {
        return nil, errorx.NewAPIError(404, "Import job not found")
    }

    if job.Status != "pending" && job.Status != "processing" {
        return nil, errorx.NewAPIError(400, "Cannot cancel job in status: "+job.Status)
    }

    updated, err := l.svcCtx.DB.UpdateImportJobStatus(l.ctx, db.UpdateImportJobStatusParams{
        ID:     req.Id,
        Status: "cancelled",
    })
    if err != nil {
        return nil, err
    }

    return jobToType(updated), nil
}
```

---

## Priority 4: Implement Housekeeping

Subscriber cleanup functionality.

### Task 4.1: HousekeepingUnconfirmed

**File:** `internal/logic/admin/housekeeping/housekeepingunconfirmedlogic.go`

```go
func (l *HousekeepingUnconfirmedLogic) HousekeepingUnconfirmed(req *types.HousekeepingUnconfirmedRequest) (*types.HousekeepingResponse, error) {
    orgID := l.ctx.Value(middleware.OrgIDKey).(string)

    cutoffDate := time.Now().AddDate(0, 0, -req.OlderThanDays)

    if req.DryRun {
        // Count only
        count, err := l.svcCtx.DB.CountUnconfirmedContacts(l.ctx, db.CountUnconfirmedContactsParams{
            OrgID:      orgID,
            CutoffDate: cutoffDate,
        })
        if err != nil {
            return nil, err
        }

        return &types.HousekeepingResponse{
            AffectedCount: int(count),
            DryRun:        true,
            Message:       fmt.Sprintf("Would delete %d unconfirmed contacts older than %d days", count, req.OlderThanDays),
        }, nil
    }

    // Actually delete
    result, err := l.svcCtx.DB.DeleteUnconfirmedContacts(l.ctx, db.DeleteUnconfirmedContactsParams{
        OrgID:      orgID,
        CutoffDate: cutoffDate,
    })
    if err != nil {
        return nil, err
    }

    affected, _ := result.RowsAffected()

    return &types.HousekeepingResponse{
        AffectedCount: int(affected),
        DryRun:        false,
        Message:       fmt.Sprintf("Deleted %d unconfirmed contacts", affected),
    }, nil
}
```

### Task 4.2: HousekeepingInactive

**File:** `internal/logic/admin/housekeeping/housekeepinginactivelogic.go`

```go
func (l *HousekeepingInactiveLogic) HousekeepingInactive(req *types.HousekeepingInactiveRequest) (*types.HousekeepingResponse, error) {
    orgID := l.ctx.Value(middleware.OrgIDKey).(string)

    cutoffDate := time.Now().AddDate(0, 0, -req.InactiveDays)

    if req.DryRun {
        count, err := l.svcCtx.DB.CountInactiveContacts(l.ctx, db.CountInactiveContactsParams{
            OrgID:      orgID,
            CutoffDate: cutoffDate,
        })
        if err != nil {
            return nil, err
        }

        return &types.HousekeepingResponse{
            AffectedCount: int(count),
            DryRun:        true,
            Message:       fmt.Sprintf("Would delete %d contacts with no activity in %d days", count, req.InactiveDays),
        }, nil
    }

    result, err := l.svcCtx.DB.DeleteInactiveContacts(l.ctx, db.DeleteInactiveContactsParams{
        OrgID:      orgID,
        CutoffDate: cutoffDate,
    })
    if err != nil {
        return nil, err
    }

    affected, _ := result.RowsAffected()

    return &types.HousekeepingResponse{
        AffectedCount: int(affected),
        DryRun:        false,
        Message:       fmt.Sprintf("Deleted %d inactive contacts", affected),
    }, nil
}
```

### Task 4.3: Add SQL Queries

**File:** `internal/db/queries/contacts.sql` (append)

```sql
-- name: CountUnconfirmedContacts :one
SELECT COUNT(*) FROM contacts
WHERE org_id = ? AND status = 'unconfirmed' AND created_at < ?;

-- name: DeleteUnconfirmedContacts :execresult
DELETE FROM contacts
WHERE org_id = ? AND status = 'unconfirmed' AND created_at < ?;

-- name: CountInactiveContacts :one
SELECT COUNT(*) FROM contacts c
WHERE c.org_id = ?
AND c.status = 'active'
AND NOT EXISTS (
    SELECT 1 FROM email_events e
    WHERE e.contact_id = c.id AND e.created_at > ?
);

-- name: DeleteInactiveContacts :execresult
DELETE FROM contacts
WHERE org_id = ?
AND status = 'active'
AND id NOT IN (
    SELECT DISTINCT contact_id FROM email_events
    WHERE created_at > ?
);
```

Then run: `make sqlc-gen`

---

## Priority 5: Implement Rules Engine Actions

The rules engine has action handlers but they're placeholders. These need real implementations.

### Task 5.1: send_email Action

**File:** `internal/services/rules/actions.go`

```go
func (r *ActionRegistry) RegisterEmailService(emailService *email.Service) {
    r.Register("send_email", func(ctx context.Context, params map[string]interface{}, dryRun bool) error {
        to, _ := params["to"].(string)
        templateID, _ := params["template_id"].(string)
        subject, _ := params["subject"].(string)
        body, _ := params["body"].(string)

        if to == "" {
            return errors.New("send_email: to is required")
        }

        if dryRun {
            slog.InfoContext(ctx, "dry-run: would send email", "to", to, "subject", subject)
            return nil
        }

        // Get org ID from context for SMTP config
        orgID, _ := ctx.Value("org_id").(string)

        if templateID != "" {
            return emailService.SendTemplateEmail(ctx, orgID, to, templateID, params)
        }

        return emailService.SendEmailWithOrgConfig(ctx, orgID, to, subject, body)
    })
}
```

### Task 5.2: tag_customer Action

```go
func (r *ActionRegistry) RegisterTagAction(db *db.Store) {
    r.Register("tag_customer", func(ctx context.Context, params map[string]interface{}, dryRun bool) error {
        contactID, _ := params["contact_id"].(string)
        tagName, _ := params["tag"].(string)

        if contactID == "" || tagName == "" {
            return errors.New("tag_customer: contact_id and tag are required")
        }

        if dryRun {
            slog.InfoContext(ctx, "dry-run: would tag customer", "contact_id", contactID, "tag", tagName)
            return nil
        }

        contact, err := db.GetContactByID(ctx, contactID)
        if err != nil {
            return err
        }

        tag, err := db.GetOrCreateTag(ctx, db.GetOrCreateTagParams{
            OrgID: contact.OrgID,
            Name:  tagName,
        })
        if err != nil {
            return err
        }

        return db.AddContactTag(ctx, db.AddContactTagParams{
            ContactID: contactID,
            TagID:     tag.ID,
        })
    })
}
```

### Task 5.3: notify_admin Action

```go
func (r *ActionRegistry) RegisterNotifyAdmin(emailService *email.Service, adminEmail string) {
    r.Register("notify_admin", func(ctx context.Context, params map[string]interface{}, dryRun bool) error {
        subject, _ := params["subject"].(string)
        message, _ := params["message"].(string)

        if subject == "" {
            subject = "Admin Notification"
        }

        if dryRun {
            slog.InfoContext(ctx, "dry-run: would notify admin", "subject", subject)
            return nil
        }

        body := fmt.Sprintf("<html><body><h2>%s</h2><p>%s</p></body></html>", subject, message)
        return emailService.SendEmail(ctx, adminEmail, subject, body)
    })
}
```

### Task 5.4: log_event Action

```go
func (r *ActionRegistry) RegisterLogEvent(db *db.Store) {
    r.Register("log_event", func(ctx context.Context, params map[string]interface{}, dryRun bool) error {
        eventType, _ := params["event_type"].(string)
        entityID, _ := params["entity_id"].(string)
        metadata, _ := params["metadata"].(map[string]interface{})

        if dryRun {
            slog.InfoContext(ctx, "dry-run: would log event", "type", eventType, "entity", entityID)
            return nil
        }

        metadataJSON, _ := json.Marshal(metadata)

        _, err := db.CreateEventLog(ctx, db.CreateEventLogParams{
            ID:        uuid.NewString(),
            EventType: eventType,
            EntityID:  entityID,
            Metadata:  string(metadataJSON),
        })
        return err
    })
}
```

---

## Priority 6: Frontend Fixes

### Task 6.1: Sequence Toggle Active

**File:** `app/src/routes/[orgSlug]/(admin)/sequences/+page.svelte`

**Line 322:** Replace empty onclick handler:

```svelte
{
    label: sequence.is_active ? 'Pause' : 'Activate',
    icon: sequence.is_active ? Pause : Play,
    onclick: async () => {
        try {
            await api.adminUpdateSequence({
                id: sequence.id,
                is_active: !sequence.is_active
            });
            // Refresh the list
            invalidateAll();
            toast.success(sequence.is_active ? 'Sequence paused' : 'Sequence activated');
        } catch (err) {
            toast.error('Failed to update sequence');
        }
    }
}
```

### Task 6.2: Rules Entity Options

**File:** `app/src/routes/[orgSlug]/(admin)/settings/automation/rules/+page.svelte`

**Line 424:** Load entity options dynamically:

```svelte
{#if formEntityType}
    <div>
        <label for="rule-entity-id" class="form-label">Entity ID</label>
        <Select id="rule-entity-id" bind:value={formEntityId}>
            <option value="">Select {formEntityType}...</option>
            {#if formEntityType === 'email_list'}
                {#each lists as list}
                    <option value={list.id}>{list.name}</option>
                {/each}
            {:else if formEntityType === 'sequence'}
                {#each sequences as seq}
                    <option value={seq.id}>{seq.name}</option>
                {/each}
            {:else if formEntityType === 'campaign'}
                {#each campaigns as campaign}
                    <option value={campaign.id}>{campaign.name}</option>
                {/each}
            {/if}
        </Select>
    </div>
{/if}
```

Also add data loading in the page's `load` function:
```typescript
// In +page.ts
export const load = async ({ parent, fetch }) => {
    const { org } = await parent();

    const [lists, sequences, campaigns] = await Promise.all([
        api.adminListEmailLists({ org_id: org.id }),
        api.adminListSequences({ org_id: org.id }),
        api.adminListCampaigns({ org_id: org.id })
    ]);

    return { lists, sequences, campaigns };
};
```

---

## Priority 7: Dashboard Stats

### Task 7.1: Complete Dashboard Stats

**File:** `internal/logic/admin/organizations/getdashboardstatslogic.go`

Add email and subscriber statistics:

```go
func (l *GetDashboardStatsLogic) GetDashboardStats(req *types.GetDashboardStatsRequest) (*types.DashboardStatsResponse, error) {
    orgID := req.Id
    resp := &types.DashboardStatsResponse{}

    // Subscriber counts
    subscriberStats, err := l.svcCtx.DB.GetSubscriberStats(l.ctx, orgID)
    if err == nil {
        resp.TotalSubscribers = int(subscriberStats.Total)
        resp.ActiveSubscribers = int(subscriberStats.Active)
        resp.UnsubscribedCount = int(subscriberStats.Unsubscribed)
    }

    // Email stats (last 30 days)
    emailStats, err := l.svcCtx.DB.GetEmailStats30Days(l.ctx, orgID)
    if err == nil {
        resp.EmailsSent = int(emailStats.Sent)
        resp.EmailsOpened = int(emailStats.Opened)
        resp.EmailsClicked = int(emailStats.Clicked)
        resp.OpenRate = calculateRate(emailStats.Opened, emailStats.Sent)
        resp.ClickRate = calculateRate(emailStats.Clicked, emailStats.Opened)
    }

    // Existing checks
    lists, err := l.svcCtx.DB.ListEmailLists(l.ctx, orgID)
    resp.HasLists = err == nil && len(lists) > 0

    org, err := l.svcCtx.DB.GetOrganizationByID(l.ctx, orgID)
    if err == nil {
        resp.EmailConfigured = org.FromEmail.Valid && org.FromEmail.String != ""
    }

    // MCP config check
    resp.HasMCPConfigured = false
    if userID, ok := l.ctx.Value("userId").(string); ok {
        mcpKeys, err := l.svcCtx.DB.ListMCPAPIKeysByUser(l.ctx, userID)
        if err == nil {
            for _, key := range mcpKeys {
                if !key.RevokedAt.Valid {
                    resp.HasMCPConfigured = true
                    break
                }
            }
        }
    }

    return resp, nil
}
```

### Task 7.2: Add SQL Queries for Stats

**File:** `internal/db/queries/contacts.sql` (append)

```sql
-- name: GetSubscriberStats :one
SELECT
    COUNT(*) as total,
    COUNT(CASE WHEN status = 'active' THEN 1 END) as active,
    COUNT(CASE WHEN status = 'unsubscribed' THEN 1 END) as unsubscribed
FROM contacts WHERE org_id = ?;

-- name: GetEmailStats30Days :one
SELECT
    COUNT(CASE WHEN event_type = 'sent' THEN 1 END) as sent,
    COUNT(CASE WHEN event_type = 'opened' THEN 1 END) as opened,
    COUNT(CASE WHEN event_type = 'clicked' THEN 1 END) as clicked
FROM email_events
WHERE org_id = ? AND created_at > datetime('now', '-30 days');
```

---

## Verification Checklist

After completing all tasks, verify:

- [ ] `go build .` succeeds
- [ ] `go test ./...` passes (all tests)
- [ ] `cd app && pnpm build` succeeds
- [ ] `cd app && pnpm test` passes
- [ ] Manual test: Create contact via SDK
- [ ] Manual test: Add/remove tags via SDK
- [ ] Manual test: Run housekeeping dry-run
- [ ] Manual test: Toggle sequence active state
- [ ] Manual test: View dashboard stats

---

## Estimated Effort

| Priority | Tasks | Effort |
|----------|-------|--------|
| P1: Fix Tests | 3 | 1-2 hours |
| P2: SDK Contacts | 6 | 2-3 hours |
| P3: Import Jobs | 4 | 2-3 hours |
| P4: Housekeeping | 3 | 1-2 hours |
| P5: Rules Actions | 4 | 2-3 hours |
| P6: Frontend Fixes | 2 | 1 hour |
| P7: Dashboard Stats | 2 | 1 hour |

**Total: ~10-15 hours**

---

## Order of Execution

1. Fix broken tests (P1) - unblocks CI
2. Implement SDK contact endpoints (P2) - needed for MCP
3. Frontend fixes (P6) - quick wins
4. Dashboard stats (P7) - user-facing
5. Housekeeping (P4) - subscriber management
6. Import jobs (P3) - bulk operations
7. Rules actions (P5) - automation

---

*Last updated: January 2026*
