
import { deleteHistory, clearHistory, addHistory } from "@/services/history";
import { loadHistoryData } from "./useHistoryData";
import type { ToastContextType } from "../use-toast";

export async function handleDeleteHistory(
    selectedHistory: any,
    setSelectedHistory: (val: any) => void,
    setShowDeleteDialog: (val: boolean) => void,
    setHistory: (data: any[]) => void,
    setLoading: (val: boolean) => void,
    toast : ToastContextType
) {
    if (selectedHistory) {
        try {
            await deleteHistory(selectedHistory.id);
            toast.success("Riwayat dihapus", {
                description: `"${selectedHistory.title}" telah dihapus`,
            });
            await loadHistoryData(setHistory, setLoading,toast);
            setShowDeleteDialog(false);
            setSelectedHistory(null);
        } catch (error) {
            toast.error("Gagal menghapus riwayat");
        }
    }
}

export async function handleClearAllHistory(
    setHistory: (data: any[]) => void,
    setLoading: (val: boolean) => void,
    toast : ToastContextType
) {
    try {
        await clearHistory();
        toast.success("Semua riwayat dihapus");
        await loadHistoryData(setHistory, setLoading, toast);
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
        keywords : string[]
    },
    setHistory: (data: any[]) => void,
    setLoading: (val: boolean) => void,
    toast : ToastContextType
) {
    try {
        await addHistory({
            ...data,
            publishedAt: new Date().toISOString(),
        });
        await loadHistoryData(setHistory, setLoading, toast);
        return true;
    } catch (error) {
        console.error("Failed to add to history:", error);
        return false;
    }
}