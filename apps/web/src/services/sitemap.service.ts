// services/sitemap.service.ts
import { apiClient } from "./api";

export interface GenerateSitemapRequest {
  baseURL: string;
  includeImages: boolean;
  limit?: number;
}

export interface GenerateSitemapResponse {
  sitemapXML: string;
  totalURLs: number;
  generatedAt: string;
}

export interface SitemapHistory {
  id: string;
  title: string;
  totalURLs: number;
  status: "success" | "failed" | "pending";
  createdAt: string;
  sitemapURL?: string;
}

export const sitemapService = {
  // Generate sitemap
  generateSitemap: async (
    params: GenerateSitemapRequest
  ): Promise<GenerateSitemapResponse> => {
    const response = await apiClient.get("/sitemap", {
      params: {
        base_url: params.baseURL,
        include_images: params.includeImages,
        limit: params.limit || 0,
      },
    });

    // Response is XML, not JSON
    return {
      sitemapXML: response.data,
      totalURLs: parseInt(response.headers["x-total-urls"] || "0"),
      generatedAt: response.headers["x-generated-at"] || new Date().toISOString(),
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

  // Get sitemap history (optional)
  getSitemapHistory: async (): Promise<SitemapHistory[]> => {
    const response = await apiClient.get("/sitemap/history");
    return response.data;
  },
};