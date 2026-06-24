// ============================================================
// FILE: types/provider.types.ts
// ============================================================

/**
 * Validation Error
 */
export interface ValidationError {
    field: string;
    message: string;
}

/**
 * API Response Wrapper
 */
export interface APIResponse<T = any> {
    data?: T;
    error?: string;
    message?: string;
    validation_errors?: ValidationError[];
}

/**
 * Connection Test Configuration
 */
export interface TestConnectionConfig {
    base_url: string;
    auth_type: string;
    auth_header: string;
    auth_prefix: string;
}

/**
 * Connection Test Result
 */
export interface TestConnectionResult {
    success: boolean;
    message?: string;
    latency?: number;
    status_code?: number;
    error?: string;
}

/**
 * Provider Statistics
 */
export interface ProviderStats {
    total_requests: number;
    total_tokens: number;
    average_latency: number;
    error_rate: number;
    last_24h_requests: number;
    active_families: number;
}

/**
 * Provider Usage Metrics
 */
export interface UsageMetrics {
    timestamp: string;
    requests: number;
    tokens: number;
    latency: number;
    errors: number;
}

/**
 * Provider Filter Options
 */
export interface ProviderFilters {
    status?: 'active' | 'inactive' | 'all';
    family?: string;
    search?: string;
    team_id?: string;
}

/**
 * Provider Sort Options
 */
export interface ProviderSort {
    field: 'name' | 'display_name' | 'created_at' | 'updated_at' | 'is_active';
    order: 'asc' | 'desc';
}


export interface APIProvider {
    id?: string;
    name: string;
    display_name: string;
    description: string | null;
    base_url: string;
    auth_type: string;
    auth_header: string;
    auth_prefix: string | null;
    default_headers: Record<string, string>;
    is_active: boolean;
    families: Family[];
    team_id?: string;
    created_at?: string;
    updated_at?: string;
}

/**
 * Provider List Response with Pagination
 */
export interface ProviderListResponse {
    data: APIProvider[];
    total: number;
    page: number;
    limit: number;
    total_pages: number;
}

/**
 * Create Provider Request
 */
export interface CreateProviderRequest {
    name: string;
    displayName: string;
    description?: string;
    baseUrl: string;
    authType: string;
    authHeader: string;
    authPrefix: string;
    textEndpoint: string;
    imageEndpoint?: string;
    defaultHeaders: Record<string, string>;
    requestTemplate: Record<string, any>;
    responseTextPath: string;
    responseImagePath?: string;
}


/**
 * Update Provider Request
 */
export type UpdateProviderRequest = Partial<CreateProviderRequest>;

/**
 * Create Response
 */
export interface CreateResponse {
    id: string;
    message: string;
}

/**
 * Update Response
 */
export interface UpdateResponse {
    message: string;
}

// ============================================================
// FILE: types/model.types.ts
// ============================================================


/**
 * AI Model Schema (Complete)
 */
export interface AIModelSchema {
    id: string;
    family_id: string | null;
    provider_id: string;
    schema_id: string | null;
    team_id: string | null;
    name: string;
    display_name: string;
    description: string | null;
    system_prompt: string | null;
    max_tokens: number | null;
    temperature: number | null;
    context_window: number | null;
    schema: RequestSchema;
    is_active: boolean;
    is_default: boolean;
    created_by: string | null;
    created_at: string;
    updated_at: string;
}

/**
 * AI Model (Simplified for lists)
 */
export interface AIModel {
    id: string;
    name: string;
    provider_id: string;
    display_name: string;
    is_active: boolean;
    is_default: boolean;
    max_tokens: number;
    team_id?: string;
    provider: string;
    created_at: string;
    updated_at: string;
}

/**
 * Model with API Key Status
 */
export interface ModelWithStatus extends AIModel {
    has_api_key: boolean;
    provider_display_name?: string;
}

/**
 * Model from API Key
 */
export interface ModelFromAPIKey {
    id: string;
    model_id: string;
    name: string;
    display_name: string;
    provider_name: string;
    provider_display_name: string;
    service: string;
}

/**
 * Create Model Request
 */
export interface CreateModelRequest {
    name: string;
    provider_id: string;
    display_name: string;
    is_active: boolean;
    is_default: boolean;
    max_tokens: number;
    system_prompt?: string;
    temperature?: number;
    context_window?: number;
    description?: string;
    family_id?: string | null;
    schema_id?: string | null;
    team_id?: string | null;
}

/**
 * Update Model Request
 */
export interface UpdateModelRequest extends Partial<CreateModelRequest> {
    id: string;
}

// ============================================================
// FILE: types/test.types.ts
// ============================================================

/**
 * Family Test Result
 */
export interface FamilyTestResult {
    family_id: string;
    family_name: string;
    success: boolean;
    endpoint: string;
    latency: number;
    error?: string;
}

/**
 * Provider Test Result
 */
export interface ProviderTestResult {
    provider_id: string;
    provider_name: string;
    success: boolean;
    families: FamilyTestResult[];
    overall_latency: number;
    timestamp: string;
}






// Form data untuk schema
export interface RequestSchemaFormData {
    id?: string;
    provider_id: string;
    name: string;
    endpoint_path: string;
    request_template: string | null;
    response_text_path: string;
    response_image_path: string | null;
    supports_temperature: boolean;
    supports_streaming: boolean;
}

/**
 * API Family Configuration (Request Schema)
 */
export interface RequestSchema {
    id: string;
    provider_id: string;
    name: string;
    endpoint_path: string;
    request_template: string ;
    response_text_path: string;
    response_image_path: string ;
    supports_temperature: boolean;
    supports_streaming: boolean;
    created_at?: string;
    updated_at?: string;
}



export interface Family {
    id: string;
    provider_id: string;
    schema_id?: string;
    name: string;
    display_name: string;
    description: string | null;
    schema?: RequestSchema;
    // Parameter baru untuk AI model configuration
    max_token?: number;        // Maximum tokens untuk response
    temperature?: number;      // Temperature untuk kreativitas (0-2)
    system_prompt?: string;    // System prompt untuk AI instructions
    created_at?: string;
    updated_at?: string;
}

/**
 * Provider Form Data (Omit fields that are auto-generated)
 */
export interface ProviderFormData{
    name: string;
    display_name: string;
    description: string | null;
    base_url: string;
    auth_type: string;
    auth_header: string;
    auth_prefix: string | null;
    default_headers: Record<string, string>;
    is_active: boolean;
    team_id?: string;
    families?: Family[];  // ← dibuat optional
}
