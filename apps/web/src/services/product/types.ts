import type {
    Product,
    AdapterConfig,
    FieldMapping,
    NestedMapping,
    WorkflowDefinition,
} from "@/types/product";

export type {
    Product,
    AdapterConfig,
    FieldMapping,
    NestedMapping,
};

export interface CreateProductRequest {
    name: Product["name"];
    platform: Product["platform"];
    api_endpoint: Product["api_endpoint"];

    api_key?: Product["api_key"];
    team_id?: Product["team_id"];
    workflow_id?: Product["workflow_id"];

    adapter_config?: Partial<AdapterConfig>;
    adapter_configs?: Partial<AdapterConfig>[];
    workflows?: Partial<WorkflowDefinition>[];
}

export interface UpdateProductRequest {
    name?: string;
    platform?: Product["platform"];
    api_endpoint?: string;
    api_key?: string;

    status?: Product["status"];
    sync_status?: Product["sync_status"];

    team_id?: string;

    adapter_config?: {
        endpoint_path?: string;
        http_method?: AdapterConfig["http_method"];
        custom_headers?: string;
        field_mapping?: string;
        response_mapping?: string;
        meta_config?: string;
        sitemap_config?: string;
        timeout_seconds?: number;
        retry_count?: number;
    };
}

export interface AddProductResponse {
    id: string;
    message: string;
}