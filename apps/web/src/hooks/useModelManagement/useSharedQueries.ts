// hooks/useModelManagement/useSharedQueries.ts

import { useQuery } from "@tanstack/react-query";
import { getProviders, getFamilies, getSchemas } from "@/services/model";

export const modelSharedKeys = {
    all: ['models'] as const,
    providers: () => [...modelSharedKeys.all, 'providers'] as const,
    families: () => [...modelSharedKeys.all, 'families'] as const,
    schemas: () => [...modelSharedKeys.all, 'schemas'] as const,
};

export function useSharedQueries() {
    const providersQuery = useQuery({
        queryKey: modelSharedKeys.providers(),
        queryFn: async () => {
            const response = await getProviders();
            return response.data || [];
        },
    });

    const familiesQuery = useQuery({
        queryKey: modelSharedKeys.families(),
        queryFn: async () => {
            const response = await getFamilies();
            return response.data || [];
        },
    });

    const schemasQuery = useQuery({
        queryKey: modelSharedKeys.schemas(),
        queryFn: async () => {
            const response = await getSchemas();
            return response.data || [];
        },
    });

    const isLoading = providersQuery.isLoading || familiesQuery.isLoading || schemasQuery.isLoading;
    const error = providersQuery.error || familiesQuery.error || schemasQuery.error;

    return {
        providers: providersQuery.data || [],
        families: familiesQuery.data || [],
        schemas: schemasQuery.data || [],
        isLoading,
        error,
        // Also return the queries if needed for refetching
        refetch: {
            providers: providersQuery.refetch,
            families: familiesQuery.refetch,
            schemas: schemasQuery.refetch,
        }
    };
}