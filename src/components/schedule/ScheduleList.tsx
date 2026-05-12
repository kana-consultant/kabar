import { Button } from "@/components/ui/button";
import { Calendar, Plus, SlidersHorizontal } from "lucide-react";
import { useNavigate } from "@tanstack/react-router";
import { ScheduleItem } from "./ScheduleItem";
import type { Draft } from "@/types/draft";
import { cn } from "@/lib/utils";

interface ScheduleListProps {
    schedules: Draft[];
    isDailySchedule: (scheduledFor?: string) => boolean;
    getScheduleDisplay: (scheduledFor?: string) => string;
    onView: (schedule: Draft) => void;
    onEdit: (schedule: Draft) => void;
    onReschedule: (schedule: Draft) => void;
    onPublishNow: (schedule: Draft) => void;
    onDelete: (schedule: Draft) => void;
}

export function ScheduleList({
    schedules,
    isDailySchedule,
    getScheduleDisplay,
    onView,
    onEdit,
    onReschedule,
    onPublishNow,
    onDelete,
}: ScheduleListProps) {
    const navigate = useNavigate();

    if (schedules.length === 0) {
        return (
            <div className={cn(
                "flex flex-col items-center justify-center rounded-2xl border py-20",
                "bg-white border-slate-200/80",
                "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
            )}>
                {/* Glow ring behind icon */}
                <div className="relative">
                    <div className={cn(
                        "absolute inset-0 rounded-full blur-xl opacity-40",
                        "bg-green-300 dark:bg-purple-500"
                    )} />
                    <div className={cn(
                        "relative flex h-14 w-14 items-center justify-center rounded-2xl",
                        "bg-green-50 text-green-600 ring-1 ring-green-200/60 shadow-sm",
                        "dark:bg-purple-500/10 dark:text-purple-400 dark:ring-purple-500/20"
                    )}>
                        <Calendar className="h-6 w-6" />
                    </div>
                </div>

                <p className="mt-5 text-sm font-medium text-slate-700 dark:text-slate-200">
                    Belum ada jadwal posting
                </p>
                <p className="mt-1 text-xs text-slate-400 dark:text-slate-500 max-w-xs text-center">
                    Buat jadwal baru untuk mulai menjadwalkan konten agar terposting otomatis.
                </p>

                <Button
                    size="sm"
                    className={cn(
                        "mt-6 gap-2 rounded-lg shadow-sm font-medium text-xs",
                        "bg-green-600 hover:bg-green-700 text-white",
                        "dark:bg-purple-600 dark:hover:bg-purple-700 dark:shadow-purple-950"
                    )}
                    onClick={() => navigate({ to: "/generate" })}
                >
                    <Plus className="h-3.5 w-3.5" />
                    Buat Jadwal Baru
                </Button>
            </div>
        );
    }

    return (
        <div className={cn(
            "overflow-hidden rounded-2xl border",
            "bg-white border-slate-200/80",
            "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
        )}>
            {/* Header */}
            <div className={cn(
                "flex items-center justify-between px-5 py-3.5 border-b",
                "border-slate-100 bg-slate-50/60",
                "dark:border-white/[0.05] dark:bg-white/[0.02]"
            )}>
                <div className="flex items-center gap-3">
                    <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                        Daftar Jadwal
                    </p>
                    <span className={cn(
                        "inline-flex h-5 min-w-5 items-center justify-center rounded-md px-1.5 text-[10px] font-semibold tabular-nums",
                        "bg-green-100 text-green-700",
                        "dark:bg-purple-500/15 dark:text-purple-300"
                    )}>
                        {schedules.length}
                    </span>
                </div>

                <div className="flex items-center gap-2">
                    <Button
                        variant="ghost"
                        size="sm"
                        className={cn(
                            "h-7 gap-1.5 px-2.5 text-[11px] rounded-lg text-slate-500",
                            "hover:text-green-600 hover:bg-green-50",
                            "dark:hover:text-purple-400 dark:hover:bg-purple-500/10"
                        )}
                    >
                        <SlidersHorizontal className="h-3 w-3" />
                        Filter
                    </Button>
                    <Button
                        size="sm"
                        className={cn(
                            "h-7 gap-1.5 px-2.5 text-[11px] rounded-lg font-medium",
                            "bg-green-600 hover:bg-green-700 text-white shadow-sm",
                            "dark:bg-purple-600 dark:hover:bg-purple-700"
                        )}
                        onClick={() => navigate({ to: "/generate" })}
                    >
                        <Plus className="h-3 w-3" />
                        Jadwal Baru
                    </Button>
                </div>
            </div>

            {/* Items */}
            <div className="p-3 space-y-1.5">
                {schedules.map((schedule) => (
                    <ScheduleItem
                        key={schedule.id}
                        schedule={schedule}
                        isDailySchedule={isDailySchedule}
                        getScheduleDisplay={getScheduleDisplay}
                        onView={onView}
                        onEdit={onEdit}
                        onReschedule={onReschedule}
                        onPublishNow={onPublishNow}
                        onDelete={onDelete}
                    />
                ))}
            </div>
        </div>
    );
}