import { apiClient } from '../api';
import type { Product } from '@/types/product';
import type { CreateProductRequest, UpdateProductRequest, AddProductResponse } from './types';
import { getProductById } from './productQueries';

// Create product (basic)
export async function createProduct(req: CreateProductRequest): Promise<{ id: string; message: string }> {
    return apiClient.post('/products', req);
}

/**
 * Add new product with full response
 */
export async function addProduct(req: CreateProductRequest): Promise<AddProductResponse> {
    const response = await apiClient.post<{ id: string; message: string }>('/products', req);
    
    // Fetch product yang baru dibuat untuk data lengkap
    let fullProduct: Product | undefined;
    try {
        fullProduct = await getProductById(response.id);
    } catch (error) {
        console.warn('Failed to fetch created product:', error);
    }
    
    return {
        id: response.id,
        message: response.message,
        product: fullProduct || {
            id: response.id,
            name: req.name,
            platform: req.platform as any,
            apiEndpoint: req.apiEndpoint,
            status: 'pending',
            syncStatus: 'idle',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            apiKey: req.apiKey,
            teamId: req.teamId,
            adapterConfig: req.adapterConfig ? {
                endpointPath: req.adapterConfig.endpointPath,
                httpMethod: req.adapterConfig.httpMethod,
                customHeaders: req.adapterConfig.customHeaders,
                fieldMapping: req.adapterConfig.fieldMapping,
                timeoutSeconds: req.adapterConfig.timeoutSeconds || 30,
                retryCount: req.adapterConfig.retryCount || 3,
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString(),
            } : undefined,
        },
    };
}

/**
 * Save product (create or update)
 */
export async function saveProduct(product: Partial<Product> & { id?: string }): Promise<Product> {
    if (product.id) {
        await updateProduct(product.id, {
            name: product.name,
            platform: product.platform,
            apiEndpoint: product.apiEndpoint,
            status: product.status,
            syncStatus: product.syncStatus,
            adapterConfig: product.adapterConfig,
        });
        return getProductById(product.id);
    } else {
        const result = await addProduct({
            name: product.name || '',
            platform: product.platform || 'custom',
            apiEndpoint: product.apiEndpoint || '',
            apiKey: product.apiKey,
            teamId: product.teamId,
            adapterConfig: product.adapterConfig ? {
                endpointPath: product.adapterConfig.endpointPath,
                httpMethod: product.adapterConfig.httpMethod,
                customHeaders: product.adapterConfig.customHeaders,
                fieldMapping: product.adapterConfig.fieldMapping,
                timeoutSeconds: product.adapterConfig.timeoutSeconds,
                retryCount: product.adapterConfig.retryCount,
            } : undefined,
        });
        return result.product;
    }
}

// Update product
export async function updateProduct(id: string, updates: UpdateProductRequest): Promise<void> {
    await apiClient.put(`/products/${id}`, updates);
}

// Delete product
export async function deleteProduct(id: string): Promise<void> {
    await apiClient.delete(`/products/${id}`);
}

// Test connection
export async function testConnection(id: string): Promise<{ success: boolean; message?: string }> {
    return apiClient.post(`/products/${id}/test`);
}

// Sync product (force sync)
export async function syncProduct(id: string): Promise<{ success: boolean; message: string }> {
    return apiClient.post(`/products/${id}/sync`);
}