// hooks/useAIManagement/useProviderManagement.ts

import { useState, useCallback, useEffect } from "react";
import { useQueryClient, useMutation } from "@tanstack/react-query";
import type { APIProvider, Family, ProviderFormData, RequestSchemaFormData } from "@/types/provider.types";
import {
  createProvider,
  updateProvider,
  deleteProvider,
  toggleProviderStatus,
} from "@/services/provider";
import { getProviders, getProviderById } from "@/services/provider/provider-api";
import { createFamily, updateFamily, deleteFamily } from "@/services/family";
import { createSchema, updateSchema, deleteSchema } from "@/services/schema";
import { useToast } from "@/hooks/use-toast";

interface ProviderManagementState {
  providers: APIProvider[];
  selectedProvider: APIProvider | null;
  loading: boolean;
  error: string | null;
}

interface UseProviderManagementOptions {
  autoLoad?: boolean;
  includeInactive?: boolean;
}

export function useProviderManagement(options: UseProviderManagementOptions = {}) {
  const { autoLoad = true, includeInactive = false } = options;
  
  const [state, setState] = useState<ProviderManagementState>({
    providers: [],
    selectedProvider: null,
    loading: false,
    error: null,
  });

  // Helpers
  const setLoading = (loading: boolean) =>
    setState((prev) => ({ ...prev, loading, error: null }));

  const setError = (error: string) =>
    setState((prev) => ({ ...prev, loading: false, error }));

  const clearError = () =>
    setState((prev) => ({ ...prev, error: null }));

  // Initial load
  useEffect(() => {
    if (!autoLoad) return;
    
    let cancelled = false;

    async function init() {
      setLoading(true);
      try {
        const providers = await getProviders();
        if (!cancelled) {
          setState((prev) => ({
            ...prev,
            providers: providers.data || [],
            loading: false,
          }));
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "Failed to load providers");
        }
      }
    }

    init();
    return () => { cancelled = true; };
  }, [autoLoad]);

  // Queries
  const loadProviders = useCallback(async () => {
    setLoading(true);
    try {
      const providers = await getProviders();
      setState((prev) => ({ ...prev, providers: providers.data || [], loading: false }));
      return providers;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load providers");
      throw err;
    }
  }, []);

  const loadActiveProviders = useCallback(async () => {
    setLoading(true);
    try {
      const providers = await getActiveProvidersList();
      setState((prev) => ({ ...prev, providers: providers || [], loading: false }));
      return providers;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load active providers");
      throw err;
    }
  }, []);

  const loadProviderById = useCallback(async (id: string) => {
    setLoading(true);
    try {
      const provider = await getProviderById(id);
      setState((prev) => ({ ...prev, selectedProvider: provider, loading: false }));
      return provider;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load provider");
      throw err;
    }
  }, []);

  // Mutations
  const createNewProvider = useCallback(async (data: ProviderFormData) => {
    setLoading(true);
    try {
      const newProvider = await createProvider(data);
      setState((prev) => ({
        ...prev,
        providers: [...prev.providers, newProvider],
        loading: false,
      }));
      return newProvider;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create provider");
      throw err;
    }
  }, []);

  const updateProviderById = useCallback(async (id: string, data: ProviderFormData) => {
    setLoading(true);
    try {
      // Hapus families dari data jika ada (karena tidak boleh dikirim ke endpoint provider)
      const updated = await updateProvider(id, data);
      setState((prev) => ({
        ...prev,
        providers: prev.providers.map((p) => (p.id === id ? updated : p)),
        selectedProvider: prev.selectedProvider?.id === id ? updated : prev.selectedProvider,
        loading: false,
      }));
      return updated;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update provider");
      throw err;
    }
  }, []);

  const deleteProviderById = useCallback(async (id: string) => {
    setLoading(true);
    try {
      await deleteProvider(id);
      setState((prev) => ({
        ...prev,
        loading: false,
        providers: prev.providers.filter((p) => p.id !== id),
        selectedProvider: prev.selectedProvider?.id === id ? null : prev.selectedProvider,
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete provider");
      throw err;
    }
  }, []);

  const toggleProviderActiveStatus = useCallback(async (id: string, isActive: boolean) => {
    setLoading(true);
    try {
      const updated = await toggleProviderStatus(id, isActive);
      setState((prev) => ({
        ...prev,
        loading: false,
        providers: prev.providers.map((p) => (p.id === id ? updated : p)),
        selectedProvider: prev.selectedProvider?.id === id ? updated : prev.selectedProvider,
      }));
      return updated;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to toggle provider status");
      throw err;
    }
  }, []);

  // Helper functions
  const getProviderById_local = useCallback((id: string) => {
    return state.providers.find(p => p.id === id);
  }, [state.providers]);

  const getActiveProvidersList = useCallback(() => {
    return state.providers.filter(p => p.is_active);
  }, [state.providers]);

  const getInactiveProvidersList = useCallback(() => {
    return state.providers.filter(p => !p.is_active);
  }, [state.providers]);

  const resetProviders = useCallback(() => {
    setState((prev) => ({
      ...prev,
      providers: [],
      selectedProvider: null,
      loading: false,
      error: null,
    }));
  }, []);

  return {
    // State
    providers: state.providers,
    selectedProvider: state.selectedProvider,
    loading: state.loading,
    error: state.error,

    // Query Actions
    loadProviders,
    loadActiveProviders,
    loadProviderById,

    // Mutation Actions
    createProvider: createNewProvider,
    updateProvider: updateProviderById,
    deleteProvider: deleteProviderById,
    toggleProviderStatus: toggleProviderActiveStatus,

    // Helper Functions
    getProviderById: getProviderById_local,
    getActiveProviders: getActiveProvidersList,
    getInactiveProviders: getInactiveProvidersList,
    clearError,
    resetProviders,
  };
}

// ============================================================
// FAMILY OPERATIONS HOOK
// ============================================================

export function useFamilyOperations(providerId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();

  const createFamilyMutation = useMutation({
    mutationFn: (data: Family) => createFamily(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['providers', providerId] });
      queryClient.invalidateQueries({ queryKey: ['families'] });
      toast.success("Family created successfully");
    },
    onError: (error: any) => {
      toast.error(error.message || "Failed to create family");
    },
  });

  const updateFamilyMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Family> }) => 
      updateFamily(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['providers', providerId] });
      queryClient.invalidateQueries({ queryKey: ['families'] });
      toast.success("Family updated successfully");
    },
    onError: (error: any) => {
      toast.error(error.message || "Failed to update family");
    },
  });

  const deleteFamilyMutation = useMutation({
    mutationFn: (id: string) => deleteFamily(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['providers', providerId] });
      queryClient.invalidateQueries({ queryKey: ['families'] });
      toast.success("Family deleted successfully");
    },
    onError: (error: any) => {
      toast.error(error.message || "Failed to delete family");
    },
  });

  return {
    createFamily: createFamilyMutation.mutateAsync,
    updateFamily: updateFamilyMutation.mutateAsync,
    deleteFamily: deleteFamilyMutation.mutateAsync,
    isCreating: createFamilyMutation.isPending,
    isUpdating: updateFamilyMutation.isPending,
    isDeleting: deleteFamilyMutation.isPending,
  };
}

// ============================================================
// SCHEMA OPERATIONS HOOK
// ============================================================

export function useSchemaOperations() {
  const queryClient = useQueryClient();
  const toast = useToast();

  const createSchemaMutation = useMutation({
    mutationFn: (data: RequestSchemaFormData) => createSchema(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['schemas'] });
      queryClient.invalidateQueries({ queryKey: ['families'] });
      queryClient.invalidateQueries({ queryKey: ['providers'] });
      toast.success("Schema created successfully");
    },
    onError: (error: any) => {
      toast.error(error.message || "Failed to create schema");
    },
  });

  const updateSchemaMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<RequestSchemaFormData> }) => 
      updateSchema(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['schemas'] });
      queryClient.invalidateQueries({ queryKey: ['families'] });
      queryClient.invalidateQueries({ queryKey: ['providers'] });
      toast.success("Schema updated successfully");
    },
    onError: (error: any) => {
      toast.error(error.message || "Failed to update schema");
    },
  });

  const deleteSchemaMutation = useMutation({
    mutationFn: (id: string) => deleteSchema(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['schemas'] });
      queryClient.invalidateQueries({ queryKey: ['families'] });
      queryClient.invalidateQueries({ queryKey: ['providers'] });
      toast.success("Schema deleted successfully");
    },
    onError: (error: any) => {
      toast.error(error.message || "Failed to delete schema");
    },
  });

  return {
    createSchema: createSchemaMutation.mutateAsync,
    updateSchema: updateSchemaMutation.mutateAsync,
    deleteSchema: deleteSchemaMutation.mutateAsync,
    isCreating: createSchemaMutation.isPending,
    isUpdating: updateSchemaMutation.isPending,
    isDeleting: deleteSchemaMutation.isPending,
  };
}