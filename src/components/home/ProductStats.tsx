import { Card, CardContent } from "@/components/ui/card";
import { Package, Wifi, WifiOff, Settings2 } from "lucide-react";
import type { Product } from "@/types/product";

interface ProductStatsProps {
    products: Product[];
}

export function ProductStats({ products }: ProductStatsProps) {
    const stats = [
        {
            title: "Total Produk",
            value: products.length,
            icon: Package,
            color: "text-blue-500",
            bg: "bg-blue-50 dark:bg-blue-950/30",
        },
        {
            title: "Terhubung",
            value: products.filter(p => p.status === "connected").length,
            icon: Wifi,
            color: "text-green-500",
            bg: "bg-green-50 dark:bg-green-950/30",
        },
        {
            title: "Menunggu",
            value: products.filter(p => p.status === "pending").length,
            icon: WifiOff,
            color: "text-yellow-500",
            bg: "bg-yellow-50 dark:bg-yellow-950/30",
        },
        {
            title: "Platform",
            value: new Set(products.map(p => p.platform)).size,
            icon: Settings2,
            color: "text-purple-500",
            bg: "bg-purple-50 dark:bg-purple-950/30",
        },
    ];

    return (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {stats.map((stat) => (
                <Card key={stat.title} className="overflow-hidden">
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-slate-500 dark:text-slate-400">{stat.title}</p>
                                <p className="text-2xl md:text-3xl font-bold mt-1">{stat.value}</p>
                            </div>
                            <div className={`rounded-full p-3 ${stat.bg}`}>
                                <stat.icon className={`h-5 w-5 ${stat.color}`} />
                            </div>
                        </div>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}