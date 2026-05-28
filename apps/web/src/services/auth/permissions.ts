import { getUserFromCookie } from '../api';
import { getUserById } from '@/services/user';
import type { UserRole } from '@/services/user';
import { getUserTeams } from '@/services/user';

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
export function getTeamId(): string {
    const user = getUserFromCookie();
    return user?.team_id || 'Viewer';
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


export function hasPermission(requiredRoles: string[]): boolean {
    return hasAnyRole(requiredRoles);
}

export function getRoleLevel(): number {
    const role = getUserRole();
    const roleLevels: Record<string, number> = {
        'viewer': 1,
        'editor': 2,
        'manager': 3,
        'admin': 4,
        'super_admin': 5,
    };
    return role ? roleLevels[role] || 0 : 0;
}

export function canAccessTeam(resourceTeamId: string | null): boolean {
    const userTeamId = getTeamId();
    const role = getUserRole();
    
    if (role === 'super_admin' || role === 'admin') return true;
    return userTeamId === resourceTeamId;
}

// Check if token exists - re-export from storage
export { hasToken } from './storage';