export interface Draft {
    id: string;
    title: string;
    topic: string;
    article: string;
    imageUrl?: string;
    imagePrompt?: string;
    createdAt: string;
    updatedAt: string;
    scheduledFor?: string;
    status: 'draft' | 'scheduled' | 'published';
    targetProducts: string[];
    hasImage: boolean;
}