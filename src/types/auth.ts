import type { User } from './user';

export interface AuthState {
    user: User | null;
    token: string | null;
    role: string | null;
    teamId: string | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    isAdmin: boolean;
    isSuperAdmin: boolean;
}

export interface AuthActions {
    login: (email: string, password: string) => Promise<boolean>;
    logout: () => Promise<void>;
    refreshUser: () => Promise<void>;
    updateUser: (user: User) => void;
    clearAuth: () => void;
}

export interface RoleChecks {
    isManagerOrAbove: boolean;
    isEditorOrAbove: boolean;
    hasRole: (role: string) => boolean;
    hasAnyRole: (roles: string[]) => boolean;
    getRoleLevel: () => number;
    canAccessTeam: (resourceTeamId: string | null) => boolean;
}

export type UseAuthReturn = AuthState & AuthActions & RoleChecks;