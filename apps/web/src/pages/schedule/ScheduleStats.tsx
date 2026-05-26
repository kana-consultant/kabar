import { Calendar, Clock, RefreshCw } from "lucide-react";
import { cn } from "@/lib/utils";

interface ScheduleStatsProps {
    totalItems: number;
    isDailySchedule?: never; // akan dipakai nanti
}

const stats = [
    {
        key: "total",
        label: "Total Jadwal",
        icon: Calendar,
        iconClass: "bg-green-50 text-green-600 ring-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:ring-green-500/20",
        dotClass: "bg-green-500",
        footer: "Aktif minggu ini",
    },
    {
        key: "oneTime",
        label: "One Time",
        icon: Clock,
        iconClass: "bg-teal-50 text-teal-600 ring-teal-200/60 dark:bg-teal-500/10 dark:text-teal-400 dark:ring-teal-500/20",
        dotClass: "bg-teal-500",
        footer: "Sekali posting",
    },
    {
        key: "daily",
        label: "Daily Schedule",
        icon: RefreshCw,
        iconClass: "bg-violet-50 text-violet-600 ring-violet-200/60 dark:bg-violet-500/10 dark:text-violet-400 dark:ring-violet-500/20",
        dotClass: "bg-violet-500 dark:bg-purple-400",
        footer: "Berulang harian",
    },
] as const;

export function ScheduleStats({ totalItems }: ScheduleStatsProps) {
    // TODO: ganti dengan data dari backend
    const values = {
        total: totalItems,
        oneTime: 0,
        daily: 0,
    };

    return (
        <div className="grid gap-3 md:grid-cols-3">
            {stats.map(({ key, label, icon: Icon, iconClass, dotClass, footer }) => (
                <div
                    key={key}
                    className={cn(
                        "group flex flex-col gap-3 rounded-xl border p-5 transition-all duration-200",
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
                        {values[key]}
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