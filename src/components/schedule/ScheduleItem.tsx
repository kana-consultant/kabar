import { Button } from "@/components/ui/button";
import { Eye, Edit, RefreshCw, Send, Trash2, Calendar } from "lucide-react";
import type { Draft } from "@/types/draft";

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
    return (
        <div className="flex items-center justify-between rounded-lg border p-4 transition-all hover:shadow-md">
            <div className="flex-1">
                <div className="flex items-center gap-3">
                    <h3 className="font-semibold">{schedule.title}</h3>
                    <span className="inline-flex items-center gap-1 rounded-full bg-blue-100 px-2 py-0.5 text-xs text-blue-700 dark:bg-blue-900 dark:text-blue-300">
                        <Calendar className="h-3 w-3" />
                        Terjadwal
                    </span>
                    {isDailySchedule(schedule.scheduledFor) && (
                        <span className="inline-flex items-center gap-1 rounded-full bg-purple-100 px-2 py-0.5 text-xs text-purple-700 dark:bg-purple-900 dark:text-purple-300">
                            <RefreshCw className="h-3 w-3" />
                            Harian
                        </span>
                    )}
                </div>
                <p className="mt-1 text-sm text-slate-500 line-clamp-2">
                    {schedule.article.replace(/<[^>]*>/g, '').substring(0, 150)}...
                </p>
                <div className="mt-2 flex flex-wrap gap-4 text-xs text-slate-400">
                    <span className="flex items-center gap-1">
                        <Calendar className="h-3 w-3" />
                        Jadwal: {getScheduleDisplay(schedule.scheduledFor)}
                    </span>
                    <span className="flex items-center gap-1">
                        🎯 Target: {schedule.targetProducts.length} produk
                    </span>
                    {schedule.hasImage && <span>🖼️ Ada gambar</span>}
                </div>
            </div>
            <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={() => onView(schedule)}>
                    <Eye className="h-4 w-4" />
                </Button>
                <Button variant="outline" size="sm" onClick={() => onEdit(schedule)}>
                    <Edit className="h-4 w-4" />
                </Button>
                <Button variant="outline" size="sm" onClick={() => onReschedule(schedule)}>
                    <RefreshCw className="h-4 w-4" />
                </Button>
                <Button variant="default" size="sm" onClick={() => onPublishNow(schedule)}>
                    <Send className="h-4 w-4" />
                    Terbitkan Sekarang
                </Button>
                <Button variant="destructive" size="sm" onClick={() => onDelete(schedule)}>
                    <Trash2 className="h-4 w-4" />
                </Button>
            </div>
        </div>
    );
}