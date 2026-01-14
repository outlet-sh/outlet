# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Outlet.sh is a self-hosted email platform (marketing + transactional) built with go-zero framework and SvelteKit. Single binary deployment with embedded frontend, SQLite database (WAL mode), and MCP (Model Context Protocol) integration for AI assistants.

## Commands

### Backend (Go)
```bash
make build        # Build main binary to bin/outlet
make run          # Build and run
make test         # Run Go tests
make gen          # Regenerate API code from outlet.api (NEVER run goctl directly)
make sqlc-gen     # Generate type-safe Go code from SQL queries
make gen-all      # Generate all code (API + sqlc + proto)
```

### Frontend (SvelteKit)
```bash
cd app
pnpm dev          # Start dev server (port 5173)
pnpm build        # Build static site to build/
pnpm check        # Type check
pnpm test         # Run tests with Vitest
```

### Database (SQLite)
Database migrations run automatically on startup via embedded Goose migrations. No manual migration commands needed.

To manually interact with the database:
```bash
sqlite3 ./data/outlet.db   # Open SQLite CLI
```

### Running Tests
```bash
go test ./...                           # All Go tests
go test ./internal/services/emailval/   # Single package
go test -run TestSyntax ./internal/...  # Single test by name
cd app && pnpm test                     # Frontend tests
```

## Architecture

### Go Backend Structure
```
outlet.go              # Main entry point - embedded config, HTTP/HTTPS/gRPC servers
outlet.api             # API definition file - defines all routes and types
internal/
├── config/            # Configuration structs (loaded from etc/outlet.yaml)
├── handler/           # Auto-generated request handlers (DO NOT EDIT)
├── logic/             # Business logic implementations
│   ├── admin/         # Admin dashboard endpoints
│   ├── auth/          # Authentication logic
│   ├── public/        # Public endpoints (signup, confirm)
│   └── sdk/           # SDK/API endpoints for external use
├── types/             # Auto-generated request/response types (DO NOT EDIT)
├── svc/               # ServiceContext - dependency injection container
├── db/                # Database layer
│   ├── migrations/    # Goose SQL migrations (auto-run on startup)
│   ├── queries/       # SQL query files for sqlc
│   └── *.sql.go       # sqlc-generated query code
├── services/          # Domain services
│   ├── email/         # Email sending (SMTP, AWS SES)
│   ├── crypto/        # Credential encryption
│   ├── emailval/      # Email validation
│   ├── rules/         # Business rules engine
│   ├── tracking/      # Email open/click tracking
│   └── webhook/       # Outbound webhook dispatcher
├── mcp/               # MCP server for AI integrations
├── middleware/        # Auth, rate limiting, API key validation
├── events/            # Event bus for internal pub/sub
├── workers/           # Background workers (email sending, retries)
└── webhook/           # Inbound webhook handlers
```

### Frontend Structure
See `app/CLAUDE.md` for detailed frontend guidance.

### Key Patterns

**API Code Generation**: The `outlet.api` file defines all routes and types. Running `make gen` generates:
- `internal/handler/` - HTTP handlers
- `internal/types/` - Request/response structs
- `app/src/lib/api/generate/` - TypeScript client

**Database Layer**: sqlc generates type-safe Go from SQL:
- Write queries in `internal/db/queries/*.sql`
- Run `make sqlc-gen` to regenerate `internal/db/*.sql.go`
- Migrations in `internal/db/migrations/` auto-run on startup

**ServiceContext**: All services are initialized in `internal/svc/servicecontext.go` and passed to handlers via dependency injection.

**Production Mode**: In production, the binary:
- Serves Let's Encrypt HTTPS on :443
- Embeds the SvelteKit static build
- Redirects www → non-www, HTTP → HTTPS
- Routes /api/*, /sdk/*, /mcp to go-zero backend

**Development Mode**: When `PRODUCTION_MODE` is unset:
- Go backend on port 9888
- Frontend dev server on port 5173 (run separately)
- gRPC on port 9889

## Critical Rules

1. **NEVER run goctl commands directly** - Always use `make gen` to regenerate API code
2. **Use pnpm** for all frontend package management
3. **No hot reload needed** - Project uses air for Go hot reloading
4. **Always build before pushing** - Run `make build` and `cd app && pnpm build`
5. **Idiomatic Go only** - One function with parameters, not multiple function variants
6. **Styles in app.css only** - No inline styles or `<style>` blocks in Svelte files
7. **Minimal changes** - Only modify code directly related to the task
8. **Never assume code is unused** - Code may be called from frontend, other services, or future features

## Environment Configuration

Configuration is embedded from `etc/outlet.yaml` with environment variable expansion. Key variables:
- `DATABASE_PATH` - SQLite database path (default: `./data/outlet.db`)
- `JWT_SECRET`, `JWT_REFRESH_SECRET` - Auth tokens
- `ANTHROPIC_API_KEY` - Claude AI integration
- `ENCRYPTION_KEY` - 32-byte hex key for credential encryption
- `PRODUCTION_MODE` - Enable HTTPS/Let's Encrypt
- `APP_DOMAIN` - Domain for production (e.g., outlet.sh)
- `ADMIN_EMAIL`, `ADMIN_PASSWORD` - Super admin seeding on first startup
