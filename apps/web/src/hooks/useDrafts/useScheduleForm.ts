import { useCallback } from "react";
import { toast } from "sonner";
import type { ScheduleConfig } from "./types";

interface UseScheduleFormParams {
    scheduleConfig: ScheduleConfig;
    setScheduleConfig: (updates: Partial<ScheduleConfig>) => void;
    onSchedule: () => Promise<void>;
    selectedDraftTitle?: string;
}

export function useScheduleForm({
    scheduleConfig,
    setScheduleConfig,
    onSchedule,
    selectedDraftTitle,
}: UseScheduleFormParams) {
    const { date, time, dailySchedule, dailyTime } = scheduleConfig;

    const validateAndSchedule = useCallback(async () => {
        if (!dailySchedule && !date) {
            toast.error("Pilih tanggal jadwal");
            return false;
        }
        
        await onSchedule();
        return true;
    }, [dailySchedule, date, onSchedule]);

    const toggleScheduleMode = useCallback((isDaily: boolean) => {
        setScheduleConfig({ dailySchedule: isDaily });
    }, [setScheduleConfig]);

    return {
        date,
        time,
        dailySchedule,
        dailyTime,
        setDate: (date: string) => setScheduleConfig({ date }),
        setTime: (time: string) => setScheduleConfig({ time }),
        setDailyTime: (dailyTime: string) => setScheduleConfig({ dailyTime }),
        toggleScheduleMode,
        validateAndSchedule,
        isFormValid: dailySchedule || !!date,
        selectedDraftTitle,
    };
}