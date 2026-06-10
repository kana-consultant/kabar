// services/model/modelQueries.ts

import { apiClient } from '../api';
import type {
    RequestSchema, Family, AIModel, ModelWithStatus,
    AIModelSchema,
    ModelFromAPIKey,
    APIResponse,
} from '@/types/provider.types';

export { getProviders } from "@/services/provider/providerMutations";
// Get all models
export async function getModels(): Promise<APIResponse<AIModel[]>> {
    const response = await apiClient.get<APIResponse<AIModel[]>> ('/families');
    return response
}

// Get single model by ID
export async function getModel(id: string): Promise<AIModel | null> {
    const response = await apiClient.get<AIModel | null>(`/models/${id}`);
    return response
}

export async function getModelSchema(id: string): Promise<AIModelSchema | null> {
    const response = await apiClient.get<AIModelSchema | null>(`/models/${id}/schema`);
    return response
}


// Get default model
export async function getDefaultModel(): Promise<AIModel> {
    const response = await apiClient.get('/models/default');
    return response as AIModel;
}


// Get families
export async function getFamilies(): Promise<APIResponse<Family[]>> {
    const response = await apiClient.get<APIResponse<Family[]>>('/families');
    return response;
}

// Get schemas
export async function getSchemas(): Promise<APIResponse<RequestSchema[]>> {
    const response = await apiClient.get<APIResponse<RequestSchema[]>>('/schemas');
    return response;
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
            model_id: key.modelId,
            name: key.model_name,
            display_name: key.model_display_name,
            provider_name: key.provider_name,
            provider_display_name: key.provider_display_name,
            service: key.service,
        }));

        return models;
    } catch (error) {
        console.error('Failed to get models from API keys:', error);
        return [];
    }
}