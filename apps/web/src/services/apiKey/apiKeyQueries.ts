import { apiClient } from '../api';
import type { APIKey, APIKeyDetail } from './types';

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