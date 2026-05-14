import { apiClient } from '../api';
import type {
    AIModel, ModelWithStatus, ModelFromAPIKey,
    UpdateResponse, UpdateProviderRequest, CreateProviderRequest, CreateResponse
} from './types';
import { type APIProvider } from './types';
import { API_ENDPOINTS } from './constants';

// Get all models
export async function getModels(): Promise<AIModel[]> {
    const response = await apiClient.get('/models');
    return response as AIModel[];
}

// Get default model
export async function getDefaultModel(): Promise<AIModel> {
    const response = await apiClient.get('/models/default');
    return response as AIModel;
}

export async function getProviders(): Promise<APIProvider[]> {
    const response = await apiClient.get<APIProvider[]>(API_ENDPOINTS.PROVIDERS);
    return response;
}

export async function getProviderById(id: string): Promise<APIProvider> {
    const response = await apiClient.get<APIProvider>(API_ENDPOINTS.PROVIDER_BY_ID(id));
    return response;
}

export async function createProvider(data: CreateProviderRequest): Promise<CreateResponse> {
    const response = await apiClient.post<CreateResponse>(API_ENDPOINTS.PROVIDERS, data);
    return response;
}

export async function updateProvider(id: string, data: UpdateProviderRequest): Promise<UpdateResponse> {
    const response = await apiClient.put<UpdateResponse>(API_ENDPOINTS.PROVIDER_BY_ID(id), data);
    return response;
}

export async function deleteProvider(id: string): Promise<void> {
    await apiClient.delete(API_ENDPOINTS.PROVIDER_BY_ID(id));
}

// Get models with status (has API key or not)
export async function getModelsWithStatus(): Promise<ModelWithStatus[]> {
    const response = await apiClient.get('/models/with-status');
    return response as ModelWithStatus[];
}

// Get models from API keys (yang sudah punya API key aktif)
export async function getModelsFromAPIKeys(): Promise<ModelFromAPIKey[]> {
    try {
        const response = await apiClient.get('/api-keys');
        const keys = response as any[];

        // Filter hanya yang service = 'text' (buat generate artikel)
        const textKeys = keys ? keys.filter(k => k.service === 'text') : [];

        // Map ke format ModelFromAPIKey
        const models: ModelFromAPIKey[] = textKeys.map(key => ({
            id: key.id,
            modelId: key.modelId,
            name: key.modelName,
            displayName: key.modelDisplayName,
            providerName: key.providerName,
            providerDisplayName: key.providerDisplayName,
            service: key.service,
        }));

        return models;
    } catch (error) {
        console.error('Failed to get models from API keys:', error);
        return [];
    }
}