---
title: 'DOSync'
slug: 'dosync'
logo: '/images/projects/dosync.svg'
website: 'https://github.com/localrivet/dosync'
founded: '2024'
role: 'Creator'
status: 'Active'
category: 'Developer Tools'
featured: true
order: 3
tagline: 'Production-grade Docker Compose orchestration for teams that choose simplicity over Kubernetes complexity.'
description: 'DOSync is a zero-downtime deployment tool for Docker Compose that automates rolling updates across single servers or multi-server fleets. Built for teams that value operational simplicity—it provides Kubernetes-level reliability without the operational burden.'
tags: ['Docker', 'DevOps', 'GoLang', 'Open Source', 'Infrastructure']
metrics:
  users: '3 GitHub Stars'
  revenue: 'Open Source'
  team: '2'
links:
  github: 'https://github.com/localrivet/dosync'
  twitter: 'https://x.com/outlet'
  linkedin: 'https://linkedin.com/in/outlet'
---

# About DOSync

DOSync is a production-grade Docker Compose orchestration tool that automates deployments across single servers or multi-server fleets. It synchronizes services with container registries, performs zero-downtime rolling updates, and provides automatic rollback on failures.

Perfect for teams that choose Docker Compose over Kubernetes for operational simplicity.

## The Problem It Solves

Kubernetes is overkill for most deployments. You don't need that complexity when your application runs fine on 5-50 servers, your team knows Docker Compose (not k8s manifests), and you value operational simplicity over theoretical scale.

DOSync gives you Kubernetes-level deployment reliability with Docker Compose simplicity. Run multiple servers with identical compose files, each with its own DOSync instance behind a load balancer. Get zero-downtime deployments, automatic rollback, and horizontal scaling without the infrastructure costs or operational burden.

## My Role

I co-built DOSync with Claude Code. I designed the architecture, wrote the core orchestration logic, implemented the registry integrations, and created the sophisticated replica detection system. I focused on making deployments feel effortless while maintaining production-grade reliability.

The goal was to prove that you don't need Kubernetes to run a reliable, scalable application—you just need better tooling around Docker Compose.

## Technology Stack

- **Go** (core orchestration engine)
- **Docker & Docker Compose** (container runtime)
- **Multi-registry support** (Docker Hub, GHCR, GCR, ACR, ECR, Harbor, Quay, DOCR)
- **SQLite** (metrics storage)
- **Systemd** (service management)

## Key Features

### Production-Grade Deployments

- Zero-downtime rolling updates with health checks
- Multi-server support (run independent instances behind a load balancer)
- Intelligent replica detection for scale-based and name-based replicas
- Service dependency management (updates in correct order)
- Multiple deployment strategies: one-at-a-time, percentage, blue-green, canary
- Automatic rollback on health check failures with full history

### Registry & Version Control

- Support for all major container registries
- Advanced tag policies: semantic versioning, numerical ordering, regex filters
- Version constraints (e.g., deploy only `>=1.0.0 <2.0.0`)
- State drift prevention (checks running containers vs compose file)

### Operations & Monitoring

- Self-updating capability
- Automatic backup management before modifications
- SQLite metrics storage
- Slack/email/webhook notifications
- Web dashboard for monitoring deployments
- Docker Compose as source of truth

## Architecture Highlights

### Replica Detection System

DOSync includes a sophisticated replica detection system that automatically identifies and manages different types of service replicas:

- **Scale-based replicas**: Handles Docker Compose `scale` and `deploy.replicas` configurations
- **Name-based replicas**: Detects blue-green deployments and custom naming patterns
- **Consistent updates**: Ensures all replicas of a service are updated correctly
- **Zero-downtime**: Rolling updates across all replicas

### Multi-Server Fleet Management

Each server runs:

- Identical `docker-compose.yml` file
- Its own DOSync instance (monitors local Docker daemon)
- Multiple replicas of each service
- Load balancer routes traffic based on health checks

**Benefits:**

- High availability (servers fail independently)
- Horizontal scaling (add servers as traffic grows)
- No single point of failure
- Standard Docker Compose (familiar tooling)
- Much lower costs than Kubernetes

### Image Policy Engine

Sophisticated tag selection policies:

```yaml
# Example: Select highest timestamp from branch-based tags
imagePolicy:
  filterTags:
    pattern: '^main-[a-fA-F0-9]+-(?P<ts>\d+)$'
    extract: 'ts'
  policy:
    numerical:
      order: desc

# Example: Semantic versioning with constraints
imagePolicy:
  policy:
    semver:
      range: '>=1.0.0 <2.0.0'
```

## Key Achievements

- Built a production-ready orchestration system that rivals Kubernetes reliability with 1/10th the complexity
- Designed intelligent replica detection that handles both Docker Compose scaling and custom naming patterns
- Implemented multi-registry support with advanced tag policies (semver, numerical, alphabetical)
- Created automatic rollback system with health check validation
- Shipped with zero external dependencies (runs as single binary)
- Achieved zero-downtime deployments across multi-server fleets without coordination mechanisms

## Use Cases

### Single Server

Perfect for:

- Side projects and MVPs
- Internal tools
- Development/staging environments
- Small businesses ($50-500k ARR)

### Multi-Server Fleet

Ideal for:

- Growing SaaS applications ($500k-5M ARR)
- High-availability web applications
- Agencies managing multiple client sites
- Teams that value "boring technology"
- 100-10,000 requests/second

## When to Use Kubernetes Instead

DOSync is designed for simplicity and operational efficiency. Use Kubernetes when you need:

- 100+ servers
- Complex multi-region deployments
- Service mesh, auto-scaling pods across nodes
- Enterprise compliance requirements

## Links & Resources

- [GitHub Repository](https://github.com/localrivet/dosync)
- [Documentation](https://github.com/localrivet/dosync#readme)
- [Troubleshooting Guide](https://github.com/localrivet/dosync/blob/main/docs/troubleshooting.md)
- [Example Configurations](https://github.com/localrivet/dosync/tree/main/examples)
