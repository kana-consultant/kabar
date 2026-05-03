import { toast } from "sonner";
import { getDrafts } from "@/services/draft";

export async function loadSchedulesData(
    setSchedules: (data: any[]) => void,
    setLoading: (val: boolean) => void
) {
    setLoading(true);
    try {
        const allDrafts = await getDrafts();
        const scheduled = allDrafts ? allDrafts.filter(d => d.status === "scheduled") : [];
        setSchedules(scheduled || []);
    } catch (error) {
        console.error("Failed to load schedules:", error);
        toast.error("Gagal memuat jadwal");
    } finally {
        setLoading(false);
    }
}