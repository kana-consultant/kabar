import { useEffect } from "react";
import { useProductsState } from "./useProductsState";
import { useProductsData } from "./useProductsData";
import { useProductsFilter } from "./useProductsFilter";
import { useProductsActions } from "./useProductsActions";

export function useProducts() {
    const {
        filteredProducts, setFilteredProducts,
        searchQuery, setSearchQuery,
        statusFilter, setStatusFilter,
        loading, setLoading,
        testingId, setTestingId,
        syncingId, setSyncingId,
        showDeleteDialog, setShowDeleteDialog,
        selectedProduct, setSelectedProduct,
        products, setProducts,
        productNames, setProductNames,
        productsLoading, setProductsLoading,
        productsError, setProductsError,
    } = useProductsState();

    const { fetchProducts, loadProducts } = useProductsData(
        setProducts, setProductNames, setProductsLoading, setProductsError,setLoading
    );

    // Load products on mount
    useEffect(() => {
        fetchProducts();
    }, [fetchProducts]);

    useEffect(() => {
        loadProducts();
    }, [loadProducts]);

    // Filter products
    useProductsFilter(products, searchQuery, statusFilter, setFilteredProducts);

    const { handleDelete, handleTestConnection, handleSync } = useProductsActions(
        loadProducts, fetchProducts, setShowDeleteDialog, setSelectedProduct,
        setTestingId, setSyncingId
    );

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