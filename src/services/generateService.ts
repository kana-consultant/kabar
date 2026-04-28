// src/services/generateService.ts
import { apiClient } from './api';

export type Tone = 'professional' | 'casual' | 'friendly' | 'formal';
export type ArticleLength = 'short' | 'medium' | 'long';
export type Language = 'id' | 'en';

export interface GenerateArticleRequest {
    topic: string;
    tone?: Tone;
    length?: ArticleLength;
    keywords?: string[];
    language?: Language;
    modelId?: string;
    autoGenerateImage : boolean
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
}

export interface GenerateImageRequest {
    prompt: string;
    modelId : string;
}

export interface GenerateImageResponse {
    imageUrl: string;
    prompt: string;
    generatedAt: string;
    model: string;
}

// Generate article
export async function generateArticle(req: GenerateArticleRequest): Promise<GenerateArticleResponse> {
    try {
        const response = await apiClient.post('/generate/article', req);
        return response as GenerateArticleResponse;
    } catch (error) {
        console.error('Generate article error:', error);
        throw error;
    }
}

// Generate image
export async function generateImage(req: GenerateImageRequest): Promise<GenerateImageResponse> {
    try {
        const response = await apiClient.post('/generate/image', req);
        return response as GenerateImageResponse;
    } catch (error) {
        console.error('Generate image error:', error);
        throw error;
    }
}