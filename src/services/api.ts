import axios, {type AxiosInstance,type InternalAxiosRequestConfig,type AxiosResponse } from 'axios';
import Cookies from 'js-cookie';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

class ApiClient {
    private client: AxiosInstance;

    constructor() {
        this.client = axios.create({
            baseURL: API_BASE_URL,
            headers: {
                'Content-Type': 'application/json',
            },
            withCredentials: true,
            timeout: 30000,
        });

        // Request interceptor
        this.client.interceptors.request.use(
            (config: InternalAxiosRequestConfig) => {
                const token = Cookies.get('auth_token');
                if (token) {
                    config.headers.Authorization = `Bearer ${token}`;
                }
                return config;
            },
            (error) => Promise.reject(error)
        );

        // Response interceptor
        this.client.interceptors.response.use(
            (response: AxiosResponse) => response,
            (error) => {
                if (error.response?.status === 401) {
                    removeAuthCookie();
                    window.location.href = '/login';
                }
                return Promise.reject(error);
            }
        );
    }

    async get<T>(url: string, params?: any): Promise<T> {
        const response = await this.client.get<T>(url, { params });
        return response.data;
    }

    async post<T>(url: string, data?: any): Promise<T> {
        const response = await this.client.post<T>(url, data);
        return response.data;
    }

    async put<T>(url: string, data?: any): Promise<T> {
        const response = await this.client.put<T>(url, data);
        return response.data;
    }

    async patch<T>(url: string, data?: any): Promise<T> {
        const response = await this.client.patch<T>(url, data);
        return response.data;
    }

    async delete<T>(url: string): Promise<T> {
        const response = await this.client.delete<T>(url);
        return response.data;
    }
}

export const apiClient = new ApiClient();

// ==================== COOKIE MANAGEMENT ====================

export const setAuthCookie = (token: string, user: any) => {
    const isSecure = import.meta.env.PROD; // true di production, false di development
    const options = { 
        expires: 7, 
        secure: isSecure,
        sameSite: 'strict' as const,
        path: '/'
    };
    
    Cookies.set('auth_token', token, options);
    Cookies.set('user', JSON.stringify(user), options);
};

export const removeAuthCookie = () => {
    Cookies.remove('auth_token');
    Cookies.remove('user');
    Cookies.remove('team_id');
};

export const getAuthToken = (): string | undefined => {
    return Cookies.get('auth_token');
};

export const getUserFromCookie = (): any => {
    const userStr = Cookies.get('user');
    if (userStr) {
        try {
            return JSON.parse(userStr);
        } catch {
            return null;
        }
    }
    return null;
};

export const isAuthenticated = (): boolean => {
    return !!getAuthToken();
};

