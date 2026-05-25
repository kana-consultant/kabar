import { Button } from "@kana-consultant/ui-kit";
import { Eye, Repeat, Trash2, CheckCircle, XCircle, Clock, Send, Calendar, FileText, Package } from "lucide-react";
import type { HistoryItem as HistoryItemType } from "@/services/history";
import { cn } from "@/lib/utils";

interface HistoryItemProps {
    item: HistoryItemType;
    onView: (item: HistoryItemType) => void;
    onRepost: (item: HistoryItemType) => void;
    onDelete: (item: HistoryItemType) => void;
    formatDate: (date: string) => string;
    getStatusData: (status: string) => { label: string; icon: string; color: string };
    getActionData: (action: string) => { label: string; icon: string };
}

const statusConfig = {
    published: {
        icon: CheckCircle,
        badge: "bg-green-50 text-green-700 border-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:border-green-500/20",
        iconWrap: "bg-green-50 text-green-600 ring-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:ring-green-500/20",
        accent: "bg-green-500",
    },
    failed: {
        icon: XCircle,
        badge: "bg-red-50 text-red-700 border-red-200/60 dark:bg-red-500/10 dark:text-red-400 dark:border-red-500/20",
        iconWrap: "bg-red-50 text-red-600 ring-red-200/60 dark:bg-red-500/10 dark:text-red-400 dark:ring-red-500/20",
        accent: "bg-red-500",
    },
    pending: {
        icon: Clock,
        badge: "bg-amber-50 text-amber-700 border-amber-200/60 dark:bg-amber-500/10 dark:text-amber-400 dark:border-amber-500/20",
        iconWrap: "bg-amber-50 text-amber-600 ring-amber-200/60 dark:bg-amber-500/10 dark:text-amber-400 dark:ring-amber-500/20",
        accent: "bg-amber-500",
    },
} as const;

const actionIconMap: Record<string, React.ElementType> = {
    published: Send,
    scheduled: Calendar,
    draft_saved: FileText,
};

export function HistoryItem({
    item, onView, onRepost, onDelete, getStatusData, getActionData,
}: HistoryItemProps) {
    const status = getStatusData(item.status);
    const action = getActionData(item.action ?? "");
    const cfg = statusConfig[item.status as "published" | "failed" | "pending"] ?? statusConfig.pending;
    const StatusIcon = cfg.icon;
    const ActionIcon = actionIconMap[item.action ?? ""] ?? FileText;

    return (
        <div className={cn(
            "group relative flex items-start gap-3 rounded-xl border p-3.5 transition-all duration-200",
            "bg-white border-slate-200/80 hover:shadow-[0_1px_8px_rgba(0,0,0,0.06)]",
            "dark:bg-[#0f0d1a] dark:border-white/[0.06]",
            item.action === "published"
                ? "hover:border-green-200/80 dark:hover:border-green-500/20"
                : item.status === "failed"
                    ? "hover:border-red-200/80 dark:hover:border-red-500/20"
                    : "hover:border-amber-200/80 dark:hover:border-amber-500/20"
        )}>
            {/* Left accent */}
            <div className={cn(
                "absolute left-0 top-3 bottom-3 w-[3px] rounded-r-full opacity-0 transition-opacity duration-200 group-hover:opacity-100",
                cfg.accent
            )} />

            {/* Status icon */}
            <div className={cn(
                "mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-xl ring-1",
                cfg.iconWrap
            )}>
                <StatusIcon className="h-4 w-4" />
            </div>

            {/* Content */}
            <div className="flex-1 min-w-0">
                <div className="flex flex-wrap items-center gap-1.5">
                    <span className="text-sm font-medium text-slate-900 dark:text-white truncate max-w-[320px]">
                        {item.title}
                    </span>

                    {/* Status badge */}
                    <span className={cn(
                        "inline-flex items-center gap-1 rounded-md px-1.5 py-0.5 text-[10px] font-medium border",
                        cfg.badge
                    )}>
                        <StatusIcon className="h-2.5 w-2.5" />
                        {status.label}
                    </span>

                    {/* Action badge */}
                    <span className={cn(
                        "inline-flex items-center gap-1 rounded-md px-1.5 py-0.5 text-[10px] font-medium border",
                        "bg-slate-50 text-slate-500 border-slate-200/60",
                        "dark:bg-white/[0.04] dark:text-slate-500 dark:border-white/[0.06]"
                    )}>
                        <ActionIcon className="h-2.5 w-2.5" />
                        {action.label}
                    </span>
                </div>

                <div className="mt-1.5 flex flex-wrap gap-3 text-[11px] text-slate-400 dark:text-slate-600">
                    <span className="flex items-center gap-1">
                        <Package className="h-3 w-3 shrink-0" />
                        {item.targetProducts?.join(", ")}
                    </span>
                    <span className="flex items-center gap-1">
                        <Clock className="h-3 w-3 shrink-0" />
                        {new Date(item.publishedAt).toLocaleString("id-ID", {
                            day: "numeric", month: "short", year: "numeric",
                            hour: "2-digit", minute: "2-digit",
                            timeZoneName: "short",
                        })}
                    </span>
                </div>
            </div>

            {/* Actions — reveal on hover */}
            <div className="flex shrink-0 items-center gap-1 opacity-0 translate-x-1 transition-all duration-150 group-hover:opacity-100 group-hover:translate-x-0">
                <Button
                    variant="ghost" size="icon" title="Lihat"
                    className="h-7 w-7 rounded-lg text-slate-400 hover:text-green-600 hover:bg-green-50 dark:hover:text-green-400 dark:hover:bg-green-500/10"
                    onClick={() => onView(item)}
                >
                    <Eye className="h-3.5 w-3.5" />
                </Button>
                <Button
                    variant="ghost" size="icon" title="Repost"
                    className="h-7 w-7 rounded-lg text-slate-400 hover:text-blue-600 hover:bg-blue-50 dark:hover:text-blue-400 dark:hover:bg-blue-500/10"
                    onClick={() => onRepost(item)}
                >
                    <Repeat className="h-3.5 w-3.5" />
                </Button>
                <Button
                    variant="ghost" size="icon" title="Hapus"
                    className="h-7 w-7 rounded-lg text-slate-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10"
                    onClick={() => onDelete(item)}
                >
                    <Trash2 className="h-3.5 w-3.5" />
                </Button>
            </div>
        </div>
    );
}