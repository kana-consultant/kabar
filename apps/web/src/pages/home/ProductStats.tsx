import { Package, Wifi, Clock, Activity } from "lucide-react";
import type { Product } from "@/services/product";
import { cn } from "@/lib/utils";

interface ProductStatsProps {
    products: Product[];
}

export function ProductStats({ products }: ProductStatsProps) {
    const totalProducts = products.length;
    const connected = products.filter(p => p.status === "connected").length;
    const pending = products.filter(p => p.status === "pending").length;
    const activeRate = totalProducts > 0 ? ((connected / totalProducts) * 100).toFixed(0) : 0;

    const statItems = [
        {
            label: "Total Produk",
            value: totalProducts,
            icon: Package,
            color: "text-emerald-600 dark:text-violet-400",
            bg: "bg-emerald-50 dark:bg-violet-950/30",
        },
        {
            label: "Terhubung",
            value: connected,
            icon: Wifi,
            color: "text-emerald-600 dark:text-violet-400",
            bg: "bg-emerald-50 dark:bg-violet-950/30",
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
        <div className="rounded-2xl border border-slate-200 dark:border-white/10 bg-white dark:bg-slate-900 p-6 shadow-sm">
            <div className="flex items-center justify-between mb-5">
                <h3 className="text-lg font-semibold text-slate-800 dark:text-white">Ringkasan Produk</h3>
                <div className="flex items-center gap-1 text-sm text-slate-500">
                    <Activity className="h-4 w-4" />
                    <span>Aktif: {activeRate}%</span>
                </div>
            </div>

            <div className="grid grid-cols-2 lg:grid-cols-3 gap-4">
                {statItems.map((item) => (
                    <div key={item.label} className="flex items-center gap-3 p-3 rounded-xl bg-slate-50 dark:bg-slate-800/50">
                        <div className={cn("rounded-xl p-2", item.bg)}>
                            <item.icon className={cn("h-4 w-4", item.color)} />
                        </div>
                        <div>
                            <p className="text-xs text-slate-500 dark:text-slate-400">{item.label}</p>
                            <p className="text-xl font-bold text-slate-800 dark:text-white">{item.value}</p>
                        </div>
                    </div>
                ))}
            </div>

            <div className="mt-4 pt-3 border-t border-slate-100 dark:border-white/10">
                <div className="flex justify-between text-xs mb-1">
                    <span className="text-slate-500">Koneksi Berhasil</span>
                    <span className="font-medium text-emerald-600 dark:text-violet-400">{activeRate}%</span>
                </div>
                <div className="h-2 w-full rounded-full bg-slate-100 dark:bg-slate-800 overflow-hidden">
                    <div
                        className="h-full rounded-full bg-gradient-to-r from-emerald-500 to-emerald-400 dark:from-violet-500 dark:to-purple-400 transition-all duration-500"
                        style={{ width: `${activeRate}%` }}
                    />
                </div>
                <div className="flex justify-between text-xs mt-2 text-slate-400">
                    <span>✓ {connected} Terhubung</span>
                    <span>⏳ {pending} Menunggu</span>
                </div>
            </div>
        </div>
    );
}