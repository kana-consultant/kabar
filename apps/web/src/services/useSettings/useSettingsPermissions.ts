import type { User } from "@/services/user";

export function useSettingsPermissions(currentUser: User | null) {
    const canManageUsers = currentUser && ['super_admin', 'admin', 'manager'].includes(currentUser.role);
    const canManageTeams = currentUser && ['super_admin', 'admin', 'manager'].includes(currentUser.role);
    const isAdmin = currentUser?.role === 'admin'

    return {
        canManageUsers,
        canManageTeams,
        isAdmin,
    };
}