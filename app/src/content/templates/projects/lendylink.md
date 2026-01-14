---
title: 'LendyLink'
slug: 'lendylink'
logo: '/images/projects/lendylink.svg'
website: 'https://lendylink.com'
founded: '2025'
role: 'Co-Founder'
status: 'Active'
category: 'Fintech'
featured: true
order: 3
tagline: 'A clean, modern ADU financing platform that helps homeowners move from idea to loan approval without the usual friction.'
description: 'LendyLink is an end-to-end ADU financing system designed to replace the slow, fragmented loan workflows homeowners struggle with. Built as a single Go binary with an embedded SvelteKit frontend, it handles qualification, routing, underwriting prep, document intelligence, and lender handoff with AI-assisted clarity. The goal is simple... help homeowners get the funding they need without the chaos.'
tags: ['Fintech', 'AI', 'GoLang', 'SvelteKit', 'Mortgages']
metrics:
  users: 'Private Beta'
  revenue: 'Bootstrapped'
  team: '2'
links:
  twitter: 'https://x.com/outlet'
  linkedin: 'https://linkedin.com/in/outlet'
---

# **About LendyLink**

LendyLink started as a simple question.
Why is it harder to _finance_ an ADU than it is to _build_ one?

Homeowners bounce between lenders, builders, half-filled PDFs, email threads, and outdated portals that feel allergic to modern software. Everyone wants to help, but nobody sees the full picture. The process drags. Momentum dies. Projects stall.

We built LendyLink to give homeowners a clean path forward.
One place to start.
One flow that removes confusion instead of adding to it.

LendyLink handles the entire front side of ADU financing. It gathers documents, checks income and property conditions, evaluates project type, and routes the loan application to the right lender with the right underwriting package. All of it runs inside a single Go binary with a SvelteKit frontend baked in at build time.

Behind the scenes, AI ties everything together. It checks for missing data, flags inconsistencies, guides homeowners through confusing parts, and helps lenders get a complete and clean file the first time. The result is a financing experience that feels modern instead of messy.

## **My Role**

I co-founded the company and built the entire platform architecture.
Every part of LendyLink — the Go backend, the SvelteKit frontend, the loan routing logic, the AI-assisted qualification flow, the embedded UI, the lender integrations, the deployment model — all of it came together the same way I build everything else. Clear code, predictable patterns, and a focus on removing friction for the people who actually have to use the system.

My job was simple.
Create an ADU financing platform that behaves like software from this decade instead of the last one.

## **Technology Stack**

- Go (go-zero framework)
- SvelteKit with adapter-static
- PostgreSQL 16
- SQLC for database models
- NATS for async coordination
- Embedded frontend through Go’s `embed`
- AI-assisted data extraction and routing
- Auto-HTTPS with autocert
- Single-binary deployment model

## **Key Achievements**

- Built a complete ADU financing flow as a single binary (API, UI, routing logic, auth, everything).
- Designed AI-assisted qualification that reduces missing data and speeds up lender reviews.
- Created a clean, step-by-step homeowner flow that dramatically improves completion rates.
- Integrated lender APIs and webhook systems for real-time status tracking.
- Shipped a secure, production-ready platform with zero external dependencies besides PostgreSQL.
- Reduced lender review time by delivering clean, consistent underwriting packages.

## **Links and Resources**

- [Product Website](https://lendylink.com)
- Documentation (internal)
