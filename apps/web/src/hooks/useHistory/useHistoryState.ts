import { useState } from "react";
import type { HistoryItem } from "@/services/history";
import type { StatusFilter, ActionFilter } from "./types";

export function useHistoryState() {
    const [history, setHistory] = useState<HistoryItem[]>([]);
    const [filteredHistory, setFilteredHistory] = useState<HistoryItem[]>([]);
    const [searchQuery, setSearchQuery] = useState<string>("");
    const [statusFilter, setStatusFilter] = useState<StatusFilter>("all");
    const [actionFilter, setActionFilter] = useState<ActionFilter>("all");
    const [loading, setLoading] = useState<boolean>(true);
    const [selectedHistory, setSelectedHistory] = useState<HistoryItem | null>(null);
    const [showDetailDialog, setShowDetailDialog] = useState<boolean>(false);
    const [showDeleteDialog, setShowDeleteDialog] = useState<boolean>(false);
    const [currentPage, setCurrentPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [totalItems, setTotalItems] = useState(0);
    const [totalSuccess, setTotalSuccess] = useState(0)
    const [totalFailed, setTotalFailed] = useState(0)

    return {
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
    };
}