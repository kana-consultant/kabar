import { useState, useEffect } from "react";
import type { ModelFormData } from "@/pages/ai-management/Model/provider.types";
import type { AIModel, AIModelSchema } from "@/types/provider.types";
import { createModel, updateModel, deleteModel } from "@/services/model/modelMutations";

export function useModelForm(initialData?: AIModelSchema, currentUserId?: string) {
    const [formData, setFormData] = useState<ModelFormData>({
        provider_id: "",
        family_id: "",
        name: "",
        display_name: "",
        description: null,
        request_template: {},
        response_text_path: null,
        response_image_path: null,
        max_tokens: 4096,
        temperature: 0.7,
        is_active: true,
        is_default: false,
        team_id: null,
        created_by: null,
    });

    console.log("KOOOOOOOOOOOOOOOOOOOOOOO")
        console.log(initialData)

    useEffect(() => {
        console.log("KOOOOOOOOOOOOOOOOOOOOOOO")
        console.log(initialData)
        if (initialData) {
            let requestTemplate: Record<string, any> = {};
            if (initialData.schema.request_template) {
                if (typeof initialData.schema.request_template === 'string') {
                    try {
                        requestTemplate = JSON.parse(initialData.schema.request_template);
                    } catch (e) {
                        requestTemplate = {};
                    }
                } else {
                    requestTemplate = initialData.schema.request_template as Record<string, any>;
                }
            }

            console.log("KFASDWBNKNFASKSNAK-----")
            console.log(initialData)
            setFormData({
                provider_id: initialData.provider_id,
                family_id: initialData.family_id as string,
                name: initialData.name,
                display_name: initialData.display_name,
                description: initialData.description,
                request_template: requestTemplate,
                response_text_path: initialData.schema.response_text_path,
                response_image_path: initialData.schema.response_image_path,
                max_tokens: initialData.max_tokens,
                temperature: initialData.temperature,
                is_active: initialData.is_active,
                is_default: initialData.is_default,
                team_id: initialData.team_id,
                created_by: initialData.created_by,
            });
        } else if (currentUserId) {
            setFormData((prev: any) => ({ ...prev, created_by: currentUserId }));
        }
    }, [initialData, currentUserId]);

    const updateField = <K extends keyof ModelFormData>(field: K, value: ModelFormData[K]) => {
        setFormData((prev: any) => ({ ...prev, [field]: value }));
    };

    const resetForm = () => {
        setFormData({
            provider_id: "",
            family_id: "",
            name: "",
            display_name: "",
            description: null,
            request_template: {},
            response_text_path: null,
            response_image_path: null,
            max_tokens: 4096,
            temperature: 0.7,
            is_active: true,
            is_default: false,
            team_id: null,
            created_by: currentUserId || null,
        });
    };

    const getSaveData = (): Omit<ModelFormData, 'created_by'> & { request_template: string } => {
        return {
            ...formData,
            request_template: JSON.stringify(formData.request_template),
        };
    };

    const isValid = () => {
        return !!formData.name && !!formData.display_name && !!formData.provider_id;
    };

    const handleSubmitModel = async (modelId?: string) => {
        try {
            const saveData = getSaveData();
            
            if (modelId) {
                // Update existing model
                const updatedModel = await updateModel(modelId, saveData as ModelFormData);
                return { success: true, data: updatedModel };
            } else {
                // Create new model
                const newModel = await createModel(saveData as any);
                return { success: true, data: newModel };
            }
        } catch (error) {
            console.error("Failed to save model:", error);
            return { success: false, error };
        }
    };

    const handleDelete = async (modelId: string) => {
        try {
            await deleteModel(modelId);
            return { success: true };
        } catch (error) {
            console.error("Failed to delete model:", error);
            return { success: false, error };
        }
    };

    return {
        formData,
        setFormData,
        updateField,
        resetForm,
        getSaveData,
        isValid,
        handleSubmitModel,
        handleDelete,
    };
}