// Package signalnet is the client for the Outlet Signal Network: pooled, hashed
// deliverability signals shared across Outlet instances through a hub run by outlet.sh.
// Emails never leave an instance in clear text; only SHA-256 of the lowercased address.
package signalnet

import (
	"context"
	"time"
)

// Signal is what the network knows about one address across all reporting instances.
type Signal struct {
	Hash          string    `json:"hash"`
	Senders       int       `json:"senders"`      // distinct instances that mailed it
	Delivered     int       `json:"delivered"`
	HardBounces   int       `json:"hard_bounces"`
	Complaints    int       `json:"complaints"`
	Engaged       int       `json:"engaged"`      // opens or clicks
	LastEngagedAt time.Time `json:"last_engaged_at"`
	TrapScore     float64   `json:"trap_score"`   // 0..1, hub heuristic
	Toxic         bool      `json:"toxic"`
}

// Event is one deliverability outcome reported by an instance.
type Event struct {
	Hash string    `json:"hash"`
	Kind string    `json:"kind"` // delivered, hard_bounce, complaint, open, click
	At   time.Time `json:"at"`
}

// Client talks to the hub. A nil *Client is valid and disabled.
type Client struct {
	HubURL      string
	InstanceKey string
}

// Hash returns the network identifier for an email address.
func Hash(email string) string { return hash(email) }

// Lookup returns the network's signal for email, or nil when disabled or unknown.
func (c *Client) Lookup(ctx context.Context, email string) (*Signal, error) { return nil, nil }

// Report sends outcomes to the hub. No-op when disabled.
func (c *Client) Report(ctx context.Context, events []Event) error { return nil }
