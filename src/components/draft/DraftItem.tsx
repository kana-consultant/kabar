import { Button } from "@/components/ui/button";
import { Eye, Edit, Calendar, Send, Trash2, Clock, FileText, CheckCircle, Package, ImageIcon } from "lucide-react";
import type { Draft } from "@/types/draft";
import { cn } from "@/lib/utils";

interface DraftItemProps {
    draft: Draft;
    onView: (draft: Draft) => void;
    onEdit: (draft: Draft) => void;
    onSchedule: (draft: Draft) => void;
    onPublishNow: (draft: Draft) => void;
    onDelete: (draft: Draft) => void;
    formatDate: (date: string) => string;
}

const statusConfig = {
    draft: {
        icon: FileText,
        label: "Draft",
        badge: "bg-amber-50 text-amber-700 border-amber-200/60 dark:bg-amber-500/10 dark:text-amber-400 dark:border-amber-500/20",
        iconWrap: "bg-amber-50 text-amber-600 ring-amber-200/60 dark:bg-amber-500/10 dark:text-amber-400 dark:ring-amber-500/20",
        accent: "bg-amber-500",
        hover: "hover:border-amber-200/80 dark:hover:border-amber-500/20",
    },
    scheduled: {
        icon: Calendar,
        label: "Terjadwal",
        badge: "bg-blue-50 text-blue-700 border-blue-200/60 dark:bg-purple-500/10 dark:text-purple-400 dark:border-purple-500/20",
        iconWrap: "bg-blue-50 text-blue-600 ring-blue-200/60 dark:bg-purple-500/10 dark:text-purple-400 dark:ring-purple-500/20",
        accent: "bg-blue-500 dark:bg-purple-500",
        hover: "hover:border-blue-200/80 dark:hover:border-purple-500/20",
    },
    published: {
        icon: CheckCircle,
        label: "Terbit",
        badge: "bg-green-50 text-green-700 border-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:border-green-500/20",
        iconWrap: "bg-green-50 text-green-600 ring-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:ring-green-500/20",
        accent: "bg-green-500",
        hover: "hover:border-green-200/80 dark:hover:border-green-500/20",
    },
} as const;

export function DraftItem({
    draft, onView, onEdit, onSchedule, onPublishNow, onDelete, formatDate,
}: DraftItemProps) {
    const cfg = statusConfig[draft.status as keyof typeof statusConfig] ?? statusConfig.draft;
    const StatusIcon = cfg.icon;

    const scheduledLabel = draft.scheduledFor
        ? draft.scheduledFor.startsWith("daily:")
            ? `Setiap hari ${draft.scheduledFor.replace("daily:", "")} `
            : draft.scheduledFor
        : null;

    return (
        <div className={cn(
            "group relative flex items-start gap-3 rounded-xl border p-3.5 transition-all duration-200",
            "bg-white border-slate-200/80",
            "dark:bg-[#0f0d1a] dark:border-white/[0.06]",
            cfg.hover
        )}>
            {/* Left accent */}
            <div className={cn(
                "absolute left-0 top-3 bottom-3 w-[3px] rounded-r-full opacity-0 transition-opacity duration-200 group-hover:opacity-100",
                cfg.accent
            )} />

            {/* Icon */}
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
                        {draft.title}
                    </span>

                    <span className={cn(
                        "inline-flex items-center gap-1 rounded-md border px-1.5 py-0.5 text-[10px] font-medium",
                        cfg.badge
                    )}>
                        <StatusIcon className="h-2.5 w-2.5" />
                        {cfg.label}
                        {draft.status === "scheduled" && scheduledLabel && (
                            <span className="opacity-70">• {scheduledLabel}</span>
                        )}
                    </span>
                </div>

                <p className="mt-1.5 text-xs text-slate-400 dark:text-slate-500 line-clamp-1 leading-relaxed">
                    {draft.article.replace(/<[^>]*>/g, '').substring(0, 140)}...
                </p>

                <div className="mt-2 flex flex-wrap gap-3 text-[11px] text-slate-400 dark:text-slate-600">
                    <span className="flex items-center gap-1">
                        <Clock className="h-3 w-3 shrink-0" />
                        {formatDate(draft.createdAt)}
                    </span>
                    <span className="flex items-center gap-1">
                        <Package className="h-3 w-3 shrink-0" />
                        {draft.targetProducts?.length} produk
                    </span>
                    {(draft.hasImage || draft.imageUrl) && (
                        <span className="flex items-center gap-1">
                            <ImageIcon className="h-3 w-3 shrink-0" />
                            Ada gambar
                        </span>
                    )}
                </div>
            </div>

            {/* Actions — reveal on hover */}
            <div className="flex shrink-0 items-center gap-1 opacity-0 translate-x-1 transition-all duration-150 group-hover:opacity-100 group-hover:translate-x-0">
                {[
                    { icon: Eye, label: "Lihat", action: onView },
                    { icon: Edit, label: "Edit", action: onEdit },
                    ...(draft.status === "draft"
                        ? [{ icon: Calendar, label: "Jadwalkan", action: onSchedule }]
                        : []),
                ].map(({ icon: Icon, label, action }) => (
                    <Button
                        key={label}
                        variant="ghost" size="icon" title={label}
                        className={cn(
                            "h-7 w-7 rounded-lg text-slate-400 transition-colors",
                            "hover:text-green-600 hover:bg-green-50",
                            "dark:hover:text-purple-400 dark:hover:bg-purple-500/10"
                        )}
                        onClick={() => action(draft)}
                    >
                        <Icon className="h-3.5 w-3.5" />
                    </Button>
                ))}

                {(draft.status === "draft" || draft.status === "scheduled") && (
                    <Button
                        size="sm"
                        className={cn(
                            "h-7 gap-1.5 px-2.5 text-[11px] font-medium rounded-lg ml-0.5",
                            "bg-green-600 hover:bg-green-700 text-white shadow-sm",
                            "dark:bg-purple-600 dark:hover:bg-purple-700"
                        )}
                        onClick={() => onPublishNow(draft)}
                    >
                        <Send className="h-3 w-3" />
                        Terbitkan
                    </Button>
                )}

                <Button
                    variant="ghost" size="icon" title="Hapus"
                    className="h-7 w-7 rounded-lg text-slate-300 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10 ml-0.5"
                    onClick={() => onDelete(draft)}
                >
                    <Trash2 className="h-3.5 w-3.5" />
                </Button>
            </div>
        </div>
    );
}