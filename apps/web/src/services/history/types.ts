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
    keywords : string[];
}

// services/types.ts
export interface PaginationParams {
    limit?: number;
    offset?: number;
}


export interface PaginatedResponse<T> {
    data: T[];
    current_page: number;
    total_pages: number;
    total_items: number;
    total_success: number;
    total_failed: number;
    limit: number;
    offset: number;
}