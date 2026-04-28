// src/services/draftService.ts
import { apiClient } from './api';
import { type ScheduleRequest } from '@/types/schedule';

export type DraftStatus = 'draft' | 'scheduled' | 'published';

export interface Draft {
    id: string;
    title: string;
    topic: string;
    article: string;
    imageUrl?: string;
    imagePrompt?: string;
    status: DraftStatus;
    scheduledFor?: string;
    targetProducts: string[];
    hasImage: boolean;
    createdBy?: string;
    teamId?: string;
    userId?: string;
    createdAt: string;
    updatedAt: string;
}

export interface CreateDraftRequest {
    Title: string;
    Topic: string;
    Article: string;
    ImageURL?: string;
    ImagePrompt?: string;
    TargetProducts: string[];
    teamId?: string;
}



export interface UpdateDraftRequest {
    title?: string;
    topic?: string;
    article?: string;
    imageUrl?: string;
    imagePrompt?: string;
    status?: DraftStatus;
    scheduledFor?: string;
    targetProducts?: string[];
    hasImage?: boolean;
}

export interface PublishResponse {
    message: string;
    status?: string;
    results?: Array<{
        product: string;
        success: boolean;
        error?: string;
        response?: string;
        product_id?: string;
    }>;
}

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

export async function draftSchedule(scheduleData: ScheduleRequest) : Promise<PublishResponse> {
    try {
        const response = await apiClient.post<PublishResponse>("/drafts/schedule ", scheduleData);
        return response;
    } catch (error) {
        console.error('Failed to update draft:', error);
        throw error;
    }
}