import type { User } from './user';

export interface AuthState {
    user: User | null;
    token: string | null;
    role: string | null;
    permissions: string[] ;
    teamId: string | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    isAdmin: boolean;
    isSuperAdmin: boolean;
}

export interface AuthSetters {
    setUser: (user: User) => void;
    setToken: (token: string | null) => void;
    setRole: (role: string | null) => void;
    setPermissions: (permissions: string[]) => void;
    setTeamId: (teamId: string | null) => void;
}

export interface AuthActions {
    login: (email: string, password: string) => Promise<boolean>;
    logout: () => Promise<void>;
    refreshUser: () => Promise<void>;
    updateUser: (user: User) => void;
    clearAuth: () => void;
    changePassword: (oldPassword: string, newPassword: string) => Promise<void>;
}

export interface RoleChecks {
    isManagerOrAbove: boolean;
    isEditorOrAbove: boolean;
    hasRole: (role: string) => boolean;
    hasAnyRole: (roles: string[]) => boolean;
    getRoleLevel: () => number;
    canAccessTeam: (resourceTeamId: string | null) => boolean;
}

export type UseAuthReturn = AuthState & AuthSetters & AuthActions & RoleChecks;