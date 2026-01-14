package tools

import (
	"testing"

	"outlet/internal/db"
	"outlet/internal/mcp/mcpctx"

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

// ========== ListCreateInput Tests ==========

func TestListCreateInput_Required(t *testing.T) {
	input := ListCreateInput{
		Name: "Newsletter",
	}

	assert.Equal(t, "Newsletter", input.Name)
	assert.Empty(t, input.Description)
	assert.Nil(t, input.DoubleOptin)
}

func TestListCreateInput_WithOptions(t *testing.T) {
	doubleOptin := true
	input := ListCreateInput{
		Name:        "Newsletter",
		Description: "Weekly updates",
		DoubleOptin: &doubleOptin,
	}

	assert.Equal(t, "Newsletter", input.Name)
	assert.Equal(t, "Weekly updates", input.Description)
	assert.NotNil(t, input.DoubleOptin)
	assert.True(t, *input.DoubleOptin)
}

func TestListCreateInput_EmptyName(t *testing.T) {
	input := ListCreateInput{
		Name: "", // Empty name should fail validation
	}

	assert.Empty(t, input.Name)
}

// ========== ListCreateOutput Tests ==========

func TestListCreateOutput_Structure(t *testing.T) {
	output := ListCreateOutput{
		ID:          "1",
		Name:        "Newsletter",
		Slug:        "newsletter",
		Description: "Weekly updates",
		DoubleOptin: true,
	}

	assert.Equal(t, "1", output.ID)
	assert.Equal(t, "Newsletter", output.Name)
	assert.Equal(t, "newsletter", output.Slug)
	assert.Equal(t, "Weekly updates", output.Description)
	assert.True(t, output.DoubleOptin)
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

// ========== SequenceCreateInput Tests ==========

func TestSequenceCreateInput_Structure(t *testing.T) {
	input := SequenceCreateInput{
		Name:         "Welcome Series",
		ListID:       "123",
		SequenceType: "lifecycle",
		TriggerEvent: "on_subscribe",
	}

	assert.Equal(t, "Welcome Series", input.Name)
	assert.Equal(t, "123", input.ListID)
	assert.Equal(t, "lifecycle", input.SequenceType)
	assert.Equal(t, "on_subscribe", input.TriggerEvent)
}

func TestSequenceCreateInput_WithActive(t *testing.T) {
	active := true
	input := SequenceCreateInput{
		Name:   "Onboarding",
		ListID: "456",
		Active: &active,
	}

	assert.Equal(t, "Onboarding", input.Name)
	assert.NotNil(t, input.Active)
	assert.True(t, *input.Active)
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
	}

	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "Welcome Series", output.Name)
	assert.Equal(t, "welcome-series", output.Slug)
	assert.Equal(t, "123", output.ListID)
	assert.True(t, output.Active)
}

// ========== SequenceListItem Tests ==========

func TestSequenceListItem_Structure(t *testing.T) {
	output := SequenceListItem{
		ID:           uuid.New().String(),
		Name:         "Welcome",
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
		Sequences: []SequenceListItem{
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

// ========== SequenceEmailAddInput Tests ==========

func TestSequenceEmailAddInput_Structure(t *testing.T) {
	input := SequenceEmailAddInput{
		SequenceID: uuid.New().String(),
		Subject:    "Welcome!",
		HTMLBody:   "<h1>Hello</h1>",
		DelayHours: 24,
	}

	assert.NotEmpty(t, input.SequenceID)
	assert.Equal(t, "Welcome!", input.Subject)
	assert.Equal(t, "<h1>Hello</h1>", input.HTMLBody)
	assert.Equal(t, 24, input.DelayHours)
}

func TestSequenceEmailAddInput_WithOptionalFields(t *testing.T) {
	active := true
	input := SequenceEmailAddInput{
		SequenceID: uuid.New().String(),
		Subject:    "Welcome!",
		HTMLBody:   "<h1>Hello</h1>",
		PlainText:  "Hello",
		Position:   1,
		Active:     &active,
	}

	assert.Equal(t, "Hello", input.PlainText)
	assert.Equal(t, 1, input.Position)
	assert.NotNil(t, input.Active)
	assert.True(t, *input.Active)
}

// ========== SequenceEmailAddOutput Tests ==========

func TestSequenceEmailAddOutput_Structure(t *testing.T) {
	output := SequenceEmailAddOutput{
		ID:         uuid.New().String(),
		SequenceID: uuid.New().String(),
		Subject:    "Welcome!",
		Position:   1,
		DelayHours: 0,
		Active:     true,
	}

	assert.NotEmpty(t, output.ID)
	assert.NotEmpty(t, output.SequenceID)
	assert.Equal(t, "Welcome!", output.Subject)
	assert.Equal(t, 1, output.Position)
	assert.True(t, output.Active)
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

// ========== OrgSelectInput Tests ==========

func TestOrgSelectInput_WithID(t *testing.T) {
	input := OrgSelectInput{
		OrgID: uuid.New().String(),
	}

	assert.NotEmpty(t, input.OrgID)
	assert.Empty(t, input.Slug)
}

func TestOrgSelectInput_WithSlug(t *testing.T) {
	input := OrgSelectInput{
		Slug: "my-org",
	}

	assert.Empty(t, input.OrgID)
	assert.Equal(t, "my-org", input.Slug)
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

// ========== OrgUpdateInput Tests ==========

func TestOrgUpdateInput_Structure(t *testing.T) {
	input := OrgUpdateInput{
		Name:      "My Company",
		FromName:  "John from My Company",
		FromEmail: "john@mycompany.com",
		ReplyTo:   "support@mycompany.com",
	}

	assert.Equal(t, "My Company", input.Name)
	assert.Equal(t, "John from My Company", input.FromName)
	assert.Equal(t, "john@mycompany.com", input.FromEmail)
	assert.Equal(t, "support@mycompany.com", input.ReplyTo)
}

func TestOrgUpdateInput_Partial(t *testing.T) {
	input := OrgUpdateInput{
		Name: "New Company Name",
	}

	assert.Equal(t, "New Company Name", input.Name)
	assert.Empty(t, input.FromName)
	assert.Empty(t, input.FromEmail)
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
