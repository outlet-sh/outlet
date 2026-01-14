package tools

import (
	"context"
	"fmt"

	"outlet/internal/db"
	"outlet/internal/mcp/mcpctx"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterOrgTools registers all organization-related MCP tools.
func RegisterOrgTools(server *mcpsdk.Server, toolCtx *mcpctx.ToolContext) {
	// org_list - always available, lists orgs the user has access to
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "org_list",
		Description: "List all organizations you have access to. For OAuth sessions, use org_select to choose which organization to work with. For API key sessions, shows only the authenticated organization.",
	}, listOrgsHandler(toolCtx))

	// org_select - for OAuth sessions to select which org to work with
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "org_select",
		Description: "Select an organization to work with (OAuth sessions only). After selecting, all other tools will operate on this organization. Use org_list to see available organizations.",
	}, selectOrgHandler(toolCtx))

	// org_get
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "org_get",
		Description: "Get the current organization's settings and configuration. Returns org name, billing status, plan, email settings, and more. Requires an organization to be selected.",
	}, getOrgHandler(toolCtx))

	// org_update
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "org_update",
		Description: "Update organization settings including name, email sender info (from_name, from_email, reply_to). Does not update billing/payment settings - use org_payment_setup for that. Requires an organization to be selected.",
	}, updateOrgHandler(toolCtx))

}

// OrgListInput defines input for org_list tool.
type OrgListInput struct{}

// OrgListItem represents an organization in the list.
type OrgListItem struct {
	ID       string `json:"id" jsonschema:"Organization UUID"`
	Name     string `json:"name" jsonschema:"Organization name"`
	Slug     string `json:"slug" jsonschema:"URL-friendly slug"`
	Role     string `json:"role,omitempty" jsonschema:"Your role in this organization"`
	Selected bool   `json:"selected" jsonschema:"Whether this org is currently selected"`
}

// OrgListOutput defines output for org_list tool.
type OrgListOutput struct {
	Organizations []OrgListItem `json:"organizations" jsonschema:"List of organizations you have access to"`
	Total         int           `json:"total" jsonschema:"Total number of organizations"`
	AuthMode      string        `json:"auth_mode" jsonschema:"Authentication mode: api_key or oauth"`
}

func listOrgsHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcpsdk.CallToolRequest, input OrgListInput) (*mcpsdk.CallToolResult, OrgListOutput, error) {
	return func(ctx context.Context, req *mcpsdk.CallToolRequest, input OrgListInput) (*mcpsdk.CallToolResult, OrgListOutput, error) {
		fmt.Printf("[MCP org_list] Called with authMode=%v, userID=%s\n", toolCtx.AuthMode(), toolCtx.UserID())

		authMode := "api_key"
		if toolCtx.AuthMode() == mcpctx.AuthModeOAuth {
			authMode = "oauth"
		}

		// Initialize as empty slice (not nil) so JSON serializes as [] not null
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
			fmt.Printf("[MCP org_list] OAuth mode - fetching ALL orgs\n")
			allOrgs, err := toolCtx.DB().ListOrganizations(ctx)
			if err != nil {
				fmt.Printf("[MCP org_list] ERROR: %v\n", err)
				return nil, OrgListOutput{}, fmt.Errorf("failed to list organizations: %w", err)
			}
			fmt.Printf("[MCP org_list] Found %d organizations\n", len(allOrgs))

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
}

// OrgSelectInput defines input for org_select tool.
type OrgSelectInput struct {
	OrgID string `json:"org_id,omitempty" jsonschema:"Organization ID to select (UUID)"`
	Slug  string `json:"slug,omitempty" jsonschema:"Organization slug to select (alternative to org_id)"`
}

// OrgSelectOutput defines output for org_select tool.
type OrgSelectOutput struct {
	ID       string `json:"id" jsonschema:"Selected organization UUID"`
	Name     string `json:"name" jsonschema:"Selected organization name"`
	Slug     string `json:"slug" jsonschema:"Selected organization slug"`
	Selected bool   `json:"selected" jsonschema:"Whether selection was successful"`
	Message  string `json:"message,omitempty" jsonschema:"Additional information"`
}

func selectOrgHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcpsdk.CallToolRequest, input OrgSelectInput) (*mcpsdk.CallToolResult, OrgSelectOutput, error) {
	return func(ctx context.Context, req *mcpsdk.CallToolRequest, input OrgSelectInput) (*mcpsdk.CallToolResult, OrgSelectOutput, error) {
		if toolCtx.AuthMode() == mcpctx.AuthModeAPIKey {
			org := toolCtx.Org()
			return nil, OrgSelectOutput{
				ID:       org.ID,
				Name:     org.Name,
				Slug:     org.Slug,
				Selected: true,
				Message:  "API key authentication is already scoped to this organization. No selection needed.",
			}, nil
		}

		if input.OrgID == "" && input.Slug == "" {
			return nil, OrgSelectOutput{}, mcpctx.NewValidationError("either org_id or slug is required", "org_id")
		}

		var fullOrg db.Organization
		var err error

		if input.OrgID != "" {
			fullOrg, err = toolCtx.DB().GetOrganizationByID(ctx, input.OrgID)
		} else {
			fullOrg, err = toolCtx.DB().GetOrganizationBySlug(ctx, input.Slug)
		}
		if err != nil {
			if input.OrgID != "" {
				return nil, OrgSelectOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("organization with ID '%s' not found", input.OrgID))
			}
			return nil, OrgSelectOutput{}, mcpctx.NewNotFoundError(fmt.Sprintf("organization with slug '%s' not found", input.Slug))
		}

		toolCtx.SelectOrg(fullOrg)

		return nil, OrgSelectOutput{
			ID:       fullOrg.ID,
			Name:     fullOrg.Name,
			Slug:     fullOrg.Slug,
			Selected: true,
			Message:  "Organization selected. All subsequent operations will use this organization.",
		}, nil
	}
}

// OrgGetInput defines input for org_get tool.
type OrgGetInput struct{}

// OrgGetOutput defines output for org_get tool.
type OrgGetOutput struct {
	ID          string `json:"id" jsonschema:"Organization UUID"`
	Name        string `json:"name" jsonschema:"Organization name"`
	Slug        string `json:"slug" jsonschema:"URL-friendly slug"`
	MaxContacts int64  `json:"max_contacts,omitempty" jsonschema:"Maximum allowed contacts"`
	FromName    string `json:"from_name,omitempty" jsonschema:"Email sender name"`
	FromEmail   string `json:"from_email,omitempty" jsonschema:"Email sender address"`
	ReplyTo     string `json:"reply_to,omitempty" jsonschema:"Reply-to email address"`
}

func getOrgHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcpsdk.CallToolRequest, input OrgGetInput) (*mcpsdk.CallToolResult, OrgGetOutput, error) {
	return func(ctx context.Context, req *mcpsdk.CallToolRequest, input OrgGetInput) (*mcpsdk.CallToolResult, OrgGetOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, OrgGetOutput{}, err
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
}

// OrgUpdateInput defines input for org_update tool.
type OrgUpdateInput struct {
	Name      string `json:"name,omitempty" jsonschema:"New organization name"`
	FromName  string `json:"from_name,omitempty" jsonschema:"Email sender name (e.g., 'John from Acme')"`
	FromEmail string `json:"from_email,omitempty" jsonschema:"Email sender address (e.g., 'john@acme.com')"`
	ReplyTo   string `json:"reply_to,omitempty" jsonschema:"Reply-to email address"`
}

// OrgUpdateOutput defines output for org_update tool.
type OrgUpdateOutput struct {
	ID        string `json:"id" jsonschema:"Organization UUID"`
	Name      string `json:"name" jsonschema:"Organization name"`
	FromName  string `json:"from_name,omitempty" jsonschema:"Email sender name"`
	FromEmail string `json:"from_email,omitempty" jsonschema:"Email sender address"`
	ReplyTo   string `json:"reply_to,omitempty" jsonschema:"Reply-to email address"`
	Updated   bool   `json:"updated" jsonschema:"Whether update was successful"`
}

func updateOrgHandler(toolCtx *mcpctx.ToolContext) func(ctx context.Context, req *mcpsdk.CallToolRequest, input OrgUpdateInput) (*mcpsdk.CallToolResult, OrgUpdateOutput, error) {
	return func(ctx context.Context, req *mcpsdk.CallToolRequest, input OrgUpdateInput) (*mcpsdk.CallToolResult, OrgUpdateOutput, error) {
		if err := toolCtx.RequireOrg(); err != nil {
			return nil, OrgUpdateOutput{}, err
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
				return nil, OrgUpdateOutput{}, fmt.Errorf("failed to update email settings: %w", err)
			}
		}

		if input.Name != "" {
			_, err := toolCtx.DB().UpdateOrganization(ctx, db.UpdateOrganizationParams{
				ID:   toolCtx.OrgID(),
				Name: input.Name,
			})
			if err != nil {
				return nil, OrgUpdateOutput{}, fmt.Errorf("failed to update organization: %w", err)
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
}
