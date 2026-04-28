import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Edit2, Trash2, Wifi, RefreshCw, Package as PackageIcon } from "lucide-react";
import type { Product } from "@/types/product";

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
            <span className="inline-flex items-center gap-1 rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700 dark:bg-green-900/30 dark:text-green-400">
                <Wifi className="h-3 w-3" />
                Terhubung
            </span>
        );
    }
    return (
        <span className="inline-flex items-center gap-1 rounded-full bg-yellow-100 px-2 py-0.5 text-xs font-medium text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400">
            <RefreshCw className="h-3 w-3" />
            Menunggu
        </span>
    );
};

export function ProductList({ products, testingId, onTest, onEdit, onDelete }: ProductListProps) {
    if (products.length === 0) {
        return (
            <Card>
                <CardContent className="py-12 text-center">
                    <PackageIcon className="mx-auto h-12 w-12 text-slate-300" />
                    <p className="mt-2 text-slate-500">Belum ada produk</p>
                </CardContent>
            </Card>
        );
    }

    return (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {products.map((product) => (
                <Card key={product.id} className="overflow-hidden transition-all hover:shadow-md">
                    <CardHeader className="pb-3">
                        <div className="flex items-start justify-between gap-2">
                            <div className="flex items-center gap-2 flex-1 min-w-0">
                                <CardTitle className="text-base truncate" title={product.name}>
                                    {product.name}
                                </CardTitle>
                            </div>
                            {getStatusBadge(product.status)}
                        </div>
                        <CardDescription className="flex items-center gap-1 mt-1">
                            {product.platform === "wordpress" && "WordPress"}
                            {product.platform === "shopify" && "Shopify"}
                            {product.platform === "custom" && "Custom API"}
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-3">
                        <div className="rounded-lg bg-slate-50 p-2 text-xs dark:bg-slate-900/50">
                            <p className="font-mono truncate" title={product.apiEndpoint}>
                                {product.apiEndpoint}
                            </p>
                        </div>
                        <div className="flex justify-between text-sm">
                            <span className="text-slate-500">Dibuat:</span>
                            <span className="font-medium">
                                {new Date(product.createdAt).toLocaleDateString("id-ID")}
                            </span>
                        </div>
                        <div className="flex gap-2 pt-2">
                            <Button
                                variant="outline"
                                size="sm"
                                className="flex-1 h-8"
                                onClick={() => onTest(product.id)}
                                disabled={testingId === product.id}
                            >
                                {testingId === product.id ? (
                                    <RefreshCw className="mr-1 h-3 w-3 animate-spin" />
                                ) : (
                                    <Wifi className="mr-1 h-3 w-3" />
                                )}
                                Test
                            </Button>
                            <Button
                                variant="outline"
                                size="sm"
                                className="flex-1 h-8"
                                onClick={() => onEdit(product)}
                            >
                                <Edit2 className="mr-1 h-3 w-3" />
                                Edit
                            </Button>
                            <Button
                                variant="destructive"
                                size="sm"
                                className="flex-1 h-8"
                                onClick={() => onDelete(product)}
                            >
                                <Trash2 className="mr-1 h-3 w-3" />
                                Hapus
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}