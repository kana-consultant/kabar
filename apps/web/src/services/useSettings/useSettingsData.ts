import { useEffect } from "react";
import { getUsers, getCurrentUser } from "@/services/user";
import type { User } from "@/services/user";
import type { ToastContextType } from "@/hooks/use-toast";

export function useSettingsData(
    setUsers: (data: User[]) => void,
    setCurrentUser: (user: User | null) => void,
    setLoading: (val: boolean) => void,
    toast : ToastContextType
) {
    const loadData = async () => {
        setLoading(true);
        try {
            const [usersData, currentUserData] = await Promise.all([
                getUsers(),
                getCurrentUser(),
            ]);
            setUsers(usersData || []);
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