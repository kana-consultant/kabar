// src/hooks/useProducts.ts
import { useState, useEffect, useCallback } from 'react';
import { toast } from 'sonner';
import { getProducts, deleteProduct, testConnection, syncProduct } from '@/services/productService';
import { type Product } from '@/types/product';

export function useProducts() {
    const [filteredProducts, setFilteredProducts] = useState<Product[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [statusFilter, setStatusFilter] = useState<string>("all");
    const [loading, setLoading] = useState(true);
    const [testingId, setTestingId] = useState<string | null>(null);
    const [syncingId, setSyncingId] = useState<string | null>(null);
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);

    // Products state
    const [products, setProducts] = useState<Product[]>([]);
    const [productNames, setProductNames] = useState<string[]>([]);
    const [productsLoading, setProductsLoading] = useState(true);
    const [productsError, setProductsError] = useState<string | null>(null);

    // Fetch products from backend
    const fetchProducts = useCallback(async () => {
        setProductsLoading(true);
        setProductsError(null);
        try {
            const productsData = await getProducts();
            setProducts(productsData || []);
            
            // Extract product names for display
            const names = productsData ? productsData.map(p => p.name) : [];
            setProductNames(names);
            
            console.log("Products loaded:", productsData);
        } catch (error) {
            console.error("Failed to fetch products:", error);
            setProductsError("Gagal memuat data produk");
            toast.error("Gagal memuat data produk", {
                description: "Silahkan refresh halaman",
            });
            // Fallback to default products if API fails
            const defaultProducts = ["TrekkingID", "CampingMart", "OutdoorGear"];
            setProductNames(defaultProducts);
        } finally {
            setProductsLoading(false);
        }
    }, []);

    // Load products on mount
    useEffect(() => {
        fetchProducts();
    }, [fetchProducts]);

    const loadProducts = useCallback(async () => {
        setLoading(true);
        try {
            const data = await getProducts();
            setProducts(data || []);
            setFilteredProducts(data || []);
            // Also update productNames
            setProductNames(data ? data.map(p => p.name) : []);
        } catch (error) {
            console.error('Failed to load products:', error);
            toast.error('Gagal memuat produk');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadProducts();
    }, [loadProducts]);

    // Filter products
    useEffect(() => {
        let filtered = [...products];
        
        if (searchQuery) {
            const query = searchQuery.toLowerCase();
            filtered = filtered.filter(p => 
                p.name.toLowerCase().includes(query) ||
                p.platform.toLowerCase().includes(query) ||
                p.apiEndpoint.toLowerCase().includes(query)
            );
        }
        
        if (statusFilter !== "all") {
            filtered = filtered.filter(p => p.status === statusFilter);
        }
        
        setFilteredProducts(filtered);
    }, [searchQuery, statusFilter, products]);

    const handleDelete = async (id: string, name: string) => {
        try {
            await deleteProduct(id);
            toast.success('Produk dihapus', { description: `"${name}" telah dihapus` });
            setShowDeleteDialog(false);
            setSelectedProduct(null);
            await loadProducts();
            await fetchProducts(); // Refresh both
        } catch (error) {
            toast.error('Gagal menghapus produk');
        }
    };

    const handleTestConnection = async (id: string) => {
        setTestingId(id);
        try {
            const result = await testConnection(id);
            if (result.success) {
                toast.success('Koneksi berhasil');
                await loadProducts();
                await fetchProducts();
            } else {
                toast.error(result.message || 'Koneksi gagal');
            }
        } catch (error) {
            toast.error('Gagal menguji koneksi');
        } finally {
            setTestingId(null);
        }
    };

    const handleSync = async (id: string) => {
        setSyncingId(id);
        try {
            const result = await syncProduct(id);
            if (result.success) {
                toast.success(result.message || 'Sinkronisasi berhasil');
                await loadProducts();
                await fetchProducts();
            } else {
                toast.error(result.message || 'Sinkronisasi gagal');
            }
        } catch (error) {
            toast.error('Gagal sinkronisasi');
        } finally {
            setSyncingId(null);
        }
    };

    return {
        // Untuk useGenerate
        products,
        productNames,
        productsLoading,
        productsError,
        
        // Untuk useProducts (filtering, dll)
        filteredProducts,
        searchQuery,
        setSearchQuery,
        statusFilter,
        setStatusFilter,
        loading,
        testingId,
        syncingId,
        showDeleteDialog,
        setShowDeleteDialog,
        selectedProduct,
        setSelectedProduct,
        
        // Actions
        loadProducts,
        handleDelete,
        handleTestConnection,
        handleSync,
        refreshProducts: fetchProducts,
    };
}