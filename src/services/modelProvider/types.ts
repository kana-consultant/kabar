export interface APIProvider {
    id: string;
    name: string;
    displayName: string;
    description: string;
    baseUrl: string;
    authType: 'bearer' | 'api_key' | 'x-api-key';
    authHeader: string;
    authPrefix: string;
    textEndpoint: string;
    imageEndpoint?: string;
    defaultHeaders: Record<string, string>;
    requestTemplate: Record<string, any>;
    responseTextPath: string;
    responseImagePath?: string;
    isActive: boolean;
    createdAt: string;
    updatedAt: string;
}


export interface CreateProviderRequest {
    name: string;
    displayName: string;
    description?: string;
    baseUrl: string;
    authType: string;
    authHeader: string;
    authPrefix: string;
    textEndpoint: string;
    imageEndpoint?: string;
    defaultHeaders: Record<string, string>;
    requestTemplate: Record<string, any>;
    responseTextPath: string;
    responseImagePath?: string;
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

export interface UpdateProviderRequest extends Partial<CreateProviderRequest> {}
export interface UpdaFteModelRequest extends Partial<CreateModelRequest> {}

export interface CreateResponse {
    id: string;
    message: string;
}

export interface UpdateResponse {
    message: string;
}


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
    provider: string;
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