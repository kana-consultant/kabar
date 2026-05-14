import { apiClient } from '../api';
import Cookies from 'js-cookie';

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