import { useEffect } from "react";
import { useScheduleState } from "./useScheduleState";
import { loadSchedulesData } from "./useScheduleData";
import { useScheduleFilter } from "./useScheduleFilter";
import { useScheduleActions } from "./useScheduleActions";
import { formatDate, getScheduleDisplay, isDailySchedule } from "./useScheduleHelpers";
import { useToast } from "../use-toast";

export function useSchedule() {
    const toast = useToast();
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
        currentPage, setCurrentPage,
        totalPages, setTotalPages,
        totalItems, setTotalItems,
    } = useScheduleState();

    const loadSchedules = (page: number = currentPage) => {
        loadSchedulesData(
            setSchedules,
            setLoading,
            toast,
            { limit: 5, offset: (page - 1) * 5 },
            setTotalItems,
            setTotalPages,
        );
    };

    // Load on mount
    useEffect(() => {
        loadSchedules(1);
    }, []);

    // Reload when page changes
    useEffect(() => {
        loadSchedules(currentPage);
    }, [currentPage]);

    // Filter logic
    useScheduleFilter(schedules, searchQuery, setFilteredSchedules);

    const { handlePublishNow, handleDelete, handleReschedule } = useScheduleActions(
        setSchedules, setLoading, setShowDeleteDialog, setSelectedSchedule,
        setShowRescheduleDialog, setNewScheduleDate, setNewScheduleTime, toast
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
        loadSchedules: () => loadSchedules(1),
        handlePublishNow,
        handleDelete: () => handleDelete(selectedSchedule),
        handleReschedule: () => handleReschedule(selectedSchedule, newScheduleDate, newScheduleTime),
        formatDate,
        getScheduleDisplay,
        isDailySchedule,
        currentPage,
        setCurrentPage,
        totalPages,
        totalItems,
    };
}