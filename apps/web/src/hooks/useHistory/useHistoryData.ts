import { toast } from "sonner";
import { getHistory } from "@/services/history";

export async function loadHistoryData(
    setHistory: (data: any[]) => void,
    setLoading: (val: boolean) => void
) {
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
}