// services/family.ts

import type { Family, FamilyFormData } from "@/types/provider.types";
import { apiClient } from "../api";

export async function createFamily(data: FamilyFormData): Promise<Family> {
    return apiClient.post<Family>('/families', data);
}

export async function updateFamily(id: string, data: Partial<FamilyFormData>): Promise<Family> {
    return apiClient.put<Family>(`/families/${id}`, data);
}

export async function deleteFamily(id: string): Promise<void> {
    return apiClient.delete(`/families/${id}`);
}

export async function getFamiliesByProvider(providerId: string): Promise<Family[]> {
    return apiClient.get(`/families/providers/${providerId}/families`);
}