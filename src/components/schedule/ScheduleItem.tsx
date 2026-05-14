import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
    Eye, Edit, RefreshCw, Send, Trash2,
    Calendar, Package, ImageIcon, Clock
} from "lucide-react";
import type { Draft } from "@/services/draft";
import { cn } from "@/lib/utils";

interface ScheduleItemProps {
    schedule: Draft;
    isDailySchedule: (scheduledFor?: string) => boolean;
    getScheduleDisplay: (scheduledFor?: string) => string;
    onView: (schedule: Draft) => void;
    onEdit: (schedule: Draft) => void;
    onReschedule: (schedule: Draft) => void;
    onPublishNow: (schedule: Draft) => void;
    onDelete: (schedule: Draft) => void;
}

export function ScheduleItem({
    schedule,
    isDailySchedule,
    getScheduleDisplay,
    onView,
    onEdit,
    onReschedule,
    onPublishNow,
    onDelete,
}: ScheduleItemProps) {
    const isDaily = isDailySchedule(schedule.scheduledFor);

    return (
        <div className={cn(
            "group relative flex items-start gap-4 rounded-xl border p-4 transition-all duration-200",
            // Light
            "bg-white border-slate-200/80 hover:border-green-300/60 hover:shadow-[0_2px_12px_rgba(22,163,74,0.08)]",
            // Dark
            "dark:bg-[#0f0d1a] dark:border-white/[0.06] dark:hover:border-purple-500/30 dark:hover:shadow-[0_2px_20px_rgba(139,92,246,0.08)]"
        )}>
            {/* Subtle left accent line */}
            <div className={cn(
                "absolute left-0 top-4 bottom-4 w-[3px] rounded-full opacity-0 transition-opacity duration-200",
                "group-hover:opacity-100",
                "bg-green-500 dark:bg-purple-500"
            )} />

            {/* Icon */}
            <div className={cn(
                "mt-0.5 flex h-10 w-10 shrink-0 items-center justify-center rounded-xl transition-colors",
                "bg-green-50 text-green-600 ring-1 ring-green-200/60",
                "dark:bg-purple-500/10 dark:text-purple-400 dark:ring-purple-500/20"
            )}>
                <Calendar className="h-4 w-4" />
            </div>

            {/* Content */}
            <div className="flex-1 min-w-0">
                <div className="flex flex-wrap items-center gap-2">
                    <h3 className="font-medium text-sm text-slate-900 dark:text-white truncate max-w-[360px]">
                        {schedule.title}
                    </h3>
                    <Badge variant="secondary" className={cn(
                        "gap-1 rounded-md px-1.5 py-0 text-[10px] font-medium h-5",
                        "bg-green-50 text-green-700 border-green-200/60",
                        "dark:bg-purple-500/10 dark:text-purple-300 dark:border-purple-500/20"
                    )}>
                        <Clock className="h-2.5 w-2.5" />
                        Terjadwal
                    </Badge>
                    {isDaily && (
                        <Badge variant="secondary" className={cn(
                            "gap-1 rounded-md px-1.5 py-0 text-[10px] font-medium h-5",
                            "bg-emerald-50 text-emerald-700 border-emerald-200/60",
                            "dark:bg-violet-500/10 dark:text-violet-300 dark:border-violet-500/20"
                        )}>
                            <RefreshCw className="h-2.5 w-2.5" />
                            Harian
                        </Badge>
                    )}
                </div>

                <p className="mt-1.5 text-xs text-slate-400 dark:text-slate-500 line-clamp-1 leading-relaxed">
                    {schedule.article.replace(/<[^>]*>/g, '').substring(0, 140)}...
                </p>

                <div className="mt-2.5 flex flex-wrap items-center gap-3 text-[11px] text-slate-400 dark:text-slate-600">
                    <span className="flex items-center gap-1.5">
                        <Calendar className="h-3 w-3 shrink-0" />
                        {getScheduleDisplay(schedule.scheduledFor)}
                    </span>
                    <span className="flex items-center gap-1.5">
                        <Package className="h-3 w-3 shrink-0" />
                        {schedule.targetProducts?.length} produk
                    </span>
                    {schedule.hasImage && (
                        <span className="flex items-center gap-1.5">
                            <ImageIcon className="h-3 w-3 shrink-0" />
                            Ada gambar
                        </span>
                    )}
                </div>
            </div>

            {/* Actions — hidden until hover */}
            <div className={cn(
                "flex shrink-0 items-center gap-1 transition-all duration-150",
                "opacity-0 translate-x-1 group-hover:opacity-100 group-hover:translate-x-0"
            )}>
                {[
                    { icon: Eye, label: "Lihat", action: onView },
                    { icon: Edit, label: "Edit", action: onEdit },
                    { icon: RefreshCw, label: "Jadwalkan ulang", action: onReschedule },
                ].map(({ icon: Icon, label, action }) => (
                    <Button
                        key={label}
                        variant="ghost"
                        size="icon"
                        title={label}
                        className={cn(
                            "h-7 w-7 rounded-lg text-slate-400 transition-colors",
                            "hover:text-green-600 hover:bg-green-50",
                            "dark:hover:text-purple-400 dark:hover:bg-purple-500/10"
                        )}
                        onClick={() => action(schedule)}
                    >
                        <Icon className="h-3.5 w-3.5" />
                    </Button>
                ))}

                <Button
                    size="sm"
                    className={cn(
                        "h-7 gap-1.5 px-2.5 text-[11px] font-medium rounded-lg ml-0.5",
                        "bg-green-600 hover:bg-green-700 text-white shadow-sm",
                        "dark:bg-purple-600 dark:hover:bg-purple-700 dark:shadow-purple-900/30"
                    )}
                    onClick={() => onPublishNow(schedule)}
                >
                    <Send className="h-3 w-3" />
                    Terbitkan
                </Button>

                <Button
                    variant="ghost"
                    size="icon"
                    title="Hapus"
                    className="h-7 w-7 rounded-lg text-slate-300 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10 ml-0.5"
                    onClick={() => onDelete(schedule)}
                >
                    <Trash2 className="h-3.5 w-3.5" />
                </Button>
            </div>
        </div>
    );
}