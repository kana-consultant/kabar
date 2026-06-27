import { useEffect } from "react";

export function useScheduleFilter(
    schedules: any,
    searchQuery: string,
    setFilteredSchedules: (data: any[]) => void
) {
    console.log(schedules)
    useEffect(() => {
        let filtered = [...schedules.data??[]];
        if (searchQuery) {
            const query = searchQuery.toLowerCase();
            filtered = filtered?.filter(s =>
                s.title.toLowerCase().includes(query) ||
                s.topic.toLowerCase().includes(query)
            );
        }
        setFilteredSchedules(filtered);
    }, [schedules, searchQuery]);
}