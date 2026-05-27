export type DraftStatus = 'draft' | 'scheduled' | 'published';

export interface Draft {
    id?: string;
    title: string;
    topic: string;
    article: string;
    image_url?: string;
    image_prompt?: string;
    status?: DraftStatus;
    scheduled_for?: string;
    target_products?: string[];
    has_image?: boolean;
    created_by?: string;
    team_id?: string;
    user_id?: string;
    created_at?: string;
    updated_at?: string;
    slug : string;
    keywords : string[] | null;
}

export interface CreateDraftRequest {
    title: string;
    topic: string;
    article: string;
    image_url?: string;
    image_prompt?: string;
    target_products: string[];
    team_id?: string;
    slug : string;
    keywords : string[] | null
}

export interface UpdateDraftRequest {
    title?: string;
    topic?: string;
    article?: string;
    image_url?: string;
    image_prompt?: string;
    status?: DraftStatus;
    scheduled_for?: string;
    target_products?: string[];
    has_image?: boolean;
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
export interface SEOScore {
    total: number;
    details: Record<string, number>;
    suggestions: string[];
}

export interface SimilarityResult {
    draft_id: string;
    similar_drafts: SimilarDraft[];
    total: number;
}

export interface SimilarDraft {
    draft_id: string;
    title: string;
    similarity: number;
}