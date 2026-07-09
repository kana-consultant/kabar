import { useState, useMemo } from "react";
import {
    CheckCircle,
    BarChart2,
    TrendingUp,
    Calendar,
    Hash,
    Target,
    Zap,
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
    type ChartOptions,
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
    cacheMetadata,
}: DraftStatsProps) {
    const [timeRange, setTimeRange] = useState<TimeRange>("30d");

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
            icon: <Hash className="h-4 w-4 text-blue-500 shrink-0" />,
            subtext: `${totalPublished} published`
        },
        {
            label: "Completion",
            value: `${completionRate.toFixed(0)}%`,
            icon: <CheckCircle className="h-4 w-4 text-green-500 shrink-0" />,
            subtext: `${totalWithImage} with image`
        },
        {
            label: "Scheduled",
            value: totalScheduled,
            icon: <Calendar className="h-4 w-4 text-purple-500 shrink-0" />,
            subtext: `${scheduledRate.toFixed(0)}% of total`
        },
        {
            label: "Avg SEO",
            value: seoScoreAvg.toFixed(1),
            icon: <Target className="h-4 w-4 text-orange-500 shrink-0" />,
            subtext: `${totalWithSEO} optimized`
        },
        {
            label: "Keywords",
            value: keywordsAvgCount.toFixed(1),
            icon: <Zap className="h-4 w-4 text-yellow-500 shrink-0" />,
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
                borderSkipped: false as any,
            },
            {
                label: "Scheduled",
                data: filteredActivity.map(a => a.scheduled),
                backgroundColor: "#7F77DD",
                borderRadius: 3,
                borderSkipped: false as any,
            },
            {
                label: "Published",
                data: filteredActivity.map(a => a.published),
                backgroundColor: "#1D9E75",
                borderRadius: 3,
                borderSkipped: false as any,
            },
        ],
    };

    const activityChartOptions: ChartOptions<"bar"> = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                display: true,
                position: "bottom" as any,
                labels: {
                    font: { size: 10, weight: "500" as any },
                    padding: 10,
                    usePointStyle: true,
                    pointStyleWidth: 8,
                }
            },
        },
        scales: {
            x: {
                stacked: true,
                ticks: { font: { size: 9 }, maxRotation: 0, autoSkip: true, maxTicksLimit: 6 },
                grid: { display: false },
                border: { display: false },
            },
            y: {
                stacked: true,
                ticks: { font: { size: 10 } },
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
                borderSkipped: false as any,
            },
        ],
    };

    const seoChartOptions: ChartOptions<"bar"> = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: {
            x: { ticks: { font: { size: 9 } }, grid: { display: false }, border: { display: false } },
            y: { ticks: { font: { size: 10 } }, border: { display: false } },
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
                borderSkipped: false as any,
            },
        ],
    };

    const productChartOptions: ChartOptions<"bar"> = {
        indexAxis: "y" as any,
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: {
            x: { ticks: { font: { size: 10 } }, border: { display: false } },
            y: { ticks: { font: { size: 10 } }, grid: { display: false }, border: { display: false } },
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
                position: "bottom" as any,
                labels: {
                    font: { size: 10, weight: "500" as any },
                    padding: 10,
                    usePointStyle: true,
                    pointStyleWidth: 8,
                }
            },
        },
        scales: {
            x: { ticks: { font: { size: 10 }, maxRotation: 0, autoSkip: true, maxTicksLimit: 6 }, grid: { display: false }, border: { display: false } },
            y: { ticks: { font: { size: 10 } }, border: { display: false } },
        },
    };

    const donutData = {
        datasets: [{
            data: [completionRate, 100 - completionRate],
            backgroundColor: ["#378ADD", "rgba(0,0,0,0.06)"],
            borderWidth: 0,
            borderRadius: 3,
        }],
    };

    const donutOptions: ChartOptions<"doughnut"> = {
        cutout: "72%",
        plugins: { legend: { display: false }, tooltip: { enabled: false } },
    };

    return (
        <div className="w-full max-w-full space-y-3 sm:space-y-4 mt-2 overflow-x-hidden">
            {/* Summary Cards */}
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-2 sm:gap-3">
                {statsCards.map((stat, index) => (
                    <div key={index} className="rounded-xl border p-2.5 sm:p-3 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06] min-w-0">
                        <div className="flex items-center justify-between mb-2 gap-1">
                            <span className="text-[10px] sm:text-[11px] text-slate-500 dark:text-slate-400 font-medium truncate">{stat.label}</span>
                            {stat.icon}
                        </div>
                        <div className="text-lg sm:text-xl font-bold text-slate-800 dark:text-slate-100 truncate">{stat.value}</div>
                        <p className="text-[9px] sm:text-[10px] text-slate-500 dark:text-slate-400 mt-1 truncate">{stat.subtext}</p>
                    </div>
                ))}
            </div>

            {/* Completion Donut */}
            <div className="rounded-xl border p-3 sm:p-4 flex flex-col sm:flex-row items-center gap-3 sm:gap-4 bg-gradient-to-r from-blue-50 to-purple-50 border-blue-100 dark:from-blue-500/5 dark:to-purple-500/5 dark:border-blue-500/20">
                <div className="relative w-20 h-20 flex-shrink-0">
                    <Doughnut data={donutData} options={donutOptions} />
                    <div className="absolute inset-0 flex items-center justify-center">
                        <span className="text-sm font-bold text-slate-700 dark:text-slate-200">{completionRate.toFixed(0)}%</span>
                    </div>
                </div>
                <div className="flex-1 w-full min-w-0">
                    <div className="flex items-center gap-2 mb-2 justify-center sm:justify-start">
                        <CheckCircle className="h-4 w-4 text-blue-500 shrink-0" />
                        <span className="text-sm font-semibold text-slate-700 dark:text-slate-300">Kelengkapan Konten</span>
                    </div>
                    <div className="space-y-2">
                        <div>
                            <div className="flex justify-between text-xs mb-1 gap-2">
                                <span className="text-slate-500 dark:text-slate-400">Gambar</span>
                                <span className="text-slate-600 dark:text-slate-300 font-medium">{totalWithImage}/{totalDraft}</span>
                            </div>
                            <div className="w-full h-1.5 bg-slate-200 dark:bg-white/10 rounded-full overflow-hidden">
                                <div className="h-full bg-blue-500 rounded-full transition-all" style={{ width: `${completionRate}%` }} />
                            </div>
                        </div>
                        <div>
                            <div className="flex justify-between text-xs mb-1 gap-2">
                                <span className="text-slate-500 dark:text-slate-400">Keywords</span>
                                <span className="text-slate-600 dark:text-slate-300 font-medium">{totalWithKeywords}/{totalDraft}</span>
                            </div>
                            <div className="w-full h-1.5 bg-slate-200 dark:bg-white/10 rounded-full overflow-hidden">
                                <div className="h-full bg-green-500 rounded-full transition-all" style={{ width: `${totalDraft > 0 ? (totalWithKeywords / totalDraft) * 100 : 0}%` }} />
                            </div>
                        </div>
                        <div>
                            <div className="flex justify-between text-xs mb-1 gap-2">
                                <span className="text-slate-500 dark:text-slate-400">SEO</span>
                                <span className="text-slate-600 dark:text-slate-300 font-medium">{totalWithSEO}/{totalDraft}</span>
                            </div>
                            <div className="w-full h-1.5 bg-slate-200 dark:bg-white/10 rounded-full overflow-hidden">
                                <div className="h-full bg-orange-500 rounded-full transition-all" style={{ width: `${totalDraft > 0 ? (totalWithSEO / totalDraft) * 100 : 0}%` }} />
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Status Breakdown */}
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 sm:gap-3">
                {Object.entries(statusBreakdown).map(([status, count]) => (
                    <div key={status} className="rounded-xl border p-2.5 sm:p-3 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06] min-w-0">
                        <div className="text-[10px] sm:text-[11px] text-slate-500 dark:text-slate-400 capitalize mb-1 truncate">{status}</div>
                        <div className="text-base sm:text-lg font-bold text-slate-800 dark:text-slate-100">{count}</div>
                        <div className="text-[9px] sm:text-[10px] text-slate-500 dark:text-slate-400 mt-1">
                            {totalDraft > 0 ? `${((count / totalDraft) * 100).toFixed(1)}%` : '0%'}
                        </div>
                    </div>
                ))}
            </div>

            {/* Charts Row 1 */}
            <div className="grid gap-3 md:grid-cols-2">
                <div className="rounded-xl border p-3 sm:p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06] min-w-0">
                    <div className="flex items-center justify-between mb-3 gap-2 flex-wrap">
                        <div className="flex items-center gap-2 min-w-0">
                            <TrendingUp className="h-4 w-4 text-blue-500 shrink-0" />
                            <span className="text-sm font-semibold text-slate-700 dark:text-slate-100 truncate">Aktivitas Harian</span>
                        </div>
                        <select value={timeRange} onChange={e => setTimeRange(e.target.value as TimeRange)} className="text-xs text-slate-500 dark:text-slate-400 bg-transparent border border-slate-200 dark:border-white/10 rounded-md px-2 py-1 cursor-pointer shrink-0">
                            <option value="7d">7 hari</option>
                            <option value="30d">30 hari</option>
                            <option value="90d">90 hari</option>
                        </select>
                    </div>
                    <div className="relative h-56 sm:h-64 w-full">
                        <Bar data={activityChartData} options={activityChartOptions} />
                    </div>
                </div>

                <div className="rounded-xl border p-3 sm:p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06] min-w-0">
                    <div className="flex items-center gap-2 mb-3">
                        <Target className="h-4 w-4 text-orange-500 shrink-0" />
                        <span className="text-sm font-semibold text-slate-700 dark:text-slate-100 truncate">Distribusi SEO Score</span>
                    </div>
                    <p className="text-xs text-slate-500 dark:text-slate-400 mb-2">Avg: {seoScoreAvg.toFixed(1)} · {totalWithSEO} optimized</p>
                    <div className="relative h-56 sm:h-64 w-full">
                        <Bar data={seoChartData} options={seoChartOptions} />
                    </div>
                </div>
            </div>

            {/* Charts Row 2 */}
            <div className="grid gap-3 md:grid-cols-2">
                <div className="rounded-xl border p-3 sm:p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06] min-w-0">
                    <div className="flex items-center gap-2 mb-3">
                        <BarChart2 className="h-4 w-4 text-purple-500 shrink-0" />
                        <span className="text-sm font-semibold text-slate-700 dark:text-slate-100 truncate">Coverage Produk</span>
                    </div>
                    <div className="relative w-full" style={{ height: `${productEntries.length * 36 + 32}px`, minHeight: "200px" }}>
                        <Bar data={productChartData} options={productChartOptions} />
                    </div>
                </div>

                <div className="rounded-xl border p-3 sm:p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06] min-w-0">
                    <div className="flex items-center gap-2 mb-3">
                        <Hash className="h-4 w-4 text-yellow-500 shrink-0" />
                        <span className="text-sm font-semibold text-slate-700 dark:text-slate-100 truncate">Top Keywords</span>
                    </div>
                    <p className="text-xs text-slate-500 dark:text-slate-400 mb-3">Avg {keywordsAvgCount.toFixed(1)} keywords/draft</p>
                    <div className="space-y-2 max-h-[300px] overflow-y-auto">
                        {topKeywords.slice(0, 15).map((kw, index) => (
                            <div key={kw.keyword} className="flex items-center justify-between p-2 rounded-lg bg-slate-50 dark:bg-white/5 gap-2">
                                <div className="flex items-center gap-2 min-w-0">
                                    <span className="text-xs font-medium text-slate-400 dark:text-slate-500 w-5 shrink-0">#{index + 1}</span>
                                    <span className="text-sm text-slate-700 dark:text-slate-300 truncate">{kw.keyword}</span>
                                </div>
                                <span className="text-xs font-medium text-slate-500 dark:text-slate-400 bg-slate-200 dark:bg-white/10 px-2 py-0.5 rounded-full shrink-0">{kw.count}x</span>
                            </div>
                        ))}
                    </div>
                </div>
            </div>

            {/* Weekly Trend */}
            {weeklyTrend.length > 0 && (
                <div className="rounded-xl border p-3 sm:p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06] min-w-0">
                    <div className="flex items-center gap-2 mb-3">
                        <TrendingUp className="h-4 w-4 text-green-500 shrink-0" />
                        <span className="text-sm font-semibold text-slate-700 dark:text-slate-100 truncate">Tren Mingguan</span>
                    </div>
                    <div className="relative h-56 sm:h-64 w-full">
                        <Line data={trendChartData} options={trendChartOptions} />
                    </div>
                </div>
            )}

            {/* Upcoming Scheduled */}
            {scheduledUpcoming.length > 0 && (
                <div className="rounded-xl border p-3 sm:p-4 bg-white border-slate-200/80 dark:bg-[#0f0d1a] dark:border-white/[0.06] min-w-0">
                    <div className="flex items-center gap-2 mb-3">
                        <Calendar className="h-4 w-4 text-purple-500 shrink-0" />
                        <span className="text-sm font-semibold text-slate-700 dark:text-slate-100 truncate">Jadwal Mendatang</span>
                    </div>
                    <div className="space-y-2">
                        {scheduledUpcoming.slice(0, 5).map((item) => (
                            <div key={item.id} className="flex items-start justify-between p-3 rounded-lg bg-slate-50 dark:bg-white/5 gap-2">
                                <div className="min-w-0 flex-1">
                                    <h4 className="text-sm font-medium text-slate-700 dark:text-slate-300 break-words">{item.title}</h4>
                                    <div className="flex items-start gap-2 mt-1">
                                        <Calendar className="h-3 w-3 text-slate-400 shrink-0 mt-0.5" />
                                        <span className="text-xs text-slate-500 dark:text-slate-400 break-words">
                                            {new Date(item.scheduled_for).toLocaleDateString('id-ID', {
                                                weekday: 'long', year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit'
                                            })}
                                        </span>
                                    </div>
                                    {item.products && item.products.length > 0 && (
                                        <div className="flex flex-wrap gap-1 mt-2">
                                            {item.products.slice(0, 3).map((product, idx) => (
                                                <span key={idx} className="text-[10px] bg-purple-100 dark:bg-purple-500/20 text-purple-700 dark:text-purple-300 px-2 py-0.5 rounded-full">
                                                    {product}
                                                </span>
                                            ))}
                                            {item.products.length > 3 && (
                                                <span className="text-[10px] text-slate-500 dark:text-slate-400">+{item.products.length - 3}</span>
                                            )}
                                        </div>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {/* Cache Metadata */}
            {cacheMetadata && (
                <div className="text-[10px] text-slate-400 dark:text-slate-500 text-right break-words">
                    Cached: {new Date(cacheMetadata.cached_at).toLocaleTimeString()} · TTL: {cacheMetadata.ttl} · {cacheMetadata.generation_time_ms.toFixed(0)}ms
                </div>
            )}
        </div>
    );
}