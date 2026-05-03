import { toast } from "sonner";
import { publishDraft, deleteDraft, updateDraft } from "@/services/draft";
import { loadSchedulesData } from "./useScheduleData";

export function useScheduleActions(
    setSchedules: (data: any[]) => void,
    setLoading: (val: boolean) => void,
    setShowDeleteDialog: (val: boolean) => void,
    setSelectedSchedule: (val: any) => void,
    setShowRescheduleDialog: (val: boolean) => void,
    setNewScheduleDate: (val: string) => void,
    setNewScheduleTime: (val: string) => void
) {
    const handlePublishNow = async (schedule: any) => {
        try {
            await publishDraft(schedule.id);
            toast.success("Konten dipublikasikan!", {
                description: `"${schedule.title}" telah dipublikasikan`,
            });
            await loadSchedulesData(setSchedules, setLoading);
        } catch (error) {
            toast.error("Gagal mempublikasikan");
        }
    };

    const handleDelete = async (selectedSchedule: any) => {
        if (selectedSchedule) {
            try {
                await deleteDraft(selectedSchedule.id);
                toast.success("Jadwal dihapus", {
                    description: `"${selectedSchedule.title}" telah dihapus dari jadwal`,
                });
                await loadSchedulesData(setSchedules, setLoading);
                setShowDeleteDialog(false);
                setSelectedSchedule(null);
            } catch (error) {
                toast.error("Gagal menghapus jadwal");
            }
        }
    };

    const handleReschedule = async (
        selectedSchedule: any,
        newScheduleDate: string,
        newScheduleTime: string
    ) => {
        if (selectedSchedule && newScheduleDate && newScheduleTime) {
            const dateTime = `${newScheduleDate}T${newScheduleTime}`;
            try {
                await updateDraft(selectedSchedule.id, { scheduledFor: dateTime });
                toast.success("Jadwal diperbarui", {
                    description: `"${selectedSchedule.title}" dijadwalkan ulang pada ${newScheduleDate} jam ${newScheduleTime}`,
                });
                await loadSchedulesData(setSchedules, setLoading);
                setShowRescheduleDialog(false);
                setSelectedSchedule(null);
                setNewScheduleDate("");
                setNewScheduleTime("09:00");
            } catch (error) {
                toast.error("Gagal memperbarui jadwal");
            }
        }
    };

    return {
        handlePublishNow,
        handleDelete,
        handleReschedule,
    };
}