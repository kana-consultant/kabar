export interface AIModel {
    id: string;
    name: string;
    providerId: string;
    displayName: string;
    description: string;
    isActive: boolean;
    isDefault: boolean;
    maxTokens: number;
    temperature: number;
    teamId?: string;
    provider: string;
    createdAt: string;
    updatedAt: string;
}

export interface ModelWithStatus extends AIModel {
    hasApiKey: boolean;
    providerDisplayName?: string;
}

export interface CreateModelRequest {
    name: string;
    providerId: string;
    displayName: string;
    description?: string;
    isDefault?: boolean;
    maxTokens?: number;
    temperature?: number;
}

export interface ModelFromAPIKey {
    id: string;
    modelId: string;
    name: string;
    displayName: string;
    providerName: string;
    providerDisplayName: string;
    service: string;
}