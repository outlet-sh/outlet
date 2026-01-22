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

// brandActions defines valid actions for each brand resource.
var brandActions = map[string][]string{
	"brand":  {"list", "select", "get", "update", "create", "delete", "dashboard_stats"},
	"domain": {"list", "create", "get", "refresh", "delete"},
}

// brandOutputSchema defines the JSON Schema for all possible brand outputs.
var brandOutputSchema = map[string]any{
	"type": "object",
	"description": `Output varies by resource and action:
- brand.list → {brands: BrandItem[], total: number, auth_mode: string}
- brand.select → {id, name, slug, selected: true, message: string}
- brand.get → {id, name, slug, max_contacts, from_name, from_email, reply_to}
- brand.update → {id, name, from_name, from_email, reply_to, updated: true}`,
	"oneOf": []map[string]any{
		{
			"title":       "BrandList",
			"description": "Returned by brand.list",
			"type":        "object",
			"properties": map[string]any{
				"brands":    map[string]any{"type": "array", "description": "List of brands"},
				"total":     map[string]any{"type": "integer", "description": "Total count"},
				"auth_mode": map[string]any{"type": "string", "description": "Authentication mode: api_key or oauth"},
			},
			"required": []string{"brands", "total", "auth_mode"},
		},
		{
			"title":       "BrandSelect",
			"description": "Returned by brand.select",
			"type":        "object",
			"properties": map[string]any{
				"id":       map[string]any{"type": "string", "description": "Brand UUID"},
				"name":     map[string]any{"type": "string", "description": "Brand name"},
				"slug":     map[string]any{"type": "string", "description": "URL-friendly slug"},
				"selected": map[string]any{"type": "boolean", "description": "Whether selection was successful"},
				"message":  map[string]any{"type": "string", "description": "Additional information"},
			},
			"required": []string{"id", "name", "slug", "selected"},
		},
		{
			"title":       "BrandGet",
			"description": "Returned by brand.get",
			"type":        "object",
			"properties": map[string]any{
				"id":           map[string]any{"type": "string", "description": "Brand UUID"},
				"name":         map[string]any{"type": "string", "description": "Brand name"},
				"slug":         map[string]any{"type": "string", "description": "URL-friendly slug"},
				"max_contacts": map[string]any{"type": "integer", "description": "Maximum allowed contacts"},
				"from_name":    map[string]any{"type": "string", "description": "Email sender name"},
				"from_email":   map[string]any{"type": "string", "description": "Email sender address"},
				"reply_to":     map[string]any{"type": "string", "description": "Reply-to email address"},
			},
			"required": []string{"id", "name", "slug"},
		},
		{
			"title":       "BrandUpdate",
			"description": "Returned by brand.update",
			"type":        "object",
			"properties": map[string]any{
				"id":         map[string]any{"type": "string", "description": "Brand UUID"},
				"name":       map[string]any{"type": "string", "description": "Brand name"},
				"from_name":  map[string]any{"type": "string", "description": "Email sender name"},
				"from_email": map[string]any{"type": "string", "description": "Email sender address"},
				"reply_to":   map[string]any{"type": "string", "description": "Reply-to email address"},
				"updated":    map[string]any{"type": "boolean", "const": true},
			},
			"required": []string{"id", "name", "updated"},
		},
	},
}

// BrandInput defines input for the unified brand tool.
type BrandInput struct {
	Resource string `json:"resource" jsonschema:"required,Resource type: brand or domain"`
	Action   string `json:"action" jsonschema:"required,Action to perform"`

	// Common
	ID string `json:"id,omitempty" jsonschema:"Resource ID (for get, update, delete, refresh)"`

	// Select-specific
	BrandID string `json:"brand_id,omitempty" jsonschema:"Brand ID to select (UUID). For brand.select."`
	Slug    string `json:"slug,omitempty" jsonschema:"Brand slug to select (alternative to brand_id). For brand.select."`

	// Create/Update-specific
	Name      string `json:"name,omitempty" jsonschema:"Brand name. For brand.create, brand.update."`
	FromName  string `json:"from_name,omitempty" jsonschema:"Email sender name (e.g., 'John from Acme'). For brand.update."`
	FromEmail string `json:"from_email,omitempty" jsonschema:"Email sender address (e.g., 'john@acme.com'). For brand.update."`
	ReplyTo   string `json:"reply_to,omitempty" jsonschema:"Reply-to email address. For brand.update."`

	// Domain-specific
	Domain string `json:"domain,omitempty" jsonschema:"Domain name (e.g., 'example.com'). For domain.create."`

	// Delete confirmation
	Confirm bool `json:"confirm,omitempty" jsonschema:"Set to true to confirm deletion. For brand.delete, domain.delete."`
}

// RegisterBrandTool registers the unified brand tool.
func RegisterBrandTool(server *mcp.Server, toolCtx *mcpctx.ToolContext) {
	mcp.AddTool(server, &mcp.Tool{
		Name:  "brand",
		Title: "Brand Management",
		Description: `Manage brands, domain identities, and select which brand to work with.

IMPORTANT: For OAuth sessions, you must use brand.select before using other tools.
For API key sessions, selection is automatic.

Resources:
- brand: Brand management
- domain: Domain identity management for email sending

BRAND RESOURCE:
- brand.list: List all brands you have access to
- brand.select: Select a brand to work with (OAuth only, requires: brand_id or slug)
- brand.get: Get current brand details
- brand.update: Update brand settings
- brand.create: Create a new brand (requires: name)
- brand.delete: Delete brand (requires: confirm=true)
- brand.dashboard_stats: Get dashboard statistics

DOMAIN RESOURCE:
- domain.list: List verified domains
- domain.create: Add a domain for verification (requires: domain)
- domain.get: Get domain details with DNS records (requires: id)
- domain.refresh: Check verification status with SES (requires: id)
- domain.delete: Remove a domain (requires: id)

Examples:
  brand(resource: brand, action: list)
  brand(resource: brand, action: select, slug: "my-company")
  brand(resource: brand, action: get)
  brand(resource: brand, action: update, from_email: "hello@example.com")
  brand(resource: brand, action: create, name: "New Company")
  brand(resource: brand, action: dashboard_stats)
  brand(resource: domain, action: list)
  brand(resource: domain, action: create, domain: "example.com")
  brand(resource: domain, action: refresh, id: "uuid")`,
		OutputSchema: brandOutputSchema,
	}, brandHandler(toolCtx))
}

func brandHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcp.CallToolRequest, input BrandInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input BrandInput) (*mcp.CallToolResult, any, error) {
		fmt.Printf("[MCP brand] Handler called - Resource: %q, Action: %q, ID: %q, BrandID: %q, Slug: %q, Name: %q\n",
			input.Resource, input.Action, input.ID, input.BrandID, input.Slug, input.Name)

		// Validate resource
		validActions, ok := brandActions[input.Resource]
		if !ok {
			fmt.Printf("[MCP brand] ERROR: invalid resource %q\n", input.Resource)
			return nil, nil, mcpctx.NewValidationError(
				fmt.Sprintf("invalid resource '%s', must be: brand or domain", input.Resource),
				"resource")
		}

		// Validate action
		if !slices.Contains(validActions, input.Action) {
			fmt.Printf("[MCP brand] ERROR: invalid action %q for resource %q\n", input.Action, input.Resource)
			return nil, nil, mcpctx.NewValidationError(
				fmt.Sprintf("invalid action '%s' for resource '%s', must be: %s",
					input.Action, input.Resource, strings.Join(validActions, ", ")),
				"action")
		}

		switch input.Resource {
		case "brand":
			return handleBrand(ctx, toolCtx, input)
		case "domain":
			return handleDomain(ctx, toolCtx, input)
		}
		return nil, nil, nil // unreachable
	}
}

// ============================================================================
// ORG HANDLERS
// ============================================================================

// BrandListItem represents an organization in the list.
type BrandListItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Role     string `json:"role,omitempty"`
	Selected bool   `json:"selected"`
}

// BrandListOutput defines output for brand.list.
type BrandListOutput struct {
	Brands   []BrandListItem `json:"brands"`
	Total    int             `json:"total"`
	AuthMode string          `json:"auth_mode"`
}

// BrandSelectOutput defines output for brand.select.
type BrandSelectOutput struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Selected bool   `json:"selected"`
	Message  string `json:"message,omitempty"`
}

// BrandGetOutput defines output for brand.get.
type BrandGetOutput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	MaxContacts int64  `json:"max_contacts,omitempty"`
	FromName    string `json:"from_name,omitempty"`
	FromEmail   string `json:"from_email,omitempty"`
	ReplyTo     string `json:"reply_to,omitempty"`
}

// BrandUpdateOutput defines output for brand.update.
type BrandUpdateOutput struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	FromName  string `json:"from_name,omitempty"`
	FromEmail string `json:"from_email,omitempty"`
	ReplyTo   string `json:"reply_to,omitempty"`
	Updated   bool   `json:"updated"`
}

func handleBrand(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	switch input.Action {
	case "list":
		return handleBrandList(ctx, toolCtx, input)
	case "select":
		return handleBrandSelect(ctx, toolCtx, input)
	case "get":
		return handleBrandGet(ctx, toolCtx, input)
	case "update":
		return handleBrandUpdate(ctx, toolCtx, input)
	case "create":
		return handleBrandCreate(ctx, toolCtx, input)
	case "delete":
		return handleBrandDelete(ctx, toolCtx, input)
	case "dashboard_stats":
		return handleBrandDashboardStats(ctx, toolCtx, input)
	}
	return nil, nil, nil
}

func handleBrandList(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	authMode := "api_key"
	if toolCtx.AuthMode() == mcpctx.AuthModeOAuth {
		authMode = "oauth"
	}

	items := make([]BrandListItem, 0)

	if toolCtx.AuthMode() == mcpctx.AuthModeAPIKey {
		org := toolCtx.Brand()
		items = append(items, BrandListItem{
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

		currentOrgID := toolCtx.BrandID()
		for _, org := range allOrgs {
			items = append(items, BrandListItem{
				ID:       org.ID,
				Name:     org.Name,
				Slug:     org.Slug,
				Role:     "admin",
				Selected: org.ID == currentOrgID,
			})
		}
	}

	return nil, BrandListOutput{
		Brands:   items,
		Total:    len(items),
		AuthMode: authMode,
	}, nil
}

func handleBrandSelect(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	// Accept either brand_id or id for the brand ID (users might use either)
	brandID := input.BrandID
	if brandID == "" && input.ID != "" {
		brandID = input.ID
	}

	fmt.Printf("[MCP brand.select] Starting - AuthMode: %v, SessionID: %s, BrandID: %q, ID: %q, Slug: %q\n",
		toolCtx.AuthMode(), toolCtx.SessionID(), input.BrandID, input.ID, input.Slug)

	if toolCtx.AuthMode() == mcpctx.AuthModeAPIKey {
		org := toolCtx.Brand()
		fmt.Printf("[MCP brand.select] API key mode - returning brand: %s\n", org.Name)
		return nil, BrandSelectOutput{
			ID:       org.ID,
			Name:     org.Name,
			Slug:     org.Slug,
			Selected: true,
			Message:  "API key authentication is already scoped to this brand. No selection needed.",
		}, nil
	}

	if brandID == "" && input.Slug == "" {
		fmt.Println("[MCP brand.select] ERROR: neither brand_id/id nor slug provided")
		return nil, nil, mcpctx.NewValidationError("either id, brand_id, or slug is required for brand.select", "id")
	}

	var fullOrg db.Organization
	var err error

	if brandID != "" {
		fmt.Printf("[MCP brand.select] Looking up by ID: %s\n", brandID)
		fullOrg, err = toolCtx.DB().GetOrganizationByID(ctx, brandID)
	} else {
		fmt.Printf("[MCP brand.select] Looking up by slug: %s\n", input.Slug)
		fullOrg, err = toolCtx.DB().GetOrganizationBySlug(ctx, input.Slug)
	}
	if err != nil {
		fmt.Printf("[MCP brand.select] ERROR: brand lookup failed: %v\n", err)
		if brandID != "" {
			return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("brand with ID '%s' not found", brandID))
		}
		return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("brand with slug '%s' not found", input.Slug))
	}

	fmt.Printf("[MCP brand.select] Found brand: %s (%s), calling SelectBrand\n", fullOrg.Name, fullOrg.ID)
	toolCtx.SelectBrand(fullOrg)
	fmt.Printf("[MCP brand.select] SelectBrand completed, HasBrand: %v\n", toolCtx.HasBrand())

	return nil, BrandSelectOutput{
		ID:       fullOrg.ID,
		Name:     fullOrg.Name,
		Slug:     fullOrg.Slug,
		Selected: true,
		Message:  "Brand selected. All subsequent operations will use this brand.",
	}, nil
}

func handleBrandGet(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireBrand(); err != nil {
		return nil, nil, err
	}
	org := toolCtx.Brand()

	return nil, BrandGetOutput{
		ID:          org.ID,
		Name:        org.Name,
		Slug:        org.Slug,
		MaxContacts: org.MaxContacts.Int64,
		FromName:    org.FromName.String,
		FromEmail:   org.FromEmail.String,
		ReplyTo:     org.ReplyTo.String,
	}, nil
}

func handleBrandUpdate(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireBrand(); err != nil {
		return nil, nil, err
	}
	org := toolCtx.Brand()

	if input.FromName != "" || input.FromEmail != "" || input.ReplyTo != "" {
		_, err := toolCtx.DB().UpdateOrgEmailSettings(ctx, db.UpdateOrgEmailSettingsParams{
			ID:        toolCtx.BrandID(),
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
			ID:   toolCtx.BrandID(),
			Name: input.Name,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update organization: %w", err)
		}
	}

	updatedOrg, err := toolCtx.DB().GetOrganizationByID(ctx, toolCtx.BrandID())
	if err != nil {
		return nil, BrandUpdateOutput{
			ID:        org.ID,
			Name:      org.Name,
			FromName:  org.FromName.String,
			FromEmail: org.FromEmail.String,
			ReplyTo:   org.ReplyTo.String,
			Updated:   true,
		}, nil
	}

	return nil, BrandUpdateOutput{
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

// BrandCreateOutput defines output for brand.create.
type BrandCreateOutput struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	APIKey  string `json:"api_key"`
	Created bool   `json:"created"`
}

// BrandDeleteOutput defines output for brand.delete.
type BrandDeleteOutput struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// BrandDashboardStatsOutput defines output for brand.dashboard_stats.
type BrandDashboardStatsOutput struct {
	TotalContacts   int64 `json:"total_contacts"`
	ActiveContacts  int64 `json:"active_contacts"`
	NewContacts30d  int64 `json:"new_contacts_30d"`
	NewContacts7d   int64 `json:"new_contacts_7d"`
	EmailsSent30d   int64 `json:"emails_sent_30d"`
	Opens30d        int64 `json:"opens_30d"`
	Clicks30d       int64 `json:"clicks_30d"`
}

func handleBrandCreate(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
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
	toolCtx.SelectBrand(org)

	return nil, BrandCreateOutput{
		ID:      org.ID,
		Name:    org.Name,
		Slug:    org.Slug,
		APIKey:  org.ApiKey,
		Created: true,
	}, nil
}

func handleBrandDelete(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireBrand(); err != nil {
		return nil, nil, err
	}

	if !input.Confirm {
		return nil, nil, mcpctx.NewValidationError("confirm=true is required to delete an organization", "confirm")
	}

	org := toolCtx.Brand()

	err := toolCtx.DB().DeleteOrganization(ctx, org.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil, BrandDeleteOutput{
		ID:      org.ID,
		Name:    org.Name,
		Success: true,
		Message: fmt.Sprintf("Organization %s deleted successfully", org.Name),
	}, nil
}

func handleBrandDashboardStats(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireBrand(); err != nil {
		return nil, nil, err
	}

	orgID := toolCtx.BrandID()

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

	return nil, BrandDashboardStatsOutput{
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

func handleDomain(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
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

func handleDomainIdentityList(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireBrand(); err != nil {
		return nil, nil, err
	}

	domains, err := toolCtx.DB().ListDomainIdentitiesByOrg(ctx, toolCtx.BrandID())
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

func handleDomainIdentityCreate(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireBrand(); err != nil {
		return nil, nil, err
	}

	if strings.TrimSpace(input.Domain) == "" {
		return nil, nil, mcpctx.NewValidationError("domain is required", "domain")
	}

	domainID := uuid.New().String()
	verificationToken := uuid.New().String()

	domain, err := toolCtx.DB().CreateDomainIdentity(ctx, db.CreateDomainIdentityParams{
		ID:                 domainID,
		OrgID:              toolCtx.BrandID(),
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

func handleDomainIdentityGet(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireBrand(); err != nil {
		return nil, nil, err
	}

	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, mcpctx.NewValidationError("id is required", "id")
	}

	domain, err := toolCtx.DB().GetDomainIdentity(ctx, input.ID)
	if err != nil {
		return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("domain %s not found", input.ID))
	}

	if domain.OrgID != toolCtx.BrandID() {
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

func handleDomainIdentityRefresh(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireBrand(); err != nil {
		return nil, nil, err
	}

	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, mcpctx.NewValidationError("id is required", "id")
	}

	domain, err := toolCtx.DB().GetDomainIdentity(ctx, input.ID)
	if err != nil {
		return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("domain %s not found", input.ID))
	}

	if domain.OrgID != toolCtx.BrandID() {
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

func handleDomainIdentityDelete(ctx context.Context, toolCtx *mcpctx.ToolContext, input BrandInput) (*mcp.CallToolResult, any, error) {
	if err := toolCtx.RequireBrand(); err != nil {
		return nil, nil, err
	}

	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, mcpctx.NewValidationError("id is required", "id")
	}

	domain, err := toolCtx.DB().GetDomainIdentity(ctx, input.ID)
	if err != nil {
		return nil, nil, mcpctx.NewNotFoundError(fmt.Sprintf("domain %s not found", input.ID))
	}

	if domain.OrgID != toolCtx.BrandID() {
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

// registerBrandToolToRegistry registers brand tool to the direct-call registry.
func registerBrandToolToRegistry(registry *ToolRegistry, toolCtx *mcpctx.ToolContext) {
	registry.Register("brand", func(ctx context.Context, args json.RawMessage) (interface{}, error) {
		var input BrandInput
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		handler := brandHandler(toolCtx)
		_, output, err := handler(ctx, nil, input)
		return output, err
	})
}
