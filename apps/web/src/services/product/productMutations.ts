import { apiClient } from '../api';
import type { Product } from '@/services/product';
import type { ProductRequest, UpdateProductRequest, AddProductResponse } from './types';
import { getProductById } from './productQueries';

// Create product (basic)
export async function createProduct(req: ProductRequest): Promise<{ id: string; message: string }> {
    return apiClient.post('/products', req);
}

/**
 * Add new product with full response
 */
export async function addProduct(req: ProductRequest): Promise<AddProductResponse> {
    console.log("payload :",req)
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
    };
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