package tools

import (
	"testing"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/mcp/mcpctx"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// createTestOrg creates a mock organization for testing
func createTestOrg(id string, name string) db.Organization {
	return db.Organization{
		ID:     id,
		Name:   name,
		Slug:   "test-org",
		ApiKey: "test-api-key",
	}
}

// createTestUser creates a mock user for testing
func createTestUser(id string, email, role string) db.User {
	return db.User{
		ID:     id,
		Email:  email,
		Name:   "Test User",
		Role:   role,
		Status: "active",
	}
}

// ========== EmailInput Tests (Unified Pattern) ==========

func TestEmailInput_ListCreate(t *testing.T) {
	doubleOptin := true
	input := EmailInput{
		Resource:    "list",
		Action:      "create",
		Name:        "Newsletter",
		Description: "Weekly updates",
		DoubleOptin: &doubleOptin,
	}

	assert.Equal(t, "list", input.Resource)
	assert.Equal(t, "create", input.Action)
	assert.Equal(t, "Newsletter", input.Name)
	assert.Equal(t, "Weekly updates", input.Description)
	assert.NotNil(t, input.DoubleOptin)
	assert.True(t, *input.DoubleOptin)
}

func TestEmailInput_ListGet(t *testing.T) {
	input := EmailInput{
		Resource: "list",
		Action:   "get",
		ID:       "123",
	}

	assert.Equal(t, "list", input.Resource)
	assert.Equal(t, "get", input.Action)
	assert.Equal(t, "123", input.ID)
}

func TestEmailInput_SequenceCreate(t *testing.T) {
	active := true
	input := EmailInput{
		Resource:     "sequence",
		Action:       "create",
		Name:         "Welcome Series",
		ListID:       "123",
		TriggerEvent: "on_subscribe",
		SequenceType: "lifecycle",
		Active:       &active,
	}

	assert.Equal(t, "sequence", input.Resource)
	assert.Equal(t, "create", input.Action)
	assert.Equal(t, "Welcome Series", input.Name)
	assert.Equal(t, "123", input.ListID)
	assert.Equal(t, "on_subscribe", input.TriggerEvent)
	assert.Equal(t, "lifecycle", input.SequenceType)
	assert.NotNil(t, input.Active)
	assert.True(t, *input.Active)
}

func TestEmailInput_TemplateCreate(t *testing.T) {
	input := EmailInput{
		Resource:   "template",
		Action:     "create",
		SequenceID: uuid.New().String(),
		Subject:    "Welcome!",
		HTMLBody:   "<h1>Hello</h1>",
		DelayHours: 24,
	}

	assert.Equal(t, "template", input.Resource)
	assert.Equal(t, "create", input.Action)
	assert.NotEmpty(t, input.SequenceID)
	assert.Equal(t, "Welcome!", input.Subject)
	assert.Equal(t, "<h1>Hello</h1>", input.HTMLBody)
	assert.Equal(t, 24, input.DelayHours)
}

// ========== ListCreateOutput Tests ==========

func TestListCreateOutput_Structure(t *testing.T) {
	output := ListCreateOutput{
		ID:          "1",
		Name:        "Newsletter",
		Slug:        "newsletter",
		Description: "Weekly updates",
		DoubleOptin: true,
		Created:     true,
	}

	assert.Equal(t, "1", output.ID)
	assert.Equal(t, "Newsletter", output.Name)
	assert.Equal(t, "newsletter", output.Slug)
	assert.Equal(t, "Weekly updates", output.Description)
	assert.True(t, output.DoubleOptin)
	assert.True(t, output.Created)
}

// ========== ListItem Tests ==========

func TestListItem_Structure(t *testing.T) {
	output := ListItem{
		ID:              "1",
		Name:            "Newsletter",
		Slug:            "newsletter",
		Description:     "Weekly updates",
		DoubleOptin:     true,
		SubscriberCount: 100,
	}

	assert.Equal(t, "1", output.ID)
	assert.Equal(t, "Newsletter", output.Name)
	assert.True(t, output.DoubleOptin)
	assert.Equal(t, int64(100), output.SubscriberCount)
}

// ========== ListListOutput Tests ==========

func TestListListOutput_Structure(t *testing.T) {
	output := ListListOutput{
		Lists: []ListItem{
			{ID: "1", Name: "Newsletter", SubscriberCount: 100},
			{ID: "2", Name: "Announcements", SubscriberCount: 50},
		},
		Total: 2,
	}

	assert.Len(t, output.Lists, 2)
	assert.Equal(t, int64(100), output.Lists[0].SubscriberCount)
	assert.Equal(t, 2, output.Total)
}

// ========== SequenceItem Tests ==========

func TestSequenceItem_Structure(t *testing.T) {
	output := SequenceItem{
		ID:           uuid.New().String(),
		Name:         "Welcome",
		Slug:         "welcome",
		ListID:       "123",
		SequenceType: "lifecycle",
		TriggerEvent: "on_subscribe",
		Active:       true,
		EmailCount:   5,
	}

	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "Welcome", output.Name)
	assert.True(t, output.Active)
	assert.Equal(t, 5, output.EmailCount)
}

// ========== SequenceListOutput Tests ==========

func TestSequenceListOutput_Structure(t *testing.T) {
	output := SequenceListOutput{
		Sequences: []SequenceItem{
			{ID: uuid.New().String(), Name: "Welcome", Active: true, EmailCount: 5},
			{ID: uuid.New().String(), Name: "Onboarding", Active: false, EmailCount: 3},
		},
		Total: 2,
	}

	assert.Len(t, output.Sequences, 2)
	assert.True(t, output.Sequences[0].Active)
	assert.Equal(t, 5, output.Sequences[0].EmailCount)
	assert.Equal(t, 2, output.Total)
}

// ========== SequenceCreateOutput Tests ==========

func TestSequenceCreateOutput_Structure(t *testing.T) {
	output := SequenceCreateOutput{
		ID:           uuid.New().String(),
		Name:         "Welcome Series",
		Slug:         "welcome-series",
		ListID:       "123",
		TriggerEvent: "on_subscribe",
		SequenceType: "lifecycle",
		Active:       true,
		Created:      true,
	}

	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "Welcome Series", output.Name)
	assert.Equal(t, "welcome-series", output.Slug)
	assert.Equal(t, "123", output.ListID)
	assert.True(t, output.Active)
	assert.True(t, output.Created)
}

// ========== TemplateCreateOutput Tests ==========

func TestTemplateCreateOutput_Structure(t *testing.T) {
	output := TemplateCreateOutput{
		ID:         uuid.New().String(),
		SequenceID: uuid.New().String(),
		Subject:    "Welcome!",
		Position:   1,
		DelayHours: 0,
		Active:     true,
		Created:    true,
	}

	assert.NotEmpty(t, output.ID)
	assert.NotEmpty(t, output.SequenceID)
	assert.Equal(t, "Welcome!", output.Subject)
	assert.Equal(t, 1, output.Position)
	assert.True(t, output.Active)
	assert.True(t, output.Created)
}

// ========== Tool Context Requirement Tests ==========

func TestToolContext_RequireOrg_WithOrg(t *testing.T) {
	orgID := uuid.New().String()
	org := createTestOrg(orgID, "Test Org")
	tc := mcpctx.NewToolContext(nil, org, "req-123", "test-agent/1.0")

	err := tc.RequireOrg()
	assert.NoError(t, err, "RequireOrg should not error when org is set")
}

func TestToolContext_RequireOrg_WithoutOrg(t *testing.T) {
	userID := uuid.New().String()
	user := createTestUser(userID, "test@example.com", "admin")
	tc := mcpctx.NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-123")

	err := tc.RequireOrg()
	assert.Error(t, err, "RequireOrg should error when no org selected")
	assert.Equal(t, mcpctx.ErrNoOrgSelected, err)
}

// ========== OrgInput Tests (Unified Pattern) ==========

func TestOrgInput_List(t *testing.T) {
	input := OrgInput{
		Resource: "org",
		Action:   "list",
	}

	assert.Equal(t, "org", input.Resource)
	assert.Equal(t, "list", input.Action)
}

func TestOrgInput_Select_WithID(t *testing.T) {
	input := OrgInput{
		Resource: "org",
		Action:   "select",
		OrgID:    uuid.New().String(),
	}

	assert.Equal(t, "org", input.Resource)
	assert.Equal(t, "select", input.Action)
	assert.NotEmpty(t, input.OrgID)
	assert.Empty(t, input.Slug)
}

func TestOrgInput_Select_WithSlug(t *testing.T) {
	input := OrgInput{
		Resource: "org",
		Action:   "select",
		Slug:     "my-org",
	}

	assert.Equal(t, "org", input.Resource)
	assert.Equal(t, "select", input.Action)
	assert.Empty(t, input.OrgID)
	assert.Equal(t, "my-org", input.Slug)
}

func TestOrgInput_Update(t *testing.T) {
	input := OrgInput{
		Resource:  "org",
		Action:    "update",
		Name:      "My Company",
		FromName:  "John from My Company",
		FromEmail: "john@mycompany.com",
		ReplyTo:   "support@mycompany.com",
	}

	assert.Equal(t, "org", input.Resource)
	assert.Equal(t, "update", input.Action)
	assert.Equal(t, "My Company", input.Name)
	assert.Equal(t, "John from My Company", input.FromName)
	assert.Equal(t, "john@mycompany.com", input.FromEmail)
	assert.Equal(t, "support@mycompany.com", input.ReplyTo)
}

// ========== OrgListItem Tests ==========

func TestOrgListItem_Structure(t *testing.T) {
	item := OrgListItem{
		ID:       uuid.New().String(),
		Name:     "My Company",
		Slug:     "my-company",
		Role:     "admin",
		Selected: true,
	}

	assert.NotEmpty(t, item.ID)
	assert.Equal(t, "My Company", item.Name)
	assert.Equal(t, "my-company", item.Slug)
	assert.Equal(t, "admin", item.Role)
	assert.True(t, item.Selected)
}

// ========== OrgListOutput Tests ==========

func TestOrgListOutput_Structure(t *testing.T) {
	output := OrgListOutput{
		Organizations: []OrgListItem{
			{ID: uuid.New().String(), Name: "Org 1", Slug: "org-1", Selected: true},
			{ID: uuid.New().String(), Name: "Org 2", Slug: "org-2", Selected: false},
		},
		Total:    2,
		AuthMode: "api_key",
	}

	assert.Len(t, output.Organizations, 2)
	assert.Equal(t, 2, output.Total)
	assert.Equal(t, "api_key", output.AuthMode)
}

// ========== OrgSelectOutput Tests ==========

func TestOrgSelectOutput_Structure(t *testing.T) {
	output := OrgSelectOutput{
		ID:       uuid.New().String(),
		Name:     "My Company",
		Slug:     "my-company",
		Selected: true,
		Message:  "Organization selected.",
	}

	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "My Company", output.Name)
	assert.True(t, output.Selected)
	assert.Contains(t, output.Message, "selected")
}

// ========== OrgGetOutput Tests ==========

func TestOrgGetOutput_Structure(t *testing.T) {
	output := OrgGetOutput{
		ID:          uuid.New().String(),
		Name:        "My Company",
		Slug:        "my-company",
		MaxContacts: 10000,
		FromName:    "Support",
		FromEmail:   "support@mycompany.com",
		ReplyTo:     "reply@mycompany.com",
	}

	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "My Company", output.Name)
	assert.Equal(t, int64(10000), output.MaxContacts)
	assert.Equal(t, "support@mycompany.com", output.FromEmail)
}

// ========== OrgUpdateOutput Tests ==========

func TestOrgUpdateOutput_Structure(t *testing.T) {
	output := OrgUpdateOutput{
		ID:        uuid.New().String(),
		Name:      "My Company",
		FromName:  "Support",
		FromEmail: "support@mycompany.com",
		ReplyTo:   "reply@mycompany.com",
		Updated:   true,
	}

	assert.NotEmpty(t, output.ID)
	assert.True(t, output.Updated)
}

// ========== generateSlug Tests ==========

func TestGenerateSlug_Simple(t *testing.T) {
	slug := generateSlug("Hello World")
	assert.Equal(t, "hello-world", slug)
}

func TestGenerateSlug_SpecialCharacters(t *testing.T) {
	slug := generateSlug("Product #1 - Special!")
	assert.Equal(t, "product-1-special", slug)
}

func TestGenerateSlug_MultipleSpaces(t *testing.T) {
	slug := generateSlug("Hello    World")
	assert.Equal(t, "hello-world", slug)
}

func TestGenerateSlug_LeadingTrailingSpaces(t *testing.T) {
	slug := generateSlug("  Hello World  ")
	assert.Equal(t, "hello-world", slug)
}

func TestGenerateSlug_Numbers(t *testing.T) {
	slug := generateSlug("Product 123")
	assert.Equal(t, "product-123", slug)
}

func TestGenerateSlug_AllCaps(t *testing.T) {
	slug := generateSlug("HELLO WORLD")
	assert.Equal(t, "hello-world", slug)
}

// ========== Action Validation Tests ==========

func TestEmailActions_ValidActions(t *testing.T) {
	// Test that list has all expected actions
	listActions := emailActions["list"]
	assert.Contains(t, listActions, "create")
	assert.Contains(t, listActions, "list")
	assert.Contains(t, listActions, "get")
	assert.Contains(t, listActions, "update")
	assert.Contains(t, listActions, "delete")

	// Test that sequence has all expected actions
	seqActions := emailActions["sequence"]
	assert.Contains(t, seqActions, "create")
	assert.Contains(t, seqActions, "list")
	assert.Contains(t, seqActions, "get")
	assert.Contains(t, seqActions, "update")
	assert.Contains(t, seqActions, "delete")

	// Test that template has all expected actions
	templateActions := emailActions["template"]
	assert.Contains(t, templateActions, "create")
	assert.Contains(t, templateActions, "list")
	assert.Contains(t, templateActions, "get")
	assert.Contains(t, templateActions, "update")
	assert.Contains(t, templateActions, "delete")
}

func TestOrgActions_ValidActions(t *testing.T) {
	orgActions := orgActions["org"]
	assert.Contains(t, orgActions, "list")
	assert.Contains(t, orgActions, "select")
	assert.Contains(t, orgActions, "get")
	assert.Contains(t, orgActions, "update")
}
