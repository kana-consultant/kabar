import { HistoryItem } from "./HistoryItem";
import { History, ChevronLeft, ChevronRight } from "lucide-react";
import type { HistoryItem as HistoryItemType } from "@/services/history";
import { cn } from "@/lib/utils";

interface HistoryListProps {
    items: HistoryItemType[];
    onView: (item: HistoryItemType) => void;
    onRepost: (item: HistoryItemType) => void;
    onDelete: (item: HistoryItemType) => void;
    formatDate: (date: string) => string;
    getStatusData: (status: string) => { label: string; icon: string; color: string };
    getActionData: (action: string) => { label: string; icon: string };
    currentPage: number;
    totalPages: number;
    totalItems: number;
    onPageChange: (page: number) => void;
}

export function HistoryList({
    items, onView, onRepost, onDelete,
    formatDate, getStatusData, getActionData,
    currentPage, totalPages, totalItems, onPageChange,
}: HistoryListProps) {

    if (items.length === 0) {
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
                        <History className="h-6 w-6" />
                    </div>
                </div>
                <p className="mt-5 text-sm font-medium text-slate-600 dark:text-slate-400">
                    Belum ada riwayat
                </p>
                <p className="mt-1 text-xs text-slate-400 dark:text-slate-600">
                    Riwayat akan muncul setelah konten dipublikasikan
                </p>
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
                        Riwayat Konten
                    </p>
                    <span className={cn(
                        "inline-flex h-5 min-w-5 items-center justify-center rounded-md px-1.5 text-[10px] font-semibold tabular-nums",
                        "bg-green-100 text-green-700",
                        "dark:bg-purple-500/15 dark:text-purple-300"
                    )}>
                        {totalItems}
                    </span>
                </div>
            </div>

            {/* Items */}
            <div className="p-3 space-y-1.5">
                {items.map((item) => (
                    <HistoryItem
                        key={item.id}
                        item={item}
                        onView={onView}
                        onRepost={onRepost}
                        onDelete={onDelete}
                        formatDate={formatDate}
                        getStatusData={getStatusData}
                        getActionData={getActionData}
                    />
                ))}
            </div>

            {/* Pagination */}
            {totalPages > 0 && (
                <div className={cn(
                    "flex items-center justify-between px-5 py-3 border-t",
                    "border-slate-100 bg-slate-50/60",
                    "dark:border-white/[0.05] dark:bg-white/[0.02]"
                )}>
                    <p className="text-xs text-slate-400 dark:text-slate-600">
                        Halaman {currentPage} dari {totalPages}
                    </p>

                    <div className="flex items-center gap-1">
                        <button
                            disabled={currentPage <= 1}
                            onClick={() => onPageChange(currentPage - 1)}
                            className={cn(
                                "flex h-7 w-7 items-center justify-center rounded-lg transition-colors",
                                "text-slate-500 hover:bg-slate-100 hover:text-slate-700",
                                "dark:text-slate-400 dark:hover:bg-white/[0.06] dark:hover:text-slate-200",
                                "disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:bg-transparent"
                            )}
                        >
                            <ChevronLeft className="h-4 w-4" />
                        </button>

                        <button
                            disabled={currentPage >= totalPages}
                            onClick={() => onPageChange(currentPage + 1)}
                            className={cn(
                                "flex h-7 w-7 items-center justify-center rounded-lg transition-colors",
                                "text-slate-500 hover:bg-slate-100 hover:text-slate-700",
                                "dark:text-slate-400 dark:hover:bg-white/[0.06] dark:hover:text-slate-200",
                                "disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:bg-transparent"
                            )}
                        >
                            <ChevronRight className="h-4 w-4" />
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}