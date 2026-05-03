import Cookies from 'js-cookie';
import { COOKIE_OPTIONS } from './config';
import type { User } from '@/types/user';

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