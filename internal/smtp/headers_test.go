package smtp

import (
	"net/mail"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOutletHeaders_Empty(t *testing.T) {
	msg := &mail.Message{
		Header: mail.Header{},
	}

	h := ParseOutletHeaders(msg)

	assert.Equal(t, "", h.ListSlug)
	assert.Empty(t, h.Tags)
	assert.Equal(t, "", h.TemplateSlug)
	assert.Equal(t, "transactional", h.Type)
	assert.True(t, h.TrackOpens)
	assert.True(t, h.TrackClicks)
	assert.Empty(t, h.Meta)
}

func TestParseOutletHeaders_ListSlug(t *testing.T) {
	msg := &mail.Message{
		Header: mail.Header{
			"X-Outlet-List": []string{"newsletter"},
		},
	}

	h := ParseOutletHeaders(msg)

	assert.Equal(t, "newsletter", h.ListSlug)
}

func TestParseOutletHeaders_Tags(t *testing.T) {
	msg := &mail.Message{
		Header: mail.Header{
			"X-Outlet-Tags": []string{"welcome, new-user, vip"},
		},
	}

	h := ParseOutletHeaders(msg)

	assert.Equal(t, []string{"welcome", "new-user", "vip"}, h.Tags)
}

func TestParseOutletHeaders_Template(t *testing.T) {
	msg := &mail.Message{
		Header: mail.Header{
			"X-Outlet-Template": []string{"order-confirmation"},
		},
	}

	h := ParseOutletHeaders(msg)

	assert.Equal(t, "order-confirmation", h.TemplateSlug)
}

func TestParseOutletHeaders_Type_Marketing(t *testing.T) {
	msg := &mail.Message{
		Header: mail.Header{
			"X-Outlet-Type": []string{"marketing"},
		},
	}

	h := ParseOutletHeaders(msg)

	assert.Equal(t, "marketing", h.Type)
}

func TestParseOutletHeaders_Type_Transactional(t *testing.T) {
	msg := &mail.Message{
		Header: mail.Header{
			"X-Outlet-Type": []string{"TRANSACTIONAL"},
		},
	}

	h := ParseOutletHeaders(msg)

	assert.Equal(t, "transactional", h.Type)
}

func TestParseOutletHeaders_Track_None(t *testing.T) {
	msg := &mail.Message{
		Header: mail.Header{
			"X-Outlet-Track": []string{"none"},
		},
	}

	h := ParseOutletHeaders(msg)

	assert.False(t, h.TrackOpens)
	assert.False(t, h.TrackClicks)
}

func TestParseOutletHeaders_Track_OpensOnly(t *testing.T) {
	msg := &mail.Message{
		Header: mail.Header{
			"X-Outlet-Track": []string{"opens"},
		},
	}

	h := ParseOutletHeaders(msg)

	assert.True(t, h.TrackOpens)
	assert.False(t, h.TrackClicks)
}

func TestParseOutletHeaders_Track_ClicksOnly(t *testing.T) {
	msg := &mail.Message{
		Header: mail.Header{
			"X-Outlet-Track": []string{"clicks"},
		},
	}

	h := ParseOutletHeaders(msg)

	assert.False(t, h.TrackOpens)
	assert.True(t, h.TrackClicks)
}

func TestParseOutletHeaders_Track_Both(t *testing.T) {
	msg := &mail.Message{
		Header: mail.Header{
			"X-Outlet-Track": []string{"opens, clicks"},
		},
	}

	h := ParseOutletHeaders(msg)

	assert.True(t, h.TrackOpens)
	assert.True(t, h.TrackClicks)
}

func TestParseOutletHeaders_Meta(t *testing.T) {
	// Test with raw email to use canonical header format
	raw := `From: sender@example.com
To: recipient@example.com
Subject: Test
X-Outlet-Meta-OrderId: 12345
X-Outlet-Meta-CustomerId: cust_abc

Body`

	msg, err := mail.ReadMessage(strings.NewReader(raw))
	assert.NoError(t, err)

	h := ParseOutletHeaders(msg)

	// Headers get canonicalized, so check with canonical form
	assert.Equal(t, "12345", h.Meta["Orderid"])
	assert.Equal(t, "cust_abc", h.Meta["Customerid"])
}

func TestParseOutletHeaders_AllFields(t *testing.T) {
	raw := `From: sender@example.com
To: recipient@example.com
Subject: Test
X-Outlet-List: newsletter
X-Outlet-Tags: welcome,vip
X-Outlet-Template: welcome-email
X-Outlet-Type: marketing
X-Outlet-Track: opens
X-Outlet-Meta-OrderId: 123

Body`

	msg, err := mail.ReadMessage(strings.NewReader(raw))
	assert.NoError(t, err)

	h := ParseOutletHeaders(msg)

	assert.Equal(t, "newsletter", h.ListSlug)
	assert.Equal(t, []string{"welcome", "vip"}, h.Tags)
	assert.Equal(t, "welcome-email", h.TemplateSlug)
	assert.Equal(t, "marketing", h.Type)
	assert.True(t, h.TrackOpens)
	assert.False(t, h.TrackClicks)
	assert.Equal(t, "123", h.Meta["Orderid"]) // canonical form
}

func TestParseOutletHeaders_FromRawEmail(t *testing.T) {
	raw := `From: sender@example.com
To: recipient@example.com
Subject: Test Email
X-Outlet-List: newsletter
X-Outlet-Tags: welcome, test
X-Outlet-Type: marketing

Hello World`

	msg, err := mail.ReadMessage(strings.NewReader(raw))
	assert.NoError(t, err)

	h := ParseOutletHeaders(msg)

	assert.Equal(t, "newsletter", h.ListSlug)
	assert.Equal(t, []string{"welcome", "test"}, h.Tags)
	assert.Equal(t, "marketing", h.Type)
}
