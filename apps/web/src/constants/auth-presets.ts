// pages/admin/ai-management/constants/auth-presets.ts

export const AUTH_PRESETS = {
    bearer_openai: {
        auth_type: "bearer",
        auth_header: "Authorization",
        auth_prefix: "Bearer",
        label: "Bearer Token (OpenAI style)",
    },
    api_key_google: {
        auth_type: "api_key",
        auth_header: "x-goog-api-key",
        auth_prefix: null,
        label: "API Key Header (Google style)",
    },
    api_key_url: {
        auth_type: "api_key",
        auth_header: "",
        auth_prefix: null,
        label: "API Key URL",
    },
    custom: {
        auth_type: "custom",
        auth_header: "",
        auth_prefix: null,
        label: "Custom",
    },
} as const;

export type AuthPreset = keyof typeof AUTH_PRESETS;