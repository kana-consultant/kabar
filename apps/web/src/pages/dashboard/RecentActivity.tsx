import { CheckCircle2, XCircle, Clock } from "lucide-react";
import { useHistory } from "@/hooks/useHistory";
import { cn } from "@/lib/utils";

const statusConfig = {
    published: {
        icon: CheckCircle2,
        badge: "bg-green-50 text-green-700 border-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:border-green-500/20",
        iconWrap: "bg-green-50 text-green-600 dark:bg-green-500/10 dark:text-green-400",
        label: "Terbit",
    },
    success: {
        icon: CheckCircle2,
        badge: "bg-green-50 text-green-700 border-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:border-green-500/20",
        iconWrap: "bg-green-50 text-green-600 dark:bg-green-500/10 dark:text-green-400",
        label: "Sukses",
    },
    failed: {
        icon: XCircle,
        badge: "bg-red-50 text-red-700 border-red-200/60 dark:bg-red-500/10 dark:text-red-400 dark:border-red-500/20",
        iconWrap: "bg-red-50 text-red-600 dark:bg-red-500/10 dark:text-red-400",
        label: "Gagal",
    },
    pending: {
        icon: Clock,
        badge: "bg-amber-50 text-amber-700 border-amber-200/60 dark:bg-amber-500/10 dark:text-amber-400 dark:border-amber-500/20",
        iconWrap: "bg-amber-50 text-amber-600 dark:bg-amber-500/10 dark:text-amber-400",
        label: "Pending",
    },
} as const;

export function RecentActivity() {
    const { history } = useHistory();
    const recent = history.slice(0, 5);

    return (
        <div className={cn(
            "overflow-hidden rounded-2xl border",
            "bg-white border-slate-200/80",
            "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
        )}>
            {/* Header */}
            <div className={cn(
                "flex items-center gap-3 px-5 py-4 border-b",
                "border-slate-100 bg-slate-50/60",
                "dark:border-white/[0.05] dark:bg-white/[0.02]"
            )}>
                <div className="flex h-8 w-8 items-center justify-center rounded-lg ring-1 bg-blue-50 text-blue-600 ring-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:ring-blue-500/20">
                    <Clock className="h-3.5 w-3.5" />
                </div>
                <div>
                    <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                        Aktivitas Terbaru
                    </p>
                    <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                        {recent.length} aktivitas terakhir
                    </p>
                </div>
            </div>

            {/* List */}
            <div className="p-3">
                {recent.length === 0 ? (
                    <div className="flex flex-col items-center justify-center py-10 text-slate-400 dark:text-slate-600">
                        <Clock className="h-8 w-8 mb-2 opacity-30" />
                        <p className="text-xs">Belum ada aktivitas</p>
                    </div>
                ) : (
                    <div className="space-y-0.5">
                        {recent.map((activity) => {
                            const cfg = statusConfig[activity.status as keyof typeof statusConfig] ?? statusConfig.pending;
                            const Icon = cfg.icon;

                            return (
                                <div
                                    key={activity.id}
                                    className="flex items-center justify-between rounded-lg px-3 py-2.5 transition-colors hover:bg-slate-50/80 dark:hover:bg-white/[0.02]"
                                >
                                    <div className="flex items-center gap-2.5 min-w-0">
                                        <div className={cn(
                                            "flex h-7 w-7 shrink-0 items-center justify-center rounded-lg",
                                            cfg.iconWrap
                                        )}>
                                            <Icon className="h-3.5 w-3.5" />
                                        </div>
                                        <div className="min-w-0">
                                            <p className="text-xs font-medium text-slate-700 dark:text-slate-200 truncate max-w-[160px]">
                                                {activity.title}
                                            </p>
                                            <p className="text-[10px] text-slate-400 dark:text-slate-600 mt-0.5">
                                                {activity.targetProducts?.join(", ")}
                                            </p>
                                        </div>
                                    </div>

                                    <span className={cn(
                                        "inline-flex items-center gap-1 rounded-md border px-1.5 py-0.5 text-[10px] font-medium shrink-0",
                                        cfg.badge
                                    )}>
                                        <Icon className="h-2.5 w-2.5" />
                                        {cfg.label}
                                    </span>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
        </div>
    );
}