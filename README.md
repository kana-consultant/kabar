# SEO Multi-Post Dashboard

**1 dashboard. N products. 1 click. Done.**

Open-source dashboard to manage SEO content (articles + images) from one place and publish it to multiple products/platforms at once.

---

## Features

- ✅ AI Article Generation (Gemini, GPT-4o, Claude, Llama, etc.)
- ✅ AI Image Generation (Nano Banana, Imagen 3.0)
- ✅ One-click publish to multiple products (WordPress, Shopify, Custom API)
- ✅ Adapter pattern for multi-platform integrations
- ✅ Multi-user team management
- ✅ Scheduled publishing
- ✅ Full publishing history tracking
- ✅ Encrypted API key storage

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go + Chi Router |
| Frontend | React + TanStack Router + Tailwind |
| Database | PostgreSQL |
| Auth | JWT + bcrypt |
| AI Providers | OpenRouter, Google Gemini |
| API Docs | Swagger |

---

# Requirements

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Redis (for queue processing)

---

# Installation

## 1. Clone Repository

```bash
git clone https://github.com/yourname/seo-dashboard.git
cd seo-dashboard
```

---

## 2. Backend Setup (Go)

```bash
cd backend
cp .env.example .env

# Edit .env with your database and JWT configuration

go mod download
go run cmd/main.go
```

Backend runs at:

```bash
http://localhost:8080
```

---

## 3. Frontend Setup (React)

```bash
cd frontend
cp .env.example .env

# Edit .env with backend API URL

npm install
npm run dev
```

Frontend runs at:

```bash
http://localhost:5173
```

---

## 4. Database Migration

```bash
psql -U postgres -d seo_db < database/schema.sql
```

---

# Environment Variables

## Backend `.env`

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=seo_db

JWT_SECRET=your_jwt_secret_key_min_32_char
ENCRYPTION_KEY=your_aes_encryption_key_32_char
```

---

## Frontend `.env`

```env
VITE_API_URL=http://localhost:8080
```

---

# API Endpoints

## Public Routes

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/auth/register | Register user |
| POST | /api/auth/login | Login user |
| POST | /api/auth/forgot-password | Forgot password |

---

## Protected Routes (Require JWT)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/dashboard/stats | Dashboard statistics |
| GET | /api/products | Get products |
| POST | /api/products | Add product |
| POST | /api/products/{id}/test | Test product API |
| GET | /api/drafts | Get drafts |
| POST | /api/drafts | Create draft |
| POST | /api/drafts/{id}/publish | Publish draft |
| POST | /api/drafts/publish | Batch publish |
| POST | /api/drafts/schedule | Schedule post |
| POST | /api/generate/article | Generate article |
| POST | /api/generate/image | Generate image |
| GET | /api/history | Publishing history |
| GET | /api/models | AI model list |
| CRUD | /api/teams | Team management |
| CRUD | /api/api-keys | API key management |

---

## Utility Routes

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | Health check |
| GET | /metrics | System metrics |
| GET | /swagger/* | Swagger documentation |

---

# Adding a New Product

## WordPress

Go to:

```text
Products → Add Product
```

Select:

```text
WordPress
```

Enter:

```text
API URL:
https://yourdomain.com/wp-json/wp/v2/posts
```

Provide:

- Application Password

Then:

```text
Test Connection → Save
```

---

## Shopify

Select:

```text
Shopify
```

Enter:

```text
Store URL:
https://store-name.myshopify.com
```

Provide:

- Access Token
- Blog ID

Then:

```text
Test Connection → Save
```

---

## Custom API

Select:

```text
Custom API
```

Configure:

- API Endpoint
- HTTP Method (POST / PUT)
- Field Mapping (drag & drop)

Then:

```text
Test Connection → Save
```

---

# Custom API Field Mapping

| Internal Field | API Field |
|----------------|-----------|
| title | post_title |
| body | content_html |
| imageUrl | thumbnail_url |
| status | is_published |

Example JSON Mapping:

```json
{
  "title": "post_title",
  "body": "content",
  "imageUrl": "featured_image",
  "status": "publish_status"
}
```

---

# Workflow

## Step 1 — Generate Content

- Open Generate menu
- Select AI model (Gemini, GPT-4o, etc.)
- Enter topic or keyword
- Click Generate Article
- Optionally generate image
- Preview result

---

## Step 2 — Save as Draft

- Edit content if needed
- Click Save to Drafts
- Select target products

---

## Step 3 — Publish

## Option A — Publish Now

- Open Drafts
- Select draft
- Click Publish Now

---

## Option B — Schedule

- Open Schedule
- Select draft and time
- Click Schedule

---

## Option C — Batch Publish

- Select multiple drafts
- Click Publish Selected
- Choose target products

```text
One click publishes everywhere 🚀
```

---

# Security

| Feature | Implementation |
|---------|----------------|
| User Passwords | bcrypt hashing |
| API Keys | AES-256 encryption |
| JWT | Access + Refresh tokens |
| CORS | Whitelisted domains |
| SQL Protection | Prepared statements |

---

# Database Structure

```text
users           - Users and authentication
teams           - Teams
team_members    - Team members

api_providers   - OpenRouter, Gemini, etc.
ai_models       - AI model catalog
api_keys        - Encrypted user keys

products        - WordPress / Shopify / Custom
adapter_configs - Custom API mappings

drafts          - Unpublished content
histories       - Publish history
```

---

# Development

## Add a New AI Model

Insert into:

```text
api_providers
ai_models
```

It will automatically appear in frontend.

---

## Add a New Platform

1. Add new enum in:

```text
products.platform
```

2. Create adapter in:

```text
internal/adapters/
```

3. Update:

```text
adapter_configs schema
```

---

# License

MIT License

---

# Contributing

Pull Requests are welcome.

For major changes:

1. Open an issue first  
2. Discuss the proposal  
3. Submit PR

---

# Contact

- Issues: GitHub Issues  
- Discord: Join Server

---

## Built with using Go + React