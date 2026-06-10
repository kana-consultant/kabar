import type { ModelFormData } from '@/pages/ai-management/Model/provider.types';
import { apiClient } from '../api';
import type { CreateModelRequest,AIModel } from '@/types/provider.types';

// Create model
// services/model/modelMutations.ts

export async function createModel(data: CreateModelRequest): Promise<AIModel> {
    const response = await apiClient.post('/models', data);
    return response as AIModel;  // Sama seperti updateModel
}

export async function updateModel(id: string, data: Partial<ModelFormData>): Promise<AIModel> {
    const response = await apiClient.put(`/models/${id}`, data);
    return response as AIModel;
}

// Delete model
export async function deleteModel(id: string): Promise<void> {
    await apiClient.delete(`/models/${id}`);
}