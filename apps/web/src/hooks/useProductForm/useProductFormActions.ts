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

    // ========== WORKFLOW ACTIONS ==========
    
    const updateWorkflowId = (workflowId: string) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            workflow_id: workflowId,
        }));
    };

    const addAdapterConfig = (adapter: AdapterConfig) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapterConfigs: [...(prev.adapterConfigs || []), adapter],
        }));
    };

    const updateAdapterConfigs = (index: number, updates: Partial<AdapterConfig>) => {
        setProduct((prev: Partial<Product>) => {
            const updatedConfigs = [...(prev.adapterConfigs || [])];
            updatedConfigs[index] = { ...updatedConfigs[index], ...updates } as AdapterConfig;
            return {
                ...prev,
                adapterConfigs: updatedConfigs,
            };
        });
    };

    const removeAdapterConfig = (index: number) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapterConfigs: (prev.adapterConfigs || []).filter((_, i) => i !== index),
        }));
    };

    return {
        updateProductInfo,
        updateAdapterConfig,
        updateFieldMapping,
        updateMetaConfig,
        updateSitemapConfig,
        // Workflow actions
        updateWorkflowId,
        addAdapterConfig,
        updateAdapterConfigs,
        removeAdapterConfig,
    };
}