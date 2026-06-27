import { apiClient } from '../api';
import type { Draft } from './types';
import type { SimilarityResult, SEOScore } from './types';

export interface DraftStats {
    total_draft: number;
    total_with_image: number;
    total_without_image: number;
    total_scheduled: number;
    product_coverage: Record<string, number>;
    daily_activity: {
        date: string;
        count: number;
    }[];
}

interface PaginationParams {
    page?: number;
    limit?: number;
    offset?: number;
    status?: string;
    search?: string;
}

export interface PaginatedResponse<T> {
    data: T[];
    current_page: number;
    total_pages: number;
    total_items: number;
    limit: number;
    offset: number;
}

export interface GetDraftsResponse {
    data: dataResponse;
    
}

export interface dataResponse {
    drafts: PaginatedResponse<Draft>;
    stats: DraftStats;
}

function buildPaginationParams(params: PaginationParams): { limit: number; offset: number } {
    const limit = params.limit ?? 10;
    const offset = params.offset ?? ((params.page ?? 1) - 1) * limit;
    return { limit, offset };
}

export async function getDrafts(params?: PaginationParams): Promise<GetDraftsResponse> {
    try {
        const { limit, offset } = buildPaginationParams(params ?? {});
        const response = await apiClient.get<GetDraftsResponse>('/drafts', {
            ...params,
            limit,
            offset,
        });
      
        return response;
    } catch (error) {
        console.error('Failed to get drafts:', error);
        throw error;
    }
}

export async function getScheduled(params?: PaginationParams): Promise<PaginatedResponse<Draft>> {
    try {
        const { limit, offset } = buildPaginationParams(params ?? {});
        const response = await apiClient.get<PaginatedResponse<Draft>>('/drafts/scheduled', {
            ...params,
            limit,
            offset,
        });

        console.log(response)
        return response;
    } catch (error) {
        console.error('Failed to get scheduled drafts:', error);
        throw error;
    }
}

export async function getDraftById(id: string): Promise<Draft | null> {
    try {
        const response = await apiClient.get<Draft>(`/drafts/${id}`);
        return response;
    } catch (error) {
        console.error('Failed to get draft:', error);
        return null;
    }
}

export async function checkSimilarity(id: string): Promise<SimilarityResult | null> {
    try {
        const response = await apiClient.get<SimilarityResult>(`/drafts/${id}/check-similarity`);
        return response;
    } catch (error) {
        console.error('Failed to get draft:', error);
        return null;
    }
}

export async function getSeoScore(id: string): Promise<SEOScore | null> {
    try {
        const response = await apiClient.get<SEOScore>(`/drafts/${id}/seo-score`);
        return response;
    } catch (error) {
        console.error('Failed to get draft:', error);
        return null;
    }
}
