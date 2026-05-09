-- 1. EXTENSIONS
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

-- 2. CORE TABLES
CREATE TABLE IF NOT EXISTS public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    email character varying(255) NOT NULL UNIQUE,
    name character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    role character varying(50) DEFAULT 'viewer',
    avatar text,
    status character varying(20) DEFAULT 'active',
    last_active timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.teams (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name character varying(255) NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by character varying(255)
);

CREATE TABLE IF NOT EXISTS public.api_providers (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name character varying(100) NOT NULL,
    display_name character varying(200) NOT NULL,
    description text,
    base_url character varying(500) NOT NULL,
    auth_type character varying(50) DEFAULT 'bearer',
    auth_header character varying(100) DEFAULT 'Authorization',
    auth_prefix character varying(50) DEFAULT 'Bearer',
    text_endpoint character varying(200) NOT NULL,
    image_endpoint character varying(200),
    default_headers jsonb DEFAULT '{}'::jsonb,
    request_template text NOT NULL,
    response_text_path character varying(500),
    response_image_path character varying(500),
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.ai_models (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name character varying(200) NOT NULL,
    provider_id uuid NOT NULL,
    display_name character varying(200) NOT NULL,
    description text,
    request_template text,
    response_text_path character varying(500),
    response_image_path character varying(500),
    is_active boolean DEFAULT true,
    is_default boolean DEFAULT false,
    max_tokens integer DEFAULT 4096,
    temperature double precision DEFAULT 0.7,
    team_id uuid,
    created_by character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.api_keys (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    service character varying(50) NOT NULL,
    provider_id uuid NOT NULL,
    model_id uuid NOT NULL,
    key_encrypted text NOT NULL,
    system_prompt text,
    team_id uuid,
    is_active boolean DEFAULT true,
    created_by character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.products (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name character varying(255) NOT NULL,
    platform character varying(50) NOT NULL,
    api_endpoint text NOT NULL,
    api_key_encrypted text,
    status character varying(20) DEFAULT 'pending',
    sync_status character varying(20) DEFAULT 'idle',
    last_sync timestamp without time zone,
    created_by uuid,
    team_id uuid,
    user_id uuid,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.adapter_configs (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    product_id uuid NOT NULL,
    endpoint_path character varying(500) NOT NULL,
    http_method character varying(10) DEFAULT 'POST',
    custom_headers jsonb DEFAULT '{}'::jsonb,
    field_mapping text NOT NULL,
    timeout_seconds integer DEFAULT 30,
    retry_count integer DEFAULT 3,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.drafts (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    title character varying(500) NOT NULL,
    topic character varying(500) NOT NULL,
    article text NOT NULL,
    image_url text,
    image_prompt text,
    status character varying(20) DEFAULT 'draft',
    scheduled_for timestamp without time zone,
    target_products jsonb DEFAULT '[]'::jsonb,
    has_image boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by uuid,
    team_id uuid,
    user_id uuid,
    published_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS public.histories (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    title character varying(500) NOT NULL,
    topic character varying(500) NOT NULL,
    content text NOT NULL,
    image_url text,
    target_products jsonb DEFAULT '[]'::jsonb,
    status character varying(20) NOT NULL,
    action character varying(20) NOT NULL,
    error_message text,
    published_at timestamp without time zone NOT NULL,
    scheduled_for timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by uuid,
    team_id uuid,
    user_id uuid
);

CREATE TABLE IF NOT EXISTS public.team_members (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    team_id uuid NOT NULL,
    user_id uuid NOT NULL,
    role character varying(50) DEFAULT 'member',
    joined_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

-- 3. FOREIGN KEYS (fk_tm_user di-comment dulu)
ALTER TABLE public.team_members ADD CONSTRAINT fk_tm_team FOREIGN KEY (team_id) REFERENCES public.teams(id) ON DELETE CASCADE;
-- ALTER TABLE public.team_members ADD CONSTRAINT fk_tm_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE public.ai_models ADD CONSTRAINT fk_ai_provider FOREIGN KEY (provider_id) REFERENCES public.api_providers(id) ON DELETE CASCADE;
ALTER TABLE public.api_keys ADD CONSTRAINT fk_key_provider FOREIGN KEY (provider_id) REFERENCES public.api_providers(id) ON DELETE CASCADE;
ALTER TABLE public.api_keys ADD CONSTRAINT fk_key_model FOREIGN KEY (model_id) REFERENCES public.ai_models(id) ON DELETE CASCADE;
ALTER TABLE public.products ADD CONSTRAINT fk_prod_team FOREIGN KEY (team_id) REFERENCES public.teams(id) ON DELETE SET NULL;
ALTER TABLE public.adapter_configs ADD CONSTRAINT fk_adapter_prod FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;
ALTER TABLE public.drafts ADD CONSTRAINT fk_draft_team FOREIGN KEY (team_id) REFERENCES public.teams(id) ON DELETE SET NULL;

-- 4. INDEXES
CREATE INDEX IF NOT EXISTS idx_users_email ON public.users(email);
CREATE INDEX IF NOT EXISTS idx_products_team ON public.products(team_id);
CREATE INDEX IF NOT EXISTS idx_drafts_team ON public.drafts(team_id);


-- Insert data ke api_providers
INSERT INTO public.api_providers (id, name, display_name, description, base_url, auth_type, auth_header, auth_prefix, text_endpoint, image_endpoint, default_headers, request_template, response_text_path, response_image_path, is_active, created_at, updated_at) VALUES
('df0bf601-acc0-418a-81f4-2c14d22d02f5', 'openrouter', 'OpenRouter', 'Unified API for multiple AI models', 'https://openrouter.ai/api/v1', 'bearer', 'Authorization', 'Bearer', '/chat/completions', NULL, '{}', '{"model":"{model}","messages":[{"role":"user","content":"{prompt}"}]}', 'choices[0].message.content', NULL, true, NOW(), NOW()),
('a07bba3f-ef56-4da7-ae37-8cdba3cacd63', 'google-gemini', 'Google Gemini', 'Google Gemini API', 'https://generativelanguage.googleapis.com/v1beta', 'bearer', 'x-goog-api-key', NULL, '/models/{model}:generateContent', NULL, '{}', '{"contents":[{"parts":[{"text":"{prompt}"}]}]}', 'candidates[0].content.parts[0].text', NULL, true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert data ke ai_models
INSERT INTO public.ai_models (id, name, provider_id, display_name, description, request_template, response_text_path, response_image_path, is_active, is_default, max_tokens, temperature, team_id, created_by, created_at, updated_at) VALUES
('7c759977-e6fe-472a-bcbe-c3ca510f112c', 'google/gemini-3.1-flash-image-preview', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Nano Banana 2', NULL, '{"model":"{model}","messages":[{"role":"user","content":"{prompt}"}],"modalities":["image","text"]}', NULL, 'choices[0].message.images[0].image_url.url', true, true, 4096, 0.7, NULL, NULL, '2026-04-26 01:51:56.050348', '2026-04-26 01:51:56.050348'),
('995629ca-c31a-4952-8af6-91692e88a6f5', 'google/gemini-2.5-pro', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Gemini 2.5 Pro (OpenRouter)', 'Most powerful Gemini model via OpenRouter. 1M context, excellent for complex research and long-form content.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, true, 1000000, 0.5, NULL, NULL, '2026-04-26 19:00:54.435157', '2026-04-26 19:00:54.435157'),
('89e6f0c4-3cff-4f53-aeb3-14c7d2c796b4', 'google/gemini-2.5-flash-lite', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Gemini 2.5 Flash-Lite (OpenRouter)', 'Lightweight version of Gemini 2.5 Flash. Most cost-effective for high-volume tasks.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 1000000, 0.7, NULL, NULL, '2026-04-26 19:00:54.435157', '2026-04-26 19:00:54.435157'),
('c131bcb3-cd3c-4cb8-a36f-4869c7dfe8d4', 'google/gemini-2.0-flash-lite', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Gemini 2.0 Flash-Lite (OpenRouter)', 'Most lightweight Gemini model. Best for simple, high-volume generation tasks.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 1000000, 0.7, NULL, NULL, '2026-04-26 19:00:54.435157', '2026-04-26 19:00:54.435157'),
('15a6707b-1e05-4337-8269-4b8e4827e4d2', 'google/gemini-2.0-flash-001', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Gemini 2.5 Flash (OpenRouter)', 'Fast and efficient Gemini model via OpenRouter. Great for daily content generation.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 1000000, 0.7, NULL, NULL, '2026-04-26 19:00:54.435157', '2026-04-26 19:00:54.435157'),
('e6e84f37-c6ac-4e4b-a54c-12662572ac66', 'openai/gpt-4o', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'GPT-4o', 'OpenAI flagship model. High intelligence, fast, multimodal. Best for complex tasks.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, true, 128000, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('b1234cef-013b-42ef-88c9-6a026d830390', 'openai/gpt-4o-mini', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'GPT-4o Mini', 'Smaller, faster, cheaper version of GPT-4o. Great for high-volume tasks.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 128000, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('822c1845-3a06-4193-a346-bbc3fccb526a', 'openai/gpt-3.5-turbo', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'GPT-3.5 Turbo', 'Fast and cost-effective model for everyday tasks.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 16384, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('4916bd4c-66aa-4611-964c-172517a75973', 'anthropic/claude-3.5-sonnet', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Claude 3.5 Sonnet', 'Most intelligent Claude model. Excellent for complex reasoning, coding, and content generation.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 200000, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('915c837a-e586-4483-a94b-278b5e1aaff9', 'anthropic/claude-3-haiku', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Claude 3 Haiku', 'Fastest Claude model. Good for simple tasks and high-volume processing.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 200000, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('b6f38b71-c215-462f-88ee-0153f1c80402', 'meta-llama/llama-3.3-70b-instruct', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Llama 3.3 70B', 'Meta most capable open model. Strong performance for general tasks.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 128000, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('b641fa63-bbba-435c-8171-ef9febafde8a', 'google/gemini-2.5-flash-image', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Nano Banana', 'Image generation with contextual understanding. Support multi-turn conversations for image editing.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}], "modalities": ["image", "text"]}', NULL, 'choices[0].message.images[0].image_url.url', true, false, 4096, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('92d4f2c6-99c5-434a-ae32-0d9e8b45c8ac', 'google/gemini-3-pro-image-preview', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Nano Banana Pro', 'Most advanced model. Support 2K/4K outputs, identity preservation up to 5 subjects.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}], "modalities": ["image", "text"], "image_config": {"aspect_ratio": "16:9", "output_quality": "4k"}}', NULL, 'choices[0].message.images[0].image_url.url', true, false, 4096, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('d2e61d86-576f-4960-9713-80e5ef441f4d', 'gemini-2.0-flash-image', 'a07bba3f-ef56-4da7-ae37-8cdba3cacd63', 'Gemini 2.0 Flash Image', 'Multimodal model that can generate and edit images.', '{"contents": [{"parts": [{"text": "{prompt}"}]}]}', NULL, 'candidates[0].content.parts[0].inlineData.data', true, true, 4096, 0.7, NULL, NULL, '2026-04-26 19:05:36.972777', '2026-04-26 19:05:36.972777'),
('b5d8196f-f1cb-4aee-9b28-330aa3b0c846', 'imagen-3.0-generate-001', 'a07bba3f-ef56-4da7-ae37-8cdba3cacd63', 'Imagen 3.0', 'Google highest quality image generation model. Photorealistic results.', '{"instances": [{"prompt": "{prompt}"}]}', NULL, 'predictions[0].bytesBase64Encoded', true, false, 4096, 0.7, NULL, NULL, '2026-04-26 19:05:36.972777', '2026-04-26 19:05:36.972777'),
('c63e09e3-43c3-44ec-b4da-24d4e0e60adc', 'google/gemini-2.0-flash', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Gemini 2.0 Flash (OpenRouter)', 'Fast and versatile Gemini model. Balanced performance for scaling across diverse tasks.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 1000000, 0.7, NULL, NULL, '2026-04-26 19:00:54.435157', '2026-04-26 19:00:54.435157'),
('b7cd3b8d-10d9-4260-b27e-58ca16a5f825', 'gemini-2.5-pro', 'a07bba3f-ef56-4da7-ae37-8cdba3cacd63', 'Gemini 2.5 Pro', 'Powerful model for complex content, research, and reasoning tasks.', '{"contents": [{"parts": [{"text": "{prompt}"}]}]}', 'candidates[0].content.parts[0].text', NULL, true, false, 32768, 0.5, NULL, NULL, '2026-04-26 19:05:36.972777', '2026-04-26 19:05:36.972777'),
('256b25ce-a127-49e6-be85-bd6c8c858af6', 'gemini-2.0-flash', 'a07bba3f-ef56-4da7-ae37-8cdba3cacd63', 'Gemini 2.0 Flash', 'Fast and versatile multimodal model for scaling across diverse tasks.', '{"contents": [{"parts": [{"text": "{prompt}"}]}]}', 'candidates[0].content.parts[0].text', NULL, true, false, 8192, 0.7, NULL, NULL, '2026-04-26 19:05:36.972777', '2026-04-26 19:05:36.972777'),
('88f46222-91a9-41d7-b95c-39e579d7413b', 'gemini-2.0-flash-lite', 'a07bba3f-ef56-4da7-ae37-8cdba3cacd63', 'Gemini 2.0 Flash-Lite', 'Lightweight version of Gemini 2.0 Flash for cost-effective high-volume tasks.', '{"contents": [{"parts": [{"text": "{prompt}"}]}]}', 'candidates[0].content.parts[0].text', NULL, true, false, 8192, 0.7, NULL, NULL, '2026-04-26 19:05:36.972777', '2026-04-26 19:05:36.972777'),
('5c03d223-ab6f-4482-984f-6c1732cbebbf', 'gemini-1.5-flash', 'a07bba3f-ef56-4da7-ae37-8cdba3cacd63', 'Gemini 1.5 Flash', 'Fast and efficient model with 1M context window.', '{"contents": [{"parts": [{"text": "{prompt}"}]}]}', 'candidates[0].content.parts[0].text', NULL, true, false, 1000000, 0.7, NULL, NULL, '2026-04-26 19:05:36.972777', '2026-04-26 19:05:36.972777'),
('488c8fe8-e7df-47f3-a275-678c3f895b70', 'gemini-1.5-pro', 'a07bba3f-ef56-4da7-ae37-8cdba3cacd63', 'Gemini 1.5 Pro', 'Powerful model with 2M context window for complex reasoning.', '{"contents": [{"parts": [{"text": "{prompt}"}]}]}', 'candidates[0].content.parts[0].text', NULL, true, false, 2000000, 0.5, NULL, NULL, '2026-04-26 19:05:36.972777', '2026-04-26 19:05:36.972777'),
('0420f294-4668-424d-a0f9-01e277adeb2d', 'meta-llama/llama-3.1-8b-instruct:free', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Llama 3.1 8B (Free)', 'Free version of Llama 3.1 8B. Good for testing and low-budget tasks.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 8192, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('e70322d1-3d64-472a-9e41-3478508b6a45', 'mistralai/mistral-large', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Mistral Large', 'Mistral top-tier model. Strong reasoning and instruction following.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 128000, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('2ea9d218-2f80-4065-9a82-f2a201f66676', 'mistralai/mistral-nemo', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Mistral Nemo', 'Small, efficient Mistral model. Good for high-volume tasks.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 128000, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('f72f9a89-c457-49cb-820d-47ed25661382', 'deepseek/deepseek-chat', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'DeepSeek-V3', 'High-performance model from DeepSeek. Excellent for coding and reasoning.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 128000, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('a7152dad-5055-446b-9ccb-7d70f985d29a', 'deepseek/deepseek-coder', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'DeepSeek-Coder', 'Specialized for code generation and programming tasks.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 128000, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507'),
('3e2bada5-3039-4687-8441-06726795e47b', 'qwen/qwen-2.5-72b-instruct', 'df0bf601-acc0-418a-81f4-2c14d22d02f5', 'Qwen 2.5 72B', 'Alibaba powerful model. Strong multilingual support.', '{"model": "{model}", "messages": [{"role": "user", "content": "{prompt}"}]}', 'choices[0].message.content', NULL, true, false, 128000, 0.7, NULL, NULL, '2026-04-26 18:41:12.050507', '2026-04-26 18:41:12.050507')
ON CONFLICT (id) DO NOTHING;