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