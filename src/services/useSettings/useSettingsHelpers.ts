import type { UserRoleType } from "./types";
import { USER_ROLES } from "./types";

// Helper function untuk mendapatkan UserRole display name
export const getRoleDisplayName = (roleType: UserRoleType): string => {
    return USER_ROLES[roleType]?.displayName || roleType;
};

// Role options untuk dropdown
export const roleOptions: { value: UserRoleType; label: string }[] = [
    { value: 'viewer', label: 'Viewer' },
    { value: 'editor', label: 'Editor' },
    { value: 'manager', label: 'Manager' },
    { value: 'admin', label: 'Admin' },
    { value: 'super_admin', label: 'Super Admin' },
];

// Filter role options berdasarkan permission current user
export const getAvailableRoles = (currentUser: any, isSuperAdmin: boolean, isAdmin: boolean): UserRoleType[] => {
    if (isSuperAdmin) {
        return ['super_admin', 'admin', 'manager', 'editor', 'viewer'];
    }
    if (isAdmin) {
        return ['admin', 'manager', 'editor', 'viewer'];
    }
    if (currentUser?.role === 'manager') {
        return ['manager', 'editor', 'viewer'];
    }
    return ['viewer'];
};