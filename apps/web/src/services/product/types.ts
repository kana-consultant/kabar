import type {
    Product,
    AdapterConfig,
    FieldMapping,
    NestedMapping,
    WorkflowDefinition,
    WorkflowNode,
} from "@/types/product";

export type {
    Product,
    AdapterConfig,
    FieldMapping,
    NestedMapping,
    WorkflowDefinition,
    WorkflowNode,
};

// Reuse type dari Product dan WorkflowDefinition
export interface CreateProductRequest {
    name: Product["name"];
    platform: Product["platform"];
    api_endpoint: Product["api_endpoint"];
    api_key?: Product["api_key"];
    team_id?: Product["team_id"];
    user_id?: Product["user_id"];
    workflow_id?: Product["workflow_id"];
    status?: Product["status"];
    sync_status?: Product["sync_status"];
    last_sync?: Product["last_sync"];
    created_by?: Product["created_by"];
    
    adapter_config?: Omit<AdapterConfig,'endpoint_path' | 'http_method' |'id' | 'created_at' | 'updated_at'> & {
        http_method?: string;
        id?: string;
        created_at?: string;
        updated_at?: string;
    };
    
    workflows?: (Omit<WorkflowDefinition, "input_mapping" | "id" |  "created_at" | "updated_at"> & {
        id? : string
        created_at?: string;
        updated_at?: string;
        nodes?: (Omit<WorkflowNode, "adapter_config_idt"| "adapter_config_idts" | "workflow_id" | "created_at" | "updated_at"> & {
            created_at?: string;
            updated_at?: string;
            adapter_config?: Omit<AdapterConfig, "created_at" | "updated_at"> & {
                created_at?: string;
                updated_at?: string;
            };
        })[];
    })[];
}

export interface UpdateProductRequest {
    name?: string;
    platform?: Product["platform"];
    api_endpoint?: string;
    api_key?: string;
    status?: Product["status"];
    sync_status?: Product["sync_status"];
    last_sync?: string;
    team_id?: string;
    user_id?: string;
    created_by?: string;
    workflow_id?: string;
    
    // Adapter config - semua field optional
    adapter_config?: Partial<AdapterConfig>;
    
    // Workflows untuk update
    workflows?: Partial<WorkflowDefinition>[];
}

export interface AddProductResponse {
    id: string;
    message: string;
}