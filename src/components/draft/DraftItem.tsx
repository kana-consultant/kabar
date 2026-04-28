import { Button } from "@/components/ui/button";
import { Eye, Edit, Calendar, Send, Trash2, Clock, FileText, CheckCircle } from "lucide-react";
import type { Draft } from "@/types/draft";

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
                    Terbit
                </span>
            );
        case "scheduled":
            return (
                <span className="inline-flex items-center gap-1 rounded-full bg-blue-100 px-2 py-0.5 text-xs text-blue-700 dark:bg-blue-900 dark:text-blue-300">
                    <Calendar className="h-3 w-3" />
                    Terjadwal {scheduledFor?.replace("daily:", "setiap hari jam ")}
                </span>
            );
        default:
            return (
                <span className="inline-flex items-center gap-1 rounded-full bg-yellow-100 px-2 py-0.5 text-xs text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300">
                    <FileText className="h-3 w-3" />
                    Draft
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
        <div className="flex items-center justify-between rounded-lg border p-4 transition-all hover:shadow-md">
            <div className="flex-1">
                <div className="flex items-center gap-3">
                    <h3 className="font-semibold">{draft.title}</h3>
                    {getStatusBadge(draft.status, draft.scheduledFor)}
                </div>
                <p className="mt-1 text-sm text-slate-500 line-clamp-2">
                    {draft.article.replace(/<[^>]*>/g, '').substring(0, 150)}...
                </p>
                <div className="mt-2 flex flex-wrap gap-4 text-xs text-slate-400">
                    <span className="flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        Dibuat: {formatDate(draft.createdAt)}
                    </span>
                    <span className="flex items-center gap-1">
                         Target: {draft.targetProducts.length} produk
                    </span>
                    {draft.hasImage && (
                        <span className="flex items-center gap-1">
                             Ada gambar
                        </span>
                    )}
                    {draft.imageUrl && !draft.hasImage && (
                        <span className="flex items-center gap-1">
                            Ada gambar
                        </span>
                    )}
                </div>
            </div>
            <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={() => onView(draft)}>
                    <Eye className="h-4 w-4" />
                </Button>
                <Button variant="outline" size="sm" onClick={() => onEdit(draft)}>
                    <Edit className="h-4 w-4" />
                </Button>
                {draft.status === "draft" && (
                    <>
                        <Button variant="outline" size="sm" onClick={() => onSchedule(draft)}>
                            <Calendar className="h-4 w-4" />
                        </Button>
                        <Button variant="default" size="sm" onClick={() => onPublishNow(draft)}>
                            <Send className="h-4 w-4" />
                        </Button>
                    </>
                )}
                {draft.status === "scheduled" && (
                    <Button variant="default" size="sm" onClick={() => onPublishNow(draft)}>
                        <Send className="h-4 w-4" />
                        Terbitkan Sekarang
                    </Button>
                )}
                <Button variant="destructive" size="sm" onClick={() => onDelete(draft)}>
                    <Trash2 className="h-4 w-4" />
                </Button>
            </div>
        </div>
    );
}