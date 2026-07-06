import { Button } from "@kana-consultant/ui-kit";
import { Eye, Edit, RefreshCw, Send, Trash2, Calendar } from "lucide-react";
import type { Draft } from "@/services/draft";
import { Can } from "@/components/ui/Can";

interface ScheduleItemProps {
    schedule: Draft;
    isDailySchedule: (scheduledFor?: string) => boolean;
    getScheduleDisplay: (scheduledFor?: string) => string;
    onView: (schedule: Draft) => void;
    onEdit?: (schedule: Draft) => void;
    onReschedule?: (schedule: Draft) => void;
    onPublishNow?: (schedule: Draft) => void;
    onDelete?: (schedule: Draft) => void;
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
    return (
        <div className="flex flex-col gap-4 rounded-lg border p-4 transition-all hover:shadow-md sm:flex-row sm:items-center sm:justify-between">
            <div className="flex-1">
                <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:gap-3">
                    <h3 className="font-semibold text-sm sm:text-base line-clamp-2">{schedule.title}</h3>
                    <div className="flex flex-wrap items-center gap-2">
                        <span className="inline-flex items-center gap-1 rounded-full bg-blue-100 px-2 py-0.5 text-xs text-blue-700 dark:bg-blue-900 dark:text-blue-300">
                            <Calendar className="h-3 w-3" />
                            <span className="hidden sm:inline">Terjadwal</span>
                        </span>
                        {isDailySchedule(schedule.scheduled_for) && (
                            <span className="inline-flex items-center gap-1 rounded-full bg-purple-100 px-2 py-0.5 text-xs text-purple-700 dark:bg-purple-900 dark:text-purple-300">
                                <RefreshCw className="h-3 w-3" />
                                <span className="hidden sm:inline">Harian</span>
                            </span>
                        )}
                    </div>
                </div>

                <p className="mt-1 text-xs sm:text-sm text-slate-500 line-clamp-2">
                    {schedule.article.replace(/<[^>]*>/g, '').substring(0, 150)}...
                </p>

                <div className="mt-2 flex flex-wrap gap-3 text-xs text-slate-400">
                    <span className="flex items-center gap-1">
                        <Calendar className="h-3 w-3" />
                        Jadwal: {getScheduleDisplay(schedule.scheduled_for)}
                    </span>
                    <span className="flex items-center gap-1">
                        🎯 Target: {schedule.target_products?.length} produk
                    </span>
                    {schedule.has_image && (
                        <span className="flex items-center gap-1">
                            🖼️ <span className="hidden sm:inline">Ada gambar</span>
                            <span className="sm:hidden">Gambar</span>
                        </span>
                    )}
                </div>
            </div>

            {/* Action Buttons */}
            <div className="flex flex-wrap gap-2 sm:flex-nowrap sm:justify-end">

                {/* view — semua role */}
                <Button variant="outline" size="sm" onClick={() => onView(schedule)}>
                    <Eye className="h-4 w-4" />
                    <span className="hidden sm:inline">Lihat</span>
                </Button>

                {/* edit */}
                {/* <Can permission="schedule:edit:team">
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => onEdit?.(schedule)}
                    >
                        <Edit className="h-4 w-4" />
                        <span className="hidden sm:inline">Edit</span>
                    </Button>
                </Can> */}

                {/* reschedule = edit */}
                <Can permission="schedule:edit:team">
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => onReschedule?.(schedule)}
                    >
                        <RefreshCw className="h-4 w-4" />
                    </Button>
                </Can>

                {/* publish */}
                {/* <Can permission="schedule:publish:team">
                    <Button
                        variant="primary"
                        size="sm"
                        onClick={() => onPublishNow?.(schedule)}
                    >
                        <Send className="h-4 w-4" />
                    </Button>
                </Can> */}

                {/* delete */}
                <Can permission="schedule:delete:team">
                    <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => onDelete?.(schedule)}
                    >
                        <Trash2 className="h-4 w-4" />
                    </Button>
                </Can>

            </div>
        </div>
    );
}