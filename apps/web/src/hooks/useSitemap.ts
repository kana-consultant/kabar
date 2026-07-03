// hooks/useSitemap.ts
import { useState, useCallback } from "react";
import { sitemapService, type GenerateSitemapRequest, type GenerateSitemapResponse, type SitemapHistory } from "@/services/sitemap.service";

interface UseSitemapReturn {
  // State
  isLoading: boolean;
  error: string | null;
  sitemapData: GenerateSitemapResponse | null;
  history: SitemapHistory[];
  isHistoryLoading: boolean;

  // Actions
  generateSitemap: (params: GenerateSitemapRequest) => Promise<void>;
  downloadSitemap: () => void;
  fetchHistory: () => Promise<void>;
  clearError: () => void;
  reset: () => void;
}

export function useSitemap(): UseSitemapReturn {
  const [isLoading, setIsLoading] = useState(false);
  const [isHistoryLoading, setIsHistoryLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [sitemapData, setSitemapData] = useState<GenerateSitemapResponse | null>(null);
  const [history, setHistory] = useState<SitemapHistory[]>([]);

  const generateSitemap = useCallback(async (params: GenerateSitemapRequest) => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await sitemapService.generateSitemap(params);
      setSitemapData(data);
    } catch (err: any) {
      setError(err.message || "Failed to generate sitemap");
      setSitemapData(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const downloadSitemap = useCallback(() => {
    if (sitemapData?.sitemapXML) {
      sitemapService.downloadSitemap(sitemapData.sitemapXML);
    }
  }, [sitemapData]);

  const fetchHistory = useCallback(async () => {
    setIsHistoryLoading(true);
    setError(null);

    try {
      const data = await sitemapService.getSitemapHistory();
      setHistory(data);
    } catch (err: any) {
      setError(err.message || "Failed to fetch history");
    } finally {
      setIsHistoryLoading(false);
    }
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  const reset = useCallback(() => {
    setSitemapData(null);
    setError(null);
    setIsLoading(false);
  }, []);

  return {
    isLoading,
    isHistoryLoading,
    error,
    sitemapData,
    history,
    generateSitemap,
    downloadSitemap,
    fetchHistory,
    clearError,
    reset,
  };
}