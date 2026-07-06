import { deleteHistory, clearHistory, addHistory } from "@/services/history";
import { loadHistoryData } from "./useHistoryData";
import type { ToastContextType } from "../use-toast";

interface PaginationState {
    currentPage: number;
    totalPages: number;
    totalItems: number;
    totalSuccess: number;
    totalFailed: number;
}

interface HistoryData {
    title: string;
    topic: string;
    content: string;
    imageUrl?: string;
    targetProducts: string[];
    status: 'success' | 'failed' | 'pending';
    action: 'published' | 'scheduled' | 'draft_saved';
    errorMessage?: string;
    scheduledFor?: string;
    keywords: string[];
    excerpt : string,
    slug : string;
}

export async function handleDeleteHistory(
    selectedHistory: any,
    setSelectedHistory: (val: any) => void,
    setShowDeleteDialog: (val: boolean) => void,
    setHistory: (data: any[]) => void,
    setLoading: (val: boolean) => void,
    setPagination: (pagination: PaginationState) => void,
    toast: ToastContextType,
    page: number = 1,
    limit: number = 5,
) {
    if (!selectedHistory) {
        toast.error("Tidak ada riwayat yang dipilih");
        return;
    }

    try {
        await deleteHistory(selectedHistory.id);

        toast.success("Riwayat dihapus", {
            description: `"${selectedHistory.title}" telah dihapus dari riwayat`,
        });

        // Reload data dengan pagination yang sama
        await loadHistoryData({
            setHistory,
            setLoading,
            setPagination,
            toast,
            page,
            limit
        });

        // Reset state
        setShowDeleteDialog(false);
        setSelectedHistory(null);

    } catch (error: any) {
        console.error("Failed to delete history:", error);

        const errorMessage = error?.response?.data?.message ||
            error?.message ||
            "Gagal menghapus riwayat";

        toast.error(errorMessage);
    }
}

export async function handleClearAllHistory(
    setHistory: (data: any[]) => void,
    setLoading: (val: boolean) => void,
    setPagination: (pagination: PaginationState) => void,
    toast: ToastContextType,
    page: number = 1,
    limit: number = 5,
) {
    try {
        await clearHistory();

        toast.success("Semua riwayat dihapus", {
            description: "Seluruh riwayat publikasi telah dibersihkan",
        });

        // Reload data setelah clear all
        await loadHistoryData({
            setHistory,
            setLoading,
            setPagination,
            toast,
            page,
            limit
        });

    } catch (error: any) {
        console.error("Failed to clear all history:", error);

        const errorMessage = error?.response?.data?.message ||
            error?.message ||
            "Gagal menghapus semua riwayat";

        toast.error(errorMessage);
    }
}

export async function addToHistory(
    data: HistoryData,
    setHistory: (data: any[]) => void,
    setLoading: (val: boolean) => void,
    setPagination: (pagination: PaginationState) => void,
    toast: ToastContextType,
    page: number = 1,
    limit: number = 5,
): Promise<boolean> {
    try {
        await addHistory(
            {
                ...data,
                publishedAt : new Date().toISOString() as string,
            }
        );

        // Reload data setelah menambah history
        await loadHistoryData({
            setHistory,
            setLoading,
            setPagination,
            toast,
            page,
            limit
        });

        return true;

    } catch (error: any) {
        console.error("Failed to add to history:", error);

        // Tidak perlu toast error di sini karena ini operasi background
        // Biarkan komponen yang memanggil yang handle error UI

        return false;
    }
}

