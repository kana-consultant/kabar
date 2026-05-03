import { useState } from "react";
import type { Product } from "@/types/product";

export function useProductsState() {
    const [filteredProducts, setFilteredProducts] = useState<Product[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [statusFilter, setStatusFilter] = useState<string>("all");
    const [loading, setLoading] = useState(true);
    const [testingId, setTestingId] = useState<string | null>(null);
    const [syncingId, setSyncingId] = useState<string | null>(null);
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
    const [products, setProducts] = useState<Product[]>([]);
    const [productNames, setProductNames] = useState<string[]>([]);
    const [productsLoading, setProductsLoading] = useState(true);
    const [productsError, setProductsError] = useState<string | null>(null);

    return {
        // Filter state
        filteredProducts, setFilteredProducts,
        searchQuery, setSearchQuery,
        statusFilter, setStatusFilter,
        loading, setLoading,
        testingId, setTestingId,
        syncingId, setSyncingId,
        showDeleteDialog, setShowDeleteDialog,
        selectedProduct, setSelectedProduct,
        
        // Products data state
        products, setProducts,
        productNames, setProductNames,
        productsLoading, setProductsLoading,
        productsError, setProductsError,
    };
}