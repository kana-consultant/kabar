import { Button } from "@/components/ui/button";
import { Eye, Edit, Calendar, Send, Trash2, Clock, FileText, CheckCircle } from "lucide-react";
import type { Draft } from "@//services/draft";

interface DraftItemProps {
    draft: Draft;
    onView: (draft: Draft) => void;
    onEdit: (draft: Draft) => void;
    onSchedule: (draft: Draft) => void;
    onPublishNow: (draft: Draft) => void;
    onDelete: (draft: Draft) => void;
    formatDate: (date: string) => string;
}

function getStatusBadge(status: string, scheduledFor?: string) {
    switch (status) {
        case "published":
            return (
                <span className="inline-flex items-center gap-1 rounded-full bg-green-100 px-2 py-0.5 text-xs text-green-700 dark:bg-green-900 dark:text-green-300">
                    <CheckCircle className="h-3 w-3" />
                    <span className="hidden sm:inline">Terbit</span>
                </span>
            );
        case "scheduled":
            return (
                <span className="inline-flex items-center gap-1 rounded-full bg-blue-100 px-2 py-0.5 text-xs text-blue-700 dark:bg-blue-900 dark:text-blue-300">
                    <Calendar className="h-3 w-3" />
                    <span className="hidden sm:inline">
                        Terjadwal {scheduledFor?.replace("daily:", "setiap hari jam ")}
                    </span>
                </span>
            );
        default:
            return (
                <span className="inline-flex items-center gap-1 rounded-full bg-yellow-100 px-2 py-0.5 text-xs text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300">
                    <FileText className="h-3 w-3" />
                    <span className="hidden sm:inline">Draft</span>
                </span>
            );
    }
}

export function DraftItem({
    draft,
    onView,
    onEdit,
    onSchedule,
    onPublishNow,
    onDelete,
    formatDate,
}: DraftItemProps) {
    return (
        <div className="flex flex-col gap-4 rounded-lg border p-4 transition-all hover:shadow-md sm:flex-row sm:items-center sm:justify-between">
            <div className="flex-1">
                <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:gap-3">
                    <h3 className="font-semibold text-sm sm:text-base line-clamp-2">{draft.title}</h3>
                    <div className="flex-shrink-0">
                        {getStatusBadge(draft.status as any, draft.scheduled_for)}
                    </div>
                </div>
                <p className="mt-1 text-xs sm:text-sm text-slate-500 line-clamp-2">
                    {draft.article.replace(/<[^>]*>/g, '').substring(0, 150)}...
                </p>
                <div className="mt-2 flex flex-wrap gap-3 text-xs text-slate-400">
                    <span className="flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        Dibuat: {formatDate(draft.created_at)}
                    </span>
                    <span className="flex items-center gap-1">
                        Target: {draft.target_products?.length} produk
                    </span>
                </div>
            </div>

            {/* Action Buttons - Responsive Grid */}
            <div className="flex flex-wrap gap-2 sm:flex-nowrap sm:justify-end">
                <Button variant="outline" size="sm" onClick={() => onView(draft)}>
                    <Eye className="h-4 w-4 sm:mr-1" />
                    <span className="hidden sm:inline">Lihat</span>
                </Button>

                <Button variant="outline" size="sm" onClick={() => onEdit(draft)}>
                    <Edit className="h-4 w-4 sm:mr-1" />
                    <span className="hidden sm:inline">Edit</span>
                </Button>

                {draft.status === "draft" && (
                    <>
                        <Button variant="outline" size="sm" onClick={() => onSchedule(draft)}>
                            <Calendar className="h-4 w-4 sm:mr-1" />
                            <span className="hidden sm:inline">Jadwal</span>
                        </Button>
                        <Button variant="default" size="sm" onClick={() => onPublishNow(draft)}>
                            <Send className="h-4 w-4 sm:mr-1" />
                            <span className="hidden sm:inline">Terbit</span>
                        </Button>
                    </>
                )}

                {draft.status === "scheduled" && (
                    <Button variant="default" size="sm" onClick={() => onPublishNow(draft)}>
                        <Send className="h-4 w-4 sm:mr-1" />
                        <span className="hidden sm:inline">Terbitkan Sekarang</span>
                        <span className="sm:hidden">Terbit</span>
                    </Button>
                )}

                <Button variant="destructive" size="sm" onClick={() => onDelete(draft)}>
                    <Trash2 className="h-4 w-4 sm:mr-1" />
                    <span className="hidden sm:inline">Hapus</span>
                </Button>
            </div>
        </div>
    );
}