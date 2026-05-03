import type { Product, AdapterConfig } from '@/types/product';

export interface CreateProductRequest {
    name: string;
    platform: string;
    apiEndpoint: string;
    apiKey?: string;
    teamId?: string;  // ← tambahan untuk assign ke team
    adapterConfig?: {
        endpointPath: string;
        httpMethod: 'POST' | 'PUT' | 'PATCH';
        customHeaders: Record<string, string>;
        fieldMapping: string;
        timeoutSeconds?: number;
        retryCount?: number;
    };
}

export interface UpdateProductRequest {
    name?: string;
    platform?: string;
    apiEndpoint?: string;
    apiKey?: string;
    status?: string;
    syncStatus?: string;  // ← tambahan
    teamId?: string;      // ← tambahan
    adapterConfig?: Partial<AdapterConfig>;
}

export interface AddProductResponse {
    id: string;
    message: string;
    product: Product;
}