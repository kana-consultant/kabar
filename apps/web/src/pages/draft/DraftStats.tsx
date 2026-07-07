import { useState, useMemo } from "react";
import {
    CheckCircle,
    BarChart2,
    TrendingUp,
    Calendar,
    Hash,
    Target,
    Zap,
    AlertCircle
} from "lucide-react";
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
    scheduled: number;
    published: number;
    with_image: number;
    with_keywords: number;
    avg_seo_score: number;
}

interface WeeklyTrend {
    week: string;
    created: number;
    scheduled: number;
    published: number;
}

interface ScheduledItem {
    id: string;
    title: string;
    scheduled_for: string;
    products?: string[];
}

interface KeywordStats {
    keyword: string;
    count: number;
}

interface DraftStatsProps {
    totalDraft: number;
    totalWithImage: number;
    totalWithoutImage: number;
    totalScheduled: number;
    totalPublished: number;
    totalWithKeywords: number;
    totalWithSEO: number;
    
    completionRate: number;
    scheduledRate: number;
    imageCoverageRate: number;
    seoScoreAvg: number;
    keywordsAvgCount: number;
    
    statusBreakdown: Record<string, number>;
    productCoverage: Record<string, number>;
    productStatus: Record<string, number>;
    topicBreakdown: Record<string, number>;
    seoScoreDistribution: Record<string, number>;
    
    dailyActivity: DailyActivity[];
    weeklyTrend?: WeeklyTrend[];
    scheduledUpcoming?: ScheduledItem[];
    
    topTopics?: Array<{ topic: string; count: number; avg_seo_score: number }>;
    topKeywords?: KeywordStats[];
    
    cacheMetadata?: { cached_at: string; ttl: string; generation_time_ms: number };
    isLoading?: boolean;
}

type TimeRange = "7d" | "30d" | "90d";

export function DraftStats({
    totalDraft,
    totalWithImage,
    totalScheduled,
    totalPublished,
    totalWithKeywords,
    totalWithSEO,
    completionRate,
    scheduledRate,
    seoScoreAvg,
    keywordsAvgCount,
    statusBreakdown,
    productCoverage,
    seoScoreDistribution,
    dailyActivity,
    weeklyTrend = [],
    scheduledUpcoming = [],
    topKeywords = [],
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

    const productEntries = Object.entries(productCoverage)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 8);

    const statsCards = [
        {
            label: "Total Draft",
            value: totalDraft,
            icon: <Hash className="h-4 w-4 text-blue-500" />,
            subtext: `${totalPublished} published`
        },
        {
            label: "Completion",
            value: `${completionRate.toFixed(0)}%`,
            icon: <CheckCircle className="h-4 w-4 text-green-500" />,
            subtext: `${totalWithImage} with image`
        },
        {
            label: "Scheduled",
            value: totalScheduled,
            icon: <Calendar className="h-4 w-4 text-purple-500" />,
            subtext: `${scheduledRate.toFixed(0)}% of total`
        },
        {
            label: "Avg SEO",
            value: seoScoreAvg.toFixed(1),
            icon: <Target className="h-4 w-4 text-orange-500" />,
            subtext: `${totalWithSEO} optimized`
        },
        {
            label: "Keywords",
            value: keywordsAvgCount.toFixed(1),
            icon: <Zap className="h-4 w-4 text-yellow-500" />,
            subtext: `${totalWithKeywords} with kw`
        }
    ];

    const activityChartData = {
        labels: filteredActivity.map(a => a.date),
        datasets: [
            {
                label: "Total",
                data: filteredActivity.map(a => a.count),
                backgroundColor: "#378ADD",
                borderRadius: 3,
                borderSkipped: false as const,
            },
            {
                label: "Scheduled",
                data: filteredActivity.map(a => a.scheduled),
                backgroundColor: "#7F77DD",
                borderRadius: 3,
                borderSkipped: false as const,
            },
            {
                label: "Published",
                data: filteredActivity.map(a => a.published),
                backgroundColor: "#1D9E75",
                borderRadius: 3,
                borderSkipped: false as const,
            },
        ],
    };

    const activityChartOptions: ChartOptions<"bar"> = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: { 
                display: true,
                position: "bottom" as const,
                labels: {
                    color: textColor,
                    font: { size: 10 },
                    padding: 10,
                    usePointStyle: true,
                }
            },
        },
        scales: {
            x: {
                stacked: true,
                ticks: { color: textColor, font: { size: 9 }, maxRotation: 0, autoSkip: true, maxTicksLimit: 8 },
                grid: { display: false },
                border: { display: false },
            },
            y: {
                stacked: true,
                ticks: { color: textColor, font: { size: 10 } },
                grid: { color: gridColor },
                border: { display: false },
            },
        },
    };

    const seoEntries = Object.entries(seoScoreDistribution).sort((a, b) => {
        const order = ["0", "1-20", "21-40", "41-60", "61-80", "81-100"];
        return order.indexOf(a[0]) - order.indexOf(b[0]);
    });
    
    const seoChartData = {
        labels: seoEntries.map(([range]) => range),
        datasets: [
            {
                label: "Draft",
                data: seoEntries.map(([, count]) => count),
                backgroundColor: ["#EF4444", "#F97316", "#F59E0B", "#84CC16", "#22C55E", "#10B981"],
                borderRadius: 3,
                borderSkipped: false as const,
            },
        ],
    };

    const seoChartOptions: ChartOptions<"bar"> = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: {
            x: { ticks: { color: textColor, font: { size: 10 } }, grid: { display: false }, border: { display: false } },
            y: { ticks: { color: textColor, font: { size: 10 } }, grid: { color: gridColor }, border: { display: false } },
        },
    };

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
        plugins: { legend: { display: false } },
        scales: {
            x: { ticks: { color: textColor, font: { size: 10 } }, grid: { color: gridColor }, border: { display: false } },
            y: { ticks: { color: textColor, font: { size: 11 } }, grid: { display: false }, border: { display: false } },
        },
    };

    const trendChartData = {
        labels: weeklyTrend.map(w => w.week),
        datasets: [
            {
                label: "Created",
                data: weeklyTrend.map(w => w.created),
                borderColor: "#378ADD",
                backgroundColor: "rgba(55,138,221,0.08)",
                borderWidth: 2,
                pointRadius: 3,
                tension: 0.4,
                fill: true,
            },
            {
                label: "Scheduled",
                data: weeklyTrend.map(w => w.scheduled),
                borderColor: "#7F77DD",
                backgroundColor: "rgba(127,119,221,0.08)",
                borderWidth: 2,
                pointRadius: 3,
                tension: 0.4,
                fill: true,
            },
            {
                label: "Published",
                data: weeklyTrend.map(w => w.published),
                borderColor: "#1D9E75",
                backgroundColor: "rgba(29,158,117,0.08)",
                borderWidth: 2,
                pointRadius: 3,
                tension: 0.4,
                fill: true,
            },
        ],
    };

    const trendChartOptions: ChartOptions<"line"> = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: { 
                display: true,
                position: "bottom" as const,
                labels: { color: textColor, font: { size: 10 }, padding: 10, usePointStyle: true }
            },
        },
        scales: {
            x: { ticks: { color: textColor, font: { size: 11 } }, grid: { display: false }, border: { display: false } },
            y: { ticks: { color: textColor, font: { size: 10 } }, grid: { color: gridColor }, border: { display: false } },
        },
    };

    const donutData = {
        datasets: [{
            data: [completionRate, 100 - completionRate],
            backgroundColor: ["#378ADD", isDark ? "rgba(255,255,255,0.08)" : "rgba(0,0,0,0.07)"],
            borderWidth: 0,
            borderRadius: 3,
        }],
    };

    const donutOptions: ChartOptions<"doughnut"> = {
        cutout: "72%",
        plugins: { legend: { display: false }, tooltip: { enabled: false } },
    };

    return (
        <div className="space-y-4 mt-2">
            {/* Summary Cards */}
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3">
                {statsCards.map((stat, index) => (
                    <div key={index} className="rounded-xl border p-3 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                        <div className="flex items-center justify-between mb-2">
                            <span className="text-[11px] text-slate-400 font-medium">{stat.label}</span>
                            {stat.icon}
                        </div>
                        <div className="text-xl font-bold text-slate-800 dark:text-slate-100">{stat.value}</div>
                        <p className="text-[10px] text-slate-400 mt-1">{stat.subtext}</p>
                    </div>
                ))}
            </div>

            {/* Completion Donut */}
            <div className="rounded-xl border p-4 flex items-center gap-4 bg-gradient-to-r from-blue-50 to-purple-50 border-blue-100 dark:from-blue-500/5 dark:to-purple-500/5 dark:border-blue-500/20">
                <div className="relative w-20 h-20 flex-shrink-0">
                    <Doughnut data={donutData} options={donutOptions} />
                    <div className="absolute inset-0 flex items-center justify-center">
                        <span className="text-sm font-bold text-slate-700 dark:text-slate-200">{completionRate.toFixed(0)}%</span>
                    </div>
                </div>
                <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                        <CheckCircle className="h-4 w-4 text-blue-500" />
                        <span className="text-sm font-semibold text-slate-700 dark:text-slate-300">Kelengkapan Konten</span>
                    </div>
                    <div className="space-y-2">
                        <div>
                            <div className="flex justify-between text-xs mb-1">
                                <span className="text-slate-500">Gambar</span>
                                <span className="text-slate-600 font-medium">{totalWithImage}/{totalDraft}</span>
                            </div>
                            <div className="w-full h-1.5 bg-slate-200 dark:bg-white/10 rounded-full overflow-hidden">
                                <div className="h-full bg-blue-500 rounded-full transition-all" style={{ width: `${completionRate}%` }} />
                            </div>
                        </div>
                        <div>
                            <div className="flex justify-between text-xs mb-1">
                                <span className="text-slate-500">Keywords</span>
                                <span className="text-slate-600 font-medium">{totalWithKeywords}/{totalDraft}</span>
                            </div>
                            <div className="w-full h-1.5 bg-slate-200 dark:bg-white/10 rounded-full overflow-hidden">
                                <div className="h-full bg-green-500 rounded-full transition-all" style={{ width: `${totalDraft > 0 ? (totalWithKeywords / totalDraft) * 100 : 0}%` }} />
                            </div>
                        </div>
                        <div>
                            <div className="flex justify-between text-xs mb-1">
                                <span className="text-slate-500">SEO</span>
                                <span className="text-slate-600 font-medium">{totalWithSEO}/{totalDraft}</span>
                            </div>
                            <div className="w-full h-1.5 bg-slate-200 dark:bg-white/10 rounded-full overflow-hidden">
                                <div className="h-full bg-orange-500 rounded-full transition-all" style={{ width: `${totalDraft > 0 ? (totalWithSEO / totalDraft) * 100 : 0}%` }} />
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Status Breakdown */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                {Object.entries(statusBreakdown).map(([status, count]) => (
                    <div key={status} className="rounded-xl border p-3 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                        <div className="text-[11px] text-slate-400 capitalize mb-1">{status}</div>
                        <div className="text-lg font-bold text-slate-800 dark:text-slate-100">{count}</div>
                        <div className="text-[10px] text-slate-400 mt-1">
                            {totalDraft > 0 ? `${((count / totalDraft) * 100).toFixed(1)}%` : '0%'}
                        </div>
                    </div>
                ))}
            </div>

            {/* Charts Row 1 */}
            <div className="grid gap-3 md:grid-cols-2">
                <div className="rounded-xl border p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                    <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center gap-2">
                            <TrendingUp className="h-4 w-4 text-blue-500" />
                            <span className="text-sm font-semibold text-slate-800 dark:text-slate-100">Aktivitas Harian</span>
                        </div>
                        <select value={timeRange} onChange={e => setTimeRange(e.target.value as TimeRange)} className="text-xs text-slate-500 bg-transparent border border-slate-200 dark:border-white/10 rounded-md px-2 py-1 cursor-pointer">
                            <option value="7d">7 hari</option>
                            <option value="30d">30 hari</option>
                            <option value="90d">90 hari</option>
                        </select>
                    </div>
                    <div className="relative h-64">
                        <Bar data={activityChartData} options={activityChartOptions} />
                    </div>
                </div>

                <div className="rounded-xl border p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                    <div className="flex items-center gap-2 mb-3">
                        <Target className="h-4 w-4 text-orange-500" />
                        <span className="text-sm font-semibold text-slate-800 dark:text-slate-100">Distribusi SEO Score</span>
                    </div>
                    <p className="text-xs text-slate-400 mb-2">Avg: {seoScoreAvg.toFixed(1)} · {totalWithSEO} optimized</p>
                    <div className="relative h-64">
                        <Bar data={seoChartData} options={seoChartOptions} />
                    </div>
                </div>
            </div>

            {/* Charts Row 2 */}
            <div className="grid gap-3 md:grid-cols-2">
                <div className="rounded-xl border p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                    <div className="flex items-center gap-2 mb-3">
                        <BarChart2 className="h-4 w-4 text-purple-500" />
                        <span className="text-sm font-semibold text-slate-800 dark:text-slate-100">Coverage Produk</span>
                    </div>
                    <div className="relative" style={{ height: `${productEntries.length * 36 + 32}px`, minHeight: "200px" }}>
                        <Bar data={productChartData} options={productChartOptions} />
                    </div>
                </div>

                <div className="rounded-xl border p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                    <div className="flex items-center gap-2 mb-3">
                        <Hash className="h-4 w-4 text-yellow-500" />
                        <span className="text-sm font-semibold text-slate-800 dark:text-slate-100">Top Keywords</span>
                    </div>
                    <p className="text-xs text-slate-400 mb-3">Avg {keywordsAvgCount.toFixed(1)} keywords/draft</p>
                    <div className="space-y-2 max-h-[300px] overflow-y-auto">
                        {topKeywords.slice(0, 15).map((kw, index) => (
                            <div key={kw.keyword} className="flex items-center justify-between p-2 rounded-lg bg-slate-50 dark:bg-white/5">
                                <div className="flex items-center gap-2">
                                    <span className="text-xs font-medium text-slate-400 w-5">#{index + 1}</span>
                                    <span className="text-sm text-slate-700 dark:text-slate-300">{kw.keyword}</span>
                                </div>
                                <span className="text-xs font-medium text-slate-500 bg-slate-200 dark:bg-white/10 px-2 py-0.5 rounded-full">{kw.count}x</span>
                            </div>
                        ))}
                    </div>
                </div>
            </div>

            {/* Weekly Trend */}
            {weeklyTrend.length > 0 && (
                <div className="rounded-xl border p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                    <div className="flex items-center gap-2 mb-3">
                        <TrendingUp className="h-4 w-4 text-green-500" />
                        <span className="text-sm font-semibold text-slate-800 dark:text-slate-100">Tren Mingguan</span>
                    </div>
                    <div className="relative h-64">
                        <Line data={trendChartData} options={trendChartOptions} />
                    </div>
                </div>
            )}

            {/* Upcoming Scheduled */}
            {scheduledUpcoming.length > 0 && (
                <div className="rounded-xl border p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                    <div className="flex items-center gap-2 mb-3">
                        <Calendar className="h-4 w-4 text-purple-500" />
                        <span className="text-sm font-semibold text-slate-800 dark:text-slate-100">Jadwal Mendatang</span>
                    </div>
                    <div className="space-y-2">
                        {scheduledUpcoming.slice(0, 5).map((item) => (
                            <div key={item.id} className="flex items-center justify-between p-3 rounded-lg bg-slate-50 dark:bg-white/5">
                                <div>
                                    <h4 className="text-sm font-medium text-slate-700 dark:text-slate-300">{item.title}</h4>
                                    <div className="flex items-center gap-2 mt-1">
                                        <Calendar className="h-3 w-3 text-slate-400" />
                                        <span className="text-xs text-slate-400">
                                            {new Date(item.scheduled_for).toLocaleDateString('id-ID', {
                                                weekday: 'long', year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit'
                                            })}
                                        </span>
                                    </div>
                                    {item.products && item.products.length > 0 && (
                                        <div className="flex gap-1 mt-2">
                                            {item.products.slice(0, 3).map((product, idx) => (
                                                <span key={idx} className="text-[10px] bg-purple-100 dark:bg-purple-500/20 text-purple-700 dark:text-purple-300 px-2 py-0.5 rounded-full">
                                                    {product}
                                                </span>
                                            ))}
                                            {item.products.length > 3 && (
                                                <span className="text-[10px] text-slate-400">+{item.products.length - 3}</span>
                                            )}
                                        </div>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}