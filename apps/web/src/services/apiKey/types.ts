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
    providerId: string;
    modelId: string;
    isActive: boolean;
    systemPrompt: string;
    createdBy: string;
    createdAt: string;
    updatedAt: string;
    providerName: string;
    providerDisplayName: string;
    modelName: string;
    modelDisplayName: string;
}

export interface CreateAPIKeyRequest {
    service: string;
    key: string;
    modelId: String;
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