import { FileText, Package, Send, TrendingUp } from "lucide-react";
import { useEffect, useState } from "react";
import { getDashboardStats, type DashboardStats } from "@/services/dashboardService";
import { cn } from "@/lib/utils";

const statConfig = [
    {
        key: "totalContent" as const,
        label: "Total Konten",
        icon: FileText,
        iconClass: "bg-blue-50 text-blue-600 ring-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:ring-blue-500/20",
        dotClass: "bg-blue-500",
        pingClass: "bg-blue-500",
        getChange: (s: DashboardStats) => s.contentChange || "0 minggu ini",
        getValue: (s: DashboardStats) => s.totalContent || 0,
    },
    {
        key: "totalProducts" as const,
        label: "Total Produk",
        icon: Package,
        iconClass: "bg-green-50 text-green-600 ring-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:ring-green-500/20",
        dotClass: "bg-green-500",
        pingClass: "bg-green-500",
        getChange: (s: DashboardStats) => s.productsChange || "0 bulan ini",
        getValue: (s: DashboardStats) => s.totalProducts || 0,
    },
    {
        key: "totalPublished" as const,
        label: "Konten Terkirim",
        icon: Send,
        iconClass: "bg-violet-50 text-violet-600 ring-violet-200/60 dark:bg-violet-500/10 dark:text-violet-400 dark:ring-violet-500/20",
        dotClass: "bg-violet-500",
        pingClass: "bg-violet-500",
        getChange: (s: DashboardStats) => `${s.publishedPercentage || 0}% sukses`,
        getValue: (s: DashboardStats) => s.totalPublished || 0,
    },
    {
        key: "averageSeoScore" as const,
        label: "SEO Score Rata-rata",
        icon: TrendingUp,
        iconClass: "bg-orange-50 text-orange-600 ring-orange-200/60 dark:bg-orange-500/10 dark:text-orange-400 dark:ring-orange-500/20",
        dotClass: "bg-orange-500",
        pingClass: "bg-orange-500",
        getChange: (s: DashboardStats) => s.seoScoreChange || "0%",
        getValue: (s: DashboardStats) => Math.round(s.averageSeoScore || 0),
    },
];

export function StatsCards() {
    const [stats, setStats] = useState<DashboardStats | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetch = async () => {
            try {
                setLoading(true);
                setStats(await getDashboardStats());
                setError(null);
            } catch {
                setError("Gagal memuat statistik");
            } finally {
                setLoading(false);
            }
        };
        fetch();
    }, []);

    if (loading) return <StatsCardsSkeleton />;

    if (error) return (
        <div className="rounded-xl border border-red-200/80 bg-red-50 p-4 text-center dark:border-red-500/20 dark:bg-red-500/10">
            <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
            <button
                onClick={() => window.location.reload()}
                className="mt-1.5 text-xs text-red-500 underline hover:no-underline"
            >
                Coba lagi
            </button>
        </div>
    );

    return (
        <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-4">
            {statConfig.map(({ label, icon: Icon, iconClass, pingClass, getChange, getValue }) => (
                <div
                    key={label}
                    className={cn(
                        "flex flex-col gap-3 rounded-xl border p-5 transition-all duration-200",
                        "bg-white border-slate-200/80 hover:border-slate-300/80 hover:shadow-sm",
                        "dark:bg-[#0f0d1a] dark:border-white/[0.06] dark:hover:border-white/[0.10]"
                    )}
                >
                    <div className="flex items-start justify-between">
                        <span className="text-xs font-medium text-slate-400 dark:text-slate-500 tracking-wide">
                            {label}
                        </span>
                        <div className={cn(
                            "flex h-8 w-8 items-center justify-center rounded-lg ring-1",
                            iconClass
                        )}>
                            <Icon className="h-3.5 w-3.5" />
                        </div>
                    </div>

                    <span className="text-3xl font-semibold tracking-tight text-slate-900 dark:text-white tabular-nums">
                        {stats ? getValue(stats) : 0}
                    </span>

                    {/* Ping dot footer */}
                    <div className="flex items-center gap-1.5">
                        <span className="relative flex h-1.5 w-1.5">
                            <span className={cn(
                                "animate-ping absolute inline-flex h-full w-full rounded-full opacity-75",
                                pingClass
                            )} />
                            <span className={cn(
                                "relative inline-flex h-1.5 w-1.5 rounded-full",
                                pingClass
                            )} />
                        </span>
                        <span className="text-[11px] text-slate-400 dark:text-slate-600">
                            {stats ? getChange(stats) : "—"}
                        </span>
                    </div>
                </div>
            ))}
        </div>
    );
}

function StatsCardsSkeleton() {
    return (
        <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-4">
            {[1, 2, 3, 4].map((i) => (
                <div
                    key={i}
                    className={cn(
                        "flex flex-col gap-3 rounded-xl border p-5",
                        "bg-white border-slate-200/80",
                        "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
                    )}
                >
                    <div className="flex items-start justify-between">
                        <div className="h-3 w-24 rounded-md bg-slate-100 dark:bg-white/[0.05] animate-pulse" />
                        <div className="h-8 w-8 rounded-lg bg-slate-100 dark:bg-white/[0.05] animate-pulse" />
                    </div>
                    <div className="h-8 w-14 rounded-md bg-slate-100 dark:bg-white/[0.05] animate-pulse" />
                    <div className="h-2.5 w-20 rounded-md bg-slate-100 dark:bg-white/[0.05] animate-pulse" />
                </div>
            ))}
        </div>
    );
}