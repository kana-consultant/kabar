import type { UserRole } from "@/services/user";

export type UserRoleType = 'superadmin' | 'admin' | 'member' | 'viewer' | 'owner';

// Predefined UserRole objects untuk TeamMember.role
export const USER_ROLES: Record<UserRoleType, UserRole> = {
    superadmin: {
        id: "role_admin",
        name: "Super Admin",
        displayName: "Super dmin",
        description: "Super access",
        scope: "global",
        level: 100
    },
    admin: {
        id: "role_admin",
        name: "admin",
        displayName: "Admin",
        description: "Admin access",
        scope: "team",
        level: 80
    },
    member: {
        id: "role_member",
        name: "member",
        displayName: "Member",
        description: "Regular team member",
        scope: "team",
        level: 50
    },
    viewer: {
        id: "role_viewer",
        name: "viewer",
        displayName: "Viewer",
        description: "View only",
        scope: "team",
        level: 20
    },
    owner: {
        id: "role_owner",
        name: "owner",
        displayName: "Owner",
        description: "Team owner with full control",
        scope: "team",
        level: 100
    }
};