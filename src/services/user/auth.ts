import { apiClient, getAuthToken, setAuthCookie, removeAuthCookie, getUserFromCookie } from '../api';
import type { User } from '@/types/user';

// Get current logged in user from cookie
export async function getCurrentUser(): Promise<User | any> {
    return getUserFromCookie();
}

// Login
export async function login(email: string, password: string): Promise<{ user: User; token: string }> {
    const response = await apiClient.post<{ token: string; user: User }>('/auth/login', { email, password });
    
    if (response.token) {
        setAuthCookie(response.token, response.user);
    }
    
    return response;
}

// Logout
export async function logout(): Promise<void> {
    removeAuthCookie();
    try {
        await apiClient.post('/auth/logout');
    } catch (error) {
        // Ignore error
    }
}

// Check if authenticated
export function isAuthenticated(): boolean {
    return !!getAuthToken();
}