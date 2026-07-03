export interface FieldMapping {
  id: string;
  source_field: string;
  target_field: string;
  is_required: boolean;
  default_value?: string;
}

export interface NestedMapping {
  id: string;
  source_field: string;
  target_field: string;
  is_required: boolean;
  default_value?: string;
  children?: FieldMapping[];
  is_expanded?: boolean;
}

export interface Product {
  id: string;
  name: string;
  domain? : string;
  platform: 'wordpress' | 'shopify' | 'custom';
  api_endpoint: string;
  api_key?: string;
  status: 'connected' | 'pending' | 'error' | 'disconnected';
  sync_status: 'idle' | 'syncing' | 'success' | 'failed';
  last_sync?: string;
  created_by?: string;
  team_id?: string;
  user_id?: string;
  created_at: string;
  updated_at: string;
  adapter_config?: AdapterConfig;
  workflow_id: string;
  workflows?: WorkflowDefinition[];
}


export interface AdapterConfig {
  id?: string;
  custom_headers?: string;
  meta_config?: string;
  sitemap_config?: string;
  timeout_seconds?: number;
  retry_count?: number;
  created_at?: string;
  updated_at?: string;
}


export interface AdapterConfigNode {
  id?: string;
  product_id?: string;
  endpoint_path?: string;
  http_method?: 'POST' | 'PUT' | 'PATCH' | 'GET' | 'DELETE';
  field_mapping?: string;
}

export interface WorkflowDefinition {
  id: string;
  product_id?: string;
  name: string;
  created_at: string;
  updated_at: string;
  nodes?: WorkflowNode[];
}

export interface WorkflowNode {
  id?: string;
  workflow_id?: string;
  adapter_config_id?: string;
  previous_node_ids?: string[] | null;
  step_order: number;
  input_mapping?: Record<string, any>;
  next_node_ids?: string[] | null;
  created_at?: string;
  adapter_config?: AdapterConfigNode;
}

export interface ProductWithDetails extends Product {
  adapter_configs: AdapterConfig[];
  workflows: WorkflowDefinition[];
  active_workflow_id?: string;
}