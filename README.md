# Kabar — Multi-Platform AI Content Management System

**Kabar** is a centralized AI-powered content management platform designed to simplify the process of creating, managing, scheduling, and publishing content across multiple platforms — all from a single dashboard.

---

## 🚀 Overview

Managing content across multiple platforms manually is time-consuming, inefficient, and error-prone. Kabar serves as a unified solution that brings the entire content workflow into one integrated system.

**Without Kabar, content teams must:**
- Log into each product dashboard one by one
- Create articles and images using separate, disconnected tools
- Publish content to each platform individually
- Manage SEO manually without centralized visibility
- Track content status across scattered systems

**With Kabar, everything happens in one place:**
- Generate high-quality articles with AI
- Generate supporting images with AI
- Publish to multiple platforms simultaneously in one click
- Schedule content for the most strategic publish times
- Optimize SEO automatically
- Generate sitemaps dynamically
- Manage meta tags across all platforms
- Collaborate as a team in a shared workspace

---

## 🎯 Problem Statement

| No | Problem | Impact |
|----|---------|--------|
| 1 | Must log into multiple product dashboards separately | Wastes team time and energy |
| 2 | Article and image creation happen in different tools | Fragmented, inefficient workflow |
| 3 | No centralized content management system | Hard to track status and maintain consistency |
| 4 | Keyword research is not integrated into the publishing flow | SEO opportunities are frequently missed |

---

## 💡 Solution

> *"From a single dashboard, users can generate articles and images, then publish them to multiple products simultaneously with just one click."*

Kabar unifies the entire content management ecosystem into one platform:

- **AI Content Generation** — Produce high-quality content with AI assistance
- **SEO Optimization** — Make content discoverable on search engines automatically
- **Multi-platform Publishing** — Deliver content to all platforms in a single action
- **Scheduling** — Set the most strategic publish times in advance
- **Analytics** — Monitor content performance from a central view
- **Team Collaboration** — Work together in one organized workspace

---

## 🏗️ System Architecture

### Layered Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ Presentation Layer                                          │
│ React 19 + TanStack Router + shadcn/ui                      │
├─────────────────────────────────────────────────────────────┤
│ API Layer                                                   │
│ Chi Router                                                  │
├─────────────────────────────────────────────────────────────┤
│ Application Layer                                           │
│ Adapter │ Generate │ Draft │ Product │ Auth │ Team          │
├─────────────────────────────────────────────────────────────┤
│ Domain Layer                                                │
│ Entities + Business Rules + Repository Interfaces           │
├─────────────────────────────────────────────────────────────┤
│ Infrastructure Layer                                        │
│ PostgreSQL │ Redis │ SMTP │ AI APIs │ Docker                │
└─────────────────────────────────────────────────────────────┘
```

Kabar is built with a clean layered architecture that separates each component's responsibility clearly. The Presentation Layer handles the user interface, the API Layer serves as the entry point for frontend requests, the Application Layer orchestrates business logic, the Domain Layer defines core entities and rules, and the Infrastructure Layer manages connections to the database, cache, email, and external services.

### 🛠️ Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | React 19, Vite, TanStack Router & Query |
| UI | Tailwind CSS, shadcn/ui |
| Backend | Go 1.25, Chi Router |
| Database | PostgreSQL / SQLite |
| Cache & Queue | Redis |
| Migration | golang-migrate |
| Auth | Session-based / JWT |
| AI Integration | AI API (pluggable) |
| Image Generation | Image Generation API |
| Infrastructure | Docker & Docker Compose |

---

## ✨ Core Features

### 1. 🤖 AI Content Generation — *Produce Publish-Ready Content with AI*

This feature allows users to create a complete, structured article simply by providing a topic or keyword. The AI automatically composes content that is informative, well-structured, and SEO-friendly — ready to publish without heavy manual editing.

**What you can do:**

- **Generate SEO-friendly articles** — The AI produces articles with a proper heading hierarchy, natural keyword density, and relevant content based on the topic provided. Every article is structured to perform well on search engines from day one.
- **Generate supporting images** — The AI generates contextually appropriate images that match the article content. No need to search stock photo sites or create visuals in a separate design tool.
- **Preview before publishing** — Users can review the full content layout before deciding to publish, ensuring quality and accuracy meet expectations before it goes live.
- **Draft management** — Unfinished content can be saved as a draft and resumed at any time. The writing process does not need to be completed in a single session, giving teams full flexibility.
- **AI model selection** — Users can select which AI model to use for content generation, balancing quality and speed depending on the use case or volume requirements.
- **Quick generate mode** — A fast-track mode for generating straightforward content without additional configuration, ideal for high-volume content production workflows.

---

### 2. 🌐 Multi-Platform Publishing — *Publish to Every Platform at Once*

Kabar is built on an **Adapter Pattern** architecture that allows a single piece of content to be published to multiple platforms simultaneously, with each platform receiving it in the correct format automatically.

**Why this matters:**
Every platform has its own API structure, payload format, and authentication method. Without Kabar, teams must manually reformat content for each platform — a repetitive, time-consuming, and error-prone process. Kabar's Adapter Layer handles all of these transformations automatically behind the scenes.

**Each adapter is capable of:**
- Targeting a different API endpoint per platform
- Reshaping the payload structure to match each platform's requirements
- Handling platform-specific authentication independently
- Applying the correct publishing strategy (immediate publish, submit as draft, etc.)

**Publishing flow:**
```
Generate Content ──► Select Target Platforms ──► Adapter Transformation ──► Deliver to All APIs
```

---

### 3. Redis-Based Scheduling — *Automate Publishing at the Right Time*

Kabar's scheduling system uses Redis as the backbone of its job queue, ensuring content is published exactly on time — even when users are offline or the server restarts.

**Why Redis:**
Redis is an in-memory data store that is extremely fast and reliable for queue management. With Redis, Kabar can handle thousands of scheduled publish jobs concurrently without overloading the primary database.

**How it works:**
1. The user sets a future publish date and time for the content
2. Kabar stores the publish job in the Redis queue
3. Background workers continuously monitor the queue and execute each job precisely at the scheduled time
4. If a publish job fails (e.g. due to a network disruption), the system automatically retries with an exponential backoff strategy

**Scheduling capabilities:**
- **Schedule future publishing** — Set a specific date and time in the future for content to go live automatically
- **Reschedule** — Change the publish time of already-scheduled content at any point before it executes
- **Cancel scheduled posts** — Remove a scheduled job before it runs if plans change
- **Retry failed jobs** — The system automatically retries failed publish attempts, reducing the risk of missed publications
- **Queue monitoring** — View the real-time status of the queue, including running, pending, and failed jobs

---

### 4. SEO Optimization — *Every Piece of Content, Search-Engine Ready*

Kabar integrates SEO optimization directly into the content creation and publishing workflow. By the time content is published, it is already optimized for search engines — no separate SEO tool required.

**SEO analysis and optimization features:**

- **SEO score analysis** — Every piece of content receives an SEO score based on factors such as keyword density, content length, heading structure, and readability. Users know exactly how well-optimized their content is before it goes live.
- **Meta title generation** — The AI automatically generates a compelling meta title that includes the primary keyword, optimized to the ideal character length for search engine display.
- **Meta description generation** — The AI produces a persuasive and informative meta description designed to increase click-through rates from search results pages.
- **Keyword integration** — Ensures the target keyword is woven naturally throughout the content body, headings, and meta tags without keyword stuffing.
- **Canonical URL support** — Prevents duplicate content issues by setting the correct canonical URL for each published page, protecting search ranking signals.

**Supported meta tags:**

- **Open Graph tags** — Controls how content appears when shared on Facebook, LinkedIn, and other social platforms, ensuring rich previews with the correct title, description, and image.
- **Twitter Cards** — Optimizes content appearance when shared on Twitter/X, enabling rich card formats that increase engagement.
- **JSON-LD** — Adds structured data markup that helps search engines understand the deeper context of the content, improving eligibility for rich results.
- **Structured Data** — Supports rich snippets in search results including article schema, breadcrumb navigation, and FAQ markup.
- **Dynamic meta generation** — Meta tags are generated dynamically per page based on its actual content, rather than using a single static template applied across all pages.

---

### 5. Sitemap System — *Always Up-to-Date, Automatically*

Kabar automatically manages XML sitemaps for all published content. A sitemap tells search engines about every page that exists on a website, enabling faster and more complete indexing of new content.

**Why this matters:**
Without a regularly updated sitemap, search engines may take days or weeks to discover newly published content. Kabar handles sitemap generation entirely automatically, so teams never need to think about it.

**What is generated automatically:**
- A primary `sitemap.xml` that aggregates all content across the system
- Per-product sitemaps for better organization and crawl efficiency
- Dynamic sitemap updates triggered every time new content is published

**Sitemap features:**

- **Auto regenerate after publish** — The sitemap is automatically refreshed every time content is published, with zero manual intervention required
- **SEO-friendly URL structure** — Content URLs follow SEO best practices to maximize ranking potential and crawlability
- **Search engine ready** — Sitemap format fully complies with the standards accepted by Google, Bing, and all major search engines
- **Multi-product support** — Manages sitemaps for multiple products and websites simultaneously within a single system

---

### 6. API Key Encryption — *Enterprise-Grade Credential Security*

Every platform connected to Kabar requires an API key or authentication credential. Kabar ensures all credentials are stored and used with the highest security standards, protecting integrations from unauthorized access.

**Encryption flow:**
```
User inputs API key ──► AES-256 Encryption ──► Store ciphertext in DB ──► Decrypt only during publishing
```

**How the security flow works:**
1. When a user enters an API key for a platform, Kabar immediately encrypts it using AES-256 encryption before writing anything to the database
2. What is stored in the database is ciphertext — completely unreadable without the encryption key
3. The original API key is only decrypted temporarily in server memory at the moment a publish job executes
4. After publishing completes, the plain-text API key is never persisted anywhere in the system

**System-wide security implementation:**

- **AES-256 encryption** — The same encryption standard used by financial institutions and military organizations to protect sensitive data at rest
- **JWT authentication** — Every request from the frontend is verified using a JSON Web Token with a limited expiry window, preventing unauthorized access
- **Session-based access** — User sessions are managed securely and can be invalidated instantly from the server side when needed
- **Protected API routes** — All API endpoints are guarded by authentication middleware, blocking any unauthenticated or unauthorized access attempts

---

##  Complete Data Flow

```
User Inputs Topic / Keyword
              │
              ▼
       Keyword Research
              │
              ▼
     Generate Article via AI
              │
              ▼
     Generate Image via AI
              │
              ▼
       Preview & Edit Content
              │
              ▼
    Select Target Platforms
              │
              ▼
          Click Publish
              │
              ▼
   Adapter Layer Transformation
  (payload reformatted per platform)
              │
              ▼
    Send to Each Platform's API
              │
              ▼
       Update Content Status
              │
              ▼
  Generate Sitemap & Meta Tags
              │
              ▼
           Done ✓
```

---

##  Project Structure

```
apps/
├── api/                          # Go backend
│   ├── internal/
│   │   ├── application/          # Application logic & use cases
│   │   ├── domain/               # Core entities & business rules
│   │   ├── infrastructure/       # DB, Redis, SMTP, AI connections
│   │   ├── presentation/         # HTTP handlers & routing
│   │   ├── middleware/           # Auth, logging, rate limiting
│   │   ├── helper/               # Reusable helper functions
│   │   ├── config/               # Application configuration
│   │   ├── container/            # Dependency injection
│   │   ├── database/             # Migrations & database connection
│   │   ├── pkg/                  # Reusable internal packages
│   │   └── utils/                # General utilities
│   └── main.go
│
├── web/                          # React frontend
│   ├── src/
│   │   ├── components/           # Reusable UI components
│   │   ├── pages/                # Application pages
│   │   ├── hooks/                # Custom React hooks
│   │   ├── services/             # API communication layer
│   │   ├── routes/               # Routing configuration
│   │   ├── contexts/             # React context & global state
│   │   ├── utils/                # Frontend utility functions
│   │   └── types/                # TypeScript type definitions
│   └── package.json
│
deployments/
├── docker-compose.dev.yml        # Docker config for development
├── docker-compose.prod.yml       # Docker config for production
└── .env.production.example       # Production environment variable template
```

---

## ⚙️ Environment Configuration

### Backend `.env`

```ini
# =========================
# APP CONFIG
# =========================
APP_BASE_URL=http://localhost:5173

# =========================
# SERVER CONFIG
# =========================
SERVER_PORT=8080

# =========================
# DATABASE CONFIG
# =========================
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_database_name
DB_SSL_MODE=disable

# =========================
# REDIS CONFIG
# =========================
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# =========================
# SECURITY
# =========================
ENCRYPTION_KEY=your_32_character_key
JWT_SECRET=your_jwt_secret
JWT_EXPIRY_HOURS=24

# =========================
# AI CONFIG
# =========================
AI_API_KEY=your_ai_api_key
IMAGE_GEN_API_KEY=your_image_api_key

# =========================
# SMTP CONFIG
# =========================
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTP_FROM_EMAIL=noreply@yourapp.com
SMTP_FROM_NAME=Kabar

SMTP_SECURE=false
SMTP_TLS=true
```

---

## 🐳 Docker Development

**Run the development stack:**
```bash
docker compose -f deployments/docker-compose.dev.yml up -d
```

**Stop all containers:**
```bash
docker compose -f deployments/docker-compose.dev.yml down
```

**Rebuild containers after changes:**
```bash
docker compose -f deployments/docker-compose.dev.yml up -d --build
```

---

## 🔧 Local Development

### Backend
```bash
cd apps/api
cp .env.example .env
go mod tidy
make migrate-up
go run main.go
```

### Frontend
```bash
cd apps/web
npm install
npm run dev
```

---

## 📖 API Documentation & Testing

Swagger documentation is available at:
```
http://localhost:8080/swagger/index.html
```

### Testing

**Backend:**
```bash
cd apps/api
make test
```

**Frontend:**
```bash
cd apps/web
npm run test
```

---

