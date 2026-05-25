import { useCallback } from 'react';
import {
    login as loginService,
    logout as logoutService,
    clearAuthData,
    updateLocalUser, changePasswordApi as changePasswordService
} from '@/services/auth';

import type { User } from '@/services/user';
import type { ToastContextType } from '../use-toast';

interface UseAuthActionsParams {
    setToken: (token: string | null) => void;
    setUser: (user: User) => void;
    setRole: (role: string | null) => void;
    setTeamId: (teamId: string | null) => void;
    setIsAdmin: (isAdmin: boolean) => void;
    setIsSuperAdmin: (isSuperAdmin: boolean) => void;
    setIsLoading: (isLoading: boolean) => void;
    loadUser: () => Promise<void>;
    clearUserState: () => void;
    toast : ToastContextType
}

export function useAuthActions({
    setToken,
    setUser,
    setRole,
    setTeamId,
    setIsAdmin,
    setIsSuperAdmin,
    setIsLoading,
    loadUser,
    clearUserState,
    toast,

}: UseAuthActionsParams) {
    const login = useCallback(async (email: string, password: string): Promise<boolean> => {
        setIsLoading(true);
        try {
            const response = await loginService(email, password);

            setToken(response.token);
            setUser(response.user);
            setRole(response.user.role);
            setTeamId(response.teamId || null);
            setIsAdmin(response.user.role === 'admin');
            // setIsSuperAdmin(response.user.role === 'super_admin');

            toast.success('Login berhasil!', {
                description: `Selamat datang, ${response.user.name}!`,
            });
            
            return true;
        } catch (error: any) {
            const status = error.response?.status;
            const message = error.response?.data?.message || error.message || 'Login gagal';

            // Handle berdasarkan status code
            switch (status) {
                case 400:
                    toast.error('Login gagal', {
                        description: 'Format email atau password tidak valid.',
                    });
                    break;
                case 403:
                    toast.error('Login gagal', {
                        description: 'Email atau password yang Anda masukkan tidak sesuai.',
                    });
                    break;
                case 401:
                    toast.error('Akses ditolak', {
                        description: message || 'Anda tidak memiliki akses untuk login.',
                    });
                    break;
                case 404:
                    toast.error('Login gagal', {
                        description: 'Akun tidak ditemukan.',
                    });
                    break;
                case 422:
                    toast.error('Validasi gagal', {
                        description: message || 'Silakan periksa kembali email dan password Anda.',
                    });
                    break;
                case 429:
                    toast.error('Terlalu banyak percobaan', {
                        description: 'Silakan coba lagi setelah beberapa saat.',
                    });
                    break;
                case 500:
                    toast.error('Server error', {
                        description: 'Terjadi kesalahan pada server. Silakan coba lagi nanti.',
                    });
                    break;
                default:
                    toast.error('Login gagal', {
                        description: message,
                    });
                    break;
            }

            return false;
        } finally {
            setIsLoading(false);
        }
    }, [setToken, setUser, setRole, setTeamId, setIsAdmin, setIsSuperAdmin, setIsLoading]);

    const logout = useCallback(async () => {
        setIsLoading(true);
        try {
            await logoutService();
            clearUserState();
            toast.success('Logout berhasil');
        } catch (error) {
            console.error('Logout failed:', error);
        } finally {
            setIsLoading(false);
        }
    }, [setIsLoading, clearUserState]);

    const refreshUser = useCallback(async () => {
        await loadUser();
    }, [loadUser]);

    const updateUser = useCallback((updatedUser: User) => {
        setUser(updatedUser);
        updateLocalUser(updatedUser);
    }, [setUser]);

    const clearAuth = useCallback(() => {
        clearAuthData();
        clearUserState();
    }, [clearUserState]);

    const changePassword = useCallback(async (oldPassword: string, newPassword: string): Promise<void> => {
        await changePasswordService(oldPassword, newPassword)
        return
    }, [changePasswordService])

    return {
        login,
        logout,
        refreshUser,
        updateUser,
        clearAuth,
        changePassword
    };
}