import { Input } from "@/components/ui/input";
import { Search } from "lucide-react";
import { cn } from "@/lib/utils";

interface DraftsHeaderProps {
    searchQuery: string;
    setSearchQuery: (value: string) => void;
    statusFilter: "all" | "draft" | "scheduled" | "published";
    setStatusFilter: (value: "all" | "draft" | "scheduled" | "published") => void;
}

const filterOptions = [
    { value: "all", label: "Semua" },
    { value: "draft", label: "Draft" },
    { value: "scheduled", label: "Terjadwal" },
    { value: "published", label: "Terbit" },
] as const;

export function DraftsHeader({
    searchQuery, setSearchQuery,
    statusFilter, setStatusFilter,
}: DraftsHeaderProps) {
    return (
        <div className="flex items-start justify-between gap-4">
            <div>
                <h2 className="text-xl font-semibold tracking-tight text-slate-900 dark:text-white">
                    Draft & Terjadwal
                </h2>
                <p className="mt-0.5 text-sm text-slate-400 dark:text-slate-500">
                    Kelola draft, jadwalkan posting, atau publikasikan langsung
                </p>
            </div>

            <div className="flex items-center gap-2 shrink-0">
                {/* Search */}
                <div className="relative">
                    <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-slate-400 pointer-events-none" />
                    <Input
                        placeholder="Cari draft..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className={cn(
                            "h-8 pl-8 w-48 text-sm rounded-lg",
                            "border-slate-200/80 bg-white placeholder:text-slate-400",
                            "dark:border-white/[0.08] dark:bg-white/[0.03] dark:placeholder:text-slate-600",
                            "focus-visible:ring-1 focus-visible:ring-green-500/40 focus-visible:border-green-400/60",
                            "dark:focus-visible:ring-purple-500/40 dark:focus-visible:border-purple-400/40"
                        )}
                    />
                </div>

                {/* Status filter — pill group */}
                <div className={cn(
                    "flex items-center gap-0.5 rounded-lg border p-0.5",
                    "bg-slate-50 border-slate-200/80",
                    "dark:bg-white/[0.03] dark:border-white/[0.08]"
                )}>
                    {filterOptions.map(({ value, label }) => (
                        <button
                            key={value}
                            onClick={() => setStatusFilter(value)}
                            className={cn(
                                "h-7 rounded-md px-3 text-xs font-medium transition-all",
                                statusFilter === value
                                    ? "bg-white text-slate-800 shadow-sm border border-slate-200/80 dark:bg-white/[0.08] dark:text-white dark:border-white/[0.10]"
                                    : "text-slate-500 hover:text-slate-700 dark:text-slate-500 dark:hover:text-slate-300"
                            )}
                        >
                            {label}
                        </button>
                    ))}
                </div>
            </div>
        </div>
    );
}