import { type ToastContextType } from "../use-toast";
import { getScheduled } from "@/services/draft";
import type { PaginationParams } from "@/services/history/types";

export async function loadSchedulesData(
    setSchedules: (data: any[]) => void,
    setLoading: (val: boolean) => void,
    toast: ToastContextType,
    params: PaginationParams,
    setTotalItems: (val: number) => void,
    setTotalPages: (val: number) => void,
) {
    setLoading(true);
    try {
        const response = await getScheduled(params);
        setSchedules(response?.data ?? []);
        setTotalItems(response?.total_items ?? 0);
        setTotalPages(response?.total_pages ?? 0);
    } catch (error) {
        console.error("Failed to load schedules:", error);
        toast.error("Gagal memuat jadwal");
    } finally {
        setLoading(false);
    }
}