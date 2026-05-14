import { apiClient } from '../api';
import type { AIModel, ModelWithStatus, ModelFromAPIKey } from './types';

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