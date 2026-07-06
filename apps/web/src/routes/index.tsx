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

            <div>
                {/* Stats Cards dengan loading prop */}
                <StatsCards isLoading={loading} />
                
                {/* Draft Stats */}
                <DraftStats
                    totalDraft={stats?.total_draft ?? 0}
                    totalWithImage={stats?.total_with_image ?? 0}
                    totalWithoutImage={stats?.total_without_image ?? 0}
                    totalScheduled={stats?.total_scheduled ?? 0}
                    productCoverage={stats?.product_coverage ?? {}}
                    dailyActivity={stats?.daily_activity ?? []}
                    isLoading={loading}
                />
                
                {/* Recent Activity */}
                <div className="grid gap-6 mt-5">
                    <RecentActivity />
                </div>
            </div>
        </div>
    );
}