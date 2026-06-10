// hooks/useModelManagement.ts

import { useState, useCallback, useEffect } from "react";
import type { AIModel, CreateModelRequest } from "@/types/provider.types";
import {
    getModels,
    getModel,
    getDefaultModel,
    getModelsWithStatus,
    getModelsFromAPIKeys,
    createModel,
    updateModel,
    deleteModel,
} from "@/services/model";

interface ModelManagementState {
    models: AIModel[];
    selectedModel: AIModel | null;
    defaultModel: AIModel | null;
    modelsWithStatus: any[];
    modelsFromAPIKeys: any[];
    loading: boolean;
    error: string | null;
}

interface UseModelManagementOptions {
    autoLoad?: boolean;
    providerId?: string;
}

export function useModelManagement(options: UseModelManagementOptions = {}) {
    const { autoLoad = true, providerId } = options;
    
    const [state, setState] = useState<ModelManagementState>({
        models: [],
        selectedModel: null,
        defaultModel: null,
        modelsWithStatus: [],
        modelsFromAPIKeys: [],
        loading: false,
        error: null,
    });

    const setLoading = (loading: boolean) =>
        setState((prev) => ({ ...prev, loading, error: null }));

    const setError = (error: string) =>
        setState((prev) => ({ ...prev, loading: false, error }));

    // Initial load
    useEffect(() => {
        if (!autoLoad) return;
        
        let cancelled = false;

        async function init() {
            setLoading(true);
            try {
                const models = await getModels() as AIModel[];
                
                let filteredModels = models as AIModel[];
                if (providerId) {
                    filteredModels = models.filter((m ) => m.provider_id as string == providerId);
                }
                
                if (!cancelled) {
                    setState((prev) => ({
                        ...prev,
                        models: filteredModels,
                        loading: false,
                    }));
                }
            } catch (err) {
                if (!cancelled) {
                    setError(err instanceof Error ? err.message : "Failed to load models");
                }
            }
        }

        init();
        return () => { cancelled = true; };
    }, [autoLoad, providerId]);

    // Queries
    const loadModels = useCallback(async () => {
        setLoading(true);
        try {
            const models = await getModels() as AIModel[];
            let filteredModels = models;
            if (providerId) {
                filteredModels = models.filter((m) => m.provider_id === providerId);
            }
            setState((prev) => ({ ...prev, models: filteredModels, loading: false }));
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load models");
        }
    }, [providerId]);

    const loadModel = useCallback(async (id: string) => {
        setLoading(true);
        try {
            const model = await getModel(id);
            setState((prev) => ({ ...prev, selectedModel: model, loading: false }));
            return model;
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load model");
            throw err;
        }
    }, []);

    const loadDefaultModel = useCallback(async () => {
        setLoading(true);
        try {
            const defaultModel = await getDefaultModel();
            setState((prev) => ({ ...prev, defaultModel, loading: false }));
            return defaultModel;
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load default model");
            throw err;
        }
    }, []);

    const loadModelsWithStatus = useCallback(async () => {
        setLoading(true);
        try {
            const modelsWithStatus = await getModelsWithStatus();
            setState((prev) => ({ ...prev, modelsWithStatus, loading: false }));
            return modelsWithStatus;
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load models with status");
            throw err;
        }
    }, []);

    const loadModelsFromAPIKeys = useCallback(async () => {
        setLoading(true);
        try {
            const modelsFromAPIKeys = await getModelsFromAPIKeys();
            setState((prev) => ({ ...prev, modelsFromAPIKeys, loading: false }));
            return modelsFromAPIKeys;
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load models from API keys");
            throw err;
        }
    }, []);

    // Mutations
    const createNewModel = useCallback(async (data: CreateModelRequest) => {
        setLoading(true);
        try {
            const newModel = await createModel(data);
            setState((prev) => ({
                ...prev,
                models: [...prev.models, newModel as AIModel],
                loading: false,
            }));
            return newModel;
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to create model");
            throw err;
        }
    }, []);

    const updateModelById = useCallback(async (id: string, data: Partial<AIModel>) => {
        setLoading(true);
        try {
            const updated = await updateModel(id, data);
            setState((prev) => ({
                ...prev,
                models: prev.models.map((m) => (m.id === id ? updated as AIModel : m)),
                selectedModel: prev.selectedModel?.id === id ? updated as AIModel : prev.selectedModel,
                loading: false,
            }));
            return updated;
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to update model");
            throw err;
        }
    }, []);

    const deleteModelById = useCallback(async (id: string) => {
        setLoading(true);
        try {
            await deleteModel(id);
            setState((prev) => ({
                ...prev,
                loading: false,
                models: prev.models.filter((m) => m.id !== id),
                selectedModel: prev.selectedModel?.id === id ? null : prev.selectedModel,
                defaultModel: prev.defaultModel?.id === id ? null : prev.defaultModel,
            }));
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to delete model");
            throw err;
        }
    }, []);

    const toggleModelStatus = useCallback(async (id: string, isActive: boolean) => {
        await updateModelById(id, { is_active: isActive } as Partial<AIModel>);
    }, [updateModelById]);

    const setDefaultModel = useCallback(async (id: string) => {
        await updateModelById(id, { is_default: true } as Partial<AIModel>);
        // Reset default lainnya untuk provider yang sama
        const model = state.models.find(m => m.id === id);
        if (model) {
            const otherModels = state.models.filter(m => m.provider_id === model.provider_id && m.id !== id);
            for (const other of otherModels) {
                if (other.is_default) {
                    await updateModelById(other.id, { is_default: false } as Partial<AIModel>);
                }
            }
        }
    }, [updateModelById, state.models]);

    return {
        // State
        models: state.models,
        selectedModel: state.selectedModel,
        defaultModel: state.defaultModel,
        modelsWithStatus: state.modelsWithStatus,
        modelsFromAPIKeys: state.modelsFromAPIKeys,
        loading: state.loading,
        error: state.error,

        // Query Actions
        loadModels,
        loadModel,
        loadDefaultModel,
        loadModelsWithStatus,
        loadModelsFromAPIKeys,

        // Mutation Actions
        createModel: createNewModel,
        updateModel: updateModelById,
        deleteModel: deleteModelById,
        toggleModelStatus,
        setDefaultModel,
    };
}