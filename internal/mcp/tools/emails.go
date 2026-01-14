package tools

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"outlet/internal/db"
	"outlet/internal/mcp/mcpctx"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func int64ToBool(i sql.NullInt64) bool {
	return i.Valid && i.Int64 == 1
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	return slug
}

// RegisterEmailTools registers all email-related MCP tools.
func RegisterEmailTools(server *mcp.Server, toolCtx *mcpctx.ToolContext) {
	// list_create
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_create",
		Description: "Create a list for organizing subscribers. Returns list ID for use with sequence_create. Lists group subscribers for targeted email campaigns.",
	}, createListHandler(toolCtx))

	// list_list
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_list",
		Description: "List all lists for the organization. Returns list IDs, names, and subscriber counts.",
	}, listListsHandler(toolCtx))

	// sequence_create
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sequence_create",
		Description: "Create an email sequence (automated drip campaign). Returns sequence ID for use with sequence_email_add. Sequences send emails automatically based on triggers. Timing: Each email's delay_hours is relative to the PREVIOUS email (or trigger event for position 1). Example: Position 1 with delay_hours=0 sends immediately on trigger; Position 2 with delay_hours=24 sends 24 hours after Position 1.",
	}, createSequenceHandler(toolCtx))

	// sequence_list
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sequence_list",
		Description: "List all email sequences for the organization. Returns sequence IDs, names, trigger events, and email counts.",
	}, listSequencesHandler(toolCtx))

	// sequence_email_add
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sequence_email_add",
		Description: "Add an email to a sequence at a specific position. Timing: delay_hours controls when this email sends AFTER the previous position completes (or after trigger event for position 1). Example: delay_hours=0 sends immediately; delay_hours=24 sends 24 hours later; delay_hours=168 sends 1 week later.",
	}, addSequenceEmailHandler(toolCtx))

	// sequence_get
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sequence_get",
		Description: "Get a sequence with all its emails, including delay_hours for each. Each email's delay_hours is relative to the previous email (not the trigger). Use this to verify email timing and content.",
	}, getSequenceHandler(toolCtx))

	// list_get
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_get",
		Description: "Get a single list by ID with subscriber count and settings.",
	}, getListHandler(toolCtx))

	// list_mutate
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_mutate",
		Description: "Apply changes to a list. Pass the fields you want to update in the 'set' object. Only provided fields are changed.",
	}, mutateListHandler(toolCtx))

	// list_delete
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_delete",
		Description: "Delete a list and all its subscriber associations. Subscribers are not deleted, only removed from the list.",
	}, deleteListHandler(toolCtx))

	// sequence_mutate
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sequence_mutate",
		Description: "Apply changes to a sequence. Pass the fields you want to update in the 'set' object. Use set.active to toggle active status. Only provided fields are changed.",
	}, mutateSequenceHandler(toolCtx))

	// sequence_delete
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sequence_delete",
		Description: "Delete a sequence and all its emails. This is permanent and cannot be undone.",
	}, deleteSequenceHandler(toolCtx))

	// sequence_email_mutate
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sequence_email_mutate",
		Description: "Apply changes to a sequence email. Pass the fields you want to update in the 'set' object. Use set.position to move the email. Note: Changing delay_hours affects when this email sends relative to the previous email.",
	}, mutateSequenceEmailHandler(toolCtx))

	// sequence_email_list
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sequence_email_list",
		Description: "List all emails in a sequence with their positions, subjects, and delay_hours. Each delay_hours is relative to the previous email position.",
	}, listSequenceEmailsHandler(toolCtx))

	// sequence_email_get
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sequence_email_get",
		Description: "Get a single sequence email by ID with full content and settings.",
	}, getSequenceEmailHandler(toolCtx))

	// sequence_email_delete
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sequence_email_delete",
		Description: "Delete an email from a sequence. Other emails' positions are not automatically adjusted.",
	}, deleteSequenceEmailHandler(toolCtx))
}

// ListCreateInput defines input for list_create tool.
type ListCreateInput struct {
	Name        string `json:"name" jsonschema:"required,List name (e.g., 'Newsletter Subscribers')"`
	Description string `json:"description,omitempty" jsonschema:"List description for internal reference"`
	DoubleOptin *bool  `json:"double_optin,omitempty" jsonschema:"Require email confirmation (default: false)"`
}

// ListCreateOutput defines output for list_create tool.
type ListCreateOutput struct {
	ID          string `json:"id" jsonschema:"List ID (opaque string)"`
	Name        string `json:"name" jsonschema:"List name"`
	Slug        string `json:"slug" jsonschema:"URL-friendly slug"`
	Description string `json:"description,omitempty" jsonschema:"List description"`
	DoubleOptin bool   `json:"double_optin" jsonschema:"Whether double opt-in is required"`
	Created     bool   `json:"created" jsonschema:"Whether the list was successfully created"`
}

func createListHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input ListCreateInput) (*mcp.CallToolResult, ListCreateOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ListCreateInput) (*mcp.CallToolResult, ListCreateOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, ListCreateOutput{}, err
		}

		// Validate name
		if strings.TrimSpace(input.Name) == "" {
			return nil, ListCreateOutput{}, mcpctx.NewValidationError("name is required", "name")
		}

		// Generate slug from name
		slug := generateSlug(input.Name)

		// Default double optin to false
		doubleOptin := false
		if input.DoubleOptin != nil {
			doubleOptin = *input.DoubleOptin
		}

		// Create list
		list, err := toolCtx.DB().CreateEmailList(ctx, db.CreateEmailListParams{
			OrgID:       toolCtx.OrgID(),
			Name:        input.Name,
			Slug:        slug,
			Description: sql.NullString{String: input.Description, Valid: input.Description != ""},
			DoubleOptin: sql.NullInt64{Int64: boolToInt64(doubleOptin), Valid: true},
		})
		if err != nil {
			return nil, ListCreateOutput{}, fmt.Errorf("failed to create list: %w", err)
		}

		return nil, ListCreateOutput{
			ID:          fmt.Sprintf("%d", list.ID),
			Name:        list.Name,
			Slug:        list.Slug,
			Description: list.Description.String,
			DoubleOptin: int64ToBool(list.DoubleOptin),
			Created:     true,
		}, nil
	}
}

// ListListInput defines input for list_list tool.
type ListListInput struct{}

// ListItem represents a single list in the list.
type ListItem struct {
	ID              string `json:"id" jsonschema:"List ID (opaque string)"`
	Name            string `json:"name" jsonschema:"List name"`
	Slug            string `json:"slug" jsonschema:"URL-friendly slug"`
	Description     string `json:"description,omitempty" jsonschema:"List description"`
	DoubleOptin     bool   `json:"double_optin" jsonschema:"Whether double opt-in is required"`
	SubscriberCount int64  `json:"subscriber_count" jsonschema:"Number of active subscribers"`
}

// ListListOutput defines output for list_list tool.
type ListListOutput struct {
	Lists []ListItem `json:"lists" jsonschema:"List of lists"`
	Total int        `json:"total" jsonschema:"Total number of lists"`
}

func listListsHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input ListListInput) (*mcp.CallToolResult, ListListOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ListListInput) (*mcp.CallToolResult, ListListOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, ListListOutput{}, err
		}

		lists, err := toolCtx.DB().ListEmailLists(ctx, toolCtx.OrgID())
		if err != nil {
			return nil, ListListOutput{}, fmt.Errorf("failed to list lists: %w", err)
		}

		items := make([]ListItem, 0, len(lists))
		for _, list := range lists {
			// Get subscriber count (active only)
			count, err := toolCtx.DB().CountListSubscribers(ctx, db.CountListSubscribersParams{
				ListID:       list.ID,
				FilterStatus: "active",
			})
			if err != nil {
				count = 0
			}

			items = append(items, ListItem{
				ID:              fmt.Sprintf("%d", list.ID),
				Name:            list.Name,
				Slug:            list.Slug,
				Description:     list.Description.String,
				DoubleOptin:     int64ToBool(list.DoubleOptin),
				SubscriberCount: count,
			})
		}

		return nil, ListListOutput{
			Lists: items,
			Total: len(items),
		}, nil
	}
}

// SequenceCreateInput defines input for sequence_create tool.
type SequenceCreateInput struct {
	Name         string `json:"name" jsonschema:"required,Sequence name (e.g., 'Welcome Series')"`
	ListID       string `json:"list_id" jsonschema:"required,List ID to attach sequence to (opaque string)"`
	TriggerEvent string `json:"trigger_event,omitempty" jsonschema:"Event that starts sequence: on_subscribe, on_purchase, manual (default: on_subscribe)"`
	SequenceType string `json:"sequence_type,omitempty" jsonschema:"Type: lifecycle, promotional, transactional (default: lifecycle)"`
	Active       *bool  `json:"active,omitempty" jsonschema:"Whether sequence is active (default: true)"`
}

// SequenceCreateOutput defines output for sequence_create tool.
type SequenceCreateOutput struct {
	ID           string `json:"id" jsonschema:"Sequence ID"`
	Name         string `json:"name" jsonschema:"Sequence name"`
	Slug         string `json:"slug" jsonschema:"URL-friendly slug"`
	ListID       string `json:"list_id" jsonschema:"Attached list ID (opaque string)"`
	TriggerEvent string `json:"trigger_event" jsonschema:"Trigger event"`
	SequenceType string `json:"sequence_type" jsonschema:"Sequence type"`
	Active       bool   `json:"active" jsonschema:"Whether sequence is active"`
	Created      bool   `json:"created" jsonschema:"Whether the sequence was successfully created"`
}

func createSequenceHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input SequenceCreateInput) (*mcp.CallToolResult, SequenceCreateOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SequenceCreateInput) (*mcp.CallToolResult, SequenceCreateOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, SequenceCreateOutput{}, err
		}

		// Validate name
		if strings.TrimSpace(input.Name) == "" {
			return nil, SequenceCreateOutput{}, mcpctx.NewValidationError("name is required", "name")
		}

		// Validate and parse list_id
		if strings.TrimSpace(input.ListID) == "" {
			return nil, SequenceCreateOutput{}, mcpctx.NewValidationError("list_id is required", "list_id")
		}
		listID, err := strconv.ParseInt(input.ListID, 10, 64)
		if err != nil {
			return nil, SequenceCreateOutput{}, mcpctx.NewValidationError("list_id must be a valid ID", "list_id")
		}

		// Generate slug
		slug := generateSlug(input.Name)

		// Default trigger event
		triggerEvent := input.TriggerEvent
		if triggerEvent == "" {
			triggerEvent = "on_subscribe"
		}

		// Default sequence type
		sequenceType := input.SequenceType
		if sequenceType == "" {
			sequenceType = "lifecycle"
		}

		// Default active to true
		active := true
		if input.Active != nil {
			active = *input.Active
		}

		// Create sequence
		orgID := toolCtx.OrgID()
		sequenceID := uuid.New().String()
		sequence, err := toolCtx.DB().CreateSequence(ctx, db.CreateSequenceParams{
			ID:           sequenceID,
			OrgID:        sql.NullString{String: orgID, Valid: true},
			ListID:       sql.NullInt64{Int64: listID, Valid: true},
			Slug:         slug,
			Name:         input.Name,
			TriggerEvent: triggerEvent,
			IsActive:     sql.NullInt64{Int64: boolToInt64(active), Valid: true},
			SequenceType: sql.NullString{String: sequenceType, Valid: true},
		})
		if err != nil {
			return nil, SequenceCreateOutput{}, fmt.Errorf("failed to create sequence: %w", err)
		}

		return nil, SequenceCreateOutput{
			ID:           sequence.ID,
			Name:         sequence.Name,
			Slug:         sequence.Slug,
			ListID:       fmt.Sprintf("%d", sequence.ListID.Int64),
			TriggerEvent: sequence.TriggerEvent,
			SequenceType: sequence.SequenceType.String,
			Active:       int64ToBool(sequence.IsActive),
			Created:      true,
		}, nil
	}
}

// SequenceListInput defines input for sequence_list tool.
type SequenceListInput struct {
	ListID string `json:"list_id,omitempty" jsonschema:"Filter by list ID (opaque string)"`
}

// SequenceListItem represents a single sequence in the list.
type SequenceListItem struct {
	ID           string `json:"id" jsonschema:"Sequence ID"`
	Name         string `json:"name" jsonschema:"Sequence name"`
	Slug         string `json:"slug" jsonschema:"URL-friendly slug"`
	ListID       string `json:"list_id" jsonschema:"Attached list ID (opaque string)"`
	ListName     string `json:"list_name,omitempty" jsonschema:"Attached list name"`
	TriggerEvent string `json:"trigger_event" jsonschema:"Trigger event"`
	SequenceType string `json:"sequence_type" jsonschema:"Sequence type"`
	Active       bool   `json:"active" jsonschema:"Whether sequence is active"`
	EmailCount   int    `json:"email_count" jsonschema:"Number of emails in sequence"`
}

// SequenceListOutput defines output for sequence_list tool.
type SequenceListOutput struct {
	Sequences []SequenceListItem `json:"sequences" jsonschema:"List of sequences"`
	Total     int                `json:"total" jsonschema:"Total number of sequences"`
}

func listSequencesHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input SequenceListInput) (*mcp.CallToolResult, SequenceListOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SequenceListInput) (*mcp.CallToolResult, SequenceListOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, SequenceListOutput{}, err
		}

		// Parse optional list_id filter
		var filterListID int64
		if strings.TrimSpace(input.ListID) != "" {
			var err error
			filterListID, err = strconv.ParseInt(input.ListID, 10, 64)
			if err != nil {
				return nil, SequenceListOutput{}, mcpctx.NewValidationError("list_id must be a valid ID", "list_id")
			}
		}

		items := make([]SequenceListItem, 0)

		orgID := toolCtx.OrgID()
		sequences, err := toolCtx.DB().ListSequencesByOrg(ctx, sql.NullString{String: orgID, Valid: true})
		if err != nil {
			return nil, SequenceListOutput{}, fmt.Errorf("failed to list sequences: %w", err)
		}

		for _, seq := range sequences {
			if filterListID > 0 && seq.ListID.Int64 != filterListID {
				continue
			}

			templates, err := toolCtx.DB().ListTemplatesBySequence(ctx, sql.NullString{String: seq.ID, Valid: true})
			emailCount := 0
			if err == nil {
				emailCount = len(templates)
			}

			items = append(items, SequenceListItem{
				ID:           seq.ID,
				Name:         seq.Name,
				Slug:         seq.Slug,
				ListID:       fmt.Sprintf("%d", seq.ListID.Int64),
				ListName:     seq.ListName.String,
				TriggerEvent: seq.TriggerEvent,
				SequenceType: seq.SequenceType.String,
				Active:       int64ToBool(seq.IsActive),
				EmailCount:   emailCount,
			})
		}

		return nil, SequenceListOutput{
			Sequences: items,
			Total:     len(items),
		}, nil
	}
}

// SequenceEmailAddInput defines input for sequence_email_add tool.
type SequenceEmailAddInput struct {
	SequenceID string `json:"sequence_id" jsonschema:"required,Sequence ID to add email to"`
	Subject    string `json:"subject" jsonschema:"required,Email subject line"`
	HTMLBody   string `json:"html_body" jsonschema:"required,HTML content of the email"`
	PlainText  string `json:"plain_text,omitempty" jsonschema:"Plain text version (auto-generated if omitted)"`
	Position   int    `json:"position,omitempty" jsonschema:"Position in sequence (1-based). Default: appends to end"`
	DelayHours int    `json:"delay_hours,omitempty" jsonschema:"Hours to wait after PREVIOUS email before sending this one. 0=immediate, 24=1 day, 168=1 week. For position 1, delay is relative to trigger event."`
	Active     *bool  `json:"active,omitempty" jsonschema:"Whether email is active (default: true)"`
}

// SequenceEmailAddOutput defines output for sequence_email_add tool.
type SequenceEmailAddOutput struct {
	ID         string `json:"id" jsonschema:"Template/email ID"`
	SequenceID string `json:"sequence_id" jsonschema:"Parent sequence ID"`
	Subject    string `json:"subject" jsonschema:"Email subject"`
	Position   int    `json:"position" jsonschema:"Position in sequence"`
	DelayHours int    `json:"delay_hours" jsonschema:"Hours delay before sending"`
	Active     bool   `json:"active" jsonschema:"Whether email is active"`
	Created    bool   `json:"created" jsonschema:"Whether the email was successfully added"`
}

func addSequenceEmailHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input SequenceEmailAddInput) (*mcp.CallToolResult, SequenceEmailAddOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SequenceEmailAddInput) (*mcp.CallToolResult, SequenceEmailAddOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, SequenceEmailAddOutput{}, err
		}

		if strings.TrimSpace(input.SequenceID) == "" {
			return nil, SequenceEmailAddOutput{}, mcpctx.NewValidationError("sequence_id is required", "sequence_id")
		}

		if strings.TrimSpace(input.Subject) == "" {
			return nil, SequenceEmailAddOutput{}, mcpctx.NewValidationError("subject is required", "subject")
		}

		if strings.TrimSpace(input.HTMLBody) == "" {
			return nil, SequenceEmailAddOutput{}, mcpctx.NewValidationError("html_body is required", "html_body")
		}

		position := input.Position
		if position <= 0 {
			templates, err := toolCtx.DB().ListTemplatesBySequence(ctx, sql.NullString{String: input.SequenceID, Valid: true})
			if err == nil {
				position = len(templates) + 1
			} else {
				position = 1
			}
		}

		active := true
		if input.Active != nil {
			active = *input.Active
		}

		templateID := uuid.New().String()
		template, err := toolCtx.DB().CreateTemplate(ctx, db.CreateTemplateParams{
			ID:           templateID,
			SequenceID:   sql.NullString{String: input.SequenceID, Valid: true},
			Position:     int64(position),
			DelayHours:   int64(input.DelayHours),
			Subject:      input.Subject,
			HtmlBody:     input.HTMLBody,
			PlainText:    sql.NullString{String: input.PlainText, Valid: input.PlainText != ""},
			TemplateType: sql.NullString{String: "sequence", Valid: true},
			IsActive:     sql.NullInt64{Int64: boolToInt64(active), Valid: true},
		})
		if err != nil {
			return nil, SequenceEmailAddOutput{}, fmt.Errorf("failed to add email to sequence: %w", err)
		}

		return nil, SequenceEmailAddOutput{
			ID:         template.ID,
			SequenceID: template.SequenceID.String,
			Subject:    template.Subject,
			Position:   int(template.Position),
			DelayHours: int(template.DelayHours),
			Active:     int64ToBool(template.IsActive),
			Created:    true,
		}, nil
	}
}

// SequenceGetInput defines input for sequence_get tool.
type SequenceGetInput struct {
	SequenceID string `json:"sequence_id" jsonschema:"required,Sequence ID to get details for"`
}

// SequenceEmail represents an email in a sequence.
type SequenceEmail struct {
	ID         string `json:"id" jsonschema:"Email/template ID"`
	Position   int    `json:"position" jsonschema:"Position in sequence (1-based)"`
	Subject    string `json:"subject" jsonschema:"Email subject line"`
	DelayHours int    `json:"delay_hours" jsonschema:"Hours after previous email (or trigger for position 1)"`
	Active     bool   `json:"active" jsonschema:"Whether email is active"`
}

// SequenceGetOutput defines output for sequence_get tool.
type SequenceGetOutput struct {
	ID           string          `json:"id" jsonschema:"Sequence ID"`
	Name         string          `json:"name" jsonschema:"Sequence name"`
	Slug         string          `json:"slug" jsonschema:"URL-friendly slug"`
	ListID       string          `json:"list_id" jsonschema:"Attached list ID (opaque string)"`
	TriggerEvent string          `json:"trigger_event" jsonschema:"Trigger event"`
	SequenceType string          `json:"sequence_type" jsonschema:"Sequence type"`
	Active       bool            `json:"active" jsonschema:"Whether sequence is active"`
	Emails       []SequenceEmail `json:"emails" jsonschema:"All emails in the sequence with timing"`
}

func getSequenceHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input SequenceGetInput) (*mcp.CallToolResult, SequenceGetOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SequenceGetInput) (*mcp.CallToolResult, SequenceGetOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, SequenceGetOutput{}, err
		}

		if strings.TrimSpace(input.SequenceID) == "" {
			return nil, SequenceGetOutput{}, mcpctx.NewValidationError("sequence_id is required", "sequence_id")
		}

		seq, err := toolCtx.DB().GetSequenceByID(ctx, input.SequenceID)
		if err != nil {
			return nil, SequenceGetOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("sequence %s not found", input.SequenceID))
		}

		if seq.OrgID.String != toolCtx.OrgID() {
			return nil, SequenceGetOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("sequence %s not found", input.SequenceID))
		}

		templates, err := toolCtx.DB().ListTemplatesBySequence(ctx, sql.NullString{String: input.SequenceID, Valid: true})
		if err != nil {
			return nil, SequenceGetOutput{}, fmt.Errorf("failed to get sequence emails: %w", err)
		}

		emails := make([]SequenceEmail, 0, len(templates))
		for _, t := range templates {
			emails = append(emails, SequenceEmail{
				ID:         t.ID,
				Position:   int(t.Position),
				Subject:    t.Subject,
				DelayHours: int(t.DelayHours),
				Active:     int64ToBool(t.IsActive),
			})
		}

		return nil, SequenceGetOutput{
			ID:           seq.ID,
			Name:         seq.Name,
			Slug:         seq.Slug,
			ListID:       fmt.Sprintf("%d", seq.ListID.Int64),
			TriggerEvent: seq.TriggerEvent,
			SequenceType: seq.SequenceType.String,
			Active:       int64ToBool(seq.IsActive),
			Emails:       emails,
		}, nil
	}
}

// ========== sequence_email_mutate ==========

// SequenceEmailMutateSet defines the fields that can be mutated on a sequence email.
type SequenceEmailMutateSet struct {
	Subject    *string `json:"subject,omitempty" jsonschema:"New subject line"`
	HTMLBody   *string `json:"html_body,omitempty" jsonschema:"New HTML body content"`
	PlainText  *string `json:"plain_text,omitempty" jsonschema:"New plain text content"`
	DelayHours *int    `json:"delay_hours,omitempty" jsonschema:"Hours after previous email (or trigger for position 1)"`
	Position   *int    `json:"position,omitempty" jsonschema:"New position in sequence (1-based) - use this to move the email"`
	Active     *bool   `json:"active,omitempty" jsonschema:"Whether email is active"`
}

// SequenceEmailMutateInput defines input for sequence_email_mutate tool.
type SequenceEmailMutateInput struct {
	ID  string                 `json:"id" jsonschema:"required,Email/template ID to mutate"`
	Set SequenceEmailMutateSet `json:"set" jsonschema:"required,Fields to update"`
}

// SequenceEmailMutateOutput defines output for sequence_email_mutate tool.
type SequenceEmailMutateOutput struct {
	ID         string `json:"id" jsonschema:"Email/template ID"`
	SequenceID string `json:"sequence_id" jsonschema:"Parent sequence ID"`
	Subject    string `json:"subject" jsonschema:"Updated subject line"`
	Position   int    `json:"position" jsonschema:"Position in sequence"`
	DelayHours int    `json:"delay_hours" jsonschema:"Delay in hours"`
	Active     bool   `json:"active" jsonschema:"Whether email is active"`
	Mutated    bool   `json:"mutated" jsonschema:"Whether the mutation was successful"`
}

func mutateSequenceEmailHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input SequenceEmailMutateInput) (*mcp.CallToolResult, SequenceEmailMutateOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SequenceEmailMutateInput) (*mcp.CallToolResult, SequenceEmailMutateOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, SequenceEmailMutateOutput{}, err
		}

		if strings.TrimSpace(input.ID) == "" {
			return nil, SequenceEmailMutateOutput{}, mcpctx.NewValidationError("id is required", "id")
		}

		template, err := toolCtx.DB().GetTemplateByID(ctx, input.ID)
		if err != nil {
			return nil, SequenceEmailMutateOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("email %s not found", input.ID))
		}

		if !template.SequenceID.Valid {
			return nil, SequenceEmailMutateOutput{}, mcpctx.NewValidationError("email is not part of a sequence", "id")
		}

		seq, err := toolCtx.DB().GetSequenceByID(ctx, template.SequenceID.String)
		if err != nil || seq.OrgID.String != toolCtx.OrgID() {
			return nil, SequenceEmailMutateOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("email %s not found", input.ID))
		}

		// Apply mutations from set
		subject := template.Subject
		if input.Set.Subject != nil {
			subject = *input.Set.Subject
		}
		htmlBody := template.HtmlBody
		if input.Set.HTMLBody != nil {
			htmlBody = *input.Set.HTMLBody
		}
		plainText := template.PlainText.String
		if input.Set.PlainText != nil {
			plainText = *input.Set.PlainText
		}
		delayHours := int(template.DelayHours)
		if input.Set.DelayHours != nil {
			delayHours = *input.Set.DelayHours
		}
		position := int(template.Position)
		if input.Set.Position != nil {
			position = *input.Set.Position
		}
		active := int64ToBool(template.IsActive)
		if input.Set.Active != nil {
			active = *input.Set.Active
		}

		err = toolCtx.DB().UpdateTemplate(ctx, db.UpdateTemplateParams{
			ID:           input.ID,
			Position:     int64(position),
			DelayHours:   int64(delayHours),
			Subject:      subject,
			HtmlBody:     htmlBody,
			PlainText:    sql.NullString{String: plainText, Valid: plainText != ""},
			TemplateType: template.TemplateType,
			IsActive:     sql.NullInt64{Int64: boolToInt64(active), Valid: true},
			DesignID:     template.DesignID,
		})
		if err != nil {
			return nil, SequenceEmailMutateOutput{}, fmt.Errorf("failed to mutate email: %w", err)
		}

		return nil, SequenceEmailMutateOutput{
			ID:         input.ID,
			SequenceID: template.SequenceID.String,
			Subject:    subject,
			Position:   position,
			DelayHours: delayHours,
			Active:     active,
			Mutated:    true,
		}, nil
	}
}

// ========== list_get ==========

// ListGetInput defines input for list_get tool.
type ListGetInput struct {
	ID string `json:"id" jsonschema:"required,List ID to retrieve"`
}

// ListGetOutput defines output for list_get tool.
type ListGetOutput struct {
	ID              string `json:"id" jsonschema:"List ID (opaque string)"`
	Name            string `json:"name" jsonschema:"List name"`
	Slug            string `json:"slug" jsonschema:"URL-friendly slug"`
	Description     string `json:"description,omitempty" jsonschema:"List description"`
	DoubleOptin     bool   `json:"double_optin" jsonschema:"Whether double opt-in is required"`
	SubscriberCount int64  `json:"subscriber_count" jsonschema:"Number of active subscribers"`
}

func getListHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input ListGetInput) (*mcp.CallToolResult, ListGetOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ListGetInput) (*mcp.CallToolResult, ListGetOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, ListGetOutput{}, err
		}

		if strings.TrimSpace(input.ID) == "" {
			return nil, ListGetOutput{}, mcpctx.NewValidationError("id is required", "id")
		}

		listID, err := strconv.ParseInt(input.ID, 10, 64)
		if err != nil {
			return nil, ListGetOutput{}, mcpctx.NewValidationError("id must be a valid ID", "id")
		}

		list, err := toolCtx.DB().GetEmailList(ctx, listID)
		if err != nil {
			return nil, ListGetOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("list %s not found", input.ID))
		}

		if list.OrgID != toolCtx.OrgID() {
			return nil, ListGetOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("list %s not found", input.ID))
		}

		count, err := toolCtx.DB().CountListSubscribers(ctx, db.CountListSubscribersParams{
			ListID:       list.ID,
			FilterStatus: "active",
		})
		if err != nil {
			count = 0
		}

		return nil, ListGetOutput{
			ID:              fmt.Sprintf("%d", list.ID),
			Name:            list.Name,
			Slug:            list.Slug,
			Description:     list.Description.String,
			DoubleOptin:     int64ToBool(list.DoubleOptin),
			SubscriberCount: count,
		}, nil
	}
}

// ========== list_mutate ==========

// ListMutateSet defines the fields that can be mutated on a list.
type ListMutateSet struct {
	Name        *string `json:"name,omitempty" jsonschema:"New list name"`
	Description *string `json:"description,omitempty" jsonschema:"New list description"`
	DoubleOptin *bool   `json:"double_optin,omitempty" jsonschema:"Double opt-in setting"`
}

// ListMutateInput defines input for list_mutate tool.
type ListMutateInput struct {
	ID  string        `json:"id" jsonschema:"required,List ID to mutate"`
	Set ListMutateSet `json:"set" jsonschema:"required,Fields to update"`
}

// ListMutateOutput defines output for list_mutate tool.
type ListMutateOutput struct {
	ID          string `json:"id" jsonschema:"List ID (opaque string)"`
	Name        string `json:"name" jsonschema:"Updated list name"`
	Slug        string `json:"slug" jsonschema:"URL-friendly slug"`
	Description string `json:"description,omitempty" jsonschema:"Updated description"`
	DoubleOptin bool   `json:"double_optin" jsonschema:"Whether double opt-in is required"`
	Mutated     bool   `json:"mutated" jsonschema:"Whether the mutation was successful"`
}

func mutateListHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input ListMutateInput) (*mcp.CallToolResult, ListMutateOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ListMutateInput) (*mcp.CallToolResult, ListMutateOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, ListMutateOutput{}, err
		}

		if strings.TrimSpace(input.ID) == "" {
			return nil, ListMutateOutput{}, mcpctx.NewValidationError("id is required", "id")
		}

		listID, err := strconv.ParseInt(input.ID, 10, 64)
		if err != nil {
			return nil, ListMutateOutput{}, mcpctx.NewValidationError("id must be a valid ID", "id")
		}

		list, err := toolCtx.DB().GetEmailList(ctx, listID)
		if err != nil {
			return nil, ListMutateOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("list %s not found", input.ID))
		}

		if list.OrgID != toolCtx.OrgID() {
			return nil, ListMutateOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("list %s not found", input.ID))
		}

		// Apply mutations from set
		name := list.Name
		if input.Set.Name != nil && *input.Set.Name != "" {
			name = *input.Set.Name
		}
		description := list.Description.String
		if input.Set.Description != nil {
			description = *input.Set.Description
		}
		doubleOptin := int64ToBool(list.DoubleOptin)
		if input.Set.DoubleOptin != nil {
			doubleOptin = *input.Set.DoubleOptin
		}

		updatedList, err := toolCtx.DB().UpdateEmailList(ctx, db.UpdateEmailListParams{
			ID:                  listID,
			Name:                name,
			Description:         sql.NullString{String: description, Valid: description != ""},
			DoubleOptin:         sql.NullInt64{Int64: boolToInt64(doubleOptin), Valid: true},
			ConfirmationSubject: list.ConfirmationEmailSubject,
			ConfirmationBody:    list.ConfirmationEmailBody,
		})
		if err != nil {
			return nil, ListMutateOutput{}, fmt.Errorf("failed to mutate list: %w", err)
		}

		return nil, ListMutateOutput{
			ID:          input.ID,
			Name:        updatedList.Name,
			Slug:        updatedList.Slug,
			Description: description,
			DoubleOptin: doubleOptin,
			Mutated:     true,
		}, nil
	}
}

// ========== list_delete ==========

// ListDeleteInput defines input for list_delete tool.
type ListDeleteInput struct {
	ID string `json:"id" jsonschema:"required,List ID to delete"`
}

// ListDeleteOutput defines output for list_delete tool.
type ListDeleteOutput struct {
	Success bool   `json:"success" jsonschema:"Whether deletion was successful"`
	Message string `json:"message" jsonschema:"Status message"`
}

func deleteListHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input ListDeleteInput) (*mcp.CallToolResult, ListDeleteOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ListDeleteInput) (*mcp.CallToolResult, ListDeleteOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, ListDeleteOutput{}, err
		}

		if strings.TrimSpace(input.ID) == "" {
			return nil, ListDeleteOutput{}, mcpctx.NewValidationError("id is required", "id")
		}

		listID, err := strconv.ParseInt(input.ID, 10, 64)
		if err != nil {
			return nil, ListDeleteOutput{}, mcpctx.NewValidationError("id must be a valid ID", "id")
		}

		list, err := toolCtx.DB().GetEmailList(ctx, listID)
		if err != nil {
			return nil, ListDeleteOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("list %s not found", input.ID))
		}

		if list.OrgID != toolCtx.OrgID() {
			return nil, ListDeleteOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("list %s not found", input.ID))
		}

		err = toolCtx.DB().DeleteEmailList(ctx, listID)
		if err != nil {
			return nil, ListDeleteOutput{}, fmt.Errorf("failed to delete list: %w", err)
		}

		return nil, ListDeleteOutput{
			Success: true,
			Message: fmt.Sprintf("List %s deleted successfully", input.ID),
		}, nil
	}
}

// ========== sequence_mutate ==========

// SequenceMutateSet defines the fields that can be mutated on a sequence.
type SequenceMutateSet struct {
	Name         *string `json:"name,omitempty" jsonschema:"New sequence name"`
	TriggerEvent *string `json:"trigger_event,omitempty" jsonschema:"Trigger event: on_subscribe, on_purchase, manual"`
	SequenceType *string `json:"sequence_type,omitempty" jsonschema:"Type: lifecycle, promotional, transactional"`
	Active       *bool   `json:"active,omitempty" jsonschema:"Active status (use this to toggle on/off)"`
}

// SequenceMutateInput defines input for sequence_mutate tool.
type SequenceMutateInput struct {
	ID  string            `json:"id" jsonschema:"required,Sequence ID to mutate"`
	Set SequenceMutateSet `json:"set" jsonschema:"required,Fields to update"`
}

// SequenceMutateOutput defines output for sequence_mutate tool.
type SequenceMutateOutput struct {
	ID           string `json:"id" jsonschema:"Sequence ID"`
	Name         string `json:"name" jsonschema:"Updated sequence name"`
	Slug         string `json:"slug" jsonschema:"URL-friendly slug"`
	TriggerEvent string `json:"trigger_event" jsonschema:"Trigger event"`
	SequenceType string `json:"sequence_type" jsonschema:"Sequence type"`
	Active       bool   `json:"active" jsonschema:"Whether sequence is active"`
	Mutated      bool   `json:"mutated" jsonschema:"Whether the mutation was successful"`
}

func mutateSequenceHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input SequenceMutateInput) (*mcp.CallToolResult, SequenceMutateOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SequenceMutateInput) (*mcp.CallToolResult, SequenceMutateOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, SequenceMutateOutput{}, err
		}

		if strings.TrimSpace(input.ID) == "" {
			return nil, SequenceMutateOutput{}, mcpctx.NewValidationError("id is required", "id")
		}

		seq, err := toolCtx.DB().GetSequenceByID(ctx, input.ID)
		if err != nil {
			return nil, SequenceMutateOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("sequence %s not found", input.ID))
		}

		if seq.OrgID.String != toolCtx.OrgID() {
			return nil, SequenceMutateOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("sequence %s not found", input.ID))
		}

		// Apply mutations from set
		name := seq.Name
		if input.Set.Name != nil && *input.Set.Name != "" {
			name = *input.Set.Name
		}
		triggerEvent := seq.TriggerEvent
		if input.Set.TriggerEvent != nil {
			triggerEvent = *input.Set.TriggerEvent
		}
		sequenceType := seq.SequenceType.String
		if input.Set.SequenceType != nil {
			sequenceType = *input.Set.SequenceType
		}
		active := int64ToBool(seq.IsActive)
		if input.Set.Active != nil {
			active = *input.Set.Active
		}

		err = toolCtx.DB().UpdateSequence(ctx, db.UpdateSequenceParams{
			ID:           input.ID,
			Name:         name,
			TriggerEvent: triggerEvent,
			SequenceType: sql.NullString{String: sequenceType, Valid: true},
			IsActive:     sql.NullInt64{Int64: boolToInt64(active), Valid: true},
			SendHour:     seq.SendHour,
			SendTimezone: seq.SendTimezone,
		})
		if err != nil {
			return nil, SequenceMutateOutput{}, fmt.Errorf("failed to mutate sequence: %w", err)
		}

		return nil, SequenceMutateOutput{
			ID:           input.ID,
			Name:         name,
			Slug:         seq.Slug,
			TriggerEvent: triggerEvent,
			SequenceType: sequenceType,
			Active:       active,
			Mutated:      true,
		}, nil
	}
}

// ========== sequence_delete ==========

// SequenceDeleteInput defines input for sequence_delete tool.
type SequenceDeleteInput struct {
	ID string `json:"id" jsonschema:"required,Sequence ID to delete"`
}

// SequenceDeleteOutput defines output for sequence_delete tool.
type SequenceDeleteOutput struct {
	Success bool   `json:"success" jsonschema:"Whether deletion was successful"`
	Message string `json:"message" jsonschema:"Status message"`
}

func deleteSequenceHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input SequenceDeleteInput) (*mcp.CallToolResult, SequenceDeleteOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SequenceDeleteInput) (*mcp.CallToolResult, SequenceDeleteOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, SequenceDeleteOutput{}, err
		}

		if strings.TrimSpace(input.ID) == "" {
			return nil, SequenceDeleteOutput{}, mcpctx.NewValidationError("id is required", "id")
		}

		seq, err := toolCtx.DB().GetSequenceByID(ctx, input.ID)
		if err != nil {
			return nil, SequenceDeleteOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("sequence %s not found", input.ID))
		}

		if seq.OrgID.String != toolCtx.OrgID() {
			return nil, SequenceDeleteOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("sequence %s not found", input.ID))
		}

		// Delete templates first by listing and deleting each
		templates, err := toolCtx.DB().ListTemplatesBySequence(ctx, sql.NullString{String: input.ID, Valid: true})
		if err != nil {
			return nil, SequenceDeleteOutput{}, fmt.Errorf("failed to list sequence emails: %w", err)
		}
		for _, t := range templates {
			if err := toolCtx.DB().DeleteTemplate(ctx, t.ID); err != nil {
				return nil, SequenceDeleteOutput{}, fmt.Errorf("failed to delete sequence email %s: %w", t.ID, err)
			}
		}

		// Delete sequence
		err = toolCtx.DB().DeleteSequence(ctx, input.ID)
		if err != nil {
			return nil, SequenceDeleteOutput{}, fmt.Errorf("failed to delete sequence: %w", err)
		}

		return nil, SequenceDeleteOutput{
			Success: true,
			Message: fmt.Sprintf("Sequence %s deleted successfully", input.ID),
		}, nil
	}
}

// ========== sequence_email_list ==========

// SequenceEmailListInput defines input for sequence_email_list tool.
type SequenceEmailListInput struct {
	SequenceID string `json:"sequence_id" jsonschema:"required,Sequence ID to list emails for"`
}

// SequenceEmailListItem represents an email in the list.
type SequenceEmailListItem struct {
	ID         string `json:"id" jsonschema:"Email/template ID"`
	Position   int    `json:"position" jsonschema:"Position in sequence (1-based)"`
	Subject    string `json:"subject" jsonschema:"Email subject line"`
	DelayHours int    `json:"delay_hours" jsonschema:"Hours after previous email (or trigger for position 1)"`
	Active     bool   `json:"active" jsonschema:"Whether email is active"`
}

// SequenceEmailListOutput defines output for sequence_email_list tool.
type SequenceEmailListOutput struct {
	SequenceID string                  `json:"sequence_id" jsonschema:"Parent sequence ID"`
	Emails     []SequenceEmailListItem `json:"emails" jsonschema:"List of emails in the sequence"`
	Total      int                     `json:"total" jsonschema:"Total number of emails"`
}

func listSequenceEmailsHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input SequenceEmailListInput) (*mcp.CallToolResult, SequenceEmailListOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SequenceEmailListInput) (*mcp.CallToolResult, SequenceEmailListOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, SequenceEmailListOutput{}, err
		}

		if strings.TrimSpace(input.SequenceID) == "" {
			return nil, SequenceEmailListOutput{}, mcpctx.NewValidationError("sequence_id is required", "sequence_id")
		}

		seq, err := toolCtx.DB().GetSequenceByID(ctx, input.SequenceID)
		if err != nil {
			return nil, SequenceEmailListOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("sequence %s not found", input.SequenceID))
		}

		if seq.OrgID.String != toolCtx.OrgID() {
			return nil, SequenceEmailListOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("sequence %s not found", input.SequenceID))
		}

		templates, err := toolCtx.DB().ListTemplatesBySequence(ctx, sql.NullString{String: input.SequenceID, Valid: true})
		if err != nil {
			return nil, SequenceEmailListOutput{}, fmt.Errorf("failed to list sequence emails: %w", err)
		}

		emails := make([]SequenceEmailListItem, 0, len(templates))
		for _, t := range templates {
			emails = append(emails, SequenceEmailListItem{
				ID:         t.ID,
				Position:   int(t.Position),
				Subject:    t.Subject,
				DelayHours: int(t.DelayHours),
				Active:     int64ToBool(t.IsActive),
			})
		}

		return nil, SequenceEmailListOutput{
			SequenceID: input.SequenceID,
			Emails:     emails,
			Total:      len(emails),
		}, nil
	}
}

// ========== sequence_email_get ==========

// SequenceEmailGetInput defines input for sequence_email_get tool.
type SequenceEmailGetInput struct {
	ID string `json:"id" jsonschema:"required,Email/template ID to retrieve"`
}

// SequenceEmailGetOutput defines output for sequence_email_get tool.
type SequenceEmailGetOutput struct {
	ID         string `json:"id" jsonschema:"Email/template ID"`
	SequenceID string `json:"sequence_id" jsonschema:"Parent sequence ID"`
	Position   int    `json:"position" jsonschema:"Position in sequence (1-based)"`
	Subject    string `json:"subject" jsonschema:"Email subject line"`
	HTMLBody   string `json:"html_body" jsonschema:"HTML content of the email"`
	PlainText  string `json:"plain_text,omitempty" jsonschema:"Plain text version"`
	DelayHours int    `json:"delay_hours" jsonschema:"Hours to wait before sending"`
	Active     bool   `json:"active" jsonschema:"Whether email is active"`
}

func getSequenceEmailHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input SequenceEmailGetInput) (*mcp.CallToolResult, SequenceEmailGetOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SequenceEmailGetInput) (*mcp.CallToolResult, SequenceEmailGetOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, SequenceEmailGetOutput{}, err
		}

		if strings.TrimSpace(input.ID) == "" {
			return nil, SequenceEmailGetOutput{}, mcpctx.NewValidationError("id is required", "id")
		}

		template, err := toolCtx.DB().GetTemplateByID(ctx, input.ID)
		if err != nil {
			return nil, SequenceEmailGetOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("email %s not found", input.ID))
		}

		if !template.SequenceID.Valid {
			return nil, SequenceEmailGetOutput{}, mcpctx.NewValidationError("email is not part of a sequence", "id")
		}

		seq, err := toolCtx.DB().GetSequenceByID(ctx, template.SequenceID.String)
		if err != nil || seq.OrgID.String != toolCtx.OrgID() {
			return nil, SequenceEmailGetOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("email %s not found", input.ID))
		}

		return nil, SequenceEmailGetOutput{
			ID:         template.ID,
			SequenceID: template.SequenceID.String,
			Position:   int(template.Position),
			Subject:    template.Subject,
			HTMLBody:   template.HtmlBody,
			PlainText:  template.PlainText.String,
			DelayHours: int(template.DelayHours),
			Active:     int64ToBool(template.IsActive),
		}, nil
	}
}

// ========== sequence_email_delete ==========

// SequenceEmailDeleteInput defines input for sequence_email_delete tool.
type SequenceEmailDeleteInput struct {
	ID string `json:"id" jsonschema:"required,Email/template ID to delete"`
}

// SequenceEmailDeleteOutput defines output for sequence_email_delete tool.
type SequenceEmailDeleteOutput struct {
	Success bool   `json:"success" jsonschema:"Whether deletion was successful"`
	Message string `json:"message" jsonschema:"Status message"`
}

func deleteSequenceEmailHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input SequenceEmailDeleteInput) (*mcp.CallToolResult, SequenceEmailDeleteOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SequenceEmailDeleteInput) (*mcp.CallToolResult, SequenceEmailDeleteOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, SequenceEmailDeleteOutput{}, err
		}

		if strings.TrimSpace(input.ID) == "" {
			return nil, SequenceEmailDeleteOutput{}, mcpctx.NewValidationError("id is required", "id")
		}

		template, err := toolCtx.DB().GetTemplateByID(ctx, input.ID)
		if err != nil {
			return nil, SequenceEmailDeleteOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("email %s not found", input.ID))
		}

		if !template.SequenceID.Valid {
			return nil, SequenceEmailDeleteOutput{}, mcpctx.NewValidationError("email is not part of a sequence", "id")
		}

		seq, err := toolCtx.DB().GetSequenceByID(ctx, template.SequenceID.String)
		if err != nil || seq.OrgID.String != toolCtx.OrgID() {
			return nil, SequenceEmailDeleteOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("email %s not found", input.ID))
		}

		err = toolCtx.DB().DeleteTemplate(ctx, input.ID)
		if err != nil {
			return nil, SequenceEmailDeleteOutput{}, fmt.Errorf("failed to delete email: %w", err)
		}

		return nil, SequenceEmailDeleteOutput{
			Success: true,
			Message: fmt.Sprintf("Email %s deleted successfully", input.ID),
		}, nil
	}
}

// registerEmailToolsToRegistry registers email tools to the direct-call registry.
func registerEmailToolsToRegistry(registry *ToolRegistry, toolCtx *mcpctx.ToolContext) {
	registry.Register("list_create", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input ListCreateInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := createListHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("list_list", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input ListListInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := listListsHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("sequence_create", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input SequenceCreateInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := createSequenceHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("sequence_list", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input SequenceListInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := listSequencesHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("sequence_email_add", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input SequenceEmailAddInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := addSequenceEmailHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("sequence_get", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input SequenceGetInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := getSequenceHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("list_get", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input ListGetInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := getListHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("list_mutate", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input ListMutateInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := mutateListHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("list_delete", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input ListDeleteInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := deleteListHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("sequence_mutate", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input SequenceMutateInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := mutateSequenceHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("sequence_delete", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input SequenceDeleteInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := deleteSequenceHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("sequence_email_list", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input SequenceEmailListInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := listSequenceEmailsHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("sequence_email_get", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input SequenceEmailGetInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := getSequenceEmailHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("sequence_email_delete", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input SequenceEmailDeleteInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := deleteSequenceEmailHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})

	registry.Register("sequence_email_mutate", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input SequenceEmailMutateInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := mutateSequenceEmailHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})
}
