import { createFileRoute } from "@tanstack/react-router";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@kana-consultant/ui-kit";
import { StatsCards } from "@/pages/dashboard/StatsCards";
import { RecentActivity } from "@/pages/dashboard/RecentActivity";
import { useDrafts } from "@/hooks/useDrafts";
import { DraftStats } from "@/pages/draft/DraftStats";

export const Route = createFileRoute("/")({
    component: Dashboard,
});

export function Dashboard() {
    const { stats } = useDrafts();

    return (
        <div className="space-y-6">
            <div>
                <h2 className="text-2xl font-bold tracking-tight">Dashboard</h2>
                <p className="text-slate-500">
                    1 dashboard, N produk, 1 klik. Selesai.
                </p>
            </div>
            <Tabs defaultValue="overview" className="space-y-1">
                <TabsList>
                    <TabsTrigger value="overview">Overview</TabsTrigger>
                    <TabsTrigger value="draft-stats">Draft Stats</TabsTrigger>
                </TabsList>

                <TabsContent value="overview">
                    <StatsCards />
                    <div className="grid gap-6 mt-5">
                        {/* <QuickGenerate /> */}
                        <RecentActivity />
                    </div>
                </TabsContent>

                <TabsContent value="draft-stats">
                    <DraftStats
                        totalDraft={stats?.total_draft ?? 0}
                        totalWithImage={stats?.total_with_image ?? 0}
                        totalWithoutImage={stats?.total_without_image ?? 0}
                        totalScheduled={stats?.total_scheduled ?? 0}
                        productCoverage={stats?.product_coverage ?? {}}
                        dailyActivity={stats?.daily_activity ?? []}
                    />
                </TabsContent>
            </Tabs>
        </div>
    );
}