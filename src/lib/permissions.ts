import type { UserRole } from "@/types/user";

export const permissions = {
    // Product permissions
    canCreateProduct: (role: UserRole) => ['admin', 'manager'].includes(role),
    canEditProduct: (role: UserRole) => ['admin', 'manager'].includes(role),
    canDeleteProduct: (role: UserRole) => ['admin'].includes(role),
    canViewProduct: (role: UserRole) => ['admin', 'manager', 'editor', 'viewer'].includes(role),

    // Content permissions
    canGenerateContent: (role: UserRole) => ['admin', 'manager', 'editor'].includes(role),
    canPublishContent: (role: UserRole) => ['admin', 'manager', 'editor'].includes(role),
    canScheduleContent: (role: UserRole) => ['admin', 'manager', 'editor'].includes(role),
    canEditContent: (role: UserRole) => ['admin', 'manager', 'editor'].includes(role),
    canDeleteContent: (role: UserRole) => ['admin', 'manager'].includes(role),

    // User permissions
    canManageUsers: (role: UserRole) => ['admin', 'manager'].includes(role),
    canManageRoles: (role: UserRole) => ['admin'].includes(role),

    // Team permissions
    canManageTeams: (role: UserRole) => ['admin', 'manager'].includes(role),
    canInviteMembers: (role: UserRole) => ['admin', 'manager'].includes(role),

    // Settings permissions
    canChangeSettings: (role: UserRole) => ['admin'].includes(role),
    canManageApiKeys: (role: UserRole) => ['admin', 'manager'].includes(role),
};