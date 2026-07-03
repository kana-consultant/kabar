// hooks/useSitemap.ts
import { useState, useCallback, useEffect } from "react";
import { sitemapService,type GenerateSitemapRequest, type GenerateSitemapResponse, type SitemapHistoryItem } from "@/services/sitemap.service";
import  { type Product, getProducts } from "@/services/product";

interface UseSitemapReturn {
    // State
    isLoading: boolean;
    isHistoryLoading: boolean;
    isProductsLoading: boolean;
    error: string | null;
    sitemapData: GenerateSitemapResponse | null;
    history: SitemapHistoryItem[];
    products: Product[];
    
    // Actions
    generateSitemap: (params: GenerateSitemapRequest) => Promise<void>;
    downloadSitemap: () => void;
    fetchHistory: () => Promise<void>;
    fetchProducts: () => Promise<void>;
    clearError: () => void;
    reset: () => void;
}

export function useSitemap(): UseSitemapReturn {
    const [isLoading, setIsLoading] = useState(false);
    const [isHistoryLoading, setIsHistoryLoading] = useState(false);
    const [isProductsLoading, setIsProductsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [sitemapData, setSitemapData] = useState<GenerateSitemapResponse | null>(null);
    const [history, setHistory] = useState<SitemapHistoryItem[]>([]);
    const [products, setProducts] = useState<Product[]>([]);

    const generateSitemap = useCallback(async (params: GenerateSitemapRequest) => {
        setIsLoading(true);
        setError(null);

        try {
            const data = await sitemapService.generateSitemap(params);
            setSitemapData(data);
            await fetchHistory();
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

    const fetchProducts = useCallback(async () => {
        setIsProductsLoading(true);
        setError(null);

        try {
            const data = await getProducts();
            setProducts(data);
        } catch (err: any) {
            setError(err.message || "Failed to fetch products");
        } finally {
            setIsProductsLoading(false);
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

    useEffect(() => {
        fetchProducts();
        fetchHistory();
    }, [fetchProducts, fetchHistory]);

    return {
        isLoading,
        isHistoryLoading,
        isProductsLoading,
        error,
        sitemapData,
        history,
        products,
        generateSitemap,
        downloadSitemap,
        fetchHistory,
        fetchProducts,
        clearError,
        reset,
    };
}