// src/services/apiKeyService.ts
import { apiClient } from './api';

export interface APIKey {
    id: string;
    service: string;
    systemPrompt?: string;
    teamId?: string;
    isActive: boolean;
    createdAt: string;
    providerId : string;
    updatedAt: string;
    modelId : string;
    
}

export interface APIKeyDetail {
		id                  : string    
		service             : string    
		providerId          : string    
		modelId             : string    
		isActive            : boolean     
		systemPrompt        : string   
		createdBy           : string   
		createdAt           : string
		updatedAt           : string 
		providerName        : string    
		providerDisplayName  : string   
		modelName           : string  
		modelDisplayName    : string   
	}

export interface CreateAPIKeyRequest {
    service: string;
    key: string;
    modelId : String;
    providerId:string;
    systemPrompt?: string;
    teamId?: string;
}

export interface UpdateAPIKeyRequest {
    service?: string;
    key?: string;
    systemPrompt?: string;
    isActive?: boolean;
}

// Get all API keys
export async function getAPIKeys(): Promise<APIKey[]> {
    try {
        const response = await apiClient.get<APIKey[]>('/api-keys');
        return response;
    } catch (error) {
        console.error('Failed to get API keys:', error);
        throw error;
    }
}

// Get API key by service
export async function getAPIKeyByService(service: string): Promise<APIKey | null> {
    try {
        const keys = await getAPIKeys();
        return keys.find(k => k.service === service) || null;
    } catch (error) {
        console.error('Failed to get API key by service:', error);
        return null;
    }
}

// Create API key
export async function createAPIKey(data: CreateAPIKeyRequest): Promise<{ id: string; message: string }> {
    try {
        const response = await apiClient.post<{ id: string; message: string }>('/api-keys', data);
        return response;
    } catch (error) {
        console.error('Failed to create API key:', error);
        throw error;
    }
}

// Update API key
export async function updateAPIKey(id: string, data: UpdateAPIKeyRequest): Promise<{ message: string }> {
    try {
        const response = await apiClient.put<{ message: string }>(`/api-keys/${id}`, data);
        return response;
    } catch (error) {
        console.error('Failed to update API key:', error);
        throw error;
    }
}

// Delete API key
export async function deleteAPIKey(id: string): Promise<void> {
    try {
        await apiClient.delete(`/api-keys/${id}`);
    } catch (error) {
        console.error('Failed to delete API key:', error);
        throw error;
    }
}

// Toggle API key active status
export async function toggleAPIKey(id: string, isActive: boolean): Promise<{ message: string }> {
    try {
        const response = await apiClient.put<{ message: string }>(`/api-keys/${id}`, { isActive });
        return response;
    } catch (error) {
        console.error('Failed to toggle API key:', error);
        throw error;
    }
}