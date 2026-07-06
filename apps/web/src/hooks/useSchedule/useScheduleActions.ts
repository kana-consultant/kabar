import { publishDraft, deleteDraft, rescheduleDraft } from "@/services/draft";
import { loadSchedulesData } from "./useScheduleData";
import type { ToastContextType } from "../use-toast";

export function useScheduleActions(
    setSchedules: (data: any[]) => void,
    setLoading: (val: boolean) => void,
    setShowDeleteDialog: (val: boolean) => void,
    setSelectedSchedule: (val: any) => void,
    setShowRescheduleDialog: (val: boolean) => void,
    setNewScheduleDate: (val: string) => void,
    setNewScheduleTime: (val: string) => void,
    toast: ToastContextType,
    // Tambahkan parameter yang dibutuhkan untuk pagination
    page: number,
    setTotalItems: (val: number) => void,
    setTotalPages: (val: number) => void,
) {
    const handlePublishNow = async (schedule: any) => {
        try {
            await publishDraft(schedule.id, null);
            toast.success("Konten dipublikasikan!", {
                description: `"${schedule.title}" telah dipublikasikan`,
            });
            await loadSchedulesData(
                setSchedules,
                setLoading,
                toast,
                { limit: 5, offset: (page - 1) * 5 },
                setTotalItems,
                setTotalPages,
            );
        } catch (error: any) {
            const errorMessage = error?.message || "Gagal mempublikasikan";
            toast.error(errorMessage);
        }
    };

    const handleDelete = async (selectedSchedule: any) => {
        if (!selectedSchedule) return;
        
        try {
            await deleteDraft(selectedSchedule.id);
            toast.success("Jadwal dihapus", {
                description: `"${selectedSchedule.title}" telah dihapus dari jadwal`,
            });
            await loadSchedulesData(
                setSchedules,
                setLoading,
                toast,
                { limit: 5, offset: (page - 1) * 5 },
                setTotalItems,
                setTotalPages,
            );
            setShowDeleteDialog(false);
            setSelectedSchedule(null);
        } catch (error: any) {
            const errorMessage = error?.message || "Gagal menghapus jadwal";
            toast.error(errorMessage);
        }
    };

    const handleReschedule = async (
        selectedSchedule: any,
        newScheduleDate: string,
        newScheduleTime: string
    ) => {
        if (!selectedSchedule || !newScheduleDate || !newScheduleTime) {
            toast.error("Mohon lengkapi tanggal dan waktu");
            return;
        }

        const dateTime = `${newScheduleDate}T${newScheduleTime}:00+07:00`; // Tambahkan detik dan timezone
        
        try {
            // ✅ Gunakan rescheduleDraft, bukan updateDraft
            await rescheduleDraft(selectedSchedule.id, { 
                scheduled_for: dateTime 
            });
            
            toast.success("Jadwal diperbarui", {
                description: `"${selectedSchedule.title}" dijadwalkan ulang pada ${newScheduleDate} jam ${newScheduleTime}`,
            });
            
            await loadSchedulesData(
                setSchedules,
                setLoading,
                toast,
                { limit: 5, offset: (page - 1) * 5 },
                setTotalItems,
                setTotalPages,
            );
            
            // Reset state
            setShowRescheduleDialog(false);
            setSelectedSchedule(null);
            setNewScheduleDate("");
            setNewScheduleTime("09:00");
        } catch (error: any) {
            const errorMessage = error?.message || "Gagal memperbarui jadwal";
            
            // Jika ada detail error dari backend
            if (error?.results) {
                console.error("Validation errors:", error.results);
            }
            
            toast.error(errorMessage);
        }
    };

    return {
        handlePublishNow,
        handleDelete,
        handleReschedule,
    };
}