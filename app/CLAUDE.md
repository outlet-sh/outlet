# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

SvelteKit 5 frontend using Svelte 5 runes, Tailwind CSS 4, and mdsvex for markdown content. Static site generation via `adapter-static`, embedded into Go binary for production.

See root `../CLAUDE.md` for full-stack architecture and backend details.

## Commands

```bash
pnpm dev          # Start dev server (port 5173)
pnpm build        # Build static site to build/
pnpm check        # Type check
pnpm test         # Run tests with Vitest
pnpm test:ui      # Run tests with Vitest UI
```

## Path Aliases

- `$lib` → `src/lib/`
- `$content` → `src/content/`

Use `$content` for importing markdown content:

```typescript
const article = await import(`$content/articles/${slug}.md`);
```

## Route Groups

Routes use SvelteKit's group syntax `(groupname)` for layout organization:

```
src/routes/
├── (www)/              # Public marketing site
│   ├── (home)/         # Homepage with unique layout
│   └── (pages)/        # Standard pages (articles, books, projects, etc.)
├── (admin)/            # Admin dashboard (auth required)
├── (auth)/             # Login/auth pages
├── (lp)/               # Landing pages
│   └── lp/10-in-100/   # Funnel landing page
├── (funnel)/           # Funnel checkout flow
├── (workshop)/         # Workshop registration
└── (thankyou)/         # Thank you pages
```

Each group can have its own `+layout.svelte` for different page structures.

## Content System

Markdown files in `src/content/` are processed by mdsvex:

```
src/content/
├── articles/           # Blog articles
├── books/              # Book pages
├── projects/           # Project showcases
└── build-in-public/    # Build in public posts
```

**Frontmatter structure** (articles):

```yaml
---
title: 'Article Title'
slug: 'article-slug'
date: '2024-01-15'
author: 'Alma Tuck'
excerpt: 'Short description'
tags: ['tag1', 'tag2']
featured: true
readTime: '5 min'
---
```

**Loading content** in `+page.ts`:

```typescript
// List all articles
const files = import.meta.glob('$content/articles/*.md', { eager: true, as: 'raw' });

// Single article by slug
const article = await import(`$content/articles/${params.slug}.md`);
```

## Generated API Client (CRITICAL)

TypeScript client generated from `outlet.api` lives in `src/lib/api/generated/`:

```typescript
import { FunnelOptin, ConfirmEmail } from '$lib/api/generated/outlet';

const response = await FunnelOptin({ email, funnel_slug });
```

Regenerate after API changes: `make gen` (from repo root)

**IMPORTANT: Generated types use snake_case (matching Go/JSON)**

Always check `outletComponents.ts` for the correct property names. DO NOT assume camelCase:

```typescript
// WRONG - assumes camelCase
stats.withStripe;
stats.newLast30Days;
customer.emailVerified;
customer.createdAt;
api.adminListCustomers({ pageSize: 10 });

// CORRECT - matches generated types
stats.with_stripe;
stats.new_last_30_days;
customer.email_verified;
customer.created_at;
api.adminListCustomers({ page_size: 10 });
```

Before using any API type, open `src/lib/api/generate/outletComponents.ts` and verify the exact property names.

## Svelte 5 Conventions

Uses Svelte 5 runes syntax:

- `$state()` for reactive state
- `$derived()` for computed values
- `$effect()` for side effects
- `$props()` for component props

```svelte
<script lang="ts">
  let { data } = $props();
  let count = $state(0);
  let doubled = $derived(count * 2);
</script>
```

## UI Components (CRITICAL)

**MUST use `$lib/components/ui` components. NEVER use raw HTML for standard UI elements.**

Available components (import from `$lib/components/ui`):

- **Layout**: Card, Modal, Drawer, Sheet, Page, PageHeader, Tabs, Collapsible, CollapsibleCard
- **Forms**: Input, Select, Textarea, Checkbox, Radio, Toggle, Slider, ValidatedInput, TagInput, SearchInput, Button, ButtonGroup, Form, DateRangePicker, DateFilterPopover
- **Feedback**: Alert, AlertDialog, Toast, Tooltip, LoadingSpinner, Spinner, Skeleton, ProgressBar, EmptyState
- **Data Display**: Badge, StatusBadge, Table, DataTable, MetricCard, Sparkline, Chart, Avatar
- **Navigation**: Breadcrumb, Dropdown, DropdownMenu
- **Editors**: EmailEditor, MarkdownEditor, ContentEditor, ContentChatPanel, Markdown
- **Misc**: Divider, Separator, ScrollArea, ThemeToggle, GradientBackground, ModelSelector, CampaignSelector

Example:

```svelte
<script>
  import { Button, Input, Card, Modal, Checkbox, Badge } from '$lib/components/ui';
</script>

<!-- CORRECT -->
<Checkbox bind:checked={value} label="Enable feature" />

<!-- WRONG - never use raw HTML for checkboxes -->
<input type="checkbox" bind:checked={value} />
```

## Styling

- Tailwind CSS 4 via `@tailwindcss/vite`
- Global styles in `src/app.css`
- No inline `<style>` blocks - use Tailwind classes only
- Typography plugin for markdown: `prose` classes

## Prerendering

Marketing pages are prerendered for SEO. Set in layout:

```typescript
export const prerender = true;
export const trailingSlash = 'never';
```

Dynamic routes (admin, auth) use SPA fallback via `200.html`.
