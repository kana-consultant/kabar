// Types
export type { 
    CreateProductRequest, 
    UpdateProductRequest, 
    AddProductResponse 
} from './types';

// Queries (GET)
export { 
    getProducts, 
    getProductById, 
    getProductsByTeam,
    getProductsByPlatform,
    getProductsByStatus,
    getProductsBySyncStatus,
    getConnectedProducts,
    getProductsNeedingSync
} from './productQueries';

// Mutations (POST, PUT, DELETE)
export { 
    createProduct, 
    addProduct, 
    saveProduct, 
    updateProduct, 
    deleteProduct,
    testConnection,
    syncProduct
} from './productMutations';