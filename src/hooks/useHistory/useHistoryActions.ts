import { toast } from "sonner";
import { deleteHistory, clearHistory, addHistory } from "@/services/history";
import { loadHistoryData } from "./useHistoryData";

export async function handleDeleteHistory(
    selectedHistory: any,
    setSelectedHistory: (val: any) => void,
    setShowDeleteDialog: (val: boolean) => void,
    setHistory: (data: any[]) => void,
    setLoading: (val: boolean) => void
) {
    if (selectedHistory) {
        try {
            await deleteHistory(selectedHistory.id);
            toast.success("Riwayat dihapus", {
                description: `"${selectedHistory.title}" telah dihapus`,
            });
            await loadHistoryData(setHistory, setLoading);
            setShowDeleteDialog(false);
            setSelectedHistory(null);
        } catch (error) {
            toast.error("Gagal menghapus riwayat");
        }
    }
}

export async function handleClearAllHistory(
    setHistory: (data: any[]) => void,
    setLoading: (val: boolean) => void
) {
    try {
        await clearHistory();
        toast.success("Semua riwayat dihapus");
        await loadHistoryData(setHistory, setLoading);
    } catch (error) {
        toast.error("Gagal menghapus riwayat");
    }
}

export async function addToHistory(
    data: {
        title: string;
        topic: string;
        content: string;
        imageUrl?: string;
        targetProducts: string[];
        status: 'success' | 'failed' | 'pending';
        action: 'published' | 'scheduled' | 'draft_saved';
        errorMessage?: string;
        scheduledFor?: string;
    },
    setHistory: (data: any[]) => void,
    setLoading: (val: boolean) => void
) {
    try {
        await addHistory({
            ...data,
            publishedAt: new Date().toISOString(),
        });
        await loadHistoryData(setHistory, setLoading);
        return true;
    } catch (error) {
        console.error("Failed to add to history:", error);
        return false;
    }
}