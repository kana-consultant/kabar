import { Calendar, RefreshCw } from "lucide-react";
import { useSchedule } from "@/hooks/useSchedule"; // sesuaikan dengan hook yang ada
import { cn } from "@/lib/utils";

export function ScheduleHeader() {
    const { schedules } = useSchedule();

    // Ambil 3 jadwal mendatang
    const upcoming = schedules
        .filter(s => s.status === "scheduled")
        .slice(0, 3);

    const getTimeDisplay = (scheduledFor?: string) => {
        if (!scheduledFor) return { time: "—", label: "" };
        if (scheduledFor.startsWith("daily:")) {
            return { time: scheduledFor.replace("daily:", ""), label: "Tiap hari" };
        }
        const d = new Date(scheduledFor);
        const today = new Date();
        const isToday = d.toDateString() === today.toDateString();
        return {
            time: d.toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit" }),
            label: isToday
                ? "Hari ini"
                : d.toLocaleDateString("id-ID", { day: "numeric", month: "short" }),
        };
    };

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
                <div className="flex h-8 w-8 items-center justify-center rounded-lg ring-1 bg-green-50 text-green-600 ring-green-200/60 dark:bg-purple-500/10 dark:text-purple-400 dark:ring-purple-500/20">
                    <Calendar className="h-3.5 w-3.5" />
                </div>
                <div>
                    <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                        Jadwal Mendatang
                    </p>
                    <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                        Konten yang akan terposting
                    </p>
                </div>
            </div>

            {/* List */}
            <div className="p-3">
                {upcoming.length === 0 ? (
                    <div className="flex flex-col items-center justify-center py-10 text-slate-400 dark:text-slate-600">
                        <Calendar className="h-8 w-8 mb-2 opacity-30" />
                        <p className="text-xs">Tidak ada jadwal mendatang</p>
                    </div>
                ) : (
                    <div className="space-y-0.5">
                        {upcoming.map((schedule) => {
                            const { time, label } = getTimeDisplay(schedule.scheduledFor);
                            const isDaily = schedule.scheduledFor?.startsWith("daily:");

                            return (
                                <div
                                    key={schedule.id}
                                    className="flex items-center gap-3 rounded-lg px-3 py-2.5 transition-colors hover:bg-slate-50/80 dark:hover:bg-white/[0.02]"
                                >
                                    {/* Time column */}
                                    <div className="flex flex-col items-center min-w-[36px] shrink-0">
                                        <span className="text-xs font-semibold text-slate-800 dark:text-slate-200 tabular-nums">
                                            {time}
                                        </span>
                                        <span className="text-[9px] text-slate-400 dark:text-slate-600 mt-0.5">
                                            {label}
                                        </span>
                                    </div>

                                    {/* Divider */}
                                    <div className="w-px self-stretch bg-slate-100 dark:bg-white/[0.05]" />

                                    {/* Info */}
                                    <div className="flex-1 min-w-0">
                                        <p className="text-xs font-medium text-slate-700 dark:text-slate-200 truncate">
                                            {schedule.title}
                                        </p>
                                        <div className="flex items-center gap-1.5 mt-0.5">
                                            <span className="text-[10px] text-slate-400 dark:text-slate-600">
                                                {schedule.targetProducts?.join(", ")}
                                            </span>
                                            {isDaily && (
                                                <span className="inline-flex items-center gap-0.5 text-[9px] font-medium text-purple-500 dark:text-purple-400">
                                                    <RefreshCw className="h-2 w-2" />
                                                    Harian
                                                </span>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
        </div>
    );
}