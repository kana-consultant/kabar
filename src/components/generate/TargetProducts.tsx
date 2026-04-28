import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import type { Product } from "@/types/product";

interface TargetProductsProps {
    products: Product[];
    selectedProducts: string[];
    postToAll: boolean;
    onToggleProduct: (product: string) => void;
    onSelectAll: () => void;
}

export function TargetProducts({
    products,
    selectedProducts,
    postToAll,
    onToggleProduct,
    onSelectAll,
}: TargetProductsProps) {
    return (
        <div>
            <div className="flex items-center justify-between mb-2">
                <Label>Target Produk</Label>
                <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={onSelectAll}
                >
                    {postToAll ? "Batalkan Semua" : "Pilih Semua"}
                </Button>
            </div>
            <div className="space-y-2">
                {products.map((product) => (
                    <label
                        key={product.id}
                        className={`flex cursor-pointer items-center gap-3 rounded-lg border p-3 transition-all ${selectedProducts.includes(product.id)
                                ? "border-blue-500 bg-blue-50 dark:bg-blue-950"
                                : "hover:bg-slate-50"
                            } ${postToAll ? "opacity-50" : ""}`}
                    >
                        <input
                            type="checkbox"
                            checked={selectedProducts.includes(product.id)}
                            onChange={() => onToggleProduct(product.id)}
                            disabled={postToAll}
                            className="h-4 w-4 rounded border-gray-300"
                        />
                        <span className="text-sm font-medium">{product.name}</span>
                    </label>
                ))}
            </div>
        </div>
    );
}