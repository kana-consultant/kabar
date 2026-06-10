// pages/admin/ai-management/hooks/useProviderForm.ts

import { useState, useCallback } from "react";
import type { APIProvider, ProviderFormData, Family, RequestSchema,  } from "@/types/provider.types";
import { DEFAULT_PROVIDER_FORM } from "@/constants/default-values";

// Helper untuk generate ID
const generateId = () => crypto.randomUUID?.() || `id-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

// Transform API Family + Schema → Form Family (flattened)
const transformFamilyToForm = (family: Family): Family => {
    
    return {
        id: family.id,
        provider_id: family.provider_id,
        schema_id: family.schema_id,
        name: family.name,
        display_name: family.display_name,
        description: family.description,
        max_token : family.max_token,
        system_prompt : family.system_prompt,
        temperature : family.temperature,
        schema: family.schema ? {
            id: family.schema.id,
            provider_id: family.schema.provider_id,
            name: family.schema.name,
            endpoint_path: family.schema.endpoint_path,
            request_template: family.schema.request_template,
            response_text_path: family.schema.response_text_path,
            response_image_path: family.schema.response_image_path,
            supports_temperature: family.schema.supports_temperature,
            supports_streaming: family.schema.supports_streaming,
        } : undefined,
    };
};

// Create empty form family
const createEmptyFamily = (): any => ({
    id: generateId(),
    provider_id: "",
    schema_id: "",
    name: "",
    display_name: "",
    description: null,
    endpoint_path: "",
    max_tokens_key: "max_tokens",
    system_role_key: "system",
    request_template: "{}",
    response_text_path: "",
    response_image_path: "",
    supports_temperature: true,
    supports_streaming: true,
});

// Transform Form Family → API Family + Schema
const transformFormToFamily = (formFamily: any, providerId: string): Family => {
    const schemaId = generateId();

    // Create RequestSchema
    const requestSchema: RequestSchema = {
        id: schemaId,
        provider_id: providerId,
        name: formFamily.name,
        endpoint_path: formFamily.endpoint_path,
        request_template: formFamily.request_template,
        response_text_path: formFamily.response_text_path,
        response_image_path: formFamily.response_image_path,
        supports_temperature: formFamily.supports_temperature,
        supports_streaming: formFamily.supports_streaming,
    };

    // Return Family with nested schema
    return {
        id: formFamily.id,
        provider_id: providerId,
        schema_id: schemaId,
        name: formFamily.name,
        display_name: formFamily.display_name,
        description: formFamily.description,
        schema: requestSchema,  // ← schema nested di dalam family
    };
};

export function useProviderForm(initialProvider: APIProvider | null) {
    const [formData, setFormData] = useState<ProviderFormData>(() => {
        if (initialProvider) {
            return {
                name: initialProvider.name,
                display_name: initialProvider.display_name,
                description: initialProvider.description,
                base_url: initialProvider.base_url,
                auth_type: initialProvider.auth_type,
                auth_header: initialProvider.auth_header,
                auth_prefix: initialProvider.auth_prefix,
                default_headers: initialProvider.default_headers || {},
                is_active: initialProvider.is_active,
                team_id: initialProvider.team_id || "",
                families: initialProvider.families?.map(family => transformFamilyToForm(family)) || [],
            };
        }

        return {
            ...DEFAULT_PROVIDER_FORM,
            families: [createEmptyFamily()],
        };
    });

    const updateFormData = useCallback((updates: Partial<ProviderFormData>) => {
        console.log(updates.auth_prefix)
        setFormData((prev) => ({ ...prev, ...updates }));
    }, []);

    const updateFamily = useCallback((familyId: string, updates: Partial<any>) => {
        setFormData((prev) => ({
            ...prev,
            families: prev.families?.map((family) =>
                family.id === familyId ? { ...family, ...updates } : family
            ),
        }));
    }, []);

    const addFamily = useCallback(() => {
        setFormData((prev) => ({
            ...prev,
            families: [...prev.families as Family[], createEmptyFamily()],
        }));
    }, []);

    const removeFamily = useCallback((familyId: string) => {
        setFormData((prev) => ({
            ...prev,
            families: prev.families?.filter((family) => family.id !== familyId),
        }));
    }, []);

    const duplicateFamily = useCallback((family: any) => {
        const duplicated = {
            ...family,
            id: generateId(),
            name: `${family.name}_copy`,
            display_name: `${family.display_name} (Copy)`,
        };
        setFormData((prev) => ({
            ...prev,
            families: [...prev.families as Family[], duplicated],
        }));
    }, []);

    const validateForm = useCallback(() => {
        const errors: Record<string, string> = {};

        if (!formData.name?.trim()) {
            errors.name = "Provider name is required";
        }
        if (!formData.display_name?.trim()) {
            errors.display_name = "Display name is required";
        }
        if (!formData.base_url?.trim()) {
            errors.base_url = "Base URL is required";
        }
        if (formData.families?.length === 0) {
            errors.families = "At least one API family is required";
        }

        formData.families?.forEach((family: any, index: number) => {
            if (!family.name?.trim()) {
                errors[`family_${family.id}_name`] = `Family ${index + 1}: Name is required`;
            }
            if (!family.display_name?.trim()) {
                errors[`family_${family.id}_display_name`] = `Family ${index + 1}: Display name is required`;
            }
            if (!family.endpoint_path?.trim()) {
                errors[`family_${family.id}_endpoint_path`] = `Family ${index + 1}: Endpoint path is required`;
            }

            if (family.request_template) {
                try {
                    JSON.parse(family.request_template);
                } catch {
                    errors[`family_${family.id}_template`] = `Family ${index + 1}: Invalid JSON format`;
                }
            }
        });

        return errors;
    }, [formData]);

    const getTestConnectionConfig = useCallback(() => {
        return {
            base_url: formData.base_url,
            auth_type: formData.auth_type,
            auth_header: formData.auth_header,
            auth_prefix: formData.auth_prefix || "",
            default_headers: formData.default_headers,
        };
    }, [formData]);

    const resetForm = useCallback(() => {
        setFormData({
            ...DEFAULT_PROVIDER_FORM,
            families: [createEmptyFamily()],
        });
    }, []);

    const getSubmitData = useCallback(() => {
        const providerId = initialProvider?.id || generateId();

        // Return API format: provider with families that have nested schema
        return {
            name: formData.name,
            display_name: formData.display_name,
            description: formData.description,
            base_url: formData.base_url,
            auth_type: formData.auth_type,
            auth_header: formData.auth_header,
            auth_prefix: formData.auth_prefix,
            default_headers: formData.default_headers,
            is_active: formData.is_active,
            team_id: formData.team_id,
            families: formData.families?.map((family: any) => transformFormToFamily(family, providerId)),
        };
    }, [formData, initialProvider]);

    return {
        formData,
        updateFormData,
        updateFamily,
        addFamily,
        removeFamily,
        duplicateFamily,
        validateForm,
        getTestConnectionConfig,
        resetForm,
        getSubmitData,
    };
}