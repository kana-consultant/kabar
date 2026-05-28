import type { User } from "@/services/user";

export function useSettingsPermissions(currentUser: User | null) {
    const canManageUsers = currentUser && ['superadmin', 'admin', 'manager'].includes(currentUser.role);
    const canManageTeams = currentUser && ['superadmin', 'admin', 'manager'].includes(currentUser.role);
    const isSuperAdmin = currentUser?.role === 'superadmin';
    const isAdmin = currentUser?.role === 'admin' || currentUser?.role === 'superadmin';

    return {
        canManageUsers,
        canManageTeams,
        isSuperAdmin,
        isAdmin,
    };
}