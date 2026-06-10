// pages/admin/ai-management/services/provider-api.ts

import { apiClient } from '@/services/api';
import type { APIProvider } from '@/types/provider.types';
import type { ProviderFormData } from "@/types/provider.types";
export { getProviders,getProviderById,updateProvider,
    deleteProvider,toggleProviderStatus,testProviderConnection,
    testConnection,testProviderFamilies,testFamilyConnection } from "@/services/provider/providerMutations";

export async function createProvider(data: ProviderFormData): Promise<APIProvider> {
    return apiClient.post<APIProvider>('/providers', data);
}

