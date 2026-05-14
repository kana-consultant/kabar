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
}