import { Button } from "@kana-consultant/ui-kit";
import { Edit2, Trash2, CheckCircle2, AlertCircle, Package as PackageIcon } from "lucide-react";
import type { Product } from "@/services/product";
import { Can } from "@/components/ui/Can";

interface ProductListProps {
    products: Product[];
    testingId: string | null;
    onTest: (id: string) => void;
    onEdit: (product: Product) => void;
    onDelete: (product: Product) => void;
}

const getStatusBadge = (status: string) => {
    if (status === "connected") {
        return (
            <span className="inline-flex items-center gap-1 rounded-full bg-emerald-100 px-2 py-0.5 text-xs font-medium text-emerald-700 dark:bg-violet-900/30 dark:text-violet-400">
                <CheckCircle2 className="h-3 w-3" />
                Terhubung
            </span>
        );
    }
    return (
        <span className="inline-flex items-center gap-1 rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-700 dark:bg-amber-900/30 dark:text-amber-400">
            <AlertCircle className="h-3 w-3" />
            Menunggu
        </span>
    );
};

export function ProductList({ products, testingId, onTest, onEdit, onDelete }: ProductListProps) {
    if (products.length === 0) {
        return (
            <div className="rounded-2xl border border-slate-200 dark:border-white/10 bg-white dark:bg-slate-900 p-12 text-center">
                <PackageIcon className="mx-auto h-12 w-12 text-slate-300 dark:text-slate-600" />
                <p className="mt-2 text-slate-500 dark:text-slate-400">Belum ada produk</p>
            </div>
        );
    }

    return (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {products.map((product) => (
                <div key={product.id} className="rounded-2xl border border-slate-200 dark:border-white/10 bg-white dark:bg-slate-900 p-5 shadow-sm hover:shadow-md transition-all">
                    <div className="flex items-start justify-between gap-2 mb-3">
                        <div>
                            <h4 className="font-semibold text-slate-800 dark:text-white truncate" title={product.name}>
                                {product.name}
                            </h4>
                            <p className="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
                                {product.platform === "wordpress" && "WordPress"}
                                {product.platform === "shopify" && "Shopify"}
                                {product.platform === "custom" && "Custom API"}
                            </p>
                        </div>
                        {getStatusBadge(product.status)}
                    </div>

                    <div className="rounded-xl bg-slate-50 dark:bg-slate-800/50 p-2.5 mb-3">
                        <p className="font-mono text-xs truncate text-slate-500 dark:text-slate-400" title={product.api_endpoint}>
                            {product.api_endpoint}
                        </p>
                    </div>

                    <div className="flex justify-between text-xs text-slate-500 mb-4">
                        <span>Dibuat:</span>
                        <span className="font-medium text-slate-700 dark:text-slate-300">
                            {new Date(product.created_at).toLocaleDateString("id-ID")}
                        </span>
                    </div>

                    <div className="flex gap-2">
                        <Can permission="product:edit:team">
                            <Button
                                variant="outline"
                                size="sm"
                                className="flex-1 h-8 rounded-xl border-slate-200 dark:border-white/10 hover:bg-emerald-50 hover:text-emerald-600 dark:hover:bg-violet-900/20 dark:hover:text-violet-400"
                                onClick={() => onEdit(product)}
                            >
                                <Edit2 className="mr-1 h-3 w-3" />
                                Edit
                            </Button>
                        </Can>

                        <Can permission="product:delete:team">
                            <Button
                                variant="destructive"
                                size="sm"
                                className="flex-1 h-8 rounded-xl"
                                onClick={() => onDelete(product)}
                            >
                                <Trash2 className="mr-1 h-3 w-3" />
                                Hapus
                            </Button>
                        </Can>
                    </div>
                </div>
            ))}
        </div>
    );
}