// hooks/useModelManagement/useEditModel.ts

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getModelSchema, updateModel } from "@/services/model";
import { useToast } from "../use-toast";
import { useNavigate } from "@tanstack/react-router";
import { useSharedQueries, modelSharedKeys } from "./useSharedQueries";

export const modelKeys = {
    all: ['models'] as const,
    lists: () => [...modelKeys.all, 'list'] as const,
    detail: (id: string) => [...modelKeys.all, 'detail', id] as const,
};

export function useEditModel(id: string) {
    const toast  = useToast();
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    const { providers, families, schemas, isLoading: isLoadingShared } = useSharedQueries();

    // Query untuk model detail
    const modelQuery = useQuery({
        queryKey: modelKeys.detail(id),
        queryFn: () => getModelSchema(id),
        enabled: !!id,
    });

    // Mutation untuk update
    const updateMutation = useMutation({
        mutationFn: async ({ id, ...data }: { id: string } & any) => {
            // Validasi data tidak kosong
            if (Object.keys(data).length === 0) {
                throw new Error('No data to update');
            }
            
            // Filter hanya field yang memiliki value
            const filteredData = Object.keys(data).reduce((acc, key) => {
                if (data[key] !== undefined && data[key] !== null) {
                    acc[key] = data[key];
                }
                return acc;
            }, {} as any);
            
            if (Object.keys(filteredData).length === 0) {
                throw new Error('No valid fields to update');
            }
            
            return updateModel(id, filteredData);
        },
        onSuccess: (_, { id }) => {
            // Invalidate queries
            queryClient.invalidateQueries({ queryKey: modelKeys.lists() });
            queryClient.invalidateQueries({ queryKey: modelKeys.detail(id) });
            
            // Show success message
            toast.success('Model updated successfully');
            
            // Navigate back to list
            navigate({ to: '/ai-management' });
        },
        onError: (error: any) => {
            console.error('Update model error:', error);
            toast.error(error.message || 'Failed to update model');
        },
    });

    // Helper function untuk update
    const handleUpdate = (updateData: any) => {
        if (!id) {
            toast.error('Model ID is required');
            return;
        }
        return updateMutation.mutateAsync({ id, ...updateData });
    };

    // Reset form function
    const resetForm = () => {
        modelQuery.refetch();
    };

    const isLoading = isLoadingShared || modelQuery.isLoading;

    return {
        // Data
        providers,
        families,
        schemas,
        model: modelQuery.data,
        
        // Status
        isLoading,
        isSubmitting: updateMutation.isPending,
        isError: modelQuery.isError || updateMutation.isError,
        error: modelQuery.error || updateMutation.error,
        
        // Actions
        updateModel: handleUpdate,
        resetForm,
        refetch: modelQuery.refetch,
        
        // Raw mutations for advanced use
        mutate: updateMutation.mutate,
        mutateAsync: updateMutation.mutateAsync,
    };
}