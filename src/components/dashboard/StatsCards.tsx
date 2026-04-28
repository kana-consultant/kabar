// src/components/dashboard/StatsCards.tsx
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { FileText, Package, Send, TrendingUp } from "lucide-react";
import { useEffect, useState } from "react";
import { getDashboardStats, type DashboardStats } from "@/services/dashboardService";
import { Skeleton } from "@/components/ui/skeleton";

interface StatItem {
    title: string;
    value: number | string;
    icon: React.ElementType;
    change: string;
    color: string;
}

export function StatsCards() {
    const [stats, setStats] = useState<DashboardStats | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchStats = async () => {
            try {
                setLoading(true);
                const data = await getDashboardStats();
                setStats(data);
                setError(null);
            } catch (err) {
                console.error("Failed to fetch stats:", err);
                setError("Gagal memuat statistik");
            } finally {
                setLoading(false);
            }
        };

        fetchStats();
    }, []);

    if (loading) {
        return <StatsCardsSkeleton />;
    }

    if (error) {
        return (
            <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-center text-red-600 dark:border-red-800 dark:bg-red-950 dark:text-red-400">
                <p>{error}</p>
                <button 
                    onClick={() => window.location.reload()} 
                    className="mt-2 text-sm underline hover:no-underline"
                >
                    Coba lagi
                </button>
            </div>
        );
    }

    const statItems: StatItem[] = [
        {
            title: "Total Konten",
            value: stats?.totalContent || 0,
            icon: FileText,
            change: stats?.contentChange || "0 minggu ini",
            color: "text-blue-600",
        },
        {
            title: "Total Produk",
            value: stats?.totalProducts || 0,
            icon: Package,
            change: stats?.productsChange || "0 bulan ini",
            color: "text-green-600",
        },
        {
            title: "Konten Terkirim",
            value: stats?.totalPublished || 0,
            icon: Send,
            change: `${stats?.publishedPercentage || 0}% sukses`,
            color: "text-purple-600",
        },
        {
            title: "SEO Score Rata-rata",
            value: Math.round(stats?.averageSeoScore || 0),
            icon: TrendingUp,
            change: stats?.seoScoreChange || "0%",
            color: "text-orange-600",
        },
    ];

    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {statItems.map((stat) => (
                <Card key={stat.title}>
                    <CardHeader className="flex flex-row items-center justify-between pb-2">
                        <CardTitle className="text-sm font-medium text-slate-500">
                            {stat.title}
                        </CardTitle>
                        <stat.icon className={`h-4 w-4 ${stat.color}`} />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{stat.value}</div>
                        <p className="text-xs text-slate-500">{stat.change}</p>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}

function StatsCardsSkeleton() {
    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {[1, 2, 3, 4].map((i) => (
                <Card key={i}>
                    <CardHeader className="flex flex-row items-center justify-between pb-2">
                        <Skeleton className="h-4 w-24" />
                        <Skeleton className="h-4 w-4 rounded-full" />
                    </CardHeader>
                    <CardContent>
                        <Skeleton className="h-8 w-16 mb-1" />
                        <Skeleton className="h-3 w-20" />
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}