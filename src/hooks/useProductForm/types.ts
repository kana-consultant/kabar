import type { Product, AdapterConfig } from "@/types/product";

export interface ProductFormState {
    name: string;
    platform: string;
    apiEndpoint: string;
    apiKey: string;
    status: string;
    lastSync: string;
    adapterConfig: {
        endpointPath: string;
        httpMethod: string;
        customHeaders: Record<string, string>;
        fieldMapping: string;
    };
}