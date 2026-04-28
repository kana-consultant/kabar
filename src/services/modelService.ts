// src/services/modelService.ts
import { apiClient } from './api';

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

// Create model
export async function createModel(data: CreateModelRequest): Promise<{ id: string; message: string }> {
    const response = await apiClient.post('/models', data);
    return response as { id: string; message: string };
}

// Update model
export async function updateModel(id: string, data: Partial<AIModel>): Promise<{ message: string }> {
    const response = await apiClient.put(`/models/${id}`, data);
    return response as { message: string };
}

// Delete model
export async function deleteModel(id: string): Promise<void> {
    await apiClient.delete(`/models/${id}`);
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