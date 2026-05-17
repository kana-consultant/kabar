export interface FieldMapping {
    id: string;
    sourceField: string;
    targetField: string;
    isRequired: boolean;
    defaultValue?: string;
}

export interface NestedMapping {
    id: string;
    sourceField: string;
    targetField: string;
    isRequired: boolean;
    defaultValue?: string;
    children?: FieldMapping[];  // ← untuk nested
    isExpanded?: boolean;       // ← untuk UI
}

// srcservicesproduct.ts
export interface AdapterConfig {
    id?: string;
    productId?: string;
    endpointPath: string;
    httpMethod: 'POST' | 'PUT' | 'PATCH';
    customHeaders: string;
    fieldMapping: string;
    metaConfig?: string;        // JSON string untuk konfigurasi meta tag
    sitemapConfig?: string;     // JSON string untuk konfigurasi sitemap
    timeoutSeconds?: number;
    retryCount?: number;
    createdAt?: string;
    updatedAt?: string;
}

export interface Product {
    id: string;
    name: string;
    platform: 'wordpress' | 'shopify' | 'custom';
    apiEndpoint: string;
    apiKey?: string;
    APIKeyEncrypted?:string;
    status: 'connected' | 'pending' | 'error' | 'disconnected';
    syncStatus: 'idle' | 'syncing' | 'success' | 'failed';
    lastSync?: string;
    createdBy?: string;
    teamId?: string;
    userId?: string;
    createdAt: string;
    updatedAt: string;
    adapterConfig?: AdapterConfig;
}