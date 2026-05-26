import { useState, useMemo, useEffect, useRef } from "react";
import {
    FileText, Calendar, CheckCircle, ImageIcon, ImageOff,
    BarChart2, TrendingUp, Sparkles, ChevronDown
} from "lucide-react";
import { cn } from "@/lib/utils";
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    LineElement,
    PointElement,
    ArcElement,
    Tooltip,
    Filler,
    type ChartOptions
} from "chart.js";
import { Bar, Line, Doughnut } from "react-chartjs-2";

ChartJS.register(
    CategoryScale,
    LinearScale,
    BarElement,
    LineElement,
    PointElement,
    ArcElement,
    Tooltip,
    Filler
);

interface DailyActivity {
    date: string;
    count: number;
}

interface WeeklyTrend {
    week: string;
    withImage: number;
    withoutImage: number;
}

interface DraftStatsProps {
    totalDraft: number;
    totalWithImage: number;
    totalWithoutImage: number;
    totalScheduled: number;
    productCoverage: Record<string, number>;
    dailyActivity: DailyActivity[];
    weeklyTrend?: WeeklyTrend[];
    onRefresh?: () => void;
    onExport?: () => void;
    isLoading?: boolean;
}

type TimeRange = "7d" | "30d" | "90d";

export function DraftStats({
    totalDraft,
    totalWithImage,
    totalWithoutImage,
    totalScheduled,
    productCoverage,
    dailyActivity,
    weeklyTrend = [],
    onRefresh,
    onExport,
    isLoading = false,
}: DraftStatsProps) {
    const [timeRange, setTimeRange] = useState<TimeRange>("30d");
    const isDark = typeof window !== "undefined"
        ? window.matchMedia("(prefers-color-scheme: dark)").matches
        : false;

    const textColor = isDark ? "rgba(255,255,255,0.45)" : "rgba(0,0,0,0.45)";
    const gridColor = isDark ? "rgba(255,255,255,0.07)" : "rgba(0,0,0,0.07)";

    const filteredActivity = useMemo(() => {
        const daysMap = { "7d": 7, "30d": 30, "90d": 90 };
        return dailyActivity.slice(-daysMap[timeRange]);
    }, [dailyActivity, timeRange]);

    const completionRate = totalDraft > 0
        ? Math.round((totalWithImage / totalDraft) * 100)
        : 0;

    const productEntries = Object.entries(productCoverage)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 8);

    const stats = [
        {
            label: "Total Draft",
            value: totalDraft,
            icon: FileText,
            trend: "+12%",
            trendUp: true,
            iconColor: "#378ADD",
            dotColor: "#378ADD",
            description: "Semua draft tersimpan",
        },
        {
            label: "Terjadwal",
            value: totalScheduled,
            icon: Calendar,
            trend: "+5%",
            trendUp: true,
            iconColor: "#7F77DD",
            dotColor: "#7F77DD",
            description: "Posting otomatis",
        },
        {
            label: "Dengan Gambar",
            value: totalWithImage,
            icon: ImageIcon,
            trend: "+8%",
            trendUp: true,
            iconColor: "#1D9E75",
            dotColor: "#1D9E75",
            description: "Sudah ada ilustrasi",
        },
        {
            label: "Tanpa Gambar",
            value: totalWithoutImage,
            icon: ImageOff,
            trend: "-3%",
            trendUp: false,
            iconColor: "#BA7517",
            dotColor: "#BA7517",
            description: "Perlu ditambah gambar",
        },
    ];

    // ─── Chart: Daily Activity ───────────────────────────────────────
    const activityChartData = {
        labels: filteredActivity.map(a => a.date),
        datasets: [
            {
                label: "Draft",
                data: filteredActivity.map(a => a.count),
                backgroundColor: "#378ADD",
                borderRadius: 3,
                borderSkipped: false as const,
            },
        ],
    };

    const activityChartOptions: ChartOptions<"bar"> = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: { label: ctx => ` ${ctx.parsed.y} draft` },
            },
        },
        scales: {
            x: {
                ticks: {
                    color: textColor,
                    font: { size: 9 },
                    maxRotation: 0,
                    autoSkip: true,
                    maxTicksLimit: 8,
                },
                grid: { display: false },
                border: { display: false },
            },
            y: {
                ticks: { color: textColor, font: { size: 10 } },
                grid: { color: gridColor },
                border: { display: false },
            },
        },
    };

    // ─── Chart: Product Coverage ─────────────────────────────────────
    const productChartData = {
        labels: productEntries.map(([name]) => name),
        datasets: [
            {
                label: "Draft",
                data: productEntries.map(([, count]) => count),
                backgroundColor: "#7F77DD",
                borderRadius: 3,
                borderSkipped: false as const,
            },
        ],
    };

    const productChartOptions: ChartOptions<"bar"> = {
        indexAxis: "y" as const,
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: { label: ctx => ` ${ctx.parsed.x} draft` },
            },
        },
        scales: {
            x: {
                ticks: { color: textColor, font: { size: 10 } },
                grid: { color: gridColor },
                border: { display: false },
            },
            y: {
                ticks: { color: textColor, font: { size: 11 } },
                grid: { display: false },
                border: { display: false },
            },
        },
    };

    // ─── Chart: Weekly Trend ─────────────────────────────────────────
    const trendChartData = {
        labels: weeklyTrend.map(w => w.week),
        datasets: [
            {
                label: "Dengan gambar",
                data: weeklyTrend.map(w => w.withImage),
                borderColor: "#1D9E75",
                backgroundColor: "rgba(29,158,117,0.08)",
                borderWidth: 2,
                pointBackgroundColor: "#1D9E75",
                pointRadius: 3,
                tension: 0.4,
                fill: true,
            },
            {
                label: "Tanpa gambar",
                data: weeklyTrend.map(w => w.withoutImage),
                borderColor: "#BA7517",
                backgroundColor: "rgba(186,117,23,0.06)",
                borderWidth: 2,
                borderDash: [5, 4],
                pointBackgroundColor: "#BA7517",
                pointStyle: "rectRot" as const,
                pointRadius: 4,
                tension: 0.4,
                fill: true,
            },
        ],
    };

    const trendChartOptions: ChartOptions<"line"> = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: ctx => ` ${ctx.dataset.label}: ${ctx.parsed.y}`,
                },
            },
        },
        scales: {
            x: {
                ticks: { color: textColor, font: { size: 11 } },
                grid: { display: false },
                border: { display: false },
            },
            y: {
                ticks: { color: textColor, font: { size: 10 } },
                grid: { color: gridColor },
                border: { display: false },
            },
        },
    };

    // ─── Chart: Donut Completion ─────────────────────────────────────
    const donutData = {
        datasets: [
            {
                data: [completionRate, 100 - completionRate],
                backgroundColor: [
                    "#378ADD",
                    isDark ? "rgba(255,255,255,0.08)" : "rgba(0,0,0,0.07)",
                ],
                borderWidth: 0,
                borderRadius: 3,
            },
        ],
    };

    const donutOptions: ChartOptions<"doughnut"> = {
        cutout: "72%",
        plugins: {
            legend: { display: false },
            tooltip: { enabled: false },
        },
    };

    const totalActivity = filteredActivity.reduce((s, a) => s + a.count, 0);
    const avgActivity = filteredActivity.length
        ? Math.round(totalActivity / filteredActivity.length)
        : 0;
    const peakActivity = filteredActivity.length
        ? Math.max(...filteredActivity.map(a => a.count))
        : 0;

    return (
        <div className="space-y-4">
            {/* Header */}
            <div>
                <div className="space-y-1">
                    <h2 className="text-2xl font-bold tracking-tight">  Dashboard Draft</h2>
                    <p className="text-sm text-muted-foreground">
                        Ringkasan performa dan aktivitas konten
                    </p>
                </div>
            </div>

            {/* Stat Cards */}
            <div className="grid grid-cols-2 gap-3 lg:grid-cols-4">
                {stats.map(({ label, value, icon: Icon, trend, trendUp, iconColor, dotColor, description }) => (
                    <div
                        key={label}
                        className="flex flex-col gap-3 rounded-xl border p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]"
                    >
                        <div className="flex items-start justify-between">
                            <div className="space-y-0.5">
                                <span className="text-[11px] text-slate-400 dark:text-slate-500">{label}</span>
                                <div className={cn(
                                    "text-[10px] font-medium",
                                    trendUp ? "text-green-600 dark:text-green-400" : "text-red-500"
                                )}>
                                    {trendUp ? "↑" : "↓"} {trend}
                                </div>
                            </div>
                            <Icon className="h-4 w-4" style={{ color: iconColor }} />
                        </div>
                        <span className="text-2xl font-medium text-slate-900 dark:text-white tabular-nums">
                            {value.toLocaleString()}
                        </span>
                        <div className="flex items-center gap-1.5">
                            <span
                                className="h-1.5 w-1.5 rounded-full"
                                style={{ backgroundColor: dotColor }}
                            />
                            <span className="text-[11px] text-slate-400 dark:text-slate-500">
                                {description}
                            </span>
                        </div>
                    </div>
                ))}
            </div>

            {/* Completion Rate with Donut */}
            <div className="rounded-xl border p-3 flex items-center gap-4 bg-slate-50 border-blue-100 dark:bg-blue-500/5 dark:border-blue-500/20">
                <div className="relative w-14 h-14 flex-shrink-0">
                    <Doughnut data={donutData} options={donutOptions} />
                    <div className="absolute inset-0 flex items-center justify-center">
                        <span className="text-[11px] font-medium text-slate-700 dark:text-slate-200">
                            {completionRate}%
                        </span>
                    </div>
                </div>
                <div className="flex-1">
                    <div className="flex items-center gap-1.5 mb-1.5">
                        <CheckCircle className="h-3.5 w-3.5 text-blue-500" />
                        <span className="text-xs font-medium text-slate-700 dark:text-slate-300">
                            Kelengkapan Gambar
                        </span>
                    </div>
                    <div className="w-full h-1.5 bg-slate-200 dark:bg-white/10 rounded-full overflow-hidden">
                        <div
                            className="h-full bg-blue-500 rounded-full transition-all duration-500"
                            style={{ width: `${completionRate}%` }}
                        />
                    </div>
                    <p className="text-[11px] text-slate-400 mt-1">
                        {totalWithImage} dari {totalDraft} draft
                    </p>
                </div>
            </div>

            {/* Charts Row */}
            <div className="grid gap-3 md:grid-cols-2">
                {/* Daily Activity */}
                <div className="rounded-xl border p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                    <div className="flex items-center justify-between mb-1">
                        <div className="flex items-center gap-1.5">
                            <TrendingUp className="h-3.5 w-3.5 text-blue-500" />
                            <span className="text-sm font-medium text-slate-800 dark:text-slate-100">
                                Aktivitas {timeRange === "7d" ? "7 Hari" : timeRange === "30d" ? "30 Hari" : "90 Hari"}
                            </span>
                        </div>
                        <select
                            value={timeRange}
                            onChange={e => setTimeRange(e.target.value as TimeRange)}
                            className="text-[11px] text-slate-500 bg-transparent border border-slate-200 dark:border-white/10 rounded-md px-1.5 py-0.5 cursor-pointer"
                        >
                            <option value="7d">7 hari</option>
                            <option value="30d">30 hari</option>
                            <option value="90d">90 hari</option>
                        </select>
                    </div>
                    <p className="text-[11px] text-slate-400 mb-3">Draft dibuat per hari</p>

                    {/* Legend */}
                    <div className="flex gap-3 mb-2 text-[11px] text-slate-400">
                        <span className="flex items-center gap-1">
                            <span className="w-2.5 h-2.5 rounded-sm inline-block" style={{ background: "#378ADD" }} />
                            Jumlah draft
                        </span>
                    </div>

                    <div className="relative h-44">
                        <Bar data={activityChartData} options={activityChartOptions} />
                    </div>

                    <div className="mt-3 pt-2.5 border-t border-slate-100 dark:border-white/[0.05] flex justify-between text-[10px] text-slate-400">
                        <span>Total: <strong className="font-medium text-slate-600 dark:text-slate-300">{totalActivity}</strong></span>
                        <span>Rata-rata: <strong className="font-medium text-slate-600 dark:text-slate-300">{avgActivity}/hari</strong></span>
                        <span>Puncak: <strong className="font-medium text-slate-600 dark:text-slate-300">{peakActivity}</strong></span>
                    </div>
                </div>

                {/* Product Coverage */}
                <div className="rounded-xl border p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                    <div className="flex items-center gap-1.5 mb-1">
                        <BarChart2 className="h-3.5 w-3.5 text-purple-500" />
                        <span className="text-sm font-medium text-slate-800 dark:text-slate-100">
                            Coverage Produk
                        </span>
                    </div>
                    <p className="text-[11px] text-slate-400 mb-3">
                        {productEntries.length} produk aktif · {totalDraft} total draft
                    </p>

                    {/* Legend */}
                    <div className="flex gap-3 mb-2 text-[11px] text-slate-400">
                        <span className="flex items-center gap-1">
                            <span className="w-2.5 h-2.5 rounded-sm inline-block" style={{ background: "#7F77DD" }} />
                            Draft per produk
                        </span>
                    </div>

                    <div
                        className="relative"
                        style={{ height: `${productEntries.length * 36 + 32}px`, minHeight: "160px" }}
                    >
                        <Bar data={productChartData} options={productChartOptions} />
                    </div>
                </div>
            </div>

            {/* Weekly Trend Line Chart */}
            {weeklyTrend.length > 0 && (
                <div className="rounded-xl border p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                    <div className="flex items-center gap-1.5 mb-1">
                        <TrendingUp className="h-3.5 w-3.5 text-green-500" />
                        <span className="text-sm font-medium text-slate-800 dark:text-slate-100">
                            Tren Mingguan
                        </span>
                    </div>
                    <p className="text-[11px] text-slate-400 mb-3">
                        Perbandingan draft dengan gambar vs tanpa gambar
                    </p>

                    {/* Legend */}
                    <div className="flex gap-4 mb-2 text-[11px] text-slate-400">
                        <span className="flex items-center gap-1">
                            <span className="w-4 h-0.5 inline-block rounded" style={{ background: "#1D9E75" }} />
                            Dengan gambar
                        </span>
                        <span className="flex items-center gap-1">
                            <span className="w-4 h-0 inline-block border-t-2 border-dashed" style={{ borderColor: "#BA7517" }} />
                            Tanpa gambar
                        </span>
                    </div>

                    <div className="relative h-44">
                        <Line data={trendChartData} options={trendChartOptions} />
                    </div>
                </div>
            )}
        </div>
    );
}