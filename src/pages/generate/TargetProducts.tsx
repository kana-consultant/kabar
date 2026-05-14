import { Button } from "@/components/ui/button";
import { Check } from "lucide-react";
import type { Product } from "@/types/product";
import { cn } from "@/lib/utils";

interface TargetProductsProps {
    products: Product[];
    selectedProducts: string[];
    postToAll: boolean;
    onToggleProduct: (product: string) => void;
    onSelectAll: () => void;
}

export function TargetProducts({
    products, selectedProducts, postToAll,
    onToggleProduct, onSelectAll,
}: TargetProductsProps) {
    return (
        <div>
            <div className="flex items-center justify-between mb-3">
                <p className="text-[11px] font-medium uppercase tracking-wide text-slate-400 dark:text-slate-600">
                    Target Produk
                </p>
                <button
                    type="button"
                    onClick={onSelectAll}
                    className="text-[11px] font-medium text-green-600 hover:text-green-700 dark:text-purple-400 dark:hover:text-purple-300 transition-colors"
                >
                    {postToAll ? "Batalkan Semua" : "Pilih Semua"}
                </button>
            </div>

            <div className="space-y-1.5">
                {products.map((product) => {
                    const checked = selectedProducts.includes(product.id);
                    return (
                        <label
                            key={product.id}
                            className={cn(
                                "flex cursor-pointer items-center gap-2.5 rounded-lg border px-3 py-2.5 transition-all duration-150",
                                checked
                                    ? "border-green-200/80 bg-green-50 dark:border-purple-500/30 dark:bg-purple-500/10"
                                    : "border-slate-200/80 hover:bg-slate-50/80 dark:border-white/[0.06] dark:hover:bg-white/[0.03]",
                                postToAll && "opacity-50 cursor-not-allowed"
                            )}
                        >
                            <div className={cn(
                                "flex h-4 w-4 shrink-0 items-center justify-center rounded border transition-colors",
                                checked
                                    ? "bg-green-600 border-green-600 dark:bg-purple-600 dark:border-purple-600"
                                    : "border-slate-300 dark:border-white/[0.15]"
                            )}>
                                {checked && <Check className="h-2.5 w-2.5 text-white" strokeWidth={3} />}
                            </div>
                            <input
                                type="checkbox"
                                checked={checked}
                                onChange={() => onToggleProduct(product.id)}
                                disabled={postToAll}
                                className="sr-only"
                            />
                            <span className="text-sm font-medium text-slate-700 dark:text-slate-200">
                                {product.name}
                            </span>
                        </label>
                    );
                })}
            </div>
        </div>
    );
}