import { createFileRoute } from "@tanstack/react-router";
import { StatsCards } from "../components/dashboard/StatsCards";
import { QuickGenerate } from "../components/dashboard/QuickGenerate";
import { RecentActivity } from "../components/dashboard/RecentActivity";
export const Route = createFileRoute("/")({
    component: Dashboard,
});

export function Dashboard() {
    return (
        <div className="space-y-6 ">
            <div>
                <h2 className="text-2xl font-bold tracking-tight">Dashboard</h2>
                <p className="text-slate-500">
                    1 dashboard, N produk, 1 klik. Selesai.
                </p>
            </div>

            <StatsCards />

            <div className="grid gap-6 md:grid-cols-2">
                <QuickGenerate />
                <RecentActivity />
            </div>
        </div>
    );
}