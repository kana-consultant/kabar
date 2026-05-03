import { useCallback } from 'react';
import { 
    login as loginService, 
    logout as logoutService,
    clearAuthData,
    updateLocalUser
} from '@/services/auth';
import { toast } from 'sonner';
import type { User } from '@/types/user';

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
    clearUserState
}: UseAuthActionsParams) {
    const login = useCallback(async (email: string, password: string): Promise<boolean> => {
        setIsLoading(true);
        try {
            const response = await loginService(email, password);
            
            setToken(response.token);
            setUser(response.user);
            setRole(response.user.role);
            setTeamId(response.teamId || null);
            setIsAdmin(response.user.role === 'admin' || response.user.role === 'super_admin');
            setIsSuperAdmin(response.user.role === 'super_admin');
            
            toast.success('Login berhasil!', {
                description: `Selamat datang, ${response.user.name}!`,
            });
            return true;
        } catch (error: any) {
            const message = error.response?.data?.message || error.message || 'Login gagal';
            toast.error('Login gagal', { description: message });
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

    return {
        login,
        logout,
        refreshUser,
        updateUser,
        clearAuth
    };
}