import { apiClient } from '../api';
import type { CreateModelRequest, AIModel } from './types';

// Create model
export async function createModel(data: CreateModelRequest): Promise<{ id: string; message: string }> {
    const response = await apiClient.post('/models', data);
    return response as { id: string; message: string };
}

// Update model
export async function updateModel(id: string, data: Partial<AIModel>): Promise<{ message: string }> {
    const response = await apiClient.put(`/models/${id}`, data);
    return response as { message: string };
}

// Delete model
export async function deleteModel(id: string): Promise<void> {
    await apiClient.delete(`/models/${id}`);
}