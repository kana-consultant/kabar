export interface APIKey {
    id: string;
    service: string;
    systemPrompt?: string;
    teamId?: string;
    isActive: boolean;
    createdAt: string;
    providerId: string;
    updatedAt: string;
    modelId: string;
}

export interface APIKeyDetail {
    id: string;
    service: string;
    provider_id: string;
    model_id: string;
    is_active: boolean;
    system_prompt: string;
    created_by: string;
    created_at: string;
    updated_at: string;
    provider_name: string;
    provider_display_name: string;
    model_name: string;
    model_display_name: string;
}

export interface CreateAPIKeyRequest {
    service: string;
    key: string;
    modelId: string;
    providerId: string;
    systemPrompt?: string;
    teamId?: string;
}

export interface UpdateAPIKeyRequest {
    service?: string;
    key?: string;
    systemPrompt?: string;
    isActive?: boolean;
}