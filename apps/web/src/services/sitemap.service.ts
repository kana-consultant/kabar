// services/sitemap.service.ts
import { apiClient } from "./api";

export interface GenerateSitemapRequest {
    productId: string;
    baseURL: string;
    includeImages: boolean;
    limit?: number;
}

export interface GenerateSitemapResponse {
    sitemapXML: string;
    totalURLs: number;
    generatedAt: string;
    productId: string;
    baseURL: string;
}

export interface SitemapHistoryItem {
    id: string;
    title: string;
    totalURLs: number;
    status: "success" | "failed" | "pending";
    createdAt: string;
    sitemapURL?: string;
    productId?: string;
    productName?: string;
    baseURL?: string;
}

export const sitemapService = {
    // Generate sitemap - menggunakan getXML untuk handle XML response
    generateSitemap: async (
        params: GenerateSitemapRequest
    ): Promise<GenerateSitemapResponse> => {
        const response = await apiClient.getXML("/sitemap", {
            product_id: params.productId,
            base_url: params.baseURL,
            include_images: params.includeImages,
            limit: params.limit || 0,
        });

        return {
            sitemapXML: response.data,
            totalURLs: parseInt(response.headers["x-total-urls"] || "0"),
            generatedAt: response.headers["x-generated-at"] || new Date().toISOString(),
            productId: response.headers["x-product-id"] || params.productId,
            baseURL: response.headers["x-base-url"] || params.baseURL,
        };
    },

    // Download sitemap as file
    downloadSitemap: (xmlContent: string, filename: string = "sitemap.xml") => {
        const blob = new Blob([xmlContent], { type: "application/xml" });
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.href = url;
        link.download = filename;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(url);
    },

    // Get sitemap history
    getSitemapHistory: async (): Promise<SitemapHistoryItem[]> => {
        return apiClient.get("/sitemap/history");
    },
};