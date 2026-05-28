import type { UserRoleType } from "./types";
import { USER_ROLES } from "./types";

// Helper function untuk mendapatkan UserRole display name
export const getRoleDisplayName = (roleType: UserRoleType): string => {
    return USER_ROLES[roleType]?.displayName || roleType;
};

// Role options untuk dropdown
export const roleOptions: { value: UserRoleType; label: string }[] = [
    { value: 'super_admin', label: 'Super Admin' },
    { value: 'admin', label: 'Admin' },
    { value: 'member', label: 'Member' },
    { value: 'viewer', label: 'Viewer' },
    { value: 'owner', label: 'Owner' },
];

// Filter role options berdasarkan permission current user
export const getAvailableRoles = (currentUserRole: UserRoleType): UserRoleType[] => {
    switch (currentUserRole) {
        case 'super_admin':
            return ['admin', 'member', 'viewer', 'owner'];
        case 'owner':
        case 'admin':
            return ['member', 'viewer'];
        default:
            return [];
    }
};

// Atau jika ingin hierarchy berdasarkan level:
export const getAvailableRolesByLevel = (currentUserRole: UserRoleType): UserRoleType[] => {
    const roleLevel: Record<UserRoleType, number> = {
        super_admin: 100,
        owner: 100,
        admin: 80,
        member: 50,
        viewer: 20
    };

    const currentLevel = roleLevel[currentUserRole];

    return roleOptions
        .filter(option => roleLevel[option.value] <= currentLevel)
        .map(option => option.value);
};