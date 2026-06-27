-- ============================================================
-- FULL SCHEMA + SEED DATA: kabar.com
-- Idempotent: aman dijalankan berulang kali tanpa duplikat data
-- ============================================================
-- 1. EXTENSIONS
-- ============================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS "pgcrypto"  WITH SCHEMA public; 

-- ============================================================
-- 2. ENUM TYPES
-- ============================================================
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'permission_scope') THEN
        CREATE TYPE permission_scope AS ENUM ('global', 'team');
    END IF;
END $$;


-- ============================================================
-- 3. FUNCTION: auto-update kolom updated_at
-- ============================================================
CREATE OR REPLACE FUNCTION public.set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ============================================================
-- 4. CORE TABLES
-- ============================================================

-- 4.1 Users
CREATE TABLE IF NOT EXISTS public.users (
    id            uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    email         varchar(255) NOT NULL UNIQUE,
    name          varchar(255) NOT NULL,
    password_hash varchar(255) NOT NULL,
    role          varchar(50)  DEFAULT 'viewer',
    avatar        text,
    status        varchar(20)  DEFAULT 'active',
    last_active   timestamp,
    created_at    timestamp    DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_users_role   CHECK (role   IN ('superadmin', 'admin', 'manager', 'editor', 'viewer')),
    CONSTRAINT chk_users_status CHECK (status IN ('active', 'inactive', 'suspended'))
);

-- 4.2 Teams
CREATE TABLE IF NOT EXISTS public.teams (
    id          uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name        varchar(255) NOT NULL,
    description text,
    created_by  varchar(255),
    created_at  timestamp    DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp    DEFAULT CURRENT_TIMESTAMP
);

-- 4.3 API Providers
CREATE TABLE IF NOT EXISTS public.api_providers (
    id              uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name            varchar(100) NOT NULL,
    display_name    varchar(200) NOT NULL,
    description     text,
    base_url        varchar(500) NOT NULL,
    auth_type       varchar(50)  DEFAULT 'bearer',
    auth_header     varchar(100) DEFAULT 'Authorization',
    auth_prefix     varchar(50)  DEFAULT 'Bearer',
    default_headers jsonb        DEFAULT '{}'::jsonb,
    team_id         uuid         REFERENCES public.teams(id) ON DELETE CASCADE,
    is_active       boolean      DEFAULT true,
    created_at      timestamp    DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp    DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (team_id, name)
);

-- 4.4 Request Schemas
-- Catatan: UNIQUE(provider_id, name) sengaja tidak dipasang (sesuai versi asli),
-- karena id (PK) sudah cukup untuk identitas row dan di-upsert via ON CONFLICT(id).
CREATE TABLE IF NOT EXISTS public.request_schemas (
    id                    uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    provider_id           uuid         NOT NULL REFERENCES public.api_providers(id) ON DELETE CASCADE,
    name                  varchar(100) NOT NULL,
    endpoint_path         varchar(500) NOT NULL,
    max_tokens_key        varchar(100),
    system_role_key       varchar(100),
    response_text_path    varchar(500),
    response_image_path   varchar(500),
    request_template      text,
    supports_temperature  boolean      DEFAULT true,
    supports_streaming    boolean      DEFAULT true,
    created_at            timestamp    DEFAULT CURRENT_TIMESTAMP,
    updated_at            timestamp    DEFAULT CURRENT_TIMESTAMP
);

-- 4.5 Model Families
CREATE TABLE IF NOT EXISTS public.model_families (
    id            uuid           DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    provider_id   uuid           NOT NULL REFERENCES public.api_providers(id) ON DELETE CASCADE,
    schema_id     uuid           NOT NULL REFERENCES public.request_schemas(id) ON DELETE CASCADE,
    name          varchar(100)   NOT NULL,
    display_name  varchar(200)   NOT NULL,
    description   text,
    max_tokens    integer        DEFAULT 1024,
    temperature   numeric(3,2)   DEFAULT 1.0,
    system_prompt text           DEFAULT 'You are a helpful assistant.',
    created_at    timestamp      DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp      DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (provider_id, name)
);

-- 4.6 API Keys
CREATE TABLE IF NOT EXISTS public.api_keys (
    id            uuid        DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    service       varchar(50) NOT NULL,
    provider_id   uuid        NOT NULL REFERENCES public.api_providers(id) ON DELETE CASCADE,
    model_id      uuid        NOT NULL REFERENCES public.model_families(id) ON DELETE CASCADE,
    key_encrypted text        NOT NULL,
    system_prompt text,
    team_id       uuid        REFERENCES public.teams(id) ON DELETE CASCADE,
    is_active     boolean     DEFAULT true,
    created_by    varchar(255),
    created_at    timestamp   DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp   DEFAULT CURRENT_TIMESTAMP
);

-- 4.7 Products
CREATE TABLE IF NOT EXISTS public.products (
    id                uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name              varchar(255) NOT NULL,
    platform          varchar(50)  NOT NULL,
    api_endpoint      text         NOT NULL,
    api_key_encrypted text,
    status            varchar(20)  DEFAULT 'pending',
    sync_status       varchar(20)  DEFAULT 'idle',
    last_sync         timestamp,
    created_by        uuid         REFERENCES public.users(id) ON DELETE SET NULL,
    team_id           uuid         REFERENCES public.teams(id) ON DELETE SET NULL,
    user_id           uuid         REFERENCES public.users(id) ON DELETE SET NULL,
    created_at        timestamp    DEFAULT CURRENT_TIMESTAMP,
    updated_at        timestamp    DEFAULT CURRENT_TIMESTAMP
);

-- 4.8 Adapter Configs
CREATE TABLE IF NOT EXISTS public.adapter_configs (
    id               uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    product_id       uuid         NOT NULL REFERENCES public.products(id) ON DELETE CASCADE,
    endpoint_path    varchar(500) NOT NULL,
    http_method      varchar(10)  DEFAULT 'POST',
    custom_headers   jsonb        DEFAULT '{}'::jsonb,
    meta_config      text,
    sitemap_config   text,
    field_mapping    text         NOT NULL,
    response_mapping jsonb        DEFAULT '{}'::jsonb,
    timeout_seconds  integer      DEFAULT 30,
    retry_count      integer      DEFAULT 3,
    created_at       timestamp    DEFAULT CURRENT_TIMESTAMP,
    updated_at       timestamp    DEFAULT CURRENT_TIMESTAMP
);

-- 4.9 Workflow Definitions
CREATE TABLE IF NOT EXISTS public.workflow_definitions (
    id          uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    product_id  uuid         NOT NULL REFERENCES public.products(id) ON DELETE CASCADE,
    name        varchar(255) NOT NULL,
    created_at  timestamp    DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp    DEFAULT CURRENT_TIMESTAMP
);

-- 4.10 Workflow Nodes
-- Catatan: endpoint_path diperbaiki dari character(30) -> varchar(500),
-- karena character(N) adalah fixed-length (di-pad spasi) dan akan
-- memotong/merusak URL path yang lebih panjang dari 30 karakter.
CREATE TABLE IF NOT EXISTS public.workflow_nodes (
    id                 uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    workflow_id        uuid         NOT NULL REFERENCES public.workflow_definitions(id) ON DELETE CASCADE,
    adapter_config_id  uuid         NOT NULL REFERENCES public.adapter_configs(id)      ON DELETE CASCADE,
    step_order         integer      NOT NULL,
    input_mapping      jsonb        DEFAULT '{}'::jsonb,
    previous_node_ids  jsonb        DEFAULT '[]'::jsonb,
    next_node_ids      jsonb        DEFAULT '[]'::jsonb,
    endpoint_path      varchar(500),
    http_method        varchar(10),
    created_at         timestamp    DEFAULT CURRENT_TIMESTAMP,
    updated_at         timestamp    DEFAULT CURRENT_TIMESTAMP
);

-- 4.11 Drafts
CREATE TABLE IF NOT EXISTS public.drafts (
    id              uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    title           varchar(500) NOT NULL,
    topic           varchar(500) NOT NULL,
    article         text         NOT NULL,
    image_url       text,
    image_prompt    text,
    status          varchar(20)  DEFAULT 'draft',
    scheduled_for   timestamp,
    target_products jsonb        DEFAULT '[]'::jsonb,
    has_image       boolean      DEFAULT false,
    excerpt         text,
    slug            text,
    seo_score       int,
    created_by      uuid         REFERENCES public.users(id) ON DELETE SET NULL,
    team_id         uuid         REFERENCES public.teams(id) ON DELETE SET NULL,
    user_id         uuid         REFERENCES public.users(id) ON DELETE SET NULL,
    published_at    timestamp,
    created_at      timestamp    DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp    DEFAULT CURRENT_TIMESTAMP
);

-- 4.12 Histories
CREATE TABLE IF NOT EXISTS public.histories (
    id              uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    title           varchar(500) NOT NULL,
    topic           varchar(500) NOT NULL,
    content         text         NOT NULL,
    image_url       text,
    target_products jsonb        DEFAULT '[]'::jsonb,
    status          varchar(20)  NOT NULL,
    action          varchar(20)  NOT NULL,
    error_message   text,
    published_at    timestamp    NOT NULL,
    scheduled_for   timestamp,
    seo_score       int,
    created_by      uuid         REFERENCES public.users(id) ON DELETE SET NULL,
    team_id         uuid         REFERENCES public.teams(id) ON DELETE SET NULL,
    user_id         uuid         REFERENCES public.users(id) ON DELETE SET NULL,
    created_at      timestamp    DEFAULT CURRENT_TIMESTAMP
);

-- 4.13 Keywords
CREATE TABLE IF NOT EXISTS public.keywords (
    id         uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    id_draft   uuid         REFERENCES public.drafts(id)    ON DELETE CASCADE ON UPDATE CASCADE,
    id_history uuid         REFERENCES public.histories(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name       varchar(255) NOT NULL,
    created_at timestamp    DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    DEFAULT CURRENT_TIMESTAMP
);

-- 4.14 Team Members
CREATE TABLE IF NOT EXISTS public.team_members (
    id        uuid        DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    team_id   uuid        NOT NULL REFERENCES public.teams(id) ON DELETE CASCADE,
    user_id   uuid        NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    role      varchar(50) DEFAULT 'member',
    joined_at timestamp   DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (team_id, user_id),
    CONSTRAINT chk_team_members_role CHECK (role IN ('owner', 'admin', 'member', 'viewer'))
);

-- 4.15 Team Invites
CREATE TABLE IF NOT EXISTS public.team_invites (
    id         uuid         DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    email      varchar(255) NOT NULL,
    team_id    uuid         NOT NULL REFERENCES public.teams(id) ON DELETE CASCADE,
    role       varchar(50)  DEFAULT 'member',
    token      varchar(255) NOT NULL UNIQUE,
    status     varchar(50)  DEFAULT 'pending',
    invited_by uuid         NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    expires_at timestamp    NOT NULL,
    created_at timestamp    DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_team_invites_status CHECK (status IN ('pending', 'accepted', 'expired', 'cancelled')),
    CONSTRAINT chk_team_invites_role   CHECK (role   IN ('admin', 'member', 'viewer', 'owner'))
);

-- 4.16 Roles
CREATE TABLE IF NOT EXISTS public.roles (
    id         serial       NOT NULL PRIMARY KEY,
    name       varchar(50)  NOT NULL UNIQUE,
    label      varchar(100) NOT NULL,
    created_at timestamptz  DEFAULT NOW()
);

-- 4.17 Permissions
CREATE TABLE IF NOT EXISTS public.permissions (
    id          serial           NOT NULL PRIMARY KEY,
    module      varchar(50)      NOT NULL,
    action      varchar(50)      NOT NULL,
    scope       permission_scope NOT NULL DEFAULT 'team',
    description text,
    created_at  timestamptz      DEFAULT NOW(),
    UNIQUE (module, action, scope)
);

-- 4.18 Role Permissions
CREATE TABLE IF NOT EXISTS public.role_permissions (
    role_id       int NOT NULL REFERENCES public.roles(id)       ON DELETE CASCADE,
    permission_id int NOT NULL REFERENCES public.permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);


-- ============================================================
-- 5. TRIGGERS — auto-update updated_at
-- ============================================================
DO $$
DECLARE
    t text;
BEGIN
    FOREACH t IN ARRAY ARRAY[
        'users', 'teams', 'api_providers', 'request_schemas', 'model_families',
        'api_keys', 'products', 'adapter_configs', 'workflow_definitions',
        'workflow_nodes', 'drafts', 'keywords', 'team_invites'
    ]
    LOOP
        EXECUTE format(
            'DROP TRIGGER IF EXISTS trg_set_updated_at ON public.%I;
             CREATE TRIGGER trg_set_updated_at
             BEFORE UPDATE ON public.%I
             FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();',
            t, t
        );
    END LOOP;
END $$;


-- ============================================================
-- 6. INDEXES
-- ============================================================

-- Users / Teams
CREATE INDEX IF NOT EXISTS idx_users_email                 ON public.users(email);
CREATE INDEX IF NOT EXISTS idx_team_members_user_id         ON public.team_members(user_id);
CREATE INDEX IF NOT EXISTS idx_team_members_team_id         ON public.team_members(team_id);

-- Providers / Schemas / Models / Keys
CREATE INDEX IF NOT EXISTS idx_api_providers_team_id        ON public.api_providers(team_id);
CREATE INDEX IF NOT EXISTS idx_request_schemas_provider_id  ON public.request_schemas(provider_id);
CREATE INDEX IF NOT EXISTS idx_model_families_provider_id   ON public.model_families(provider_id);
CREATE INDEX IF NOT EXISTS idx_model_families_schema_id     ON public.model_families(schema_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_provider_id         ON public.api_keys(provider_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_model_id            ON public.api_keys(model_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_team_id             ON public.api_keys(team_id);

-- Products / Adapter / Workflow
CREATE INDEX IF NOT EXISTS idx_products_team               ON public.products(team_id);
CREATE INDEX IF NOT EXISTS idx_products_created_by          ON public.products(created_by);
CREATE INDEX IF NOT EXISTS idx_products_user_id             ON public.products(user_id);
CREATE INDEX IF NOT EXISTS idx_adapter_configs_product_id   ON public.adapter_configs(product_id);
CREATE INDEX IF NOT EXISTS idx_workflow_definitions_product ON public.workflow_definitions(product_id);
CREATE INDEX IF NOT EXISTS idx_workflow_nodes_workflow_id   ON public.workflow_nodes(workflow_id);
CREATE INDEX IF NOT EXISTS idx_workflow_nodes_adapter_id    ON public.workflow_nodes(adapter_config_id);

-- Drafts / Histories / Keywords
CREATE INDEX IF NOT EXISTS idx_drafts_team                 ON public.drafts(team_id);
CREATE INDEX IF NOT EXISTS idx_drafts_created_by            ON public.drafts(created_by);
CREATE INDEX IF NOT EXISTS idx_drafts_user_id               ON public.drafts(user_id);
CREATE INDEX IF NOT EXISTS idx_histories_team_id            ON public.histories(team_id);
CREATE INDEX IF NOT EXISTS idx_histories_created_by         ON public.histories(created_by);
CREATE INDEX IF NOT EXISTS idx_histories_user_id            ON public.histories(user_id);
CREATE INDEX IF NOT EXISTS idx_keywords_id_draft            ON public.keywords(id_draft);
CREATE INDEX IF NOT EXISTS idx_keywords_id_history          ON public.keywords(id_history);

-- Team invites
CREATE INDEX IF NOT EXISTS idx_team_invites_email           ON public.team_invites(email);
CREATE INDEX IF NOT EXISTS idx_team_invites_team_id         ON public.team_invites(team_id);
CREATE INDEX IF NOT EXISTS idx_team_invites_token           ON public.team_invites(token);
CREATE INDEX IF NOT EXISTS idx_team_invites_status          ON public.team_invites(status);
CREATE INDEX IF NOT EXISTS idx_team_invites_expires_at      ON public.team_invites(expires_at);
CREATE INDEX IF NOT EXISTS idx_team_invites_email_status    ON public.team_invites(email, status);

CREATE UNIQUE INDEX IF NOT EXISTS idx_team_invites_email_team_pending
    ON public.team_invites(email, team_id)
    WHERE status = 'pending';

-- RBAC reverse lookup
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON public.role_permissions(permission_id);


-- ============================================================
-- 7. SEED DATA — Users & Teams
-- ============================================================

INSERT INTO public.users (id, email, name, password_hash, role, status, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'superadmin@kabar.com', 'Super Admin User', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'superadmin', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'admin@kabar.com',      'Admin User',       '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin',      'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'manager@kabar.com',    'Manager User',     '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'manager',    'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'editor@kabar.com',     'Editor User',      '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'editor',     'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'viewer@kabar.com',     'Viewer User',      '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'viewer',     'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (email) DO NOTHING;

INSERT INTO public.teams (id, name, description, created_by, created_at, updated_at)
SELECT
    gen_random_uuid(),
    'SEO Management Team',
    'Main team for SEO content management and scheduling',
    'admin@kabar.com',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
WHERE NOT EXISTS (
    SELECT 1 FROM public.teams WHERE name = 'SEO Management Team'
);

INSERT INTO public.team_members (id, team_id, user_id, role, joined_at)
SELECT gen_random_uuid(), t.id, u.id, 'owner', CURRENT_TIMESTAMP
FROM   public.teams t, public.users u
WHERE  t.name = 'SEO Management Team' AND u.email = 'superadmin@kabar.com'
ON CONFLICT (team_id, user_id) DO NOTHING;

INSERT INTO public.team_members (id, team_id, user_id, role, joined_at)
SELECT gen_random_uuid(), t.id, u.id, 'admin', CURRENT_TIMESTAMP
FROM   public.teams t, public.users u
WHERE  t.name = 'SEO Management Team' AND u.email = 'admin@kabar.com'
ON CONFLICT (team_id, user_id) DO NOTHING;

INSERT INTO public.team_members (id, team_id, user_id, role, joined_at)
SELECT gen_random_uuid(), t.id, u.id, 'member', CURRENT_TIMESTAMP
FROM   public.teams t, public.users u
WHERE  t.name = 'SEO Management Team' AND u.email = 'manager@kabar.com'
ON CONFLICT (team_id, user_id) DO NOTHING;

INSERT INTO public.team_members (id, team_id, user_id, role, joined_at)
SELECT gen_random_uuid(), t.id, u.id, 'member', CURRENT_TIMESTAMP
FROM   public.teams t, public.users u
WHERE  t.name = 'SEO Management Team' AND u.email = 'editor@kabar.com'
ON CONFLICT (team_id, user_id) DO NOTHING;

INSERT INTO public.team_members (id, team_id, user_id, role, joined_at)
SELECT gen_random_uuid(), t.id, u.id, 'member', CURRENT_TIMESTAMP
FROM   public.teams t, public.users u
WHERE  t.name = 'SEO Management Team' AND u.email = 'viewer@kabar.com'
ON CONFLICT (team_id, user_id) DO NOTHING;


-- ============================================================
-- 8. SEED DATA — API Providers
-- ============================================================

INSERT INTO public.api_providers (
    id, name, display_name, description,
    base_url, auth_type, auth_header, auth_prefix,
    default_headers, team_id, is_active,
    created_at, updated_at
) VALUES
(
    'a1000000-0000-0000-0000-000000000001',
    'anthropic', 'Anthropic',
    'Anthropic AI — provider of the Claude model family.',
    'https://api.anthropic.com',
    'bearer', 'x-api-key', '',
    '{"anthropic-version": "2023-06-01", "content-type": "application/json"}'::jsonb,
    NULL, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
(
    'a1000000-0000-0000-0000-000000000002',
    'openai', 'OpenAI',
    'OpenAI — provider of GPT and o-series reasoning models.',
    'https://api.openai.com',
    'bearer', 'Authorization', 'Bearer',
    '{"content-type": "application/json"}'::jsonb,
    NULL, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
(
    'a1000000-0000-0000-0000-000000000003',
    'google_ai', 'Google AI (Gemini)',
    'Google AI Studio — provider of the Gemini model family.',
    'https://generativelanguage.googleapis.com',
    'ApiKey', 'Authorization', 'ApiKey',
    '{"content-type": "application/json"}'::jsonb,
    NULL, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
(
    'a1000000-0000-0000-0000-000000000004',
    'groq', 'Groq',
    'Groq Cloud — ultra-fast inference for open-source LLMs.',
    'https://api.groq.com/openai',
    'bearer', 'Authorization', 'Bearer',
    '{"content-type": "application/json"}'::jsonb,
    NULL, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
(
    'a1000000-0000-0000-0000-000000000005',
    'openrouter', 'OpenRouter',
    'OpenRouter — unified gateway to 100+ LLM providers.',
    'https://openrouter.ai/api',
    'bearer', 'Authorization', 'Bearer',
    '{"content-type": "application/json", "HTTP-Referer": "https://kabar.com", "X-Title": "kabar.com"}'::jsonb,
    NULL, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
-- Catatan: upsert pakai ON CONFLICT (id), bukan (team_id, name).
-- team_id bernilai NULL untuk seluruh baris seed ini, dan UNIQUE(team_id, name)
-- TIDAK pernah ter-trigger saat team_id IS NULL (NULL tidak dianggap sama
-- dengan NULL lain di pengecekan UNIQUE biasa), sehingga run berulang akan
-- gagal pada id (primary key) yang sudah hardcoded. id dipakai sebagai
-- target konflik agar re-run benar-benar idempotent.
ON CONFLICT (id) DO UPDATE SET
    display_name    = EXCLUDED.display_name,
    description     = EXCLUDED.description,
    base_url        = EXCLUDED.base_url,
    auth_type       = EXCLUDED.auth_type,
    auth_header     = EXCLUDED.auth_header,
    auth_prefix     = EXCLUDED.auth_prefix,
    default_headers = EXCLUDED.default_headers,
    is_active       = EXCLUDED.is_active,
    updated_at      = CURRENT_TIMESTAMP;


-- ============================================================
-- 9. SEED DATA — Request Schemas
-- ============================================================

INSERT INTO public.request_schemas (
    id, provider_id, name, endpoint_path,
    max_tokens_key, system_role_key,
    response_text_path, response_image_path,
    request_template,
    supports_temperature, supports_streaming,
    created_at, updated_at
) VALUES

-- Anthropic Messages API
(
    'b1000000-0000-0000-0000-000000000001',
    'a1000000-0000-0000-0000-000000000001',
    'anthropic_messages_v1', '/v1/messages',
    'max_tokens', 'system',
    'content[0].text', NULL,
    '{
        "model": "{model}",
        "max_tokens": 1024,
        "temperature": 1,
        "system": "{system_prompt}",
        "messages": [
            { "role": "user", "content": "{prompt}" }
        ]
    }',
    true, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),

-- OpenAI Chat Completions
(
    'b1000000-0000-0000-0000-000000000002',
    'a1000000-0000-0000-0000-000000000002',
    'openai_chat_completions_v1', '/v1/chat/completions',
    'max_completion_tokens', 'messages',
    'choices[0].message.content', NULL,
    '{
        "model": "{model}",
        "max_completion_tokens": 1024,
        "temperature": 1,
        "messages": [
            { "role": "system", "content": "{system_prompt}" },
            { "role": "user",   "content": "{prompt}" }
        ]
    }',
    true, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),

-- OpenAI Responses API
(
    'b1000000-0000-0000-0000-000000000003',
    'a1000000-0000-0000-0000-000000000002',
    'openai_responses_v1', '/v1/responses',
    'max_output_tokens', 'instructions',
    'output[0].content[0].text', NULL,
    '{
        "model": "{model}",
        "max_output_tokens": 1024,
        "temperature": 1,
        "instructions": "{system_prompt}",
        "input": "{prompt}"
    }',
    true, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),

-- Google Gemini generateContent
(
    'b1000000-0000-0000-0000-000000000004',
    'a1000000-0000-0000-0000-000000000003',
    'gemini_generate_content_v1', '/v1beta/models/{model}:generateContent',
    'generationConfig.maxOutputTokens', 'system_instruction',
    'candidates[0].content.parts[0].text', NULL,
    '{
        "system_instruction": {
            "parts": [{ "text": "{system_prompt}" }]
        },
        "contents": [
            {
                "role": "user",
                "parts": [{ "text": "{prompt}" }]
            }
        ],
        "generationConfig": {
            "maxOutputTokens": 1024,
            "temperature": 1.0
        }
    }',
    true, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),

-- Groq
(
    'b1000000-0000-0000-0000-000000000005',
    'a1000000-0000-0000-0000-000000000004',
    'groq_chat_completions_v1', '/openai/v1/chat/completions',
    'max_tokens', 'messages',
    'choices[0].message.content', NULL,
    '{
        "model": "{model}",
        "max_tokens": 1024,
        "temperature": 1,
        "messages": [
            { "role": "system", "content": "{system_prompt}" },
            { "role": "user",   "content": "{prompt}" }
        ]
    }',
    true, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),

-- OpenRouter
(
    'b1000000-0000-0000-0000-000000000006',
    'a1000000-0000-0000-0000-000000000005',
    'openrouter_chat_completions_v1', '/api/v1/chat/completions',
    'max_tokens', 'messages',
    'choices[0].message.content', NULL,
    '{
        "model": "{model}",
        "max_tokens": 1024,
        "temperature": 1,
        "messages": [
            { "role": "system", "content": "{system_prompt}" },
            { "role": "user",   "content": "{prompt}" }
        ]
    }',
    true, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
ON CONFLICT (id) DO UPDATE SET
    endpoint_path        = EXCLUDED.endpoint_path,
    max_tokens_key       = EXCLUDED.max_tokens_key,
    system_role_key      = EXCLUDED.system_role_key,
    response_text_path   = EXCLUDED.response_text_path,
    response_image_path  = EXCLUDED.response_image_path,
    request_template     = EXCLUDED.request_template,
    supports_temperature = EXCLUDED.supports_temperature,
    supports_streaming   = EXCLUDED.supports_streaming,
    updated_at           = CURRENT_TIMESTAMP;


-- ============================================================
-- 10. SEED DATA — Model Families
-- ============================================================

INSERT INTO public.model_families (
    id, provider_id, schema_id,
    name, display_name, description,
    max_tokens, temperature, system_prompt,
    created_at, updated_at
) VALUES
-- Anthropic
(
    'c1000000-0000-0000-0000-000000000001',
    'a1000000-0000-0000-0000-000000000001',
    'b1000000-0000-0000-0000-000000000001',
    'claude_3_5', 'Claude 3.5',
    'Anthropic Claude 3.5 family — Haiku dan Sonnet.',
    1024, 1.0, 'You are a helpful assistant.',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
(
    'c1000000-0000-0000-0000-000000000002',
    'a1000000-0000-0000-0000-000000000001',
    'b1000000-0000-0000-0000-000000000001',
    'claude_3_7', 'Claude 3.7',
    'Anthropic Claude 3.7 Sonnet — extended thinking support.',
    1024, 1.0, 'You are a helpful assistant.',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
(
    'c1000000-0000-0000-0000-000000000003',
    'a1000000-0000-0000-0000-000000000001',
    'b1000000-0000-0000-0000-000000000001',
    'claude_4', 'Claude 4',
    'Anthropic Claude 4 family — Sonnet dan Opus.',
    1024, 1.0, 'You are a helpful assistant.',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
-- OpenAI
(
    'c1000000-0000-0000-0000-000000000004',
    'a1000000-0000-0000-0000-000000000002',
    'b1000000-0000-0000-0000-000000000002',
    'gpt_4o', 'GPT-4o',
    'OpenAI GPT-4o multimodal family.',
    1024, 1.0, 'You are a helpful assistant.',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
(
    'c1000000-0000-0000-0000-000000000005',
    'a1000000-0000-0000-0000-000000000002',
    'b1000000-0000-0000-0000-000000000002',
    'gpt_4_1', 'GPT-4.1',
    'OpenAI GPT-4.1 family dengan context window 1M token.',
    1024, 1.0, 'You are a helpful assistant.',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
(
    'c1000000-0000-0000-0000-000000000006',
    'a1000000-0000-0000-0000-000000000002',
    'b1000000-0000-0000-0000-000000000003',
    'o_series', 'o-series (Reasoning)',
    'OpenAI o-series reasoning models — pakai Responses API.',
    1024, 1.0, 'You are a helpful assistant.',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
-- Google
(
    'c1000000-0000-0000-0000-000000000007',
    'a1000000-0000-0000-0000-000000000003',
    'b1000000-0000-0000-0000-000000000004',
    'gemini_2_0', 'Gemini 2.0',
    'Google Gemini 2.0 family.',
    1024, 1.0, 'You are a helpful assistant.',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
(
    'c1000000-0000-0000-0000-000000000008',
    'a1000000-0000-0000-0000-000000000003',
    'b1000000-0000-0000-0000-000000000004',
    'gemini_2_5', 'Gemini 2.5',
    'Google Gemini 2.5 family — Flash dan Pro.',
    1024, 1.0, 'You are a helpful assistant.',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
-- Groq
(
    'c1000000-0000-0000-0000-000000000009',
    'a1000000-0000-0000-0000-000000000004',
    'b1000000-0000-0000-0000-000000000005',
    'llama_4', 'Llama 4 (Groq)',
    'Meta Llama 4 models served via Groq.',
    1024, 1.0, 'You are a helpful assistant.',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
(
    'c1000000-0000-0000-0000-000000000010',
    'a1000000-0000-0000-0000-000000000004',
    'b1000000-0000-0000-0000-000000000005',
    'deepseek_groq', 'DeepSeek (Groq)',
    'DeepSeek models served via Groq.',
    1024, 1.0, 'You are a helpful assistant.',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
),
-- OpenRouter
(
    'c1000000-0000-0000-0000-000000000011',
    'a1000000-0000-0000-0000-000000000005',
    'b1000000-0000-0000-0000-000000000006',
    'openrouter_mixed', 'OpenRouter Mixed',
    'Berbagai model agregat via gateway OpenRouter.',
    1024, 1.0, 'You are a helpful assistant.',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
ON CONFLICT (provider_id, name) DO UPDATE SET
    display_name  = EXCLUDED.display_name,
    description   = EXCLUDED.description,
    max_tokens    = EXCLUDED.max_tokens,
    temperature   = EXCLUDED.temperature,
    system_prompt = EXCLUDED.system_prompt,
    updated_at    = CURRENT_TIMESTAMP;


-- ============================================================
-- 11. SEED DATA — Roles & Permissions
-- ============================================================

INSERT INTO public.roles (name, label) VALUES
    ('superadmin', 'Super Admin'),
    ('admin',      'Admin'),
    ('manager',    'Manager'),
    ('editor',     'Editor'),
    ('viewer',     'Viewer')
ON CONFLICT (name) DO NOTHING;

INSERT INTO public.permissions (module, action, scope, description) VALUES
    ('draft',           'view',        'global', 'Lihat semua draft lintas team'),
    ('draft',           'view',        'team',   'Lihat draft milik team sendiri'),
    ('draft',           'create',      'team',   'Buat draft baru'),
    ('draft',           'edit',        'team',   'Edit draft'),
    ('draft',           'delete',      'team',   'Hapus draft'),
    ('draft',           'publish',     'team',   'Publish draft'),
    ('histories',       'view',        'global', 'Lihat semua histori lintas team'),
    ('histories',       'view',        'team',   'Lihat histori milik team sendiri'),
    ('histories',       'delete',      'global', 'Hapus histori (hanya superadmin)'),
    ('product',         'view',        'global', 'Lihat semua product lintas team'),
    ('product',         'view',        'team',   'Lihat product milik team sendiri'),
    ('product',         'create',      'team',   'Tambah product baru'),
    ('product',         'edit',        'team',   'Edit product'),
    ('product',         'delete',      'team',   'Hapus product'),
    ('product',         'inject',      'team',   'Inject product ke blog post'),
    ('user_management', 'view',        'global', 'Lihat semua user lintas team'),
    ('user_management', 'manage',      'global', 'Kelola semua user lintas team'),
    ('user_management', 'manage',      'team',   'Kelola user di team sendiri'),
    ('user_management', 'assign_role', 'team',   'Assign role ke user'),
    ('schedule',        'view',        'global', 'Lihat semua schedule lintas team'),
    ('schedule',        'view',        'team',   'Lihat schedule milik team sendiri'),
    ('schedule',        'create',      'team',   'Buat schedule baru'),
    ('schedule',        'edit',        'team',   'Edit schedule'),
    ('schedule',        'delete',      'team',   'Hapus schedule'),
    ('schedule',        'publish',     'team',   'Publish schedule'),
    ('api_keys',        'view',        'global', 'Lihat semua API keys lintas team'),
    ('api_keys',        'view',        'team',   'Lihat API keys milik team sendiri'),
    ('api_keys',        'create',      'team',   'Tambah API key baru'),
    ('api_keys',        'edit',        'team',   'Edit API key'),
    ('api_keys',        'delete',      'team',   'Hapus API key')
ON CONFLICT (module, action, scope) DO NOTHING;

-- Superadmin: semua permission
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'superadmin'
ON CONFLICT DO NOTHING;

-- Admin
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'admin'
  AND (
      (p.module = 'draft'           AND p.action IN ('view','create','edit','delete','publish') AND p.scope = 'team') OR
      (p.module = 'histories'       AND p.action = 'view'                                       AND p.scope = 'team') OR
      (p.module = 'product'         AND p.action IN ('view','create','edit','delete','inject')   AND p.scope = 'team') OR
      (p.module = 'user_management' AND p.action IN ('manage','assign_role')                    AND p.scope = 'team') OR
      (p.module = 'schedule'        AND p.action IN ('view','create','edit','delete','publish')  AND p.scope = 'team') OR
      (p.module = 'api_keys'        AND p.action IN ('view','create','edit','delete')            AND p.scope = 'team')
  )
ON CONFLICT DO NOTHING;

-- Manager
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'manager'
  AND (
      (p.module = 'draft'     AND p.action IN ('view','create','edit','delete','publish') AND p.scope = 'team') OR
      (p.module = 'histories' AND p.action = 'view'                                       AND p.scope = 'team') OR
      (p.module = 'product'   AND p.action IN ('view','create','edit','inject')           AND p.scope = 'team') OR
      (p.module = 'schedule'  AND p.action IN ('view','create','edit','publish')          AND p.scope = 'team') OR
      (p.module = 'api_keys'  AND p.action IN ('view','create','edit')                   AND p.scope = 'team')
  )
ON CONFLICT DO NOTHING;

-- Editor
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'editor'
  AND (
      (p.module = 'draft'     AND p.action IN ('view','create','edit') AND p.scope = 'team') OR
      (p.module = 'histories' AND p.action = 'view'                    AND p.scope = 'team') OR
      (p.module = 'product'   AND p.action IN ('view','inject')        AND p.scope = 'team') OR
      (p.module = 'schedule'  AND p.action IN ('view','create','edit') AND p.scope = 'team') OR
      (p.module = 'api_keys'  AND p.action = 'view'                   AND p.scope = 'team')
  )
ON CONFLICT DO NOTHING;

-- Viewer: view-only
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'viewer'
  AND p.module IN ('draft','histories','product','schedule','api_keys')
  AND p.action = 'view'
  AND p.scope  = 'team'
ON CONFLICT DO NOTHING;