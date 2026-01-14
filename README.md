# Outlet.sh

The modern, self-hosted email platform for indie hackers and creators. Marketing campaigns, transactional emails, and AI-powered automation—all in a single binary.

## Why Outlet.sh?

- **Single binary deployment** — No PHP, no external database. Just download and run.
- **SQLite with WAL mode** — Zero configuration, handles 100k+ subscribers.
- **Marketing + Transactional** — One platform for newsletters and application emails.
- **AI-native MCP integration** — Build email sequences by talking to Claude.
- **Built-in backup system** — One-click download or automated S3 backups.
- **Full source code included** — No obfuscation, modify anything you need.

## Quick Start

### Option 1: Direct Binary

```bash
# Download the binary for your platform
wget https://outlet.sh/download/outlet-linux-amd64
chmod +x outlet-linux-amd64

# Set required environment variables
export JWT_SECRET="your-secret-key-here"
export JWT_REFRESH_SECRET="your-refresh-secret-here"
export ADMIN_EMAIL="admin@example.com"
export ADMIN_PASSWORD="your-admin-password"

# Run it
./outlet-linux-amd64
```

Open `http://localhost:9888` and log in with your admin credentials.

### Option 2: Docker Compose

```yaml
services:
  outlet:
    image: outlet/outlet:latest
    ports:
      - "443:443"
      - "80:80"
    volumes:
      - ./data:/data
    environment:
      - PRODUCTION_MODE=true
      - APP_DOMAIN=mail.yourdomain.com
      - JWT_SECRET=your-secret-key-here
      - JWT_REFRESH_SECRET=your-refresh-secret-here
      - ADMIN_EMAIL=admin@example.com
      - ADMIN_PASSWORD=your-admin-password
    restart: always
```

```bash
docker compose up -d
```

### Option 3: Systemd Service

Create `/etc/systemd/system/outlet.service`:

```ini
[Unit]
Description=Outlet.sh Email Platform
After=network.target

[Service]
Type=simple
User=outlet
WorkingDirectory=/opt/outlet
ExecStart=/opt/outlet/outlet
Restart=always
RestartSec=5
EnvironmentFile=/opt/outlet/.env

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable outlet
sudo systemctl start outlet
```

## Features

### Subscriber Management
- CSV import/export
- Custom fields and tags
- Segmentation and filtering
- Double opt-in flow
- Automatic bounce/complaint handling
- GDPR compliance tools

### Marketing Campaigns
- Rich text and HTML editor
- Template library
- Broadcast and scheduled campaigns
- RSS-to-email campaigns
- Open/click tracking
- Per-campaign analytics

### Transactional Email API
- REST API for application emails
- Template-based messages with merge fields
- Dedicated sending queue (separate from marketing)
- Per-template analytics
- Webhook callbacks for delivery events
- Auto-generated SDKs (TypeScript, Go)

### Automation
- Autoresponder sequences
- Trigger rules (tag added, link clicked, date-based)
- Transactional-to-marketing handoff

### AI-Powered Email (MCP Integration)

Outlet.sh includes a built-in MCP server, enabling AI assistants like Claude to manage your email operations:

```
You: "Create a 7-day welcome sequence for my SaaS. Day 1 introduces
     the product, Day 3 covers the main feature, Day 5 shares a case
     study, Day 7 offers a discount."

Claude: [Creates sequence, adds emails with proper delays, writes
        copy, sets triggers — all via MCP tools]
```

MCP capabilities include:
- List and subscriber management
- Campaign creation and scheduling
- Sequence building with email content
- Analytics queries and insights
- Transactional email sending

### Backup System

- **Dashboard download** — One-click database backup
- **S3 automated backup** — Scheduled backups with retention policies
- **CLI backup** — `outlet backup` or `outlet backup --s3`
- **Safe live backups** — Uses SQLite backup API (no corruption)
- **Easy restore** — Replace file and restart, or use dashboard

## Configuration

Configuration uses environment variables. Create a `.env` file or set them directly:

### Required

| Variable | Description |
|----------|-------------|
| `JWT_SECRET` | Secret key for access tokens |
| `JWT_REFRESH_SECRET` | Secret key for refresh tokens |
| `ADMIN_EMAIL` | Super admin email (created on first startup) |
| `ADMIN_PASSWORD` | Super admin password |

### Production

| Variable | Description |
|----------|-------------|
| `PRODUCTION_MODE` | Set to `true` for HTTPS/Let's Encrypt |
| `APP_DOMAIN` | Your domain (e.g., `mail.yourdomain.com`) |
| `APP_BASE_URL` | Full URL (e.g., `https://mail.yourdomain.com`) |

### Optional

| Variable | Description |
|----------|-------------|
| `DATABASE_PATH` | SQLite path (default: `./data/outlet.db`) |
| `ENCRYPTION_KEY` | 32-byte hex key for credential encryption |
| `ANTHROPIC_API_KEY` | Claude AI integration |
| `STRIPE_SECRET_KEY` | Payment processing |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook verification |

## Amazon SES Setup

Outlet.sh works best with Amazon SES for cost-effective sending (~$0.10 per 1,000 emails).

1. **Create SES credentials** in AWS Console → SES → SMTP Settings
2. **Verify your domain** in SES → Verified Identities
3. **Configure in Outlet.sh** dashboard → Settings → Email Provider
4. **Request production access** if still in SES sandbox

## API & SDKs

Outlet.sh provides a REST API with auto-generated clients:

### Transactional Email Example (TypeScript)

```typescript
import { OutletSDK } from '@outlet/sdk';

const outlet = new OutletSDK({
  baseUrl: 'https://mail.yourdomain.com',
  apiKey: 'your-api-key'
});

await outlet.transactional.send({
  to: 'user@example.com',
  template_id: 'welcome-email',
  variables: {
    name: 'John',
    activation_link: 'https://...'
  }
});
```

### Transactional Email Example (Go)

```go
import "github.com/outlet/outlet-go"

client := outlet.NewClient("https://mail.yourdomain.com", "your-api-key")

err := client.Transactional.Send(ctx, outlet.SendRequest{
    To:         "user@example.com",
    TemplateID: "welcome-email",
    Variables: map[string]string{
        "name":            "John",
        "activation_link": "https://...",
    },
})
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Single Binary                           │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  SvelteKit  │  │   go-zero   │  │    MCP Server       │  │
│  │  Frontend   │  │   REST API  │  │  (AI Integration)   │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Email     │  │   Rules     │  │     Webhook         │  │
│  │   Workers   │  │   Engine    │  │    Dispatcher       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                    SQLite (WAL Mode)                        │
│                    ./data/outlet.db                         │
└─────────────────────────────────────────────────────────────┘
```

- **Frontend**: Embedded SvelteKit 5 static build
- **Backend**: go-zero framework with auto-generated handlers
- **Database**: SQLite with WAL mode, auto-migrations on startup
- **Workers**: Background email sending with rate limiting
- **MCP**: Model Context Protocol server with embedded OAuth

## Comparison

| Feature | Sendy | ListMonk | Outlet.sh |
|---------|-------|----------|-----------|
| Price | $69 | Free | $99 |
| Stack | PHP/MySQL | Go/Postgres | Go/SQLite |
| Deployment | Complex | Medium | Single binary |
| Marketing emails | Yes | Yes | Yes |
| Transactional | No | Yes | Yes |
| Auto-updates | Manual | Manual | Built-in |
| Source code | Obfuscated | Open | Full source |
| AI/MCP integration | No | No | Yes |
| Built-in backup | No | No | Yes (+ S3) |
| External DB required | Yes | Yes | No |

## Development

### Prerequisites

- Go 1.21+
- Node.js 20+
- pnpm

### Setup

```bash
# Clone the repository
git clone https://github.com/outlet/outlet.sh
cd outlet.sh

# Install dependencies
go mod download
cd app && pnpm install && cd ..

# Create .env file
cp .env.example .env
# Edit .env with your values

# Run backend (uses air for hot reload)
make run

# In another terminal, run frontend
cd app && pnpm dev
```

### Build

```bash
# Build Go binary
make build

# Build frontend
cd app && pnpm build

# The binary embeds the frontend automatically
```

### Code Generation

```bash
# After modifying outlet.api
make gen          # Regenerates handlers, types, and TS client

# After modifying SQL queries
make sqlc-gen     # Regenerates type-safe Go from SQL

# All code generation
make gen-all
```

### Testing

```bash
# Go tests
go test ./...
go test ./internal/services/emailval/   # Single package
go test -run TestName ./internal/...    # Single test

# Frontend tests
cd app && pnpm test
```

## License

Outlet.sh is a commercial product with full source code included. See [LICENSE](LICENSE) for details.

- Personal and commercial use allowed
- Instance limits per license tier (honor system)
- Self-host forever (works without license server)
- All v1.x updates included free

## Support

- Documentation: https://outlet.sh/docs
- Issues: https://github.com/outlet/outlet.sh/issues
- Email: support@outlet.sh

---

Built with care for indie hackers who value simplicity, ownership, and modern tooling.
