import { FileText, Clock, Calendar, CheckCircle } from "lucide-react";
import type { Draft } from "@/services/draft";
import { cn } from "@/lib/utils";

interface DraftStatsProps {
    drafts: Draft[];
}

export function DraftStats({ drafts }: DraftStatsProps) {
    const stats = [
        {
            label: "Total Draft",
            value: drafts.length,
            icon: FileText,
            iconClass: "bg-blue-50 text-blue-600 ring-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:ring-blue-500/20",
            dotClass: "bg-blue-500",
            footer: "Semua konten",
        },
        {
            label: "Belum Terbit",
            value: drafts.filter(d => d.status === "draft").length,
            icon: Clock,
            iconClass: "bg-amber-50 text-amber-600 ring-amber-200/60 dark:bg-amber-500/10 dark:text-amber-400 dark:ring-amber-500/20",
            dotClass: "bg-amber-500",
            footer: "Menunggu publish",
        },
        {
            label: "Terjadwal",
            value: drafts.filter(d => d.status === "scheduled").length,
            icon: Calendar,
            iconClass: "bg-blue-50 text-blue-600 ring-blue-200/60 dark:bg-purple-500/10 dark:text-purple-400 dark:ring-purple-500/20",
            dotClass: "bg-blue-500 dark:bg-purple-500",
            footer: "Akan terposting",
        },
        {
            label: "Sudah Terbit",
            value: drafts.filter(d => d.status === "published").length,
            icon: CheckCircle,
            iconClass: "bg-green-50 text-green-600 ring-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:ring-green-500/20",
            dotClass: "bg-green-500",
            footer: "Berhasil dipublikasi",
        },
    ];

    return (
        <div className="grid gap-3 md:grid-cols-4">
            {stats.map(({ label, value, icon: Icon, iconClass, dotClass, footer }) => (
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
                        {value}
                    </span>

                    <div className="flex items-center gap-1.5">
                        <span className={cn("h-1.5 w-1.5 rounded-full", dotClass)} />
                        <span className="text-[11px] text-slate-400 dark:text-slate-600">
                            {footer}
                        </span>
                    </div>
                </div>
            ))}
        </div>
    );
}