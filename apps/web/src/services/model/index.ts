// services/model/index.ts

// Types
export type { 
    AIModel, 
    ModelWithStatus, 
    CreateModelRequest, 
    ModelFromAPIKey 
} from '@/types/provider.types';

// Queries (GET)
export { 
    getModels, 
    getDefaultModel, 
    getModelsWithStatus, 
    getModelsFromAPIKeys,
    getProviders,
    getFamilies,
    getSchemas,
    getModel,
    getModelSchema
} from './modelQueries';

// Mutations (POST, PUT, DELETE)
export { createModel, updateModel, deleteModel } from './modelMutations';