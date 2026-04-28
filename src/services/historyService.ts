// src/services/historyService.ts
import { apiClient } from './api';

export interface HistoryItem {
    id: string;
    title: string;
    topic: string;
    content: string;
    imageUrl?: string;
    targetProducts: string[];
    status: 'success' | 'failed' | 'pending';
    action: 'published' | 'scheduled' | 'draft_saved';
    errorMessage?: string;
    publishedAt: string;
    scheduledFor?: string;
    createdBy?: string;
    teamId?: string;
    createdAt: string;
}

// Get all history
export async function getHistory(params?: {
    page?: number;
    limit?: number;
    status?: string;
    action?: string;
}): Promise<HistoryItem[]> {
    return apiClient.get<HistoryItem[]>('/history', { params });
}

// Add to history
export async function addHistory(data: Omit<HistoryItem, 'id' | 'createdAt'>): Promise<{ id: string; message: string }> {
    return apiClient.post('/history', data);
}

export async function getHistoryById(id: string): Promise<HistoryItem> {
    return apiClient.get(`/history/${id}`);
}


// Delete history
export async function deleteHistory(id: string): Promise<void> {
    await apiClient.delete(`/history/${id}`);
}

// Clear all history
export async function clearHistory(): Promise<void> {
    await apiClient.delete('/history');
}