// lib/api.ts
import axios, { 
  type AxiosInstance, 
  type InternalAxiosRequestConfig, 
  type AxiosResponse,
  type AxiosRequestConfig 
} from 'axios';
import Cookies from 'js-cookie';
import { jwtDecode } from "jwt-decode";

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
            timeout: 300000,
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
            (response: AxiosResponse) => {
                // Jika response adalah XML, biarkan sebagai string
                const contentType = response.headers['content-type'] || '';
                if (contentType.includes('xml')) {
                    return response;
                }
                return response;
            },
            (error) => {
                if (error.response?.status === 401) {
                    removeAuthCookie();
                    window.location.href = '/login';
                }
                return Promise.reject(error);
            }
        );
    }

    // GET untuk JSON response
    async get<T>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.client.get<T>(url, { 
            params, 
            ...config,
        });
        return response.data;
    }

    // GET khusus untuk XML response - returns full response with headers
    async getXML(url: string, params?: any): Promise<{ data: string; headers: any; status: number }> {
        const response = await this.client.get(url, {
            params,
            // Prevent axios from parsing XML
            transformResponse: [(data) => data],
            // Important: don't set Content-Type for request
            headers: {
                'Accept': 'application/xml',
            },
        });

        return {
            data: response.data as string,
            headers: response.headers,
            status: response.status,
        };
    }

    async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.client.post<T>(url, data, config);
        return response.data;
    }

    async put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.client.put<T>(url, data, config);
        return response.data;
    }

    async patch<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.client.patch<T>(url, data, config);
        return response.data;
    }

    async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.client.delete<T>(url, config);
        return response.data;
    }
}

export const apiClient = new ApiClient();

// ==================== COOKIE MANAGEMENT ====================

export const setAuthCookie = (token: string, user: any) => {
    const isSecure = import.meta.env.PROD;
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
};

export const getAuthToken = (): string | undefined => {
    return Cookies.get('auth_token');
};

// ==================== JWT DECODE ====================

interface JwtPayload {
    user_id: string;
    team_id: string;
    email: string;
    name: string;
    role: string;
    exp: number;
    iat?: number;
    [key: string]: any;
}

export const getUserFromCookie = (): JwtPayload | null => {
    const token = Cookies.get("auth_token");
    if (!token) return null;

    try {
        const decoded = jwtDecode<JwtPayload>(token);
        return decoded;
    } catch {
        return null;
    }
};

export const isAuthenticated = (): boolean => {
    return !!getAuthToken();
};