-- =====================================================
-- 1. USERS TABLE (RBAC)
-- =====================================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'manager', 'editor', 'viewer')),
    avatar TEXT,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    last_active TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);

-- =====================================================
-- 2. TEAMS TABLE
-- =====================================================
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_teams_name ON teams(name);

-- =====================================================
-- 3. TEAM MEMBERS TABLE (Many-to-Many)
-- =====================================================
CREATE TABLE team_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('manager', 'editor', 'viewer')),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(team_id, user_id)
);

CREATE INDEX idx_team_members_team_id ON team_members(team_id);
CREATE INDEX idx_team_members_user_id ON team_members(user_id);

-- =====================================================
-- 4. PRODUCTS TABLE
-- =====================================================
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    platform VARCHAR(50) NOT NULL CHECK (platform IN ('wordpress', 'shopify', 'custom')),
    api_endpoint TEXT NOT NULL,
    api_key_encrypted TEXT, -- Encrypted API key
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('connected', 'pending')),
    last_sync TIMESTAMP,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_platform ON products(platform);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_team_id ON products(team_id);
CREATE INDEX idx_products_created_by ON products(created_by);

-- =====================================================
-- 5. ADAPTER CONFIGS TABLE (One-to-One dengan Product)
-- =====================================================
CREATE TABLE adapter_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL UNIQUE REFERENCES products(id) ON DELETE CASCADE,
    endpoint_path VARCHAR(500) NOT NULL,
    http_method VARCHAR(10) DEFAULT 'POST' CHECK (http_method IN ('POST', 'PUT', 'PATCH')),
    custom_headers JSONB DEFAULT '{}'::jsonb,
    field_mapping TEXT NOT NULL, -- JSON string (bisa format apapun)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_adapter_configs_product_id ON adapter_configs(product_id);

-- =====================================================
-- 6. DRAFTS TABLE
-- =====================================================
CREATE TABLE drafts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    topic VARCHAR(500) NOT NULL,
    article TEXT NOT NULL,
    image_url TEXT,
    image_prompt TEXT,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'scheduled', 'published')),
    scheduled_for TIMESTAMP,
    target_products JSONB DEFAULT '[]'::jsonb, -- Array of product IDs or names
    has_image BOOLEAN DEFAULT FALSE,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_drafts_status ON drafts(status);
CREATE INDEX idx_drafts_created_by ON drafts(created_by);
CREATE INDEX idx_drafts_team_id ON drafts(team_id);
CREATE INDEX idx_drafts_scheduled_for ON drafts(scheduled_for);
CREATE INDEX idx_drafts_title ON drafts(title);

-- =====================================================
-- 7. HISTORIES TABLE
-- =====================================================
CREATE TABLE histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    topic VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    target_products JSONB DEFAULT '[]'::jsonb,
    status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'failed', 'pending')),
    action VARCHAR(20) NOT NULL CHECK (action IN ('published', 'scheduled', 'draft_saved')),
    error_message TEXT,
    published_at TIMESTAMP NOT NULL,
    scheduled_for TIMESTAMP,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_histories_status ON histories(status);
CREATE INDEX idx_histories_action ON histories(action);
CREATE INDEX idx_histories_created_by ON histories(created_by);
CREATE INDEX idx_histories_team_id ON histories(team_id);
CREATE INDEX idx_histories_published_at ON histories(published_at);

-- =====================================================
-- 8. SCHEDULES TABLE (Optional - bisa dari draft dengan status scheduled)
-- =====================================================
CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    draft_id UUID NOT NULL REFERENCES drafts(id) ON DELETE CASCADE,
    scheduled_for TIMESTAMP NOT NULL,
    is_daily BOOLEAN DEFAULT FALSE,
    daily_time TIME,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processed', 'failed')),
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_schedules_draft_id ON schedules(draft_id);
CREATE INDEX idx_schedules_scheduled_for ON schedules(scheduled_for);
CREATE INDEX idx_schedules_status ON schedules(status);

-- =====================================================
-- 9. API KEYS TABLE
-- =====================================================
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service VARCHAR(50) NOT NULL CHECK (service IN ('gemini', 'image', 'openai')),
    key_encrypted TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_api_keys_service ON api_keys(service);
CREATE INDEX idx_api_keys_is_active ON api_keys(is_active);

-- =====================================================
-- 10. SETTINGS TABLE (Global settings)
-- =====================================================
CREATE TABLE settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(255) NOT NULL UNIQUE,
    value JSONB NOT NULL,
    description TEXT,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_settings_key ON settings(key);

-- =====================================================
-- DEFAULT SETTINGS DATA
-- =====================================================
INSERT INTO settings (key, value, description) VALUES
    ('default_post_status', '"draft"', 'Default status for new posts'),
    ('auto_save_draft', 'true', 'Auto save draft while editing'),
    ('default_timezone', '"Asia/Jakarta"', 'Default timezone for scheduling'),
    ('email_notifications', 'true', 'Enable email notifications');

-- =====================================================
-- DEFAULT API KEYS (placeholder)
-- =====================================================
INSERT INTO api_keys (service, key_encrypted, is_active) VALUES
    ('gemini', '', false),
    ('image', '', false);

-- =====================================================
-- TRIGGER untuk updated_at (PostgreSQL)
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger ke semua tabel yang punya updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_teams_updated_at BEFORE UPDATE ON teams FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_adapter_configs_updated_at BEFORE UPDATE ON adapter_configs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_drafts_updated_at BEFORE UPDATE ON drafts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_schedules_updated_at BEFORE UPDATE ON schedules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_api_keys_updated_at BEFORE UPDATE ON api_keys FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_settings_updated_at BEFORE UPDATE ON settings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- DEFAULT USERS (untuk development)
-- =====================================================
INSERT INTO users (id, email, name, role, status) VALUES
    (gen_random_uuid(), 'admin@seo.com', 'Admin User', 'admin', 'active'),
    (gen_random_uuid(), 'manager@seo.com', 'Manager User', 'manager', 'active'),
    (gen_random_uuid(), 'editor@seo.com', 'Editor User', 'editor', 'active'),
    (gen_random_uuid(), 'viewer@seo.com', 'Viewer User', 'viewer', 'active');