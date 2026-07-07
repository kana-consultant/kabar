import { createFileRoute } from "@tanstack/react-router";
import { StatsCards } from "@/pages/dashboard/StatsCards";
import { RecentActivity } from "@/pages/dashboard/RecentActivity";
import { useDrafts } from "@/hooks/useDrafts";
import { DraftStats } from "@/pages/draft/DraftStats";

export const Route = createFileRoute("/")({
    component: Dashboard,
});

export function Dashboard() {
    const { stats, loading } = useDrafts();

    return (
        <div className="space-y-6">
            {/* Header */}
            <div>
                <h2 className="text-2xl font-bold tracking-tight">Dashboard</h2>
                <p className="text-slate-500">
                    1 dashboard, N produk, 1 klik. Selesai.
                </p>
            </div>

            {/* Stats Cards dengan loading prop */}
            <StatsCards isLoading={loading} />

            {/* Draft Stats */}
            <DraftStats
                // Basic metrics
                totalDraft={stats?.total_draft ?? 0}
                totalWithImage={stats?.total_with_image ?? 0}
                totalWithoutImage={stats?.total_without_image ?? 0}
                totalScheduled={stats?.total_scheduled ?? 0}
                totalPublished={stats?.total_published ?? 0}
                totalWithKeywords={stats?.total_with_keywords ?? 0}
                totalWithSEO={stats?.total_with_seo ?? 0}

                // Derived metrics
                completionRate={stats?.completion_rate ?? 0}
                scheduledRate={stats?.scheduled_rate ?? 0}
                imageCoverageRate={stats?.image_coverage_rate ?? 0}
                seoScoreAvg={stats?.seo_score_avg ?? 0}
                keywordsAvgCount={stats?.keywords_avg_count ?? 0}

                // Breakdowns
                statusBreakdown={stats?.status_breakdown ?? {}}
                productCoverage={stats?.product_coverage ?? {}}
                productStatus={stats?.product_status ?? {}}
                topicBreakdown={stats?.topic_breakdown ?? {}}
                seoScoreDistribution={stats?.seo_score_distribution ?? {}}

                // Time series
                dailyActivity={stats?.daily_activity ?? []}
                weeklyTrend={stats?.weekly_trend ?? []}
                scheduledUpcoming={stats?.scheduled_upcoming ?? []}

                // Content quality
                topTopics={stats?.top_topics ?? []}
                topKeywords={stats?.top_keywords ?? []}

               
                cacheMetadata={stats?.cache_metadata}

                // UI state
                isLoading={loading}
            />

            {/* Recent Activity */}
            <RecentActivity />
        </div>
    );
}