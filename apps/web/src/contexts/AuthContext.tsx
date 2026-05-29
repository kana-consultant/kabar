// src/contexts/AuthContext.tsx
import React, { createContext, useContext } from 'react';
import { useEffect } from 'react';

import type { User } from '@/services/user';
import { useAuth as authService } from '@/hooks/auth/useAuth';
import { getPermissions } from '@/services/user/permissions';


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

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const auth = authService();

    const can = (permission: string): boolean => {
        const allowedRoles = getPermissions();
        return allowedRoles.includes(permission);
    }

    // Load user on mount
    useEffect(() => {
        auth.refreshUser();
    }, []); // Hanya jalan sekali saat mount

    const value: AuthContextType = {
        // User data
        user: auth.user,
        token: auth.token,
        role: auth.role,
        teamId: auth.teamId,

        // Status
        isAuthenticated: auth.isAuthenticated,
        isLoading: auth.isLoading,
        isAdminUser: auth.isAdmin,
        isSuperAdmin: auth.isSuperAdmin,
        hasPermission: auth.hasAnyRole,

        // Permissions
        can,

        // Actions
        login: auth.login,
        logout: auth.logout,
        refreshUser: auth.refreshUser,
        updateUser: auth.updateUser,
    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuthContext() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuthContext must be used within an AuthProvider');
    }
    return context;
}


// Re-export untuk kemudahan (bisa juga di hooks/useAuth.ts)
export const useAuth = useAuthContext; 