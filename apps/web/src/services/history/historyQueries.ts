import { apiClient } from '../api';
import type { HistoryItem } from './types';

// Get all history
export async function getHistory(params?: {
    page?: number;
    limit?: number;
    status?: string;
    action?: string;
}): Promise<HistoryItem[]> {
    return apiClient.get<HistoryItem[]>('/history', { params });
}

export async function getHistoryById(id: string): Promise<HistoryItem> {
    return apiClient.get(`/history/${id}`);
}