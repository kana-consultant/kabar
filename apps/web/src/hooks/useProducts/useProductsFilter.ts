import { useEffect } from "react";

export function useProductsFilter(
    products: any[],
    searchQuery: string,
    statusFilter: string,
    setFilteredProducts: (data: any[]) => void
) {
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
}