import { apiClient } from '../api';
import type { HistoryItem } from './types';

// Add to history
export async function addHistory(data: Omit<HistoryItem, 'id' | 'createdAt'>): Promise<{ id: string; message: string }> {
    return apiClient.post('/history', data);
}

// Delete history
export async function deleteHistory(id: string): Promise<void> {
    await apiClient.delete(`/history/${id}`);
}

// Clear all history
export async function clearHistory(): Promise<void> {
    await apiClient.delete('/history');
}