import { apiClient } from '../api';
import type { HistoryItem } from './types';
import type { PaginatedResponse, PaginationParams } from './types';

export async function getHistory(params?: PaginationParams & {
    status?: string;
    action?: string;
}): Promise<PaginatedResponse<HistoryItem>> {
    return apiClient.get<PaginatedResponse<HistoryItem>>('/history', params);
}

export async function getHistoryById(id: string): Promise<HistoryItem> {
    return apiClient.get(`/history/${id}`);
}