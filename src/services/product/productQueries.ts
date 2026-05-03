import { apiClient } from '../api';
import type { Product } from '@/types/product';

// Get all products (backend akan otomatis filter berdasarkan team dari token)
export async function getProducts(params?: {
    page?: number;
    limit?: number;
    platform?: string;
    status?: string;
    syncStatus?: string;  // ← tambahan
    teamId?: string;      // ← untuk admin filter by team
}): Promise<Product[]> {
    return apiClient.get<Product[]>('/products', { params });
}

// Get product by id
export async function getProductById(id: string): Promise<Product> {
    return apiClient.get<Product>(`/products/${id}`);
}

// Get products by team (admin only)
export async function getProductsByTeam(teamId: string): Promise<Product[]> {
    return apiClient.get<Product[]>(`/teams/${teamId}/products`);
}

// Get products by platform
export async function getProductsByPlatform(platform: string): Promise<Product[]> {
    return getProducts({ platform });
}

// Get products by status
export async function getProductsByStatus(status: string): Promise<Product[]> {
    return getProducts({ status });
}

// Get products by sync status
export async function getProductsBySyncStatus(syncStatus: string): Promise<Product[]> {
    return getProducts({ syncStatus });
}

// Get connected products
export async function getConnectedProducts(): Promise<Product[]> {
    return getProducts({ status: 'connected' });
}

// Get products that need sync
export async function getProductsNeedingSync(): Promise<Product[]> {
    return getProducts({ syncStatus: 'failed' });
}