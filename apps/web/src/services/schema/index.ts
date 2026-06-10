// services/schema.ts

import type { RequestSchema, RequestSchemaFormData } from "@/types/provider.types";
import { apiClient } from "../api";


export async function createSchema(data: RequestSchemaFormData): Promise<RequestSchema> {
    return apiClient.post<RequestSchema>('/schemas', data);
}

export async function updateSchema(id: string, data: Partial<RequestSchemaFormData>): Promise<RequestSchema> {
    return apiClient.put<RequestSchema>(`/schemas/${id}`, data);
}

export async function deleteSchema(id: string): Promise<void> {
    return apiClient.delete(`/schemas/${id}`);
}