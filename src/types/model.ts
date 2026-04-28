// src/types/model.ts
export interface ModelWithStatus {
    id: string;
    name: string;
    providerId: string;
    displayName: string;
    description: string;
    isActive: boolean;
    isDefault: boolean;
    maxTokens: number;
    temperature: number;
    createdAt: string;
    updatedAt: string;
    hasApiKey: boolean;
    provider?: string;
    providerDisplayName?: string;
}