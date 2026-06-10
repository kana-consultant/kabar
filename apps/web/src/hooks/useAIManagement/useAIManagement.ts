// hooks/useAIManagement.ts

import { useState, useCallback, useEffect } from "react";
import type { APIProvider, AIModel } from "@/types/provider.types";
import { getProviders, getActiveProviders } from "@/services/provider/providerMutations";
import {
    deleteProvider as deleteProviderApi,
} from "../../services/provider";
import { getModels } from "@/services/modelProvider";

// TODO: ganti dengan import dari services/model ketika sudah dibuat
// import { getModels, deleteModel as deleteModelApi, toggleModelStatus as toggleModelStatusApi, setDefaultModel as setDefaultModelApi } from "../../services/model";

// ─────────────────────────────────────────────
// State shape
// ─────────────────────────────────────────────

interface AIManagementState {
    providers: APIProvider[];
    models: AIModel[];
    loading: boolean;
    error: string | null;
}

// ─────────────────────────────────────────────
// Hook
// ─────────────────────────────────────────────

export function useAIManagement() {
    const [state, setState] = useState<AIManagementState>({
        providers: [],
        models: [],
        loading: false,
        error: null,
    });

    // ── Helpers ────────────────────────────────

    const setLoading = (loading: boolean) =>
        setState((prev) => ({ ...prev, loading, error: null }));

    const setError = (error: string) =>
        setState((prev) => ({ ...prev, loading: false, error }));

    // ── Initial load ───────────────────────────

    useEffect(() => {
        let cancelled = false;

        async function init() {
            setLoading(true);
            try {
                const [providers,models] = await Promise.all([
                    getProviders(),
                    getModels(),
                ]);

                if (!cancelled) {
                    setState((prev : any) => ({
                        ...prev,
                        providers : providers.data,
                        models : models.data,
                        loading: false,
                    }));
                }
            } catch (err) {
                if (!cancelled) {
                    setError(err instanceof Error ? err.message : "Failed to load data");
                }
            }
        }

        init();
        return () => { cancelled = true; };
    }, []);

    // ── Providers ──────────────────────────────

    const loadProviders = useCallback(async () => {
        setLoading(true);
        try {
            const resp = await getProviders();
            setState((prev) => ({ ...prev, providers: resp.data as APIProvider[], loading: false }));
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load providers");
        }
    }, []);

    const loadActiveProviders = useCallback(async () => {
        setLoading(true);
        try {
            const resp = await getActiveProviders();
            setState((prev) => ({ ...prev, providers: resp.data as APIProvider[], loading: false }));
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load active providers");
        }
    }, []);

    const deleteProvider = useCallback(async (id: string) => {
        setLoading(true);
        try {
            await deleteProviderApi(id);
            setState((prev) => ({
                ...prev,
                loading: false,
                providers: prev.providers.filter((p) => p.id !== id),
                models: prev.models.filter((m) => m.provider_id  !== id),
            }));
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to delete provider");
        }
    }, []);

  
    // ── Models ─────────────────────────────────

    const loadModels = useCallback(async () => {
        setLoading(true);
        try {
            // const data = await getModels();
            // setState((prev) => ({ ...prev, models: data, loading: false }));
            setState((prev) => ({ ...prev, loading: false }));
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load models");
        }
    }, []);

    const deleteModel = useCallback(async (id: string) => {
        setLoading(true);
        try {
            // await deleteModelApi(id);
            setState((prev) => ({
                ...prev,
                loading: false,
                models: prev.models.filter((m) => m.id !== id),
            }));
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to delete model");
        }
    }, []);

    const toggleModelStatus = useCallback(async (id: string, isActive: boolean) => {
        setLoading(true);
        try {
            // const updated = await toggleModelStatusApi(id, isActive);
            setState((prev) => ({
                ...prev,
                loading: false,
                models: prev.models.map((m) =>
                    m.id === id ? { ...m, is_active: isActive } : m
                ),
            }));
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to toggle model status");
        }
    }, []);

    const setDefaultModel = useCallback(async (id: string) => {
        setLoading(true);
        try {
            // await setDefaultModelApi(id);
            setState((prev) => ({
                ...prev,
                loading: false,
                models: prev.models.map((m) => ({ ...m, is_default: m.id === id })),
            }));
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to set default model");
        }
    }, []);

    // ── Return ─────────────────────────────────

    return {
        // State
        providers: state.providers,
        models: state.models,
        loading: state.loading,
        error: state.error,

        // Provider actions
        loadProviders,
        loadActiveProviders,
        deleteProvider,

        // Model actions
        loadModels,
        deleteModel,
        toggleModelStatus,
        setDefaultModel,
    };
}