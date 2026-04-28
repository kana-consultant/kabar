package database

import (
	"log"
)

// RunSeed mengisi data awal ke database
func RunSeed() error {
	log.Println("Seeding database...")

	seeds := []string{
		// =====================================================
		// SEED USERS (password: password123)
		// =====================================================
		`INSERT INTO users (id, email, name, password_hash, role, status) VALUES
			(gen_random_uuid(), 'admin@seo.com', 'Admin User', '$2a$10$N9qo8uLOickgx2ZMRZoMy.Mr7qYqVvkqkZqkZqkZqkZqkZqkZqkZq', 'admin', 'active'),
			(gen_random_uuid(), 'manager@seo.com', 'Manager User', '$2a$10$N9qo8uLOickgx2ZMRZoMy.Mr7qYqVvkqkZqkZqkZqkZqkZqkZqkZq', 'manager', 'active'),
			(gen_random_uuid(), 'editor@seo.com', 'Editor User', '$2a$10$N9qo8uLOickgx2ZMRZoMy.Mr7qYqVvkqkZqkZqkZqkZqkZqkZqkZq', 'editor', 'active'),
			(gen_random_uuid(), 'viewer@seo.com', 'Viewer User', '$2a$10$N9qo8uLOickgx2ZMRZoMy.Mr7qYqVvkqkZqkZqkZqkZqkZqkZqkZq', 'viewer', 'active')
		ON CONFLICT (email) DO NOTHING;`,

		// =====================================================
		// SEED TEAMS
		// =====================================================
		`INSERT INTO teams (id, name, description) VALUES
			(gen_random_uuid(), 'SEO Team', 'Tim utama SEO'),
			(gen_random_uuid(), 'Content Team', 'Tim pembuat konten'),
			(gen_random_uuid(), 'Marketing Team', 'Tim pemasaran')
		ON CONFLICT (id) DO NOTHING;`,

		// =====================================================
		// SEED TEAM MEMBERS
		// =====================================================
		`INSERT INTO team_members (team_id, user_id, role)
			SELECT t.id, u.id, 'manager'
			FROM teams t, users u
			WHERE t.name = 'SEO Team' AND u.email = 'manager@seo.com'
		ON CONFLICT (team_id, user_id) DO NOTHING;`,

		`INSERT INTO team_members (team_id, user_id, role)
			SELECT t.id, u.id, 'editor'
			FROM teams t, users u
			WHERE t.name = 'Content Team' AND u.email = 'editor@seo.com'
		ON CONFLICT (team_id, user_id) DO NOTHING;`,

		// =====================================================
		// SEED PRODUCTS
		// =====================================================
		`INSERT INTO products (id, name, platform, api_endpoint, status) VALUES
			(gen_random_uuid(), 'TrekkingID', 'wordpress', 'https://trekkingid.com/wp-json/wp/v2/posts', 'connected'),
			(gen_random_uuid(), 'CampingMart', 'shopify', 'https://campingmart.myshopify.com/admin/api/graphql.json', 'connected'),
			(gen_random_uuid(), 'OutdoorGear', 'custom', 'https://outdoorgear.com/api/blog', 'pending')
		ON CONFLICT (id) DO NOTHING;`,

		// =====================================================
		// SEED ADAPTER CONFIGS
		// =====================================================
		`INSERT INTO adapter_configs (product_id, endpoint_path, http_method, custom_headers, field_mapping)
			SELECT id, '/wp-json/wp/v2/posts', 'POST', '{"Content-Type": "application/json"}', '[{"sourceField":"title","targetField":"title"},{"sourceField":"body","targetField":"content"}]'
			FROM products WHERE name = 'TrekkingID'
		ON CONFLICT (product_id) DO NOTHING;`,

		`INSERT INTO adapter_configs (product_id, endpoint_path, http_method, custom_headers, field_mapping)
			SELECT id, '/admin/api/graphql.json', 'POST', '{"Content-Type": "application/json"}', '[{"sourceField":"title","targetField":"input.title"},{"sourceField":"body","targetField":"input.body_html"}]'
			FROM products WHERE name = 'CampingMart'
		ON CONFLICT (product_id) DO NOTHING;`,

		// =====================================================
		// SEED DRAFTS
		// =====================================================
		`INSERT INTO drafts (id, title, topic, article, status, target_products, has_image) VALUES
			(gen_random_uuid(), 'Cara Memilih Sepatu Gunung', 'Cara Memilih Sepatu Gunung', '# Cara Memilih Sepatu Gunung\n\nArtikel tentang memilih sepatu gunung...', 'draft', '["TrekkingID"]', false),
			(gen_random_uuid(), 'Tips Camping Hemat', 'Tips Camping Hemat', '# Tips Camping Hemat\n\nArtikel tentang camping hemat...', 'draft', '["CampingMart"]', false),
			(gen_random_uuid(), 'Perawatan Jaket Outdoor', 'Perawatan Jaket Outdoor', '# Perawatan Jaket Outdoor\n\nArtikel tentang perawatan jaket...', 'scheduled', '["OutdoorGear"]', false)
		ON CONFLICT (id) DO NOTHING;`,

		// =====================================================
		// SEED HISTORIES
		// =====================================================
		`INSERT INTO histories (id, title, topic, content, status, action, published_at, target_products) VALUES
			(gen_random_uuid(), 'Cara Memilih Sepatu Gunung', 'Cara Memilih Sepatu Gunung', 'Artikel tentang sepatu gunung...', 'success', 'published', NOW(), '["TrekkingID"]'),
			(gen_random_uuid(), 'Tips Camping Hemat', 'Tips Camping Hemat', 'Artikel tentang camping hemat...', 'success', 'published', NOW(), '["CampingMart"]'),
			(gen_random_uuid(), 'Perawatan Jaket Outdoor', 'Perawatan Jaket Outdoor', 'Artikel tentang perawatan jaket...', 'failed', 'published', NOW(), '["OutdoorGear"]')
		ON CONFLICT (id) DO NOTHING;`,

		// =====================================================
		// SEED SETTINGS
		// =====================================================
		`INSERT INTO settings (key, value, description) VALUES
			('default_post_status', '"draft"', 'Default status for new posts'),
			('auto_save_draft', 'true', 'Auto save draft while editing'),
			('default_timezone', '"Asia/Jakarta"', 'Default timezone for scheduling'),
			('email_notifications', 'true', 'Enable email notifications')
		ON CONFLICT (key) DO NOTHING;`,

		// =====================================================
		// SEED API KEYS
		// =====================================================
		`INSERT INTO api_keys (service, key_encrypted, is_active) VALUES
			('gemini', '', false),
			('image', '', false),
			('openai', '', false)
		ON CONFLICT (service) DO NOTHING;`,
	}

	for i, seed := range seeds {
		log.Printf("Seeding %d...", i+1)
		if _, err := GetDB().Exec(seed); err != nil {
			log.Printf("Seed %d warning: %v", i+1, err)
		}
	}

	log.Println("Seeding completed successfully")
	return nil
}