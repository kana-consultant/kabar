// pages/admin/ai-management/hooks/useFamilyManagement.ts

import { useState, useCallback } from "react";
import type { Family, RequestSchema } from "@/types/provider.types";

const generateId = () => crypto.randomUUID?.() || `id-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

// Clean family data - remove flattened fields
const cleanFamily = (family: any): Family => {
    return {
        id: family.id || generateId(),
        provider_id: family.provider_id || "",
        schema_id: family.schema?.id || family.schema_id || generateId(),
        name: family.name || "",
        display_name: family.display_name || "",
        system_prompt : family.system_prompt,
        max_token : family.max_token,
        temperature : family.temperature,
        description: family.description || null,
        schema: {
            id: family.schema?.id || generateId(),
            provider_id: family.provider_id || "",
            name: family.schema?.name || family.name || "",
            endpoint_path: family.schema?.endpoint_path || family.endpoint_path || "",
            request_template: family.schema?.request_template || family.request_template || "{}",
            response_text_path: family.schema?.response_text_path || family.response_text_path || "",
            response_image_path: family.schema?.response_image_path || family.response_image_path || "",
            supports_temperature: family.schema?.supports_temperature ?? family.supports_temperature ?? true,
            supports_streaming: family.schema?.supports_streaming ?? family.supports_streaming ?? true,
        },
    };
};

// Create empty schema
const createEmptySchema = (providerId?: string): RequestSchema => ({
    id: generateId(),
    provider_id: providerId || "",
    name: "",
    endpoint_path: "",
    request_template: "{}",
    response_text_path: "",
    response_image_path: "",
    supports_temperature: true,
    supports_streaming: true,
});

// Create empty family dengan schema (NO flattened fields)
const createEmptyFamily = (providerId?: string): Family => ({
    id: generateId(),
    provider_id: providerId || "",
    schema_id: generateId(),
    name: "",
    display_name: "",
    description: null,
    schema: createEmptySchema(providerId),
});

export function useFamilyManagement(initialFamilies: Family[] = [], providerId?: string) {
    const [families, setFamilies] = useState<Family[]>(() => {
        if (initialFamilies && initialFamilies.length > 0) {
            // Clean each family to remove flattened fields
            return initialFamilies.map(family => cleanFamily(family));
        }
        return [createEmptyFamily(providerId)];
    });

    const addFamily = useCallback((): Family => {
        const newFamily = createEmptyFamily(providerId);
        setFamilies(prev => [...prev, newFamily]);
        return newFamily;
    }, [providerId]);

    const removeFamily = useCallback((id: string) => {
        setFamilies(prev => {
            if (prev.length <= 1) {
                console.warn("Cannot remove the last family");
                return prev;
            }
            return prev.filter(family => family.id !== id);
        });
    }, []);

    const updateFamily = useCallback((id: string, updates: Partial<Family>) => {
        setFamilies(prev => prev.map(family =>
            family.id === id ? { ...family, ...updates } : family
        ));
    }, []);

    const updateFamilySchema = useCallback((familyId: string, schemaUpdates: Partial<RequestSchema>) => {
        setFamilies((prev  : any) => prev.map((family : any) =>
            family.id === familyId 
                ? { 
                    ...family, 
                    schema: { ...family.schema, ...schemaUpdates }
                  }
                : family
        ));
    }, []);

    const duplicateFamily = useCallback((family: Family): Family => {
        const newId = generateId();
        const newSchemaId = generateId();
        
        return {
            id: newId,
            provider_id: providerId || family.provider_id,
            schema_id: newSchemaId,
            name: `${family.name}-copy`,
            display_name: `${family.display_name} (Copy)`,
            description: family.description,
            schema: {
                ...family.schema,
                id: newSchemaId,
                name: `${family.schema.name}-copy`,
                provider_id: providerId || family.schema.provider_id,
            },
        };
    }, [providerId]);

    const validateFamilies = useCallback(() => {
        const errors: Record<string, string> = {};
        
        families.forEach((family, index) => {
            if (!family.name?.trim()) {
                errors[`family_${family.id}_name`] = `Family ${index + 1}: Name is required`;
            }
            if (!family.display_name?.trim()) {
                errors[`family_${family.id}_display_name`] = `Family ${index + 1}: Display name is required`;
            }
            if (!family.schema?.endpoint_path?.trim()) {
                errors[`family_${family.id}_endpoint_path`] = `Family ${index + 1}: Endpoint path is required`;
            }
            if (!family.schema?.response_text_path?.trim()) {
                errors[`family_${family.id}_response_text_path`] = `Family ${index + 1}: Response text path is required`;
            }
            if (family.schema?.request_template) {
                try {
                    JSON.parse(family.schema.request_template);
                } catch {
                    errors[`family_${family.id}_template`] = `Family ${index + 1}: Invalid JSON format`;
                }
            }
        });

        return errors;
    }, [families]);

  

    const resetFamilies = useCallback(() => {
        setFamilies([createEmptyFamily(providerId)]);
    }, [providerId]);

    const getFamilyById = useCallback((id: string) => {
        return families.find(family => family.id === id);
    }, [families]);

    return {
        families,
        addFamily,
        removeFamily,
        updateFamily,
        updateFamilySchema,
        duplicateFamily,
        validateFamilies,
        resetFamilies,
        getFamilyById,
    };
}