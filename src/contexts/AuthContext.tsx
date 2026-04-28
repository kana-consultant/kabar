// src/contexts/AuthContext.tsx
import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { 
    login as loginService, 
    logout as logoutService, 
    getCurrentUser, 
    getUserRole, 
    getTeamId, 
    isAdmin as checkIsAdmin,
    isSuperAdmin as checkIsSuperAdmin
} from '@/services/authService';
import type { User } from '@/types/user';
import { toast } from 'sonner';

interface AuthContextType {
    // User data
    user: User | null;
    token: string | null;
    role: string | null;
    teamId: string | null;

    // Status
    isAuthenticated: boolean;
    isLoading: boolean;
    isAdminUser: boolean;      // ← renamed to avoid conflict
    isSuperAdmin: boolean;      // ← added

    // Permissions
    hasPermission: (roles: string[]) => boolean;
    can: (permission: string) => boolean;

    // Actions
    login: (email: string, password: string) => Promise<boolean>;
    logout: () => Promise<void>;
    refreshUser: () => Promise<void>;
    updateUser: (user: User) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Permission mapping (gunakan role yang sesuai dengan backend)
const permissionsMap: Record<string, string[]> = {
    // Draft permissions
    'draft:create': ['editor', 'manager', 'admin', 'super_admin'],
    'draft:read': ['viewer', 'editor', 'manager', 'admin', 'super_admin'],
    'draft:update': ['editor', 'manager', 'admin', 'super_admin'],
    'draft:delete': ['manager', 'admin', 'super_admin'],
    'draft:publish': ['manager', 'admin', 'super_admin'],
    'draft:schedule': ['manager', 'admin', 'super_admin'],
    'draft:review': ['manager', 'admin', 'super_admin'],
    'draft:approve': ['manager', 'admin', 'super_admin'],
    
    // Product permissions
    'product:create': ['manager', 'admin', 'super_admin'],
    'product:read': ['viewer', 'editor', 'manager', 'admin', 'super_admin'],
    'product:update': ['manager', 'admin', 'super_admin'],
    'product:delete': ['admin', 'super_admin'],
    'product:sync': ['manager', 'admin', 'super_admin'],
    
    // Team permissions
    'team:create': ['manager', 'admin', 'super_admin'],
    'team:read': ['viewer', 'editor', 'manager', 'admin', 'super_admin'],
    'team:update': ['manager', 'admin', 'super_admin'],
    'team:delete': ['admin', 'super_admin'],
    'team:manage_members': ['manager', 'admin', 'super_admin'],
    
    // User permissions
    'user:create': ['admin', 'super_admin'],
    'user:read': ['admin', 'super_admin'],
    'user:update': ['admin', 'super_admin'],
    'user:delete': ['admin', 'super_admin'],
    'user:manage_roles': ['admin', 'super_admin'],
    
    // Settings permissions
    'settings:read': ['admin', 'super_admin'],
    'settings:update': ['admin', 'super_admin'],
    
    // Report permissions
    'report:read': ['manager', 'admin', 'super_admin'],
    'report:export': ['manager', 'admin', 'super_admin'],
};

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [role, setRole] = useState<string | null>(null);
    const [teamId, setTeamId] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isAdminUser, setIsAdminUser] = useState(false);
    const [isSuperAdminUser, setIsSuperAdminUser] = useState(false);

    // Load user from localStorage on mount
    useEffect(() => {
        const loadUser = async () => {
            setIsLoading(true);
            try {
                const storedToken = localStorage.getItem('auth_token');
                const storedUser = localStorage.getItem('user');

                if (storedToken && storedUser) {
                    setToken(storedToken);
                    setUser(JSON.parse(storedUser));
                    setRole(getUserRole());
                    setTeamId(getTeamId());
                    setIsAdminUser(checkIsAdmin());
                    setIsSuperAdminUser(checkIsSuperAdmin());
                } else {
                    const currentUser = await getCurrentUser();
                    if (currentUser) {
                        setUser(currentUser);
                        setRole(currentUser.role);
                        setTeamId(getTeamId());
                        setIsAdminUser(currentUser.role === 'admin');
                        setIsSuperAdminUser(currentUser.role === 'super_admin');
                    }
                }
            } catch (error) {
                console.error('Failed to load user:', error);
            } finally {
                setIsLoading(false);
            }
        };

        loadUser();
    }, []);

    const login = useCallback(async (email: string, password: string): Promise<boolean> => {
        setIsLoading(true);
        try {
            const response = await loginService(email, password);
            setToken(response.token);
            setUser(response.user);
            setRole(response.user.role);
            setTeamId(localStorage.getItem('team_id'));
            setIsAdminUser(response.user.role === 'admin');
            setIsSuperAdminUser(response.user.role === 'super_admin');

            toast.success('Login berhasil', {
                description: `Selamat datang, ${response.user.name}!`,
            });
            return true;
        } catch (error) {
            console.error('Login failed:', error);
            toast.error('Login gagal', {
                description: 'Email atau password salah',
            });
            return false;
        } finally {
            setIsLoading(false);
        }
    }, []);

    const logout = useCallback(async () => {
        setIsLoading(true);
        try {
            await logoutService();
            setToken(null);
            setUser(null);
            setRole(null);
            setTeamId(null);
            setIsAdminUser(false);
            setIsSuperAdminUser(false);
            toast.success('Logout berhasil');
        } catch (error) {
            console.error('Logout failed:', error);
        } finally {
            setIsLoading(false);
        }
    }, []);

    const refreshUser = useCallback(async () => {
        try {
            const currentUser = await getCurrentUser();
            if (currentUser) {
                setUser(currentUser);
                setRole(currentUser.role);
                setIsAdminUser(currentUser.role === 'admin');
                setIsSuperAdminUser(currentUser.role === 'super_admin');
            }
        } catch (error) {
            console.error('Failed to refresh user:', error);
        }
    }, []);

    const updateUser = useCallback((updatedUser: User) => {
        setUser(updatedUser);
        localStorage.setItem('user', JSON.stringify(updatedUser));
        setRole(updatedUser.role);
        setIsAdminUser(updatedUser.role === 'admin');
        setIsSuperAdminUser(updatedUser.role === 'super_admin');
    }, []);

    const hasPermission = useCallback((requiredRoles: string[]): boolean => {
        if (!role) return false;
        if (role === 'super_admin') return true;
        return requiredRoles.includes(role);
    }, [role]);

    const can = useCallback((permission: string): boolean => {
        const allowedRoles = permissionsMap[permission] || [];
        return hasPermission(allowedRoles);
    }, [hasPermission]);

    const value: AuthContextType = {
        user,
        token,
        role,
        teamId,
        isAuthenticated: !!user && !!token,
        isLoading,
        isAdminUser,
        isSuperAdmin: isSuperAdminUser,
        hasPermission,
        can,
        login,
        logout,
        refreshUser,
        updateUser,
    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
}

// Custom hook untuk menggunakan AuthContext
export function useAuthContext() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuthContext must be used within an AuthProvider');
    }
    return context;
}

// Re-export untuk kemudahan (bisa juga di hooks/useAuth.ts)
export const useAuth = useAuthContext; 