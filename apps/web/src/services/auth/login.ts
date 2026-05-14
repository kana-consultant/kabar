import { apiClient } from '../api';
import Cookies from 'js-cookie';
import { COOKIE_OPTIONS } from './config';
import type { LoginResponse } from './types';

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