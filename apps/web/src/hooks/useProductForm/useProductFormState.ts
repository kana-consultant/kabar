import { useState } from "react";
import type { Product, AdapterConfig } from "@/services/product";

const createEmptyAdapterConfig = (): AdapterConfig => ({
    id: "",
    product_id: "",
    endpoint_path: "",
    http_method: "POST",
    custom_headers: "",
    field_mapping: "",
    response_mapping: "",
    meta_config: "",
    sitemap_config: "",
    timeout_seconds: 30,
    retry_count: 3,
    created_at: "",
    updated_at: "",
});

const createEmptyProduct = (): Product => ({
    id: "",
    name: "",
    platform: "wordpress",
    api_endpoint: "",
    sync_status: "idle",
    created_at: "",
    updated_at: "",
    api_key: "",
    status: "pending",
    adapter_config: createEmptyAdapterConfig(),
    workflow_id: "",
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