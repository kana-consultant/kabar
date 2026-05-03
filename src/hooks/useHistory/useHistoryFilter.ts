import { useEffect } from "react";

export function useHistoryFilter(
    history: any[],
    searchQuery: string,
    statusFilter: string,
    actionFilter: string,
    setFilteredHistory: (data: any[]) => void
) {
    useEffect(() => {
        let filtered = [...history];
        
        if (statusFilter !== "all") {
            filtered = filtered.filter(h => h.status === statusFilter);
        }
        
        if (actionFilter !== "all") {
            filtered = filtered.filter(h => h.action === actionFilter);
        }
        
        if (searchQuery) {
            const query = searchQuery.toLowerCase();
            filtered = filtered.filter(h =>
                h.title.toLowerCase().includes(query) ||
                h.topic.toLowerCase().includes(query)
            );
        }
        
        setFilteredHistory(filtered);
    }, [history, searchQuery, statusFilter, actionFilter]);
}