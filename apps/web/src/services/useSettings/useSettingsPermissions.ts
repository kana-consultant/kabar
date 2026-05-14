import type { User } from "@/services/user";

export function useSettingsPermissions(currentUser: User | null) {
    const canManageUsers = currentUser && ['super_admin', 'admin', 'manager'].includes(currentUser.role);
    const canManageTeams = currentUser && ['super_admin', 'admin', 'manager'].includes(currentUser.role);
    const isSuperAdmin = currentUser?.role === 'super_admin';
    const isAdmin = currentUser?.role === 'admin' || currentUser?.role === 'super_admin';

    return {
        canManageUsers,
        canManageTeams,
        isSuperAdmin,
        isAdmin,
    };
}