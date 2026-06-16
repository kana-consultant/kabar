// Types
export type { 
    ProductRequest, 
    UpdateProductRequest, 
    AddProductResponse,
    Product,
    AdapterConfig,
    FieldMapping
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
    updateProduct, 
    deleteProduct,
    testConnection,
    syncProduct
} from './productMutations';