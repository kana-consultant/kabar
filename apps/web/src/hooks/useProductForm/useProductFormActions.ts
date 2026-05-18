import type { Product, AdapterConfig } from "@/services/product";

export function useProductFormActions(
    setProduct: (updater: any) => void
) {
    const updateProductInfo = (updates: Partial<Product>) => {
        setProduct((prev: Partial<Product>) => ({ ...prev, ...updates }));
    };

    const updateAdapterConfig = (updates: Partial<AdapterConfig>) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapterConfig: { ...prev.adapterConfig, ...updates } as AdapterConfig,
        }));
    };

    const updateFieldMapping = (value: string) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapterConfig: {
                ...prev.adapterConfig,
                fieldMapping: value,
            } as AdapterConfig,
        }));
    };

    const updateMetaConfig = (metaConfig: string) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapterConfig: {
                ...prev.adapterConfig,
                metaConfig: metaConfig,
            } as AdapterConfig,
        }));
    };

    const updateSitemapConfig = (sitemapConfig: string) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapterConfig: {
                ...prev.adapterConfig,
                sitemapConfig: sitemapConfig,
            } as AdapterConfig,
        }));
    };

    return {
        updateProductInfo,
        updateAdapterConfig,
        updateFieldMapping,
        updateMetaConfig,
        updateSitemapConfig
    };
}