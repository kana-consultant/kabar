import Cookies from 'js-cookie';
import { getUserFromCookie } from '../api';
import { getUserById } from './userQueries';
import type { UserRole } from '@/types/user';
import { getUserTeams } from './teamQueries';

// Get user role from cookie
export function getUserRole(): UserRole | any {
    const user = getUserFromCookie();
    return user?.role || null;
}

// Get role level
export function getUserRoleLevel(): number {
    const role = getUserRole();
    return role?.level || 0;
}

// Get role name
export function getUserRoleName(): string {
    const role = getUserRole();
    return role?.name || 'viewer';
}

// Get role display name
export function getUserRoleDisplayName(): string {
    const role = getUserRole();
    return role?.displayName || 'Viewer';
}

// Get team ID from cookie
export function getTeamId(): string | null {
    return Cookies.get('team_id') || null;
}

// Check if user is super admin
export function isSuperAdmin(): boolean {
    const roleName = getUserRoleName();
    return roleName === 'super_admin';
}

// Check if user is admin (includes super_admin)
export function isAdmin(): boolean {
    const roleName = getUserRoleName();
    return roleName === 'admin' || roleName === 'super_admin';
}

// Check if user is manager or above
export function isManagerOrAbove(): boolean {
    const roleLevel = getUserRoleLevel();
    return roleLevel >= 70; // manager level = 70
}

// Check if user is editor or above
export function isEditorOrAbove(): boolean {
    const roleLevel = getUserRoleLevel();
    return roleLevel >= 50; // editor level = 50
}

// Check if user is viewer or above (everyone)
export function isViewerOrAbove(): boolean {
    const roleLevel = getUserRoleLevel();
    return roleLevel >= 10; // viewer level = 10
}

// Check if user has specific role by name
export function hasRole(roleName: string): boolean {
    const userRoleName = getUserRoleName();
    if (userRoleName === 'super_admin') return true;
    return userRoleName === roleName;
}

// Check if user has any of the required role names
export function hasAnyRole(roleNames: string[]): boolean {
    const userRoleName = getUserRoleName();
    if (userRoleName === 'super_admin') return true;
    return roleNames.includes(userRoleName);
}

// Check if user has permission based on role level
export function hasMinLevel(requiredLevel: number): boolean {
    const userLevel = getUserRoleLevel();
    return userLevel >= requiredLevel;
}

// Check if user can access resource based on team
export async function canAccessResource(userId: string, resourceTeamId: string): Promise<boolean> {
    const user = await getUserById(userId);
    if (!user) return false;
    
    const roleName = user.role || "viewer";
    
    // Super admin and admin can access everything
    if (roleName === 'super_admin' || roleName === 'admin') return true;
    
    // Check if user is in the team
    const userTeams = await getUserTeams(userId);
    return userTeams.some(team => team.id === resourceTeamId);
}