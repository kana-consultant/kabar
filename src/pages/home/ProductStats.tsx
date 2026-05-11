import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Package, Wifi,  Clock,  Activity } from "lucide-react";
import type { Product } from "@/services/product";
import { cn } from "@/lib/utils";

interface ProductStatsProps {
    products: Product[];
}

export function ProductStats({ products }: ProductStatsProps) {
    const totalProducts = products.length;
    const connected = products.filter(p => p.status === "connected").length;
    const pending = products.filter(p => p.status === "pending").length
    const activeRate = totalProducts > 0 ? ((connected / totalProducts) * 100).toFixed(0) : 0;

    const statItems = [
        {
            label: "Total Produk",
            value: totalProducts,
            icon: Package,
            color: "text-indigo-500",
            bg: "bg-indigo-50 dark:bg-indigo-950/30",
        },
        {
            label: "Terhubung",
            value: connected,
            icon: Wifi,
            color: "text-emerald-500",
            bg: "bg-emerald-50 dark:bg-emerald-950/30",
        },
        {
            label: "Menunggu",
            value: pending,
            icon: Clock,
            color: "text-amber-500",
            bg: "bg-amber-50 dark:bg-amber-950/30",
        },
    ];

    return (
        <Card className="shadow-sm">
            <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                    <CardTitle className="text-lg font-semibold">
                        Ringkasan Produk
                    </CardTitle>
                    <div className="flex items-center gap-1 text-sm text-slate-500">
                        <Activity className="h-4 w-4" />
                        <span>Aktif: {activeRate}%</span>
                    </div>
                </div>
            </CardHeader>
            <CardContent>
                {/* Grid statistik */}
                <div className="grid grid-cols-2 lg:grid-cols-3 gap-4">
                    {statItems.map((item) => (
                        <div
                            key={item.label}
                            className="flex items-center gap-3 p-3 rounded-lg bg-slate-50 dark:bg-slate-900/50"
                        >
                            <div className={cn("rounded-full p-2", item.bg)}>
                                <item.icon className={cn("h-4 w-4", item.color)} />
                            </div>
                            <div>
                                <p className="text-xs text-slate-500 dark:text-slate-400">
                                    {item.label}
                                </p>
                                <p className="text-xl font-bold">
                                    {item.value}
                                </p>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Progress bar overall */}
                <div className="mt-4 pt-3 border-t">
                    <div className="flex justify-between text-xs mb-1">
                        <span className="text-slate-500">Koneksi Berhasil</span>
                        <span className="font-medium text-emerald-600 dark:text-emerald-400">
                            {activeRate}%
                        </span>
                    </div>
                    <div className="h-2 w-full rounded-full bg-slate-100 dark:bg-slate-800 overflow-hidden">
                        <div 
                            className="h-full rounded-full bg-gradient-to-r from-emerald-500 to-emerald-400 transition-all duration-500"
                            style={{ width: `${activeRate}%` }}
                        />
                    </div>
                    <div className="flex justify-between text-xs mt-2 text-slate-400">
                        <span>✓ {connected} Terhubung</span>
                        <span>⏳ {pending} Menunggu</span>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}