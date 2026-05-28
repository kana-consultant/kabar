import Cookies from 'js-cookie';
import { getUserFromCookie } from '../api';
import { getUserById } from './userQueries';
import { getUserTeams } from './teamQueries';

// ─── User ────────────────────────────────────────────────
export function getUserFromStorage() {
    return getUserFromCookie();
}

// ─── Role ────────────────────────────────────────────────
export function getUserRole(): string {
    const user = getUserFromCookie();
    return user?.role || 'viewer';
}

export function getUserRoleLevel(): number {
    const levelMap: Record<string, number> = {
        super_admin: 100,
        owner: 100,
        admin: 80,
        manager: 70,
        editor: 50,
        member: 30,
        viewer: 10,
    };
    return levelMap[getUserRole()] ?? 0;
}

// ─── Permissions ─────────────────────────────────────────
export function getPermissions(): string[] {
    // ✅ fix typo: permssions → permissions
    const user = getUserFromCookie();
    if (user?.permissions) return user.permissions;

    // fallback dari cookie terpisah
    try {
        const raw = Cookies.get('permissions');
        return raw ? JSON.parse(raw) : [];
    } catch {
        return [];
    }
}

export function hasPermission(permission: string): boolean {
    const role = getUserRole();
    if (role === 'super_admin') return true;
    return getPermissions().includes(permission);
}

// ─── Team ─────────────────────────────────────────────────
export function getTeamId(): string | null {
    const user = getUserFromCookie();
    // ✅ fix: sebelumnya return 'Viewer' kalau kosong
    return user?.team_id || Cookies.get('team_id') || null;
}

// ─── Role Checks ──────────────────────────────────────────
export function isSuperAdmin(): boolean {
    return getUserRole() === 'super_admin';
}

export function isAdmin(): boolean {
    return ['admin', 'super_admin'].includes(getUserRole());
}

export function isManagerOrAbove(): boolean {
    return getUserRoleLevel() >= 70;
}

export function isEditorOrAbove(): boolean {
    return getUserRoleLevel() >= 50;
}

export function isViewerOrAbove(): boolean {
    return getUserRoleLevel() >= 10;
}

export function hasRole(roleName: string): boolean {
    if (isSuperAdmin()) return true;
    return getUserRole() === roleName;
}

export function hasAnyRole(roleNames: string[]): boolean {
    if (isSuperAdmin()) return true;
    return roleNames.includes(getUserRole());
}

export function hasMinLevel(requiredLevel: number): boolean {
    return getUserRoleLevel() >= requiredLevel;
}

// ─── Async ────────────────────────────────────────────────
export async function canAccessResource(
    userId: string,
    resourceTeamId: string
): Promise<boolean> {
    const user = await getUserById(userId);
    if (!user) return false;

    const role = user.role || 'viewer';
    if (['super_admin', 'admin'].includes(role)) return true;

    const userTeams = await getUserTeams(userId);
    return userTeams.some(team => team.id === resourceTeamId);
}