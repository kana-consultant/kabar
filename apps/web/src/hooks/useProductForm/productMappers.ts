// productMappers.ts
import type { Product } from "@/services/product";
import type { ProductRequest, UpdateProductRequest } from "@/services/product";

const toJsonString = (value: unknown, fallback: Record<string, unknown> = {}): string => {
    if (typeof value === 'string') return value;
    return JSON.stringify(value || fallback);
};

export const mapProductToCreateRequest = (product: Partial<Product>): ProductRequest => {
    let adapterConfig = undefined;
    if (product.adapter_config) {
        adapterConfig = {
            custom_headers: toJsonString(product.adapter_config.custom_headers),
            meta_config: toJsonString(product.adapter_config.meta_config),
            sitemap_config: toJsonString(product.adapter_config.sitemap_config),
            timeout_seconds: product.adapter_config.timeout_seconds || 30,
            retry_count: product.adapter_config.retry_count || 3,
        };
    }

    let workflows = undefined;
    if (product.workflows && product.workflows.length > 0) {
        workflows = product.workflows.map(workflow => ({
            name: workflow.name,
            nodes: workflow.nodes?.map(node => ({
                id: node.id,
                workflow_id: node.workflow_id || workflow.id,
                step_order: node.step_order,
                next_node_id: node.next_node_id,
                previous_node_ids: node.previous_node_ids,
                adapter_config: node.adapter_config ? {
                    id: node.adapter_config.id,
                    endpoint_path: node.adapter_config.endpoint_path,
                    http_method: node.adapter_config.http_method,
                    field_mapping: toJsonString(node.input_mapping),
                    retry_count: node.adapter_config.retry_count || 3,
                } : undefined,
            })) || [],
        }));
    }

    return {
        name: product.name!,
        platform: product.platform!,
        api_endpoint: product.api_endpoint!,
        api_key: product.api_key,
        team_id: product.team_id,
        user_id: product.user_id,
        workflow_id: product.workflow_id,
        status: product.status || "pending",
        sync_status: product.sync_status || "idle",
        last_sync: product.last_sync,
        created_by: product.created_by,
        adapter_config: adapterConfig,
        ...(workflows && { workflows }),
    };
};

export const mapProductToUpdateRequest = (product: Partial<Product>): UpdateProductRequest => {
    let adapter_config = undefined;
    if (product.adapter_config) {
        adapter_config = {
            endpoint_path: product.adapter_config.endpoint_path,
            http_method: product.adapter_config.http_method,
            custom_headers: toJsonString(product.adapter_config.custom_headers),
            field_mapping: toJsonString(product.adapter_config.field_mapping),
            response_mapping: toJsonString(product.adapter_config.response_mapping),
            meta_config: toJsonString(product.adapter_config.meta_config),
            sitemap_config: toJsonString(product.adapter_config.sitemap_config),
            timeout_seconds: product.adapter_config.timeout_seconds,
            retry_count: product.adapter_config.retry_count,
        };
    }

    let workflows = undefined;
    if (product.workflows && product.workflows.length > 0) {
        workflows = product.workflows.map(workflow => ({
            id: workflow.id,
            product_id: workflow.product_id || product.id || "",
            name: workflow.name,
            nodes: workflow.nodes?.map(node => ({
                id: node.id,
                adapter_config_id: node.adapter_config_id,
                step_order: node.step_order,
                next_node_id: node.next_node_id,
                previous_node_ids: node.previous_node_ids || [],
                adapter_config: node.adapter_config ? {
                    id: node.adapter_config.id,
                    endpoint_path: node.adapter_config.endpoint_path,
                    http_method: node.adapter_config.http_method,
                    field_mapping: toJsonString(node.adapter_config.field_mapping),
                } : undefined,
            })) || [],
        }));
    }

    return {
        name: product.name,
        platform: product.platform,
        api_endpoint: product.api_endpoint,
        api_key: product.api_key,
        status: product.status,
        sync_status: product.sync_status,
        team_id: product.team_id,
        user_id: product.user_id,
        created_by: product.created_by,
        workflow_id: product.workflow_id,
        adapter_config: adapter_config,
        workflows: workflows as any,
    };
};