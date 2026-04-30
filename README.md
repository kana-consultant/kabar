# Kana Automation Builder for Article executoR (KABAR)

**1 dashboard. N products. 1 click. Done.**

KABAR is an open-source automation dashboard designed to manage, generate, and distribute SEO content (articles and images) across multiple platforms from a single unified system.

---

## What is KABAR?

Kana Automation Builder for Article executoR (KABAR) is a fully autonomous AI-powered content automation platform that streamlines the entire SEO content lifecycle — from generation to multi-platform publishing.

Instead of manually writing articles and publishing them one by one to different platforms, KABAR centralizes and automates the entire workflow:

1. Generate — AI produces SEO-optimized articles and images
2. Save — Store content as drafts for editing and review
3. Publish — Instantly publish to multiple platforms
4. Schedule — Plan and automate future publishing
5. Track — Monitor publishing history and activity

KABAR is built for scalability, enabling teams and individuals to manage high-volume content production efficiently.

---

## Features

### Content Generation

* AI-powered article generation using multiple models
* AI image generation support
* SEO optimization with structured output
* Ready-to-publish HTML content

### Publishing

* One-click multi-platform publishing
* Integration with WordPress, Shopify, and custom APIs
* Scheduled publishing with timezone support
* Batch publishing for multiple drafts

### Management

* Multi-user system with role-based access control
* Draft management system
* Product/platform configuration
* Full publishing history tracking

### Automation

* Background job processing with scheduler
* Automated publishing workflows
* Extensible system for future automation features

---

## Tech Stack

| Layer            | Technology                         |
| ---------------- | ---------------------------------- |
| Backend          | Go (Chi Router)                    |
| Frontend         | React + TanStack Router + Tailwind |
| Database         | PostgreSQL                         |
| Cache / Queue    | Redis                              |
| Authentication   | JWT + bcrypt                       |


---

## Architecture Overview

KABAR is built using a clean and modular architecture approach, separating concerns into multiple layers:

```
Frontend (React)
│
▼
Backend API (Go)
│
┌──────┴──────────────┐
│ Application Layer   │
│ - Draft             │
│ - Generate          │
│ - Product           │
│ - Team              │
│ - User              │
│ - History           │
│ - Scheduler         │
└──────┬──────────────┘
▼
Domain Layer (Entities + Interfaces)
▼
Infrastructure Layer
(PostgreSQL, Redis, AI Services)
```

This architecture ensures scalability, maintainability, and clear separation between business logic and infrastructure.

---

## Pipeline Flow

### Content Generation Flow

```
[ User Input ]
       ↓
[ AI Generation Engine ]
       ↓
[ SEO Optimization & Structuring ]
       ↓
[ Draft Storage ]
       ↓
[ Manual Review / Editing ]
       ↓
[ Publish Trigger ]
       ↓
[ Distribution Engine ]
       ↓
[ History & Tracking ]
```

---

### Scheduled Publishing Flow

```
[ Scheduler (Cron + Redis) ]
               ↓
[ Scan Scheduled Drafts ]
               ↓
[ Time Validation (scheduled_for <= now) ]
               ↓
[ Auto Publish Execution ]
               ↓
[ Status Update (published) ]
               ↓
[ History Logging ]
```

---

### Multi-Platform Publishing Flow

```
[ Draft Selected ]
        ↓
[ Product Selection ]
        ↓
[ Adapter Layer (Transform Content) ]
        ↓
[ API Dispatch (Parallel Requests) ]
        ↓
[ Response Collection ]
        ↓
[ Success / Failed Mapping ]
        ↓
[ History Recording ]
```


---

## Requirements

* Go 1.21+
* Node.js 20+
* pnpm
* PostgreSQL
* Redis
* Docker (optional)

---

## Getting Started

### Clone Repository

```bash
git clone https://github.com/yourname/kabar.git
cd kabar
```

### Install Dependencies

```bash
cd src/app/api/go mod download
cd ../frontend && pnpm install
```

### Run Infrastructure

```bash
docker-compose up -d
```

### Setup Environment

```bash
cp src/app/api/.env.example backend/.env
cp .env.example frontend/.env
```

### Run Application

```bash
cd src/app/api/ go run cmd/api/main.go
cd frontend && pnpm dev
```

---

## API Overview

### Public Endpoints

* POST /api/auth/register
* POST /api/auth/login
* GET /health

### Protected Endpoints

* Draft management
* Content generation
* Product integration
* Publishing system
* History tracking

---

## Product Integration

KABAR supports multiple platform integrations:

### WordPress

* REST API integration
* Application password authentication

### Shopify

* Store API integration
* Access token authentication

### Custom API

* Fully configurable endpoint
* Field mapping support

---

## Security

* Password hashing using bcrypt
* API key encryption using AES-256
* JWT-based authentication
* Input validation and sanitization
* Protection against SQL injection

---

## Role-Based Access Control

| Role        | Description         |
| ----------- | ------------------- |
| super_admin | Full system access  |
| admin       | Team-level access   |
| manager     | Manage team members |
| viewer      | Read-only access    |

---

## Development

```bash
go test ./...
pnpm test
```

---

## Docker

```bash
docker-compose up -d
```

---

## Documentation

```
http://localhost:8080/swagger/index.html
```

---

## Contributing

* Open an issue
* Discuss changes
* Submit pull request

---

## License

MIT License

---

## Credits

* Built with Go and React
* Powered by OpenRouter
* Infrastructure using PostgreSQL and Redis

---

## Contact

* GitHub Issues
* Community channel (Discord, etc.)

---

Kana Automation Builder for Article executoR (KABAR)

1 dashboard. N products. 1 click. Done.
