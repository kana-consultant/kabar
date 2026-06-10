// hooks/useModelManagement/useCreateModel.ts

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createModel } from "@/services/model";
import { modelKeys } from "@/services/modelProvider/constants";
import { useSharedQueries } from "./useSharedQueries";

export function useCreateModel() {
    const queryClient = useQueryClient();
    const { providers, families, schemas, isLoading: isLoadingShared } = useSharedQueries();

    // Mutation
    const createMutation = useMutation({
        mutationFn: createModel,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: modelKeys.lists() });
        },
        onError: (error: any) => {
            console.error('Create model error:', error);
        },
    });

    return {
        // Data
        providers,
        families,
        schemas,
        
        // States
        isLoading: isLoadingShared,
        isSubmitting: createMutation.isPending,
        
        // Actions
        createModel: createMutation.mutateAsync,
        
        // Errors
        error: createMutation.error,
    };
}