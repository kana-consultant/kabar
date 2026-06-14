export interface ProductFormState {
    name: string;
    platform: string;
    api_endpoint: string;
    api_key: string;
    status: string;
    lastSync: string;
    adapter_config: {
        endpointPath: string;
        httpMethod: string;
        customHeaders: Record<string, string>;
        fieldMapping: string;
    };
}