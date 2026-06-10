// /services/provider/providerMutations.ts

import type { APIProvider, APIResponse,ProviderFormData,
    TestConnectionConfig,
    TestConnectionResult,
    ProviderTestResult } from '@/types/provider.types';
import { apiClient } from '../api';


export async function getProviders(params?: { isActive?: boolean }): Promise<APIResponse<APIProvider[]>> {
    const response = await apiClient.get<APIResponse<APIProvider[]>>(`/providers`, { params });
    return response;
}

export async function createProvider(data: ProviderFormData): Promise<APIProvider> {
    return apiClient.post<APIProvider>('/providers', data);
}

export async function updateProvider(id: string, data: ProviderFormData): Promise<APIProvider> {
    return apiClient.put<APIProvider>(`/providers/${id}`, data);
}

export async function deleteProvider(id: string): Promise<void> {
    await apiClient.delete(`/providers/${id}`);
}

export async function toggleProviderStatus(id: string, isActive: boolean): Promise<APIProvider> {
    return apiClient.patch<APIProvider>(`/providers/${id}/status`, { is_active: isActive });
}

// Get provider by id
export async function getProviderById(id: string): Promise<APIProvider> {
    return apiClient.get<APIProvider>(`/providers/${id}`);
}

// Get providers by team (admin only)
export async function getProvidersByTeam(teamId: string): Promise<APIProvider[]> {
    return apiClient.get<APIProvider[]>(`/teams/${teamId}/providers`);
}

// Get active providers only
export async function getActiveProviders(): Promise<APIResponse<APIProvider[]>> {
    return await getProviders({ isActive: true });
}


/**
 * Test koneksi provider yang sudah tersimpan (by ID)
 */
export async function testProviderConnection(id: string): Promise<TestConnectionResult> {
    const startTime = Date.now();
    try {
        const data = await apiClient.get<TestConnectionResult>(`/providers/${id}/test`);
        return {
            ...data,
            success: true,
            message: "Connection successful",
            latency: Date.now() - startTime,
        };
    } catch (error: any) {
        return {
            success: false,
            message: error?.response?.data?.message || "Connection test failed",
            status_code: error?.response?.status,
            latency: Date.now() - startTime,
            error: error instanceof Error ? error.message : "Unknown error",
        };
    }
}

/**
 * Test koneksi dari konfigurasi form (sebelum disimpan)
 */
export async function testConnection(config: TestConnectionConfig): Promise<TestConnectionResult> {
    const startTime = Date.now();
    try {
        const data = await apiClient.post<TestConnectionResult>('/test-connection', config);
        return {
            ...data,
            success: true,
            message: "Connection successful",
            latency: Date.now() - startTime,
        };
    } catch (error: any) {
        return {
            success: false,
            message: error?.response?.data?.message || "Connection test failed",
            status_code: error?.response?.status,
            latency: Date.now() - startTime,
            error: error instanceof Error ? error.message : "Unknown error",
        };
    }
}

/**
 * Test semua families milik sebuah provider
 */
export async function testProviderFamilies(id: string): Promise<ProviderTestResult> {
    return apiClient.post<ProviderTestResult>(`/providers/${id}/test-families`, {});
}

/**
 * Test koneksi satu family spesifik milik sebuah provider
 */
export async function testFamilyConnection(
    providerId: string,
    familyId: string,
    testData?: Record<string, unknown>
): Promise<TestConnectionResult> {
    const startTime = Date.now();
    try {
        const data = await apiClient.post<TestConnectionResult>(
            `/providers/${providerId}/families/${familyId}/test`,
            { test_data: testData }
        );
        return {
            ...data,
            success: true,
            message: "Family test successful",
            latency: Date.now() - startTime,
        };
    } catch (error: any) {
        return {
            success: false,
            message: error?.response?.data?.message || "Family test failed",
            status_code: error?.response?.status,
            latency: Date.now() - startTime,
            error: error instanceof Error ? error.message : "Unknown error",
        };
    }
}