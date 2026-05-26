import { useEffect } from "react";
import { getUsers, getTeams, getCurrentUser } from "@/services/user";
import type { User, Team } from "@/services/user";
import type { ToastContextType } from "@/hooks/use-toast";

export function useSettingsData(
    setUsers: (data: User[]) => void,
    setTeams: (data: Team[]) => void,
    setCurrentUser: (user: User | null) => void,
    setLoading: (val: boolean) => void,
    toast : ToastContextType
) {
    const loadData = async () => {
        setLoading(true);
        try {
            const [usersData, teamsData, currentUserData] = await Promise.all([
                getUsers(),
                getTeams(),
                getCurrentUser(),
            ]);
            setUsers(usersData || []);
            setTeams(teamsData || []);
            setCurrentUser(currentUserData);
        } catch (error) {
            console.error("Failed to load settings data:", error);
            toast.error("Gagal memuat data pengaturan");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadData();
    }, []);

    return { loadData };
}   