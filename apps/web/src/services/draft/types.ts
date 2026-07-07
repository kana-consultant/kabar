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
    excerpt : string;
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
    excerpt : string
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
    excerpt : string
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
    // Basic metrics
    total_draft: number;
    total_with_image: number;
    total_without_image: number;
    total_scheduled: number;
    total_published: number;
    total_with_keywords: number;
    total_with_seo: number;
    
    // Derived metrics
    completion_rate: number;
    scheduled_rate: number;
    image_coverage_rate: number;
    seo_score_avg: number;
    keywords_avg_count: number;
    
    // Breakdowns
    status_breakdown: Record<string, number>;
    product_coverage: Record<string, number>;
    product_status: Record<string, number>;
    topic_breakdown: Record<string, number>;
    seo_score_distribution: Record<string, number>;
    
    // Time series
    daily_activity: DailyActivity[];
    weekly_trend?: WeeklyTrend[];
    scheduled_upcoming?: ScheduledItem[];
    
    // Content quality metrics
    top_topics?: TopicStats[];
    top_keywords?: KeywordStats[];
    
    // Performance metrics
    average_completion_time: number;
    cache_metadata?: CacheMetadata;
}

export interface DailyActivity {
    date: string;
    count: number;
    scheduled: number;
    published: number;
    with_image: number;
    with_keywords: number;
    avg_seo_score: number;
}

export interface WeeklyTrend {
    week: string;
    created: number;
    scheduled: number;
    published: number;
}

export interface TopicStats {
    topic: string;
    count: number;
    avg_seo_score: number;
}

export interface KeywordStats {
    keyword: string;
    count: number;
}

export interface ScheduledItem {
    id: string;
    title: string;
    scheduled_for: string; // ISO date string
    products?: string[];
}

export interface CacheMetadata {
    cached_at: string; // ISO date string
    ttl: string;
    generation_time_ms: number;
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