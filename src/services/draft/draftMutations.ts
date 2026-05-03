import { apiClient } from '../api';
import { type ScheduleRequest } from '@/types/schedule';
import type { CreateDraftRequest, UpdateDraftRequest, PublishResponse } from './types';

// Create draft
export async function createDraft(draftData: CreateDraftRequest): Promise<{ id: string; message: string }> {
    try {
        const response = await apiClient.post<{ id: string; message: string }>('/drafts', draftData);
        return response;
    } catch (error) {
        console.error('Failed to create draft:', error);
        throw error;
    }
}

// Update draft
export async function updateDraft(id: string, draftData: UpdateDraftRequest): Promise<{ message: string }> {
    try {
        const response = await apiClient.put<{ message: string }>(`/drafts/${id}`, draftData);
        return response;
    } catch (error) {
        console.error('Failed to update draft:', error);
        throw error;
    }
}

// Delete draft
export async function deleteDraft(id: string): Promise<void> {
    try {
        await apiClient.delete(`/drafts/${id}`);
    } catch (error) {
        console.error('Failed to delete draft:', error);
        throw error;
    }
}

// Publish draft
export async function publishDraft(id: string, scheduledFor?: string): Promise<PublishResponse> {
    try {
        const payload = scheduledFor ? { scheduledFor } : {};
        const response = await apiClient.post<PublishResponse>(`/drafts/${id}/publish`, payload);
        return response;
    } catch (error: any) {
        console.error('Failed to publish draft:', error);

        // If error response contains results, pass them through
        if (error?.response?.data?.results) {
            throw {
                message: error.response.data.message || 'Failed to publish draft',
                results: error.response.data.results,
                response: error.response
            };
        }

        throw error;
    }
}

export async function publishDraftInstant(draftData: CreateDraftRequest): Promise<PublishResponse> {
    try {
        const response = await apiClient.post<PublishResponse>(`/drafts/publish`, draftData);
        return response;
    } catch (error) {
        console.error('Failed to update draft:', error);
        throw error;
    }
}

export async function draftSchedule(scheduleData: ScheduleRequest): Promise<PublishResponse> {
    try {
        const response = await apiClient.post<PublishResponse>("/drafts/schedule ", scheduleData);
        return response;
    } catch (error) {
        console.error('Failed to update draft:', error);
        throw error;
    }
}