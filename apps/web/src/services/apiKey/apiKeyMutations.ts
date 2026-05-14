import { apiClient } from '../api';
import type { CreateAPIKeyRequest, UpdateAPIKeyRequest } from './types';

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