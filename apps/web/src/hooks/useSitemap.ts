// hooks/useSitemap.ts
import { useState, useCallback, useEffect } from "react";
import { sitemapService,type GenerateSitemapRequest, type GenerateSitemapResponse } from "@/services/sitemap.service";
import  { type Product, getProducts } from "@/services/product";

interface UseSitemapReturn {
    // State
    isLoading: boolean;
    isProductsLoading: boolean;
    error: string | null;
    sitemapData: GenerateSitemapResponse | null;
    products: Product[];
    
    // Actions
    generateSitemap: (params: GenerateSitemapRequest) => Promise<void>;
    downloadSitemap: () => void;
    clearError: () => void;
    reset: () => void;
}

export function useSitemap(): UseSitemapReturn {
    const [isLoading, setIsLoading] = useState(false);
    const [isProductsLoading, setIsProductsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [sitemapData, setSitemapData] = useState<GenerateSitemapResponse | null>(null);
    const [products, setProducts] = useState<Product[]>([]);

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
    }, [fetchProducts]);

    return {
        isLoading,
        isProductsLoading,
        error,
        sitemapData,
        products,
        generateSitemap,
        downloadSitemap,
        clearError,
        reset,
    };
}