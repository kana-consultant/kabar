// pages/admin/ai-management/hooks/useProviderForm.ts

import { useState, useCallback } from "react";
import type { APIProvider, ProviderFormData, Family } from "@/types/provider.types";
import { DEFAULT_PROVIDER_FORM } from "@/constants/default-values";

// Helper untuk generate ID
const generateId = () => crypto.randomUUID?.() || `id-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

// Default family template - sesuai dengan tipe Family
const createDefaultFamily = (): Family => ({
    id: generateId(),
    provider_id: "",
    name: "",
    display_name: "",
    description: null,
    max_token: undefined,
    temperature: undefined,
    system_prompt: "",
    schema: {
        id: generateId(),
        provider_id: "",
        name: "",
        endpoint_path: "",
        request_template: "{}",
        response_text_path: "choices[0].message.content",
        response_image_path: "",
        supports_temperature: true,
        supports_streaming: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
    }
});

export function useProviderForm(initialProvider: APIProvider | null) {
    const [formData, setFormData] = useState<ProviderFormData>(() => {
        if (initialProvider) {
            // Gunakan langsung data dari API tanpa transformasi
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
                families: initialProvider.families || [],
            };
        }

        return {
            ...DEFAULT_PROVIDER_FORM,
            families: [createDefaultFamily()],
        };
    });

    const updateFormData = useCallback((updates: Partial<ProviderFormData>) => {
        setFormData((prev) => {
            // Deep merge untuk nested objects jika diperlukan
            if (updates.default_headers) {
                return {
                    ...prev,
                    ...updates,
                    default_headers: {
                        ...prev.default_headers,
                        ...updates.default_headers,
                    }
                };
            }
            return { ...prev, ...updates };
        });
    }, []);

    const updateFamily = useCallback((familyId: string, updates: Partial<Family>) => {
        setFormData((prev) => ({
            ...prev,
            families: prev.families?.map((family) =>
                family.id === familyId 
                    ? {
                        ...family,
                        ...updates,
                        // Deep merge schema jika ada
                        schema: updates.schema 
                            ? { ...family.schema, ...updates.schema }
                            : family.schema
                    }
                    : family
            ),
        }));
    }, []);

    const addFamily = useCallback(() => {
        setFormData((prev) => ({
            ...prev,
            families: [...(prev.families || []), createDefaultFamily()],
        }));
    }, []);

    const removeFamily = useCallback((familyId: string) => {
        setFormData((prev) => ({
            ...prev,
            families: prev.families?.filter((family) => family.id !== familyId) || [],
        }));
    }, []);

    const duplicateFamily = useCallback((family: Family) => {
        const duplicated: Family = {
            ...family,
            id: generateId(),
            name: `${family.name}_copy`,
            display_name: `${family.display_name || family.name} (Copy)`,
            schema: family.schema ? {
                ...family.schema,
                id: generateId(),
                name: `${family.schema.name}_copy`,
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
            } : undefined,
        };
        
        setFormData((prev) => ({
            ...prev,
            families: [...(prev.families || []), duplicated],
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
        if (!formData.families || formData.families.length === 0) {
            errors.families = "At least one API family is required";
        }

        formData.families?.forEach((family, index) => {
            if (!family.name?.trim()) {
                errors[`family_${family.id}_name`] = `Family ${index + 1}: Name is required`;
            }
            if (!family.display_name?.trim()) {
                errors[`family_${family.id}_display_name`] = `Family ${index + 1}: Display name is required`;
            }
            if (!family.schema?.endpoint_path?.trim()) {
                errors[`family_${family.id}_endpoint_path`] = `Family ${index + 1}: Endpoint path is required`;
            }

            // Validasi request template jika ada
            const template = family.schema?.request_template;
            if (template && template !== "{}") {
                try {
                    JSON.parse(template);
                } catch {
                    errors[`family_${family.id}_template`] = `Family ${index + 1}: Invalid JSON format in request template`;
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
            families: [createDefaultFamily()],
        });
    }, []);

    const getSubmitData = useCallback(() => {
        // Tidak perlu transformasi lagi karena struktur sudah sama
        return {
            ...formData,
            id: initialProvider?.id, // Include ID jika edit mode
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