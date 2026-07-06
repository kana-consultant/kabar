import { getHistory } from "@/services/history";
import type { ToastContextType } from "../use-toast";

interface PaginationState {
    currentPage: number;
    totalPages: number;
    totalItems: number;
    totalSuccess: number;
    totalFailed: number;
}

interface LoadHistoryParams {
    setHistory: (data: any[]) => void;
    setLoading: (val: boolean) => void;
    setPagination: (pagination: PaginationState) => void;
    toast: ToastContextType;
    page?: number;
    limit?: number;
}

export async function loadHistoryData({
    setHistory,
    setLoading,
    setPagination,
    toast,
    page = 1,
    limit = 5
}: LoadHistoryParams) {
    setLoading(true);
    try {
        const offset = (page - 1) * limit;
        const response = await getHistory({ limit, offset });

        setHistory(response.data || []);
        setPagination({
            currentPage: response.current_page,
            totalPages: response.total_pages,
            totalItems: response.total_items,
            totalSuccess: response.total_success ?? 0,
            totalFailed: response.total_failed ?? 0,
        });
    } catch (error) {
        console.error("Failed to load history:", error);
        toast.error("Gagal memuat riwayat");
    } finally {
        setLoading(false);
    }
}