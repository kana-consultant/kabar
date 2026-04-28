import { Button } from "@/components/ui/button";
import { Eye, Repeat, Trash2 } from "lucide-react";
import type { HistoryItem as HistoryItemType } from "@/types/history";

interface HistoryItemProps {
    item: HistoryItemType;
    onView: (item: HistoryItemType) => void;
    onRepost: (item: HistoryItemType) => void;
    onDelete: (item: HistoryItemType) => void;
    formatDate: (date: string) => string;
    getStatusData: (status: string) => { label: string; icon: string; color: string };
    getActionData: (action: string) => { label: string; icon: string };
}

export function HistoryItem({
    item,
    onView,
    onRepost,
    onDelete,
    formatDate,
    getStatusData,
    getActionData,
}: HistoryItemProps) {
    const status = getStatusData(item.status);

    const statusColorClass = {
        green: "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400",
        red: "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400",
        yellow: "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400",
    };

    return (
        <div className="flex items-center justify-between border-b pb-4 last:border-0">
            <div className="space-y-1 flex-1">
                <div className="flex items-center gap-2 flex-wrap">
                    <p className="font-medium">{item.title}</p>
                    <span className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium ${statusColorClass[status.color as keyof typeof statusColorClass]}`}>
                        {status.icon} {status.label}
                    </span>
                </div>
                <div className="flex flex-wrap gap-2 text-sm text-slate-500">
                    <span>Dikirim ke: {item.targetProducts.join(", ")}</span>
                    <span>•</span>
                    <span>{formatDate(item.publishedAt)}</span>
                </div>
                {item.errorMessage && (
                    <p className="text-xs text-red-500">{item.errorMessage}</p>
                )}
            </div>
            <div className="flex gap-2 ml-4">
                <Button variant="ghost" size="sm" onClick={() => onView(item)}>
                    <Eye className="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="sm" onClick={() => onRepost(item)}>
                    <Repeat className="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="sm" className="text-red-600" onClick={() => onDelete(item)}>
                    <Trash2 className="h-4 w-4" />
                </Button>
            </div>
        </div>
    );
}