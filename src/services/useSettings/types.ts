import type { User, Team, UserRole } from "@/types/user";

export type UserRoleType = 'super_admin' | 'admin' | 'manager' | 'editor' | 'viewer';

// Predefined UserRole objects untuk TeamMember.role
export const USER_ROLES: Record<UserRoleType, UserRole> = {
    super_admin: {
        id: "role_super_admin",
        name: "super_admin",
        displayName: "Super Admin",
        description: "Full access to everything",
        scope: "system",
        level: 100
    },
    admin: {
        id: "role_admin",
        name: "admin",
        displayName: "Admin",
        description: "Admin access",
        scope: "global",
        level: 80
    },
    manager: {
        id: "role_manager",
        name: "manager",
        displayName: "Manager",
        description: "Can manage team",
        scope: "team",
        level: 60
    },
    editor: {
        id: "role_editor",
        name: "editor",
        displayName: "Editor",
        description: "Can edit content",
        scope: "team",
        level: 40
    },
    viewer: {
        id: "role_viewer",
        name: "viewer",
        displayName: "Viewer",
        description: "View only",
        scope: "team",
        level: 20
    }
};