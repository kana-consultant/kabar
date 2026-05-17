// src/services/dashboardService.ts
import { apiClient } from './api';

export interface DashboardStats {
    totalContent: number;
    totalProducts: number;
    totalPublished: number;
    averageSeoScore: number;
    contentChange: string;
    productsChange: string;
    publishedPercentage: number;
    seoScoreChange: string;
}

export async function getDashboardStats(): Promise<DashboardStats> {
    try {
        const response = await apiClient.get<DashboardStats>('/dashboard/stats');
        return response;
    } catch (error) {
        console.error('Failed to get dashboard stats:', error);
        throw error;
    }
}