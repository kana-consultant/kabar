import { useEffect } from "react";
import { useHistoryState } from "./useHistoryState";
import { loadHistoryData } from "./useHistoryData";
import { useHistoryFilter } from "./useHistoryFilter";
import { handleDeleteHistory, handleClearAllHistory, addToHistory } from "./useHistoryActions";
import { formatDate, getStatusData, getActionData } from "./useHistoryHelpers";
import { useToast } from "../use-toast";

export function useHistory() {
    const toast = useToast();
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
        currentPage, setCurrentPage,
        totalPages, setTotalPages,
        totalItems, setTotalItems,
        totalSuccess, setTotalSuccess,
        totalFailed, setTotalFailed,
    } = useHistoryState();

    const setPagination = ({
        currentPage,
        totalPages,
        totalItems,
        totalSuccess,
        totalFailed,
    }: {
        currentPage: number;
        totalPages: number;
        totalItems: number;
        totalSuccess: number;
        totalFailed: number;
    }) => {
        setCurrentPage(currentPage);
        setTotalPages(totalPages);
        setTotalItems(totalItems);
        setTotalSuccess(totalSuccess);
        setTotalFailed(totalFailed);
    };

    const load = (page: number = 1) => {
        loadHistoryData({ setHistory, setLoading, setPagination, toast, page });
    };

    // Load history on mount
    useEffect(() => {
        load(1);
    }, []);

    // Filter logic
    useHistoryFilter(history, searchQuery, statusFilter, actionFilter, setFilteredHistory);

    const handlePageChange = (page: number) => {
        setCurrentPage(page);
        load(page);
    };

    const handleDelete = () => handleDeleteHistory(
        selectedHistory, setSelectedHistory, setShowDeleteDialog, setHistory, setLoading, setPagination, toast, 1
    );

    const handleClearAll = () => handleClearAllHistory(setHistory, setLoading, setPagination, toast, 1);

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
        keywords: string[];
        excerpt: string;
        slug: string;
    }) => {
        return addToHistory(data, setHistory, setLoading, setPagination, toast, 1);
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
        // Pagination
        currentPage,
        totalPages,
        totalItems,
        handlePageChange,
        // Functions
        handleDelete,
        handleClearAll,
        addToHistory: addToHistoryWrapper,
        formatDate,
        // Helpers
        getStatusData,
        getActionData,
        loadHistory: load,
        totalSuccess, setTotalSuccess,
        totalFailed, setTotalFailed,
    };
}