import { getUserRole } from './user';

export function isAdmin(): boolean {
    const role = getUserRole();
    return role === 'admin';
}

export function isSuperAdmin(): boolean {
    const role = getUserRole();
    return role === 'super_admin';
}

export function isAnyAdmin(): boolean {
    const role = getUserRole();
    return role === 'admin' || role === 'super_admin';
}

export function isManagerOrAbove(): boolean {
    const role = getUserRole();
    return role === 'manager' || role === 'admin' || role === 'super_admin';
}

export function isEditorOrAbove(): boolean {
    const role = getUserRole();
    return role === 'editor' || role === 'manager' || role === 'admin' || role === 'super_admin';
}