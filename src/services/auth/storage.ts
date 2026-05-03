import Cookies from 'js-cookie';
import { COOKIE_OPTIONS } from './config';
import type { User } from '@/types/user';

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