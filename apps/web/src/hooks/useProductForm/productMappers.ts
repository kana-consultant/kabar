// productMappers.ts
import type { Product } from "@/services/product";
import type { CreateProductRequest, UpdateProductRequest } from "@/services/product";

export const mapProductToCreateRequest = (product: Partial<Product>): CreateProductRequest => ({
    name: product.name!,
    platform: product.platform!,
    api_endpoint: product.api_endpoint!,
    api_key: product.api_key,
    team_id: product.team_id,
    workflow_id: product.workflow_id,
    adapter_config: product.adapter_config ? {
        endpoint_path: product.adapter_config.endpoint_path,
        http_method: product.adapter_config.http_method,
        custom_headers: typeof product.adapter_config.custom_headers === 'string'
            ? product.adapter_config.custom_headers
            : JSON.stringify(product.adapter_config.custom_headers),
        field_mapping: typeof product.adapter_config.field_mapping === 'string'
            ? product.adapter_config.field_mapping
            : JSON.stringify(product.adapter_config.field_mapping),
        response_mapping: typeof product.adapter_config.response_mapping === 'string'
            ? product.adapter_config.response_mapping
            : JSON.stringify(product.adapter_config.response_mapping),
        meta_config: product.adapter_config.meta_config,
        sitemap_config: product.adapter_config.sitemap_config,
        timeout_seconds: product.adapter_config.timeout_seconds,
        retry_count: product.adapter_config.retry_count,
    } : undefined,
});

export const mapProductToUpdateRequest = (product: Partial<Product>): UpdateProductRequest => ({
    name: product.name,
    platform: product.platform,
    api_endpoint: product.api_endpoint,
    api_key: product.api_key,
    status: product.status,
    sync_status: product.sync_status,
    team_id: product.team_id,
    adapter_config: product.adapter_config ? {
        endpoint_path: product.adapter_config.endpoint_path,
        http_method: product.adapter_config.http_method,
        custom_headers: typeof product.adapter_config.custom_headers === 'string'
            ? product.adapter_config.custom_headers
            : JSON.stringify(product.adapter_config.custom_headers),
        field_mapping: typeof product.adapter_config.field_mapping === 'string'
            ? product.adapter_config.field_mapping
            : JSON.stringify(product.adapter_config.field_mapping),
        response_mapping: typeof product.adapter_config.response_mapping === 'string'
            ? product.adapter_config.response_mapping
            : JSON.stringify(product.adapter_config.response_mapping),
        meta_config: product.adapter_config.meta_config,
        sitemap_config: product.adapter_config.sitemap_config,
        timeout_seconds: product.adapter_config.timeout_seconds,
        retry_count: product.adapter_config.retry_count,
    } : undefined,
});