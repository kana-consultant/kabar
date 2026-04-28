import { useState, useEffect } from "react";
import { toast } from "sonner";
import { 
    getHistory, 
    deleteHistory, 
    clearHistory, 
    addHistory 
} from "@/services/historyService";
import type { HistoryItem } from "@/types/history";

type StatusFilter = "all" | "success" | "failed" | "pending";
type ActionFilter = "all" | "published" | "scheduled" | "draft_saved";

export function useHistory() {
    const [history, setHistory] = useState<HistoryItem[]>([]);
    const [filteredHistory, setFilteredHistory] = useState<HistoryItem[]>([]);
    const [searchQuery, setSearchQuery] = useState<string>("");
    const [statusFilter, setStatusFilter] = useState<StatusFilter>("all");
    const [actionFilter, setActionFilter] = useState<ActionFilter>("all");
    const [loading, setLoading] = useState<boolean>(true);
    const [selectedHistory, setSelectedHistory] = useState<HistoryItem | null>(null);
    const [showDetailDialog, setShowDetailDialog] = useState<boolean>(false);
    const [showDeleteDialog, setShowDeleteDialog] = useState<boolean>(false);

    const loadHistory = async () => {
        setLoading(true);
        try {
            const data = await getHistory();
            setHistory(data || []);
        } catch (error) {
            console.error("Failed to load history:", error);
            toast.error("Gagal memuat riwayat");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadHistory();
    }, []);

    useEffect(() => {
        let filtered = [...history];
        
        if (statusFilter !== "all") {
            filtered = filtered.filter(h => h.status === statusFilter);
        }
        
        if (actionFilter !== "all") {
            filtered = filtered.filter(h => h.action === actionFilter);
        }
        
        if (searchQuery) {
            const query = searchQuery.toLowerCase();
            filtered = filtered.filter(h =>
                h.title.toLowerCase().includes(query) ||
                h.topic.toLowerCase().includes(query)
            );
        }
        
        setFilteredHistory(filtered);
    }, [history, searchQuery, statusFilter, actionFilter]);

    const handleDelete = async () => {
        if (selectedHistory) {
            try {
                await deleteHistory(selectedHistory.id);
                toast.success("Riwayat dihapus", {
                    description: `"${selectedHistory.title}" telah dihapus`,
                });
                await loadHistory();
                setShowDeleteDialog(false);
                setSelectedHistory(null);
            } catch (error) {
                toast.error("Gagal menghapus riwayat");
            }
        }
    };

    const handleClearAll = async () => {
        try {
            await clearHistory();
            toast.success("Semua riwayat dihapus");
            await loadHistory();
        } catch (error) {
            toast.error("Gagal menghapus riwayat");
        }
    };

    const addToHistory = async (data: {
        title: string;
        topic: string;
        content: string;
        imageUrl?: string;
        targetProducts: string[];
        status: 'success' | 'failed' | 'pending';
        action: 'published' | 'scheduled' | 'draft_saved';
        errorMessage?: string;
        scheduledFor?: string;
    }) => {
        try {
            await addHistory({
                ...data,
                publishedAt: new Date().toISOString(),
            });
            await loadHistory();
            return true;
        } catch (error) {
            console.error("Failed to add to history:", error);
            return false;
        }
    };

    const formatDate = (dateString: string): string => {
        return new Date(dateString).toLocaleDateString("id-ID", {
            day: "numeric",
            month: "long",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    // Hanya data mentah, bukan JSX!
    const getStatusData = (status: string): { label: string; icon: string; color: string } => {
        switch (status) {
            case "published":
                return { label: "Berhasil", icon: "✅", color: "green" };
            case "failed":
                return { label: "Gagal", icon: "❌", color: "red" };
            default:
                return { label: "Pending", icon: "⏳", color: "yellow" };
        }
    };

    const getActionData = (action: string): { label: string; icon: string } => {
        switch (action) {
            case "published":
                return { label: "Publikasi", icon: "🚀" };
            case "scheduled":
                return { label: "Terjadwal", icon: "📅" };
            default:
                return { label: "Draft", icon: "📝" };
        }
    };

    return {
        // Data states
        history,
        filteredHistory,
        searchQuery,
        setSearchQuery,
        statusFilter,
        setStatusFilter,
        actionFilter,
        setActionFilter,
        loading,
        selectedHistory,
        setSelectedHistory,
        showDetailDialog,
        setShowDetailDialog,
        showDeleteDialog,
        setShowDeleteDialog,
        // Functions (operasi)
        handleDelete,
        handleClearAll,
        addToHistory,
        formatDate,
        // Data helpers (bukan JSX!)
        getStatusData,
        getActionData,
        loadHistory,
    };
}