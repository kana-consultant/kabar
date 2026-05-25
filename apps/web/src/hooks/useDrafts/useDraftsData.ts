import { useState, useEffect, useMemo, useCallback } from "react";
import { useToast } from "@/hooks/use-toast"; //   Import useToast
import { getDrafts } from "@/services/draft";
import type { Draft } from "@/services/draft";
import type { StatusFilter } from "./types";

export function useDraftsData() {
    const toast = useToast(); //   Gunakan hook toast
    const [drafts, setDrafts] = useState<Draft[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState("");
    const [statusFilter, setStatusFilter] = useState<StatusFilter>("all");

    const loadDrafts = useCallback(async () => {
        setLoading(true);
        try {
            const allDrafts = await getDrafts();
            setDrafts(allDrafts ?? []);
        } catch (error) {
            console.error("Failed to load drafts:", error);
            toast.error("Gagal memuat draft"); //   Gunakan toast.error
        } finally {
            setLoading(false);
        }
    }, [toast]); //   Tambahkan toast ke dependencies

    // Filter drafts based on search and status
    const filteredDrafts = useMemo(() => {
        let filtered = [...drafts];
        
        if (statusFilter !== "all") {
            filtered = filtered.filter(d => d.status === statusFilter);
        }
        
        if (searchQuery) {
            const query = searchQuery.toLowerCase();
            filtered = filtered.filter(d =>
                d.title.toLowerCase().includes(query) ||
                d.topic.toLowerCase().includes(query)
            );
        }
        
        return filtered;
    }, [drafts, searchQuery, statusFilter]);

    // Initial load
    useEffect(() => {
        loadDrafts();
    }, [loadDrafts]);

    return {
        drafts,
        filteredDrafts,
        loading,
        searchQuery,
        setSearchQuery,
        statusFilter,
        setStatusFilter,
        loadDrafts,
    };
}