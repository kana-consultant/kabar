<div align="center">


<br/>


### Kana Automation Builder for Article executoR

**AI-Powered Content Automation Platform**

<br/>

<a href="#overview">Overview</a> • <a href="#features">Features</a> • <a href="#architecture">Architecture</a> • <a href="#pipeline-flow">Pipeline Flow</a> • <a href="#getting-started">Getting Started</a> • <a href="#tech-stack">Tech Stack</a> • <a href="#security">Security</a> • <a href="#license">License</a>

<br/><br/>

![Go](https://img.shields.io/badge/Go-1.21-blue?logo=go)
![React](https://img.shields.io/badge/React-19-blue?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue?logo=typescript)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue?logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-7-red?logo=redis)
![Docker](https://img.shields.io/badge/Docker-ready-blue?logo=docker)

</div>

---

## Overview

Kana Automation Builder for Article executoR (KABAR) is a fully autonomous AI-powered platform designed to automate the entire SEO content lifecycle — from generation to multi-platform publishing.

KABAR eliminates repetitive manual workflows by providing a centralized system for:

* AI content generation
* Draft management
* Multi-platform publishing
* Scheduled automation
* Activity tracking and analytics

---

## Features

### Content Generation

* AI-powered article generation (multi-model support)
* AI image generation integration
* SEO-optimized structured output
* Ready-to-publish HTML content

### Publishing

* One-click multi-platform publishing
* WordPress, Shopify, and Custom API support
* Scheduled publishing with timezone support
* Batch publishing

### Management

* Role-based access control (RBAC)
* Draft and content lifecycle management
* Product/platform configuration
* Full publishing history tracking

### Automation

* Cron-based scheduler with Redis
* Background job processing
* Extensible automation workflows

---

## Architecture

```text
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

---

## Pipeline Flow

### Content Generation Flow

```text
User Input
   ↓
AI Generation
   ↓
SEO Optimization
   ↓
Draft Storage
   ↓
Review
   ↓
Publish
   ↓
History Tracking
```

### Scheduled Publishing Flow

```text
Scheduler
   ↓
Scan Drafts
   ↓
Validate Time
   ↓
Auto Publish
   ↓
Update Status
   ↓
Save History
```

### Multi-Platform Publishing Flow

```text
Draft Selected
   ↓
Select Products
   ↓
Transform Content
   ↓
API Dispatch
   ↓
Collect Result
   ↓
Save History
```

---

## Getting Started

### Clone Repository

```bash
git clone https://github.com/yourname/kabar.git
cd kabar
```

### Install Dependencies

```bash
cd backend
go mod download

cd ../frontend
pnpm install
```

### Run Infrastructure

```bash
docker-compose up -d
```

### Run Application

```bash
cd backend
go run cmd/api/main.go

cd frontend
pnpm dev
```

---

## Tech Stack

| Layer            | Technology               |
| ---------------- | ------------------------ |
| Backend          | Go (Chi Router)          |
| Frontend         | React + TanStack Router  |
| UI / Styling     | Tailwind CSS + shadcn/ui |
| Database         | PostgreSQL               |
| Cache / Queue    | Redis                    |
| Authentication   | JWT + bcrypt             |


---

## Security

* bcrypt password hashing
* AES-256 encryption for API keys
* JWT authentication
* Input validation & sanitization
* SQL injection protection

---

## License

MIT License

---

<div align="center">

Kana Automation Builder for Article executoR (KABAR)

**1 dashboard. N products. 1 click. Done.**

</div>
