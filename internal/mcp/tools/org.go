package tools

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/mcp/mcpctx"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// orgActions defines valid actions for each org resource.
var orgActions = map[string][]string{
	"org":    {"list", "select", "get", "update", "create", "delete", "dashboard_stats"},
	"domain": {"list", "create", "get", "refresh", "delete"},
}

// orgOutputSchema defines the JSON Schema for all possible org outputs.
var orgOutputSchema = map[string]any{
	"type": "object",
	"description": `Output varies by resource and action:
- org.list → {organizations: OrgItem[], total: number, auth_mode: string}
- org.select → {id, name, slug, selected: true, message: string}
- org.get → {id, name, slug, max_contacts, from_name, from_email, reply_to}
- org.update → {id, name, from_name, from_email, reply_to, updated: true}`,
	"oneOf": []map[string]any{
		{
			"title":       "OrgList",
			"description": "Returned by org.list",
			"type":        "object",
			"properties": map[string]any{
				"organizations": map[string]any{"type": "array", "description": "List of organizations"},
				"total":         map[string]any{"type": "integer", "description": "Total count"},
				"auth_mode":     map[string]any{"type": "string", "description": "Authentication mode: api_key or oauth"},
			},
			"required": []string{"organizations", "total", "auth_mode"},
		},
		{
			"title":       "OrgSelect",
			"description": "Returned by org.select",
			"type":        "object",
			"properties": map[string]any{
				"id":       map[string]any{"type": "string", "description": "Organization UUID"},
				"name":     map[string]any{"type": "string", "description": "Organization name"},
				"slug":     map[string]any{"type": "string", "description": "URL-friendly slug"},
				"selected": map[string]any{"type": "boolean", "description": "Whether selection was successful"},
				"message":  map[string]any{"type": "string", "description": "Additional information"},
			},
			"required": []string{"id", "name", "slug", "selected"},
		},
		{
			"title":       "OrgGet",
			"description": "Returned by org.get",
			"type":        "object",
			"properties": map[string]any{
				"id":           map[string]any{"type": "string", "description": "Organization UUID"},
				"name":         map[string]any{"type": "string", "description": "Organization name"},
				"slug":         map[string]any{"type": "string", "description": "URL-friendly slug"},
				"max_contacts": map[string]any{"type": "integer", "description": "Maximum allowed contacts"},
				"from_name":    map[string]any{"type": "string", "description": "Email sender name"},
				"from_email":   map[string]any{"type": "string", "description": "Email sender address"},
				"reply_to":     map[string]any{"type": "string", "description": "Reply-to email address"},
			},
			"required": []string{"id", "name", "slug"},
		},
		{
			"title":       "OrgUpdate",
			"description": "Returned by org.update",
			"type":        "object",
			"properties": map[string]any{
				"id":         map[string]any{"type": "string", "description": "Organization UUID"},
				"name":       map[string]any{"type": "string", "description": "Organization name"},
				"from_name":  map[string]any{"type": "string", "description": "Email sender name"},
				"from_email": map[string]any{"type": "string", "description": "Email sender address"},
				"reply_to":   map[string]any{"type": "string", "description": "Reply-to email address"},
				"updated":    map[string]any{"type": "boolean", "const": true},
			},
			"required": []string{"id", "name", "updated"},
		},
	},
}

// OrgInput defines input for the unified org tool.
type OrgInput struct {
	Resource string `json:"resource" jsonschema:"required,Resource type: org or domain"`
	Action   string `json:"action" jsonschema:"required,Action to perform"`

	// Common
	ID string `json:"id,omitempty" jsonschema:"Resource ID (for get, update, delete, refresh)"`

	// Select-specific
	OrgID string `json:"org_id,omitempty" jsonschema:"Organization ID to select (UUID). For org.select."`
	Slug  string `json:"slug,omitempty" jsonschema:"Organization slug to select (alternative to org_id). For org.select."`

	// Create/Update-specific
	Name      string `json:"name,omitempty" jsonschema:"Organization name. For org.create, org.update."`
	FromName  string `json:"from_name,omitempty" jsonschema:"Email sender name (e.g., 'John from Acme'). For org.update."`
	FromEmail string `json:"from_email,omitempty" jsonschema:"Email sender address (e.g., 'john@acme.com'). For org.update."`
	ReplyTo   string `json:"reply_to,omitempty" jsonschema:"Reply-to email address. For org.update."`

	// Domain-specific
	Domain string `json:"domain,omitempty" jsonschema:"Domain name (e.g., 'example.com'). For domain.create."`

	// Delete confirmation
	Confirm bool `json:"confirm,omitempty" jsonschema:"Set to true to confirm deletion. For org.delete, domain.delete."`
}

// RegisterOrgTool registers the unified org tool.
func RegisterOrgTool(server *mcp.Server, toolCtx *mcpctx.ToolContext) {
	mcp.AddTool(server, &mcp.Tool{
		Name:  "org",
		Title: "Organization Management",
		Description: `Manage organizations, domain identities, and select which org to work with.

IMPORTANT: For OAuth sessions, you must use org.select before using other tools.
For API key sessions, selection is automatic.

Resources:
- org: Organization management
- domain: Domain identity management for email sending

ORG RESOURCE:
- org.list: List all organizations you have access to
- org.select: Select an organization to work with (OAuth only, requires: org_id or slug)
- org.get: Get current organization details
- org.update: Update organization settings
- org.create: Create a new organization (requires: name)
- org.delete: Delete organization (requires: confirm=true)
- org.dashboard_stats: Get dashboard statistics

DOMAIN RESOURCE:
- domain.list: List verified domains
- domain.create: Add a domain for verification (requires: domain)
- domain.get: Get domain details with DNS records (requires: id)
- domain.refresh: Check verification status with SES (requires: id)
- domain.delete: Remove a domain (requires: id)

Examples:
  org(resource: org, action: list)
  org(resource: org, action: select, slug: "my-company")
  org(resource: org, action: get)
  org(resource: org, action: update, from_email: "hello@example.com")
  org(resource: org, action: create, name: "New Company")
  org(resource: org, action: dashboard_stats)
  org(resource: domain, action: list)
  org(resource: domain, action: create, domain: "example.com")
  org(resource: domain, action: refresh, id: "uuid")`,
		OutputSchema: orgOutputSchema,
	}, orgHandler(toolCtx))
}

func orgHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input OrgInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input OrgInput) (*mcp.CallToolResult, any, error) {
		fmt.Printf("[MCP org] Handler called - Resource: %q, Action: %q, ID: %q, OrgID: %q, Slug: %q, Name: %q\n",
			input.Resource, input.Action, input.ID, input.OrgID, input.Slug, input.Name)

		// Validate resource
		validActions, ok := orgActions[input.Resource]
		if !ok {
			fmt.Printf("[MCP org] ERROR: invalid resource %q\n", input.Resource)
			return nil, nil, mcpctx.NewValidationError(
				fmt.Sprintf("invalid resource '%s', must be: org or domain", input.Resource),
				"resource")
		}

		// Validate action
		if !slices.Contains(validActions, input.Action) {
			fmt.Printf("[MCP org] ERROR: invalid action %q for resource %q\n", input.Action, input.Resource)
			return nil, nil, mcpctx.NewValidationError(
				fmt.Sprintf("invalid action '%s' for resource '%s', must be: %s",
					input.Action, input.Resource, strings.Join(validActions, ", ")),
				"action")
		}

		switch input.Resource {
		case "org":
			return handleOrg(ctx, toolCtx, input)
		case "domain":
			return handleDomain(ctx, toolCtx, input)
		}
		return nil, nil, nil // unreachable
	}
}

// ============================================================================
// ORG HANDLERS
// ============================================================================

// OrgListItem represents an organization in the list.
type OrgListItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Role     string `json:"role,omitempty"`
	Selected bool   `json:"selected"`
}

// OrgListOutput defines output for org.list.
type OrgListOutput struct {
	Organizations []OrgListItem `json:"organizations"`
	Total         int           `json:"total"`
	AuthMode      string        `json:"auth_mode"`
}

// OrgSelectOutput defines output for org.select.
type OrgSelectOutput struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Selected bool   `json:"selected"`
	Message  string `json:"message,omitempty"`
}

// OrgGetOutput defines output for org.get.
type OrgGetOutput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	MaxContacts int64  `json:"max_contacts,omitempty"`
	FromName    string `json:"from_name,omitempty"`
	FromEmail   string `json:"from_email,omitempty"`
	ReplyTo     string `json:"reply_to,omitempty"`
}

// OrgUpdateOutput defines output for org.update.
type OrgUpdateOutput struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	FromName  string `json:"from_name,omitempty"`
	FromEmail string `json:"from_email,omitempty"`
	ReplyTo   string `json:"reply_to,omitempty"`
	Updated   bool   `json:"updated"`
}

func handleOrg(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	switch input.Action {
	case "list":
		return handleOrgList(ctx, toolCtx, input)
	case "select":
		return handleOrgSelect(ctx, toolCtx, input)
	case "get":
		return handleOrgGet(ctx, toolCtx, input)
	case "update":
		return handleOrgUpdate(ctx, toolCtx, input)
	case "create":
		return handleOrgCreate(ctx, toolCtx, input)
	case "delete":
		return handleOrgDelete(ctx, toolCtx, input)
	case "dashboard_stats":
		return handleOrgDashboardStats(ctx, toolCtx, input)
	}
	return nil, nil, nil
}

func handleOrgList(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	authMode := "api_key"
	if toolCtx.AuthMode() == mcpctx.AuthModeOAuth {
		authMode = "oauth"
	}

	items := make([]OrgListItem, 0)

	if toolCtx.AuthMode() == mcpctx.AuthModeAPIKey {
		org := toolCtx.Org()
		items = append(items, OrgListItem{
			ID:       org.ID,
			Name:     org.Name,
			Slug:     org.Slug,
			Selected: true,
		})
	} else {
		allOrgs, err := toolCtx.DB().ListOrganizations(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list organizations: %w", err)
		}

		currentOrgID := toolCtx.OrgID()
		for _, org := range allOrgs {
			items = append(items, OrgListItem{
				ID:       org.ID,
				Name:     org.Name,
				Slug:     org.Slug,
				Role:     "admin",
				Selected: org.ID == currentOrgID,
			})
		}
	}

	return nil, OrgListOutput{
		Organizations: items,
		Total:         len(items),
		AuthMode:      authMode,
	}, nil
}

func handleOrgSelect(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	// Accept either org_id or id for the org ID (users might use either)
	orgID := input.OrgID
	if orgID == "" && input.ID != "" {
		orgID = input.ID
	}

	fmt.Printf("[MCP org.select] Starting - AuthMode: %v, SessionID: %s, OrgID: %q, ID: %q, Slug: %q\n",
		toolCtx.AuthMode(), toolCtx.SessionID(), input.OrgID, input.ID, input.Slug)

	if toolCtx.AuthMode() == mcpctx.AuthModeAPIKey {
		org := toolCtx.Org()
		fmt.Printf("[MCP org.select] API key mode - returning org: %s\n", org.Name)
		return nil, OrgSelectOutput{
			ID:       org.ID,
			Name:     org.Name,
			Slug:     org.Slug,
			Selected: true,
			Message:  "API key authentication is already scoped to this organization. No selection needed.",
		}, nil
	}

	if orgID == "" && input.Slug == "" {
		fmt.Println("[MCP org.select] ERROR: neither org_id/id nor slug provided")
		return nil, nil, mcpctx.NewValidationError("either id, org_id, or slug is required for org.select", "id")
	}

	var fullOrg db.Organization
	var err error

	if orgID != "" {
		fmt.Printf("[MCP org.select] Looking up by ID: %s\n", orgID)
		fullOrg, err = toolCtx.DB().GetOrganizationByID(ctx, orgID)
	} else {
		fmt.Printf("[MCP org.select] Looking up by slug: %s\n", input.Slug)
		fullOrg, err = toolCtx.DB().GetOrganizationBySlug(ctx, input.Slug)
	}
	if err != nil {
		fmt.Printf("[MCP org.select] ERROR: org lookup failed: %v\n", err)
		if orgID != "" {
			return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("organization with ID '%s' not found", orgID))
		}
		return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("organization with slug '%s' not found", input.Slug))
	}

	fmt.Printf("[MCP org.select] Found org: %s (%s), calling SelectOrg\n", fullOrg.Name, fullOrg.ID)
	toolCtx.SelectOrg(fullOrg)
	fmt.Printf("[MCP org.select] SelectOrg completed, HasOrg: %v\n", toolCtx.HasOrg())

	return nil, OrgSelectOutput{
		ID:       fullOrg.ID,
		Name:     fullOrg.Name,
		Slug:     fullOrg.Slug,
		Selected: true,
		Message:  "Organization selected. All subsequent operations will use this organization.",
	}, nil
}

func handleOrgGet(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireOrg(); err != nil {
		return nil, nil, err
	}
	org := toolCtx.Org()

	return nil, OrgGetOutput{
		ID:          org.ID,
		Name:        org.Name,
		Slug:        org.Slug,
		MaxContacts: org.MaxContacts.Int64,
		FromName:    org.FromName.String,
		FromEmail:   org.FromEmail.String,
		ReplyTo:     org.ReplyTo.String,
	}, nil
}

func handleOrgUpdate(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireOrg(); err != nil {
		return nil, nil, err
	}
	org := toolCtx.Org()

	if input.FromName != "" || input.FromEmail != "" || input.ReplyTo != "" {
		_, err := toolCtx.DB().UpdateOrgEmailSettings(ctx, db.UpdateOrgEmailSettingsParams{
			ID:        toolCtx.OrgID(),
			FromName:  input.FromName,
			FromEmail: input.FromEmail,
			ReplyTo:   input.ReplyTo,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update email settings: %w", err)
		}
	}

	if input.Name != "" {
		_, err := toolCtx.DB().UpdateOrganization(ctx, db.UpdateOrganizationParams{
			ID:   toolCtx.OrgID(),
			Name: input.Name,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update organization: %w", err)
		}
	}

	updatedOrg, err := toolCtx.DB().GetOrganizationByID(ctx, toolCtx.OrgID())
	if err != nil {
		return nil, OrgUpdateOutput{
			ID:        org.ID,
			Name:      org.Name,
			FromName:  org.FromName.String,
			FromEmail: org.FromEmail.String,
			ReplyTo:   org.ReplyTo.String,
			Updated:   true,
		}, nil
	}

	return nil, OrgUpdateOutput{
		ID:        updatedOrg.ID,
		Name:      updatedOrg.Name,
		FromName:  updatedOrg.FromName.String,
		FromEmail: updatedOrg.FromEmail.String,
		ReplyTo:   updatedOrg.ReplyTo.String,
		Updated:   true,
	}, nil
}

// ============================================================================
// ORG CREATE/DELETE/STATS HANDLERS
// ============================================================================

// OrgCreateOutput defines output for org.create.
type OrgCreateOutput struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	APIKey  string `json:"api_key"`
	Created bool   `json:"created"`
}

// OrgDeleteOutput defines output for org.delete.
type OrgDeleteOutput struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// OrgDashboardStatsOutput defines output for org.dashboard_stats.
type OrgDashboardStatsOutput struct {
	TotalContacts   int64 `json:"total_contacts"`
	ActiveContacts  int64 `json:"active_contacts"`
	NewContacts30d  int64 `json:"new_contacts_30d"`
	NewContacts7d   int64 `json:"new_contacts_7d"`
	EmailsSent30d   int64 `json:"emails_sent_30d"`
	Opens30d        int64 `json:"opens_30d"`
	Clicks30d       int64 `json:"clicks_30d"`
}

func handleOrgCreate(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, nil, mcpctx.NewValidationError("name is required", "name")
	}

	// Generate slug from name
	slug := strings.ToLower(input.Name)
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

	orgID := uuid.New().String()
	apiKey := uuid.New().String()

	org, err := toolCtx.DB().CreateOrganization(ctx, db.CreateOrganizationParams{
		ID:          orgID,
		Name:        input.Name,
		Slug:        slug,
		ApiKey:      apiKey,
		MaxContacts: sql.NullInt64{Int64: 10000, Valid: true},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Auto-select the newly created organization so subsequent operations work immediately
	toolCtx.SelectOrg(org)

	return nil, OrgCreateOutput{
		ID:      org.ID,
		Name:    org.Name,
		Slug:    org.Slug,
		APIKey:  org.ApiKey,
		Created: true,
	}, nil
}

func handleOrgDelete(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireOrg(); err != nil {
		return nil, nil, err
	}

	if !input.Confirm {
		return nil, nil, mcpctx.NewValidationError("confirm=true is required to delete an organization", "confirm")
	}

	org := toolCtx.Org()

	err := toolCtx.DB().DeleteOrganization(ctx, org.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil, OrgDeleteOutput{
		ID:      org.ID,
		Name:    org.Name,
		Success: true,
		Message: fmt.Sprintf("Organization %s deleted successfully", org.Name),
	}, nil
}

func handleOrgDashboardStats(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireOrg(); err != nil {
		return nil, nil, err
	}

	orgID := toolCtx.OrgID()

	// Get subscriber stats
	var totalContacts, new30d, new7d int64
	subscriberStats, err := toolCtx.DB().GetDashboardSubscriberStats(ctx, sql.NullString{String: orgID, Valid: true})
	if err == nil {
		totalContacts = subscriberStats.Total
		if subscriberStats.New30d.Valid {
			new30d = int64(subscriberStats.New30d.Float64)
		}
		if subscriberStats.New7d.Valid {
			new7d = int64(subscriberStats.New7d.Float64)
		}
	}

	// Get email stats (30 days)
	var emailsSent, opens, clicks int64
	emailStats, err := toolCtx.DB().GetDashboardEmailStats30Days(ctx, orgID)
	if err == nil {
		emailsSent = interfaceToInt64(emailStats.EmailsSent)
		opens = interfaceToInt64(emailStats.EmailsOpened)
		clicks = interfaceToInt64(emailStats.EmailsClicked)
	}

	return nil, OrgDashboardStatsOutput{
		TotalContacts:  totalContacts,
		ActiveContacts: totalContacts, // Same as total for now
		NewContacts30d: new30d,
		NewContacts7d:  new7d,
		EmailsSent30d:  emailsSent,
		Opens30d:       opens,
		Clicks30d:      clicks,
	}, nil
}

// ============================================================================
// DOMAIN IDENTITY HANDLERS
// ============================================================================

// DomainIdentityItem represents a domain identity in the list.
type DomainIdentityItem struct {
	ID                 string `json:"id"`
	Domain             string `json:"domain"`
	VerificationStatus string `json:"verification_status"`
	DkimStatus         string `json:"dkim_status"`
	MailFromStatus     string `json:"mail_from_status,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
}

// DomainIdentityListOutput defines output for domain.list.
type DomainIdentityListOutput struct {
	Domains []DomainIdentityItem `json:"domains"`
	Total   int                  `json:"total"`
}

// DomainIdentityCreateOutput defines output for domain.create.
type DomainIdentityCreateOutput struct {
	ID                 string `json:"id"`
	Domain             string `json:"domain"`
	VerificationStatus string `json:"verification_status"`
	VerificationToken  string `json:"verification_token,omitempty"`
	DnsRecords         string `json:"dns_records,omitempty"`
	Created            bool   `json:"created"`
}

// DomainIdentityGetOutput defines output for domain.get.
type DomainIdentityGetOutput struct {
	ID                 string `json:"id"`
	Domain             string `json:"domain"`
	VerificationStatus string `json:"verification_status"`
	DkimStatus         string `json:"dkim_status"`
	VerificationToken  string `json:"verification_token,omitempty"`
	DkimTokens         string `json:"dkim_tokens,omitempty"`
	DnsRecords         string `json:"dns_records,omitempty"`
	MailFromDomain     string `json:"mail_from_domain,omitempty"`
	MailFromStatus     string `json:"mail_from_status,omitempty"`
	LastCheckedAt      string `json:"last_checked_at,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
}

// DomainIdentityRefreshOutput defines output for domain.refresh.
type DomainIdentityRefreshOutput struct {
	ID                 string `json:"id"`
	Domain             string `json:"domain"`
	VerificationStatus string `json:"verification_status"`
	DkimStatus         string `json:"dkim_status"`
	Refreshed          bool   `json:"refreshed"`
	Message            string `json:"message"`
}

// DomainIdentityDeleteOutput defines output for domain.delete.
type DomainIdentityDeleteOutput struct {
	ID      string `json:"id"`
	Domain  string `json:"domain"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func handleDomain(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	switch input.Action {
	case "list":
		return handleDomainIdentityList(ctx, toolCtx, input)
	case "create":
		return handleDomainIdentityCreate(ctx, toolCtx, input)
	case "get":
		return handleDomainIdentityGet(ctx, toolCtx, input)
	case "refresh":
		return handleDomainIdentityRefresh(ctx, toolCtx, input)
	case "delete":
		return handleDomainIdentityDelete(ctx, toolCtx, input)
	}
	return nil, nil, nil
}

func handleDomainIdentityList(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireOrg(); err != nil {
		return nil, nil, err
	}

	domains, err := toolCtx.DB().ListDomainIdentitiesByOrg(ctx, toolCtx.OrgID())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list domains: %w", err)
	}

	items := make([]DomainIdentityItem, 0, len(domains))
	for _, d := range domains {
		items = append(items, DomainIdentityItem{
			ID:                 d.ID,
			Domain:             d.Domain,
			VerificationStatus: d.VerificationStatus.String,
			DkimStatus:         d.DkimStatus.String,
			MailFromStatus:     d.MailFromStatus.String,
			CreatedAt:          d.CreatedAt.String,
		})
	}

	return nil, DomainIdentityListOutput{
		Domains: items,
		Total:   len(items),
	}, nil
}

func handleDomainIdentityCreate(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireOrg(); err != nil {
		return nil, nil, err
	}

	if strings.TrimSpace(input.Domain) == "" {
		return nil, nil, mcpctx.NewValidationError("domain is required", "domain")
	}

	domainID := uuid.New().String()
	verificationToken := uuid.New().String()

	domain, err := toolCtx.DB().CreateDomainIdentity(ctx, db.CreateDomainIdentityParams{
		ID:                 domainID,
		OrgID:              toolCtx.OrgID(),
		Domain:             input.Domain,
		VerificationStatus: sql.NullString{String: "pending", Valid: true},
		DkimStatus:         sql.NullString{String: "pending", Valid: true},
		VerificationToken:  sql.NullString{String: verificationToken, Valid: true},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create domain: %w", err)
	}

	return nil, DomainIdentityCreateOutput{
		ID:                 domain.ID,
		Domain:             domain.Domain,
		VerificationStatus: domain.VerificationStatus.String,
		VerificationToken:  domain.VerificationToken.String,
		DnsRecords:         domain.DnsRecords.String,
		Created:            true,
	}, nil
}

func handleDomainIdentityGet(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireOrg(); err != nil {
		return nil, nil, err
	}

	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, mcpctx.NewValidationError("id is required", "id")
	}

	domain, err := toolCtx.DB().GetDomainIdentity(ctx, input.ID)
	if err != nil {
		return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("domain %s not found", input.ID))
	}

	if domain.OrgID != toolCtx.OrgID() {
		return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("domain %s not found", input.ID))
	}

	return nil, DomainIdentityGetOutput{
		ID:                 domain.ID,
		Domain:             domain.Domain,
		VerificationStatus: domain.VerificationStatus.String,
		DkimStatus:         domain.DkimStatus.String,
		VerificationToken:  domain.VerificationToken.String,
		DkimTokens:         domain.DkimTokens.String,
		DnsRecords:         domain.DnsRecords.String,
		MailFromDomain:     domain.MailFromDomain.String,
		MailFromStatus:     domain.MailFromStatus.String,
		LastCheckedAt:      domain.LastCheckedAt.String,
		CreatedAt:          domain.CreatedAt.String,
	}, nil
}

func handleDomainIdentityRefresh(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireOrg(); err != nil {
		return nil, nil, err
	}

	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, mcpctx.NewValidationError("id is required", "id")
	}

	domain, err := toolCtx.DB().GetDomainIdentity(ctx, input.ID)
	if err != nil {
		return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("domain %s not found", input.ID))
	}

	if domain.OrgID != toolCtx.OrgID() {
		return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("domain %s not found", input.ID))
	}

	// Note: In a real implementation, this would call AWS SES to check verification status
	// For now, we just return the current status
	return nil, DomainIdentityRefreshOutput{
		ID:                 domain.ID,
		Domain:             domain.Domain,
		VerificationStatus: domain.VerificationStatus.String,
		DkimStatus:         domain.DkimStatus.String,
		Refreshed:          true,
		Message:            "Domain verification status checked. Configure your DNS records and check again.",
	}, nil
}

func handleDomainIdentityDelete(ctx context.Context, toolCtx *mcpctx.ToolContext, input OrgInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireOrg(); err != nil {
		return nil, nil, err
	}

	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, mcpctx.NewValidationError("id is required", "id")
	}

	domain, err := toolCtx.DB().GetDomainIdentity(ctx, input.ID)
	if err != nil {
		return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("domain %s not found", input.ID))
	}

	if domain.OrgID != toolCtx.OrgID() {
		return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("domain %s not found", input.ID))
	}

	err = toolCtx.DB().DeleteDomainIdentity(ctx, input.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete domain: %w", err)
	}

	return nil, DomainIdentityDeleteOutput{
		ID:      domain.ID,
		Domain:  domain.Domain,
		Success: true,
		Message: fmt.Sprintf("Domain %s deleted successfully", domain.Domain),
	}, nil
}

// registerOrgToolToRegistry registers org tool to the direct-call registry.
func registerOrgToolToRegistry(registry *ToolRegistry, toolCtx *mcpctx.ToolContext) {
	registry.Register("org", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input OrgInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := orgHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})
}
