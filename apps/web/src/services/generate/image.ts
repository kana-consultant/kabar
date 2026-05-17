import { apiClient } from '../api';
import type { GenerateImageRequest, GenerateImageResponse } from './types';

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