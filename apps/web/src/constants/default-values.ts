// constants/default-values.ts

export const DEFAULT_REQUEST_SCHEMA = {
    id: "",
    provider_id: "",
    name: "",
    endpoint_path: "",
    max_tokens_key: "max_tokens",
    system_role_key: "system",
    request_template: "{}",
    response_text_path: "",
    response_image_path: "",
    supports_temperature: true,
    supports_streaming: true,
};

export const DEFAULT_FAMILY = {
    id: "",
    provider_id: "",
    schema_id: "",
    name: "",
    display_name: "",
    description: null,
    schema: DEFAULT_REQUEST_SCHEMA,
};

export const DEFAULT_PROVIDER_FORM = {
    name: "",
    display_name: "",
    description: null,
    base_url: "",
    auth_type: "bearer",
    auth_header: "Authorization",
    auth_prefix: "Bearer",
    default_headers: {},
    is_active: true,
    families: [],
    team_id: "",
};