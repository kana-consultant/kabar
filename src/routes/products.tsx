import { useNavigate } from "@tanstack/react-router";
import { useProducts } from "@/hooks/useProducts";
import { ProductHeader } from "../components/home/ProductHeader";
import { ProductStats } from "../components/home/ProductStats";
import { ProductList } from "../components/home/ProductList";
import { DeleteProductDialog } from "@/components/home/DeleteProductDialog";
import { LoadingCard } from "@/components/ui/LoadingCard";
import type { Product } from "@/types/product";
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute("/products")({
    component: Products,
});

export default function Products() {
    const navigate = useNavigate();
    const {
        products,
        filteredProducts,
        searchQuery,
        setSearchQuery,
        loading,
        showDeleteDialog,
        setShowDeleteDialog,
        selectedProduct,
        handleDelete,
        handleTestConnection,
        testingId,
        setSelectedProduct
    } = useProducts();

    // Loading state with full page loader
    if (loading) {
        return <LoadingCard />;
    }

    return (
        <div className="max-w-7xl mx-auto space-y-6 ">
            <ProductHeader
                searchQuery={searchQuery}
                setSearchQuery={setSearchQuery}
                onAddProduct={() => navigate({ to: "/products/add" })}
            />

            <ProductStats products={products} />

            <ProductList
                products={filteredProducts }
                testingId={testingId}
                onTest={handleTestConnection}
                onEdit={(product: Product) => navigate({ to: `/products/${product.id}/edit` })}
                onDelete={(product: Product) => {
                    setSelectedProduct(product);
                    setShowDeleteDialog(true);
                }}
            />

            <DeleteProductDialog
                open={showDeleteDialog}
                onOpenChange={setShowDeleteDialog}
                product={selectedProduct}
                onDelete={handleDelete}
            />
        </div>
    );
}