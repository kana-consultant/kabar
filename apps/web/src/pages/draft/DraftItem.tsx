import { cn } from "@/lib/utils";
import { Button } from "@kana-consultant/ui-kit"
import { Clock, Package, ImageIcon, Eye, Edit, Calendar, Send, Trash2, GitCompare, BarChart2 } from "lucide-react";
import type { Draft } from "@/services/draft";
import { Can } from "@/components/ui/Can";
import { useAuth } from "@/contexts/AuthContext";

const statusConfig = {
    draft: {
        icon: Edit,
        label: "Draft",
        badge: "border-slate-200 text-slate-500 bg-slate-50 dark:border-white/10 dark:text-slate-400 dark:bg-white/5",
        iconWrap: "bg-slate-50 text-slate-400 ring-slate-200/60 dark:bg-white/5 dark:text-slate-500 dark:ring-white/10",
        accent: "bg-slate-300 dark:bg-slate-600",
        hover: "hover:border-slate-300/80 dark:hover:border-white/[0.10]",
    },
    scheduled: {
        icon: Calendar,
        label: "Terjadwal",
        badge: "border-blue-200 text-blue-600 bg-blue-50 dark:border-blue-500/20 dark:text-blue-400 dark:bg-blue-500/10",
        iconWrap: "bg-blue-50 text-blue-500 ring-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:ring-blue-500/20",
        accent: "bg-blue-400 dark:bg-blue-500",
        hover: "hover:border-blue-300/80 dark:hover:border-blue-500/20",
    },
    published: {
        icon: Send,
        label: "Terbit",
        badge: "border-green-200 text-green-600 bg-green-50 dark:border-green-500/20 dark:text-green-400 dark:bg-green-500/10",
        iconWrap: "bg-green-50 text-green-500 ring-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:ring-green-500/20",
        accent: "bg-green-400 dark:bg-green-500",
        hover: "hover:border-green-300/80 dark:hover:border-green-500/20",
    },
};

interface DraftItemProps {
    draft: Draft;
    onView: (draft: Draft) => void;
    onEdit: (draft: Draft) => void;
    onSchedule: (draft: Draft) => void;
    onPublishNow: (draft: Draft) => void;
    onDelete: (draft: Draft) => void;
    formatDate: (date: string) => string;
    checkSimilarity: (draft: Draft) => void;
    getSeoScore: (draft: Draft) => void;
}

export function DraftItem({
    draft, onView, onEdit, onSchedule, onPublishNow, onDelete, formatDate,
    checkSimilarity, getSeoScore,
}: DraftItemProps) {
    const cfg = statusConfig[draft.status as keyof typeof statusConfig] ?? statusConfig.draft;
    const StatusIcon = cfg.icon;

    const scheduledLabel = draft.scheduled_for
        ? draft.scheduled_for.startsWith("daily:")
            ? `Setiap hari ${draft.scheduled_for.replace("daily:", "")} `
            : draft.scheduled_for
        : null;

    const actions = [
        { icon: Eye, label: "Lihat", action: onView, permission: "draft:view:team" },
        { icon: Edit, label: "Edit", action: onEdit, permission: "draft:edit:team" },
        { icon: GitCompare, label: "Cek Kemiripan", action: checkSimilarity, permission: "draft:view:team" },
        { icon: BarChart2, label: "SEO Score", action: getSeoScore, permission: "draft:view:team" },
        ...(draft.status === "draft"
            ? [{ icon: Calendar, label: "Jadwalkan", action: onSchedule, permission: "draft:edit:team" }]
            : []),
    ];

    const {can} = useAuth()

    const visibleActions = actions.filter(a => can(a.permission));

    return (
        <div className={cn(
            "group relative rounded-xl border transition-all duration-200",
            "bg-white border-slate-200/80",
            "dark:bg-[#0f0d1a] dark:border-white/[0.06]",
            cfg.hover
        )}>
            {/* Left accent */}
            <div className={cn(
                "absolute left-0 top-3 bottom-3 w-[3px] rounded-r-full opacity-0 transition-opacity duration-200 group-hover:opacity-100",
                cfg.accent
            )} />

            {/* Main row */}
            <div className="flex items-start gap-3 p-3.5">
                {/* Icon */}
                <div className={cn(
                    "mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-xl ring-1",
                    cfg.iconWrap
                )}>
                    <StatusIcon className="h-4 w-4" />
                </div>

                {/* Content */}
                <div className="min-w-0 flex-1 overflow-hidden">
                    <div className="flex flex-wrap items-center gap-1.5">
                        <span className="truncate text-sm font-medium text-slate-900 dark:text-white">
                            {draft.title}
                        </span>
                        <span className={cn(
                            "inline-flex shrink-0 items-center gap-1 rounded-md border px-1.5 py-0.5 text-[10px] font-medium",
                            cfg.badge
                        )}>
                            <StatusIcon className="h-2.5 w-2.5" />
                            {cfg.label}
                            {draft.status === "scheduled" && scheduledLabel && (
                                <span className="opacity-70">• {scheduledLabel}</span>
                            )}
                        </span>
                    </div>

                    <p className="mt-1.5 line-clamp-1 text-xs leading-relaxed text-slate-400 dark:text-slate-500">
                        {draft.article.replace(/<[^>]*>/g, '').substring(0, 120)}...
                    </p>

                    <div className="mt-2 flex flex-wrap gap-3 text-[11px] text-slate-400 dark:text-slate-600">
                        <span className="flex items-center gap-1">
                            <Clock className="h-3 w-3 shrink-0" />
                            {formatDate(draft.created_at ?? "")}
                        </span>
                        <span className="flex items-center gap-1">
                            <Package className="h-3 w-3 shrink-0" />
                            {draft.target_products?.length} produk
                        </span>
                        {(draft.has_image || draft.image_url) && (
                            <span className="flex items-center gap-1">
                                <ImageIcon className="h-3 w-3 shrink-0" />
                                Ada gambar
                            </span>
                        )}
                    </div>
                </div>

                {/* Actions desktop */}
                <div className="hidden shrink-0 items-center gap-1 opacity-0 translate-x-1 transition-all duration-150 group-hover:opacity-100 group-hover:translate-x-0 sm:flex">
                    <Can permission="draft:view:team">
                        {visibleActions.map(({ icon: Icon, label, action }) => (
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
                    </Can>

                    <Can permission="draft:publish:team">
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
                    </Can>

                    <Can permission="draft:delete:team">
                        <Button
                            variant="ghost" size="icon" title="Hapus"
                            className="h-7 w-7 rounded-lg text-slate-300 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10 ml-0.5"
                            onClick={() => onDelete(draft)}
                        >
                            <Trash2 className="h-3.5 w-3.5" />
                        </Button>
                    </Can>
                </div>
            </div>

            {/* Actions mobile */}
            <div className="flex flex-wrap items-center gap-1 border-t border-slate-100 px-3.5 py-2 dark:border-white/[0.04] sm:hidden">
                <Can permission="draft:view:team">
                    {actions.map(({ icon: Icon, label, action }) => (
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
                </Can>

                <Can permission="draft:publish:team">
                    {(draft.status === "draft" || draft.status === "scheduled") && (
                        <Button
                            size="sm"
                            className={cn(
                                "h-7 gap-1.5 px-2.5 text-[11px] font-medium rounded-lg",
                                "bg-green-600 hover:bg-green-700 text-white shadow-sm",
                                "dark:bg-purple-600 dark:hover:bg-purple-700"
                            )}
                            onClick={() => onPublishNow(draft)}
                        >
                            <Send className="h-3 w-3" />
                            Terbitkan
                        </Button>
                    )}
                </Can>

                <Can permission="draft:delete:team">
                    <Button
                        variant="ghost" size="icon" title="Hapus"
                        className="h-7 w-7 rounded-lg text-slate-300 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10"
                        onClick={() => onDelete(draft)}
                    >
                        <Trash2 className="h-3.5 w-3.5" />
                    </Button>
                </Can>
            </div>
        </div>
    );
}