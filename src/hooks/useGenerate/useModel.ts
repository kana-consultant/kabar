// src/hooks/useGenerate/useModels.ts
import { useCallback } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { getAPIKeys } from "@/services/apiKey";
import { useState } from "react";
import type { APIKeyDetail } from "@/services/apiKey";

// Query keys
export const modelKeys = {
    all: ["models"] as const,
    lists: () => [...modelKeys.all, "list"] as const,
    details: () => [...modelKeys.all, "detail"] as const,
    byService: (service: string) => [...modelKeys.all, "service", service] as const,
    textModels: () => [...modelKeys.all, "text"] as const,
    imageModels: () => [...modelKeys.all, "image"] as const,
};

// Cache configuration
const CACHE_CONFIG = {
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000,   // 10 minutes
};

export function useModels() {
    const queryClient = useQueryClient();
    
   const [selectedModelId, setSelectedModelId] = useState("");

    // Query untuk mendapatkan models dari API keys
    const {
        data: models = [],
        isLoading: loadingModels,
        isError,
        error,
        refetch: refetchModels,
    } = useQuery({
        queryKey: modelKeys.lists(),
        queryFn: async () => {
            const data = await getAPIKeys();
            console.log("Models loaded:", data);
            return data as APIKeyDetail[];
        },
        staleTime: CACHE_CONFIG.staleTime,
        gcTime: CACHE_CONFIG.gcTime,
    });

    // Get text and image models separately
    const getTextAndImageModels = useCallback(async () => {
        // Try to get from cache first
        const cachedModels = queryClient.getQueryData<APIKeyDetail[]>(modelKeys.lists());
        
        if (cachedModels) {
            const articleModel = cachedModels.find((model: any) => model.service === 'text');
            const imageModel = cachedModels.find((model: any) => model.service === 'image');
            return { articleModel, imageModel, availableModels: cachedModels };
        }
        
        // Fallback to API call
        const data = await getAPIKeys();
        const availableModels = data as any;
        
        const articleModel = availableModels.find((model: any) => model.service === 'text');
        const imageModel = availableModels.find((model: any) => model.service === 'image');
        
        return { articleModel, imageModel, availableModels };
    }, [queryClient]);

    // Helper function untuk get provider color
    const getProviderColor = useCallback((providerName?: string) => {
        switch (providerName?.toLowerCase()) {
            case "openai":
                return "bg-green-100 text-green-700";
            case "anthropic":
                return "bg-purple-100 text-purple-700";
            case "google":
                return "bg-blue-100 text-blue-700";
            case "openrouter":
                return "bg-orange-100 text-orange-700";
            case "cohere":
                return "bg-teal-100 text-teal-700";
            default:
                return "bg-gray-100 text-gray-700";
        }
    }, []);

    // Helper function untuk get service badge color
    const getServiceBadgeColor = useCallback((service: string) => {
        switch (service) {
            case "text":
                return "bg-blue-100 text-blue-700";
            case "image":
                return "bg-purple-100 text-purple-700";
            case "embedding":
                return "bg-green-100 text-green-700";
            default:
                return "bg-gray-100 text-gray-700";
        }
    }, []);

    // Get text models only
    const getTextModels = useCallback(() => {
        return models.filter((model: any) => model.service === 'text');
    }, [models]);

    // Get image models only
    const getImageModels = useCallback(() => {
        return models.filter((model: any) => model.service === 'image');
    }, [models]);

    // Get active models
    const getActiveModels = useCallback(() => {
        return models.filter((model: any) => model.isActive);
    }, [models]);

    // Get model by ID
    const getModelById = useCallback((id: string) => {
        return models.find((model: any) => model.id === id);
    }, [models]);

    // Refetch function
    const refetch = useCallback(() => {
        return refetchModels();
    }, [refetchModels]);

    // Manual reload with toast
    const reloadModels = useCallback(async () => {
        try {
            await refetchModels();
            toast.success("Models refreshed", {
                description: "AI models have been reloaded",
            });
        } catch (error) {
            console.error("Failed to reload models:", error);
            toast.error("Failed to reload models");
        }
    }, [refetchModels]);

    // Computed values
    const textModels = getTextModels();
    const imageModels = getImageModels();
    const activeModels = getActiveModels();

    return {
        // Data
        models,
        
        // Loading & Error states
        loadingModels,
        isLoading: loadingModels,
        isError,
        error: error as Error | null,

        selectedModelId,
        setSelectedModelId,
        
        // Helper functions for UI
        getProviderColor,
        getServiceBadgeColor,
        
        // Filter functions
        getTextModels,
        getImageModels,
        getActiveModels,
        getModelById,
        
        // Actions
        refetch,
        reloadModels,
        getTextAndImageModels,
        
        // Computed values
        textModels,
        imageModels,
        activeModels,
        hasModels: models.length > 0,
        hasTextModels: textModels.length > 0,
        hasImageModels: imageModels.length > 0,
        
        // Default models
        defaultTextModel: textModels.find((m: any) => m.isDefault) || textModels[0],
        defaultImageModel: imageModels.find((m: any) => m.isDefault) || imageModels[0],
    };
}

// Hook untuk text models only
export function useTextModels() {
    const { textModels, loadingModels, ...rest } = useModels();
    
    return {
        textModels: textModels,
        loadingModels,
        ...rest,
    };
}

// Hook untuk image models only
export function useImageModels() {
    const { imageModels, loadingModels, ...rest } = useModels();
    
    return {
        imageModels: imageModels,
        loadingModels,
        ...rest,
    };
}

// Hook untuk active models only
export function useActiveModels() {
    const { activeModels, loadingModels, ...rest } = useModels();
    
    return {
        activeModels: activeModels,
        loadingModels,
        ...rest,
    };
}