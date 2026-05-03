import { apiClient } from '../api';
import type { Draft } from './types';

// Get all drafts
export async function getDrafts(params?: { status?: string; search?: string }): Promise<Draft[]> {
    try {
        const response = await apiClient.get<Draft[]>('/drafts', params);
        return response;
    } catch (error) {
        console.error('Failed to get drafts:', error);
        throw error;
    }
}

// Get draft by ID
export async function getDraftById(id: string): Promise<Draft | null> {
    try {
        const response = await apiClient.get<Draft>(`/drafts/${id}`);
        return response;
    } catch (error) {
        console.error('Failed to get draft:', error);
        return null;
    }
}