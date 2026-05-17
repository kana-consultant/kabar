import { Button } from "@/components/ui/button";
import { FileText, Plus } from "lucide-react";
import { useNavigate } from "@tanstack/react-router";
import { DraftItem } from "./DraftItem";
import type { Draft } from "@/services/draft";
import { cn } from "@/lib/utils";

interface DraftListProps {
    drafts: Draft[];
    onView: (draft: Draft) => void;
    onEdit: (draft: Draft) => void;
    onSchedule: (draft: Draft) => void;
    onPublishNow: (draft: Draft) => void;
    onDelete: (draft: Draft) => void;
    formatDate: (date: string) => string;
}

export function DraftList({
    drafts, onView, onEdit, onSchedule,
    onPublishNow, onDelete, formatDate,
}: DraftListProps) {
    const navigate = useNavigate();

    if (drafts.length === 0) {
        return (
            <div className={cn(
                "flex flex-col items-center justify-center rounded-2xl border py-20",
                "bg-white border-slate-200/80",
                "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
            )}>
                <div className="relative">
                    <div className="absolute inset-0 rounded-full blur-xl opacity-30 bg-slate-300 dark:bg-slate-600" />
                    <div className={cn(
                        "relative flex h-14 w-14 items-center justify-center rounded-2xl",
                        "bg-slate-100 text-slate-400 ring-1 ring-slate-200/60",
                        "dark:bg-white/[0.04] dark:text-slate-600 dark:ring-white/[0.06]"
                    )}>
                        <FileText className="h-6 w-6" />
                    </div>
                </div>
                <p className="mt-5 text-sm font-medium text-slate-600 dark:text-slate-400">
                    Belum ada draft
                </p>
                <p className="mt-1 text-xs text-slate-400 dark:text-slate-600 text-center max-w-xs">
                    Draft akan muncul setelah kamu generate konten
                </p>
                <Button
                    size="sm"
                    className={cn(
                        "mt-6 gap-2 rounded-lg text-xs font-medium shadow-sm",
                        "bg-green-600 hover:bg-green-700 text-white",
                        "dark:bg-purple-600 dark:hover:bg-purple-700"
                    )}
                    onClick={() => navigate({ to: "/generate" })}
                >
                    <Plus className="h-3.5 w-3.5" />
                    Buat Draft Baru
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
                        Daftar Draft
                    </p>
                    <span className={cn(
                        "inline-flex h-5 min-w-5 items-center justify-center rounded-md px-1.5 text-[10px] font-semibold tabular-nums",
                        "bg-green-100 text-green-700",
                        "dark:bg-purple-500/15 dark:text-purple-300"
                    )}>
                        {drafts.length}
                    </span>
                </div>

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
                    Draft Baru
                </Button>
            </div>

            {/* Items */}
            <div className="p-3 space-y-1.5">
                {drafts.map((draft) => (
                    <DraftItem
                        key={draft.id}
                        draft={draft}
                        onView={onView}
                        onEdit={onEdit}
                        onSchedule={onSchedule}
                        onPublishNow={onPublishNow}
                        onDelete={onDelete}
                        formatDate={formatDate}
                    />
                ))}
            </div>
        </div>
    );
}