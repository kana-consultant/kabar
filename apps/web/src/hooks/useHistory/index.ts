import { useEffect } from "react";
import { useHistoryState } from "./useHistoryState";
import { loadHistoryData } from "./useHistoryData";
import { useHistoryFilter } from "./useHistoryFilter";
import { handleDeleteHistory, handleClearAllHistory, addToHistory } from "./useHistoryActions";
import { formatDate, getStatusData, getActionData } from "./useHistoryHelpers";

export function useHistory() {
    const {
        history, setHistory,
        filteredHistory, setFilteredHistory,
        searchQuery, setSearchQuery,
        statusFilter, setStatusFilter,
        actionFilter, setActionFilter,
        loading, setLoading,
        selectedHistory, setSelectedHistory,
        showDetailDialog, setShowDetailDialog,
        showDeleteDialog, setShowDeleteDialog,
    } = useHistoryState();

    // Load history on mount
    useEffect(() => {
        loadHistoryData(setHistory, setLoading);
    }, []);

    // Filter logic
    useHistoryFilter(history, searchQuery, statusFilter, actionFilter, setFilteredHistory);

    const handleDelete = () => handleDeleteHistory(
        selectedHistory, setSelectedHistory, setShowDeleteDialog, setHistory, setLoading
    );

    const handleClearAll = () => handleClearAllHistory(setHistory, setLoading);

    const addToHistoryWrapper = async (data: {
        title: string;
        topic: string;
        content: string;
        imageUrl?: string;
        targetProducts: string[];
        status: 'success' | 'failed' | 'pending';
        action: 'published' | 'scheduled' | 'draft_saved';
        errorMessage?: string;
        scheduledFor?: string;
    }) => {
        return addToHistory(data, setHistory, setLoading);
    };

    return {
        // Data states
        history,
        filteredHistory,
        searchQuery,
        setSearchQuery,
        statusFilter,
        setStatusFilter,
        actionFilter,
        setActionFilter,
        loading,
        selectedHistory,
        setSelectedHistory,
        showDetailDialog,
        setShowDetailDialog,
        showDeleteDialog,
        setShowDeleteDialog,
        // Functions (operasi)
        handleDelete,
        handleClearAll,
        addToHistory: addToHistoryWrapper,
        formatDate,
        // Data helpers (bukan JSX!)
        getStatusData,
        getActionData,
        loadHistory: () => loadHistoryData(setHistory, setLoading),
    };
}