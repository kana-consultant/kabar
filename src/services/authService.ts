// src/services/authService.ts
import { apiClient } from './api';
import type { User } from '@/types/user';
import Cookies from 'js-cookie';

export interface LoginResponse {
    token: string;
    user: User;
    teamId?: string;
}

// Cookie configuration
const COOKIE_OPTIONS = {
    expires: 7, // 7 days
    secure: import.meta.env.PROD, // true in production
    sameSite: 'strict' as const,
    path: '/'
};

// ==================== LOGIN / LOGOUT ====================

// Login - simpan ke localStorage dan cookie
export async function login(email: string, password: string): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>('/auth/login', { email, password });

    console.log(response.token)
    
    // Save to localStorage (for backward compatibility)
    localStorage.setItem('auth_token', response.token);
    localStorage.setItem('user', JSON.stringify(response.user));
    if (response.teamId) {
        localStorage.setItem('team_id', response.teamId);
    }
    
    // Save to cookie (for API interceptor)
    Cookies.set('auth_token', response.token, COOKIE_OPTIONS);
    Cookies.set('user', JSON.stringify(response.user), COOKIE_OPTIONS);
    if (response.teamId) {
        Cookies.set('team_id', response.teamId, COOKIE_OPTIONS);
    }
    
    return response;
}

// Logout - hapus dari localStorage dan cookie
export async function logout(): Promise<void> {
    // Remove from localStorage
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user');
    localStorage.removeItem('team_id');
    
    // Remove from cookie
    Cookies.remove('auth_token');
    Cookies.remove('user');
    Cookies.remove('team_id');
    
    // Optional: call logout endpoint
    try {
        await apiClient.post('/auth/logout');
    } catch (error) {
        // Ignore error
    }
}

// ==================== GET USER DATA ====================

// Get current user - cek dari cookie dulu, lalu localStorage
export async function getCurrentUser(): Promise<User | null> {
    // Try to get from cookie first
    const userFromCookie = Cookies.get('user');
    if (userFromCookie) {
        try {
            return JSON.parse(userFromCookie);
        } catch {
            // Fall through to localStorage
        }
    }
    
    // Fallback to localStorage
    const userStr = localStorage.getItem('user');
    if (userStr) {
        try {
            return JSON.parse(userStr);
        } catch {
            return null;
        }
    }
    return null;
}

// Get user role - dari cookie atau localStorage
export function getUserRole(): string | null {
    // Try cookie first
    const userFromCookie = Cookies.get('user');
    if (userFromCookie) {
        try {
            const user = JSON.parse(userFromCookie);
            return user.role;
        } catch {
            // Fall through
        }
    }
    
    // Fallback to localStorage
    const userStr = localStorage.getItem('user');
    if (userStr) {
        try {
            const user = JSON.parse(userStr);
            return user.role;
        } catch {
            return null;
        }
    }
    return null;
}

// Get team ID - dari cookie atau localStorage
export function getTeamId(): string | null {
    // Try cookie first
    const teamIdFromCookie = Cookies.get('team_id');
    if (teamIdFromCookie) {
        return teamIdFromCookie;
    }
    
    // Fallback to localStorage
    return localStorage.getItem('team_id');
}

// Get auth token - dari cookie atau localStorage
export function getToken(): string | null {
    // Try cookie first
    const tokenFromCookie = Cookies.get('auth_token');
    if (tokenFromCookie) {
        return tokenFromCookie;
    }
    
    // Fallback to localStorage
    return localStorage.getItem('auth_token');
}

// Check if token exists
export function hasToken(): boolean {
    return !!getToken();
}

// Clear all auth data
export function clearAuthData(): void {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user');
    localStorage.removeItem('team_id');
    
    Cookies.remove('auth_token');
    Cookies.remove('user');
    Cookies.remove('team_id');
}

// ==================== ADMIN CHECKS ====================

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

// ==================== PERMISSION CHECKS ====================

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

// ==================== USER DATA ====================

export function getUserName(): string {
    const user = getCurrentUserSync();
    return user?.name || '';
}

export function getUserEmail(): string {
    const user = getCurrentUserSync();
    return user?.email || '';
}

export function getUserAvatar(): string | undefined {
    const user = getCurrentUserSync();
    return user?.avatar;
}

// Helper sync version
function getCurrentUserSync(): User | null {
    const userFromCookie = Cookies.get('user');
    if (userFromCookie) {
        try {
            return JSON.parse(userFromCookie);
        } catch {
            // Fall through
        }
    }
    
    const userStr = localStorage.getItem('user');
    if (userStr) {
        try {
            return JSON.parse(userStr);
        } catch {
            return null;
        }
    }
    return null;
}

export function updateLocalUser(updates: Partial<User>): void {
    const user = getCurrentUserSync();
    if (user) {
        const updatedUser = { ...user, ...updates };
        
        // Update localStorage
        localStorage.setItem('user', JSON.stringify(updatedUser));
        
        // Update cookie
        Cookies.set('user', JSON.stringify(updatedUser), COOKIE_OPTIONS);
    }
}

// ==================== TOKEN MANAGEMENT ====================

export function setAuthCookie(token: string, user: User, teamId?: string): void {
    Cookies.set('auth_token', token, COOKIE_OPTIONS);
    Cookies.set('user', JSON.stringify(user), COOKIE_OPTIONS);
    if (teamId) {
        Cookies.set('team_id', teamId, COOKIE_OPTIONS);
    }
}

export function removeAuthCookie(): void {
    Cookies.remove('auth_token');
    Cookies.remove('user');
    Cookies.remove('team_id');
}