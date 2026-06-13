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

// types/product.ts
export interface Product {
    id: string;
    name: string;
    platform: 'wordpress' | 'shopify' | 'custom';
    apiEndpoint: string;
    apiKey?: string;
    status: 'connected' | 'pending' | 'error' | 'disconnected';
    syncStatus: 'idle' | 'syncing' | 'success' | 'failed';
    lastSync?: string;
    createdBy?: string;
    teamId?: string;
    userId?: string;
    createdAt: string;
    updatedAt: string;
    adapterConfig?: AdapterConfig;  // semua adapter milik product
    workflow_id: string;
    // Relations
    adapterConfigs?: AdapterConfig[];  // semua adapter milik product
    workflows?: WorkflowDefinition[];   // semua workflow milik product
}

// types/adapter.ts
export interface AdapterConfig {
    id: string;
    productId: string;
    endpointPath: string;
    httpMethod: 'POST' | 'PUT' | 'PATCH' | 'GET' | 'DELETE';
    customHeaders: string;
    fieldMapping: string;  // JSON string
    responseMapping: Record<string, any> | string;  // JSON string untuk extract response
    metaConfig?: string;
    sitemapConfig?: string;
    timeoutSeconds?: number;
    retryCount?: number;
    createdAt?: string;
    updatedAt?: string;
}

// types/workflow.ts
export interface WorkflowDefinition {
    id: string;
    productId: string;
    name: string;
    createdAt: string;
    updatedAt: string;
    nodes?: WorkflowNode[];
}

export interface WorkflowNode {
    id: string;
    workflowId: string;
    adapterConfigId: string;
    stepOrder: number;
    inputMapping: Record<string, any>;
    nextNodeId?: string | null;
    createdAt: string;

    // Relations (dari join)
    adapterConfig?: AdapterConfig;
}

// Product dengan data lengkap (response dari GET)
export interface ProductWithDetails extends Product {
    adapterConfigs: AdapterConfig[];
    workflows: WorkflowDefinition[];
    activeWorkflowId?: string;  // workflow yang sedang aktif
}