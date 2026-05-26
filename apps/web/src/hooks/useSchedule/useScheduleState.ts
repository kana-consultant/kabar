import { useState } from "react";
import type { Draft } from "@/services/draft";

export function useScheduleState() {
    const [schedules, setSchedules] = useState<Draft[]>([]);
    const [filteredSchedules, setFilteredSchedules] = useState<Draft[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [loading, setLoading] = useState(true);
    const [selectedSchedule, setSelectedSchedule] = useState<Draft | null>(null);
    const [showDetailDialog, setShowDetailDialog] = useState(false);
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);
    const [showRescheduleDialog, setShowRescheduleDialog] = useState(false);
    const [newScheduleDate, setNewScheduleDate] = useState("");
    const [newScheduleTime, setNewScheduleTime] = useState("09:00");
    const [currentPage, setCurrentPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [totalItems, setTotalItems] = useState(0);

    return {
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
    };
}