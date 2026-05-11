export type DraftStatus = 'draft' | 'scheduled' | 'published';

export interface Draft {
    id: string;
    title: string;
    topic: string;
    article: string;
    image_url?: string;
    image_prompt?: string;
    status?: DraftStatus;
    scheduled_for?: string;
    target_products: string[];
    has_image?: boolean;
    created_by?: string;
    team_id?: string;
    user_id?: string;
    created_at: string;
    updated_at: string;
}

export interface CreateDraftRequest {
    title: string;
    topic: string;
    article: string;
    image_url?: string;
    image_prompt?: string;
    target_products: string[];
    team_id?: string;
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