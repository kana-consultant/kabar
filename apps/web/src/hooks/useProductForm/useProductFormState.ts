import { useState } from "react";
import type { Product, AdapterConfig } from "@/services/product";

const createEmptyAdapterConfig = (): Omit<AdapterConfig, 'id' | 'productId'> => ({
    endpointPath: "",
    httpMethod: "POST",
    customHeaders: JSON.stringify({ "Content-Type": "application/json" }),
    fieldMapping: "{}",
    responseMapping: {},
    metaConfig: "{}",
    sitemapConfig: "{}",
    timeoutSeconds: 30,
    retryCount: 3,
});

const createEmptyProduct = (): Product => ({
    id: "",
    name: "",
    platform: "wordpress",
    apiEndpoint: "",
    apiKey: "",
    status: "pending",
    syncStatus: "idle",
    lastSync: undefined,
    createdBy: undefined,
    teamId: undefined,
    userId: undefined,
    createdAt: "",
    updatedAt: "",
    adapterConfig: undefined,
    workflow_id: "",
    adapterConfigs: [],
    workflows: [],
});

export function useProductFormState() {
    const [loading, setLoading] = useState(false);
    const [testing, setTesting] = useState(false);
    const [product, setProduct] = useState<Product>(createEmptyProduct());

    return {
        loading,
        setLoading,
        testing,
        setTesting,
        product,
        setProduct,
    };
}