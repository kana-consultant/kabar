import { useEffect } from "react";
import { useScheduleState } from "./useScheduleState";
import { loadSchedulesData } from "./useScheduleData";
import { useScheduleFilter } from "./useScheduleFilter";
import { useScheduleActions } from "./useScheduleActions";
import { formatDate, getScheduleDisplay, isDailySchedule } from "./useScheduleHelpers";

export function useSchedule() {
    const {
        schedules, setSchedules,
        filteredSchedules, setFilteredSchedules,
        searchQuery, setSearchQuery,
        loading, setLoading,
        selectedSchedule, setSelectedSchedule,
        showDetailDialog, setShowDetailDialog,
        showDeleteDialog, setShowDeleteDialog,
        showRescheduleDialog, setShowRescheduleDialog,
        newScheduleDate, setNewScheduleDate,
        newScheduleTime, setNewScheduleTime,
    } = useScheduleState();

    // Load schedules on mount
    useEffect(() => {
        loadSchedulesData(setSchedules, setLoading);
    }, []);

    // Filter logic
    useScheduleFilter(schedules, searchQuery, setFilteredSchedules);

    const { handlePublishNow, handleDelete, handleReschedule } = useScheduleActions(
        setSchedules, setLoading, setShowDeleteDialog, setSelectedSchedule,
        setShowRescheduleDialog, setNewScheduleDate, setNewScheduleTime
    );

    return {
        schedules,
        filteredSchedules,
        searchQuery,
        setSearchQuery,
        loading,
        selectedSchedule,
        setSelectedSchedule,
        showDetailDialog,
        setShowDetailDialog,
        showDeleteDialog,
        setShowDeleteDialog,
        showRescheduleDialog,
        setShowRescheduleDialog,
        newScheduleDate,
        setNewScheduleDate,
        newScheduleTime,
        setNewScheduleTime,
        loadSchedules: () => loadSchedulesData(setSchedules, setLoading),
        handlePublishNow,
        handleDelete: () => handleDelete(selectedSchedule),
        handleReschedule: () => handleReschedule(selectedSchedule, newScheduleDate, newScheduleTime),
        formatDate,
        getScheduleDisplay,
        isDailySchedule,
    };
}