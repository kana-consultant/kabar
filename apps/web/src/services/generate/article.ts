import { apiClient } from '../api';
import type { GenerateArticleRequest, GenerateArticleResponse } from './types';

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