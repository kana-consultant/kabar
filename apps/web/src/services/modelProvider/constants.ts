export const API_ENDPOINTS = {
    PROVIDERS: '/providers',
    PROVIDER_BY_ID: (id: string) => `/providers/${id}`,
    MODELS: '/models',
    MODEL_BY_ID: (id: string) => `/models/${id}`,
    MODELS_DEFAULT: '/models/default',
} as const;

export const AUTH_TYPES = ['bearer', 'api_key', 'x-api-key'] as const;

export const providerKeys = {
    all: ['providers'] as const,
    lists: () => [...providerKeys.all, 'list'] as const,
    details: () => [...providerKeys.all, 'detail'] as const,
    detail: (id: string) => [...providerKeys.details(), id] as const,
};

export const modelKeys = {
    all: ['models'] as const,
    lists: () => [...modelKeys.all, 'list'] as const,
    details: () => [...modelKeys.all, 'detail'] as const,
    detail: (id: string) => [...modelKeys.details(), id] as const,
    default: () => [...modelKeys.all, 'default'] as const,
};

export const CACHE_CONFIG = {
    staleTime: 5 * 60 * 1000,
    gcTime: 10 * 60 * 1000,
};