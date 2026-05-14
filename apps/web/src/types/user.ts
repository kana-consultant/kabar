export interface User {
    id: string;
    email: string;
    name: string;
    role: 'super_admin' | 'admin' | 'manager' | 'editor' | 'viewer';
    avatar?: string;
    status: 'active' | 'inactive' | 'suspended';
    emailVerified: boolean;
    twoFactorEnabled: boolean;
    lastActive?: string;
    lastLogin?: string;
    createdAt: string;
    updatedAt: string;
}

export interface UserRole {
    id: string;
    name: string;
    displayName: string;
    description?: string;
    scope: 'global' | 'team' | 'system';
    level: number;
}

export interface TeamMember {
    userId: string;
    userEmail: string;
    userName: string;
    role: UserRole;
    joinedAt: string;
}

export interface Team {
    id: string;
    name: string;
    description: string;
    members: TeamMember[];
    createdAt: string;
    updatedAt: string;
}

export interface Permission {
    id: string;
    name: string;
    description: string;
    roles: UserRole[];
}