export type Tone = 'professional' | 'casual' | 'friendly' | 'formal';
export type ArticleLength = 'short' | 'medium' | 'long';
export type Language = 'id' | 'en';

export interface GenerateArticleRequest {
    topic: string;
    tone?: Tone;
    length?: ArticleLength;
    autoGenerateImage? : boolean;
    language?: Language;
    modelId?: string;
    imageModelId? : string;
}

export interface GenerateArticleResponse {
    title: string;
    content: string;
    excerpt: string;
    keywords: string[];
    imagePrompt: string;
    imageUrl: string;
    wordCount: number;
    readabilityScore: number;
    seoScore: number;
    slug: string;
}

export interface GenerateImageRequest {
    prompt: string;
    modelId: string;
}

export interface GenerateImageResponse {
    imageUrl: string;
    prompt: string;
    generatedAt: string;
    model: string;
}