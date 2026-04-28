import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Search, Trash2 } from "lucide-react";

interface HistoryHeaderProps {
    searchQuery: string;
    setSearchQuery: (value: string) => void;
    statusFilter: string;
    setStatusFilter: (value: "all" | "success" | "failed" | "pending") => void;
    actionFilter: string;
    setActionFilter: (value: "all" | "published" | "scheduled" | "draft_saved") => void;
    onClearAll: () => void;
}

export default function HistoryHeader({
    searchQuery,
    setSearchQuery,
    statusFilter,
    setStatusFilter,
    actionFilter,
    setActionFilter,
    onClearAll,
}: HistoryHeaderProps) {
    return (
        <div className="space-y-4">
            <div>
                <h2 className="text-2xl font-bold tracking-tight">History</h2>
                <p className="text-slate-500">Riwayat generate dan publish konten</p>
            </div>

            <div className="flex flex-wrap gap-4">
                <div className="relative flex-1 min-w-[200px]">
                    <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
                    <Input
                        placeholder="Cari riwayat..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-9"
                    />
                </div>

                <select
                    value={statusFilter}
                    onChange={(e) => setStatusFilter(e.target.value as any)}
                    className="rounded-md border px-3 py-2 text-sm"
                >
                    <option value="all">Semua Status</option>
                    <option value="success">Berhasil</option>
                    <option value="failed">Gagal</option>
                    <option value="pending">Pending</option>
                </select>

                <select
                    value={actionFilter}
                    onChange={(e) => setActionFilter(e.target.value as any)}
                    className="rounded-md border px-3 py-2 text-sm"
                >
                    <option value="all">Semua Aksi</option>
                    <option value="published">Publikasi</option>
                    <option value="scheduled">Terjadwal</option>
                    <option value="draft_saved">Draft</option>
                </select>

                <Button variant="destructive" onClick={onClearAll}>
                    <Trash2 className="mr-2 h-4 w-4" />
                    Hapus Semua
                </Button>
            </div>
        </div>
    );
}