// src/services/modelService.ts
import { apiClient } from './api';

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

export interface AIModel {
    id: string;
    name: string;
    providerId: string;
    provider?: APIProvider;
    displayName: string;
    description: string;
    isActive: boolean;
    isDefault: boolean;
    maxTokens: number;
    temperature: number;
    teamId?: string;
    createdAt: string;
    updatedAt: string;
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

// Get all models
export async function getModels(): Promise<AIModel[]> {
    const response = await apiClient.get('/models');
    return response as AIModel[];
}

// Get all providers
export async function getProviders(): Promise<APIProvider[]> {
    const response = await apiClient.get('/providers');
    return response as APIProvider[];
}

// Create provider
export async function createProvider(data: CreateProviderRequest): Promise<{ id: string; message: string }> {
    const response = await apiClient.post('/providers', data);
    return response as { id: string; message: string };
}

// Update provider
export async function updateProvider(id: string, data: Partial<CreateProviderRequest>): Promise<{ message: string }> {
    const response = await apiClient.put(`/providers/${id}`, data);
    return response as { message: string };
}

// Delete provider
export async function deleteProvider(id: string): Promise<void> {
    await apiClient.delete(`/providers/${id}`);
}

// Create model
export async function createModel(data: CreateModelRequest): Promise<{ id: string; message: string }> {
    const response = await apiClient.post('/models', data);
    return response as { id: string; message: string };
}

// Update model
export async function updateModel(id: string, data: Partial<CreateModelRequest>): Promise<{ message: string }> {
    const response = await apiClient.put(`/models/${id}`, data);
    return response as { message: string };
}

// Delete model
export async function deleteModel(id: string): Promise<void> {
    await apiClient.delete(`/models/${id}`);
}

// Get default model
export async function getDefaultModel(): Promise<AIModel> {
    const response = await apiClient.get('/models/default');
    return response as AIModel;
}