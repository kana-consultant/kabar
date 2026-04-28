import { Input } from "@/components/ui/input";
import { Search } from "lucide-react";

interface DraftsHeaderProps {
    searchQuery: string;
    setSearchQuery: (value: string) => void;
    statusFilter: "all" | "draft" | "scheduled" | "published";
    setStatusFilter: (value: "all" | "draft" | "scheduled" | "published") => void;
}

export function DraftsHeader({
    searchQuery,
    setSearchQuery,
    statusFilter,
    setStatusFilter,
}: DraftsHeaderProps) {
    return (
        <div className="flex items-center justify-between">
            <div>
                <h2 className="text-2xl font-bold tracking-tight">Draft & Terjadwal</h2>
                <p className="text-slate-500">
                    Kelola draft, jadwalkan posting, atau publikasikan langsung
                </p>
            </div>
            <div className="flex gap-2">
                <div className="relative">
                    <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
                    <Input
                        placeholder="Cari draft..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-8 w-64"
                    />
                </div>
                <select
                    value={statusFilter}
                    onChange={(e) => setStatusFilter(e.target.value as any)}
                    className="rounded-md border px-3 py-2 text-sm"
                >
                    <option value="all">Semua</option>
                    <option value="draft">Draft</option>
                    <option value="scheduled">Terjadwal</option>
                    <option value="published">Terbit</option>
                </select>
            </div>
        </div>
    );
}