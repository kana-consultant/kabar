import { getUserRole } from './user';
import { getTeamId } from '../auth';

export function hasRole(requiredRole: string): boolean {
    const role = getUserRole();
    if (!role) return false;
    if (role === 'super_admin') return true;
    return role === requiredRole;
}

export function hasAnyRole(requiredRoles: string[]): boolean {
    const role = getUserRole();
    if (!role) return false;
    if (role === 'super_admin') return true;
    return requiredRoles.includes(role);
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