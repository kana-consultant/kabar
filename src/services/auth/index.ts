// Types
export type { LoginResponse } from './types';

// Config
export { COOKIE_OPTIONS } from './config';

// Login/Logout
export { login } from './login';
export { logout } from './logout';

// Storage
export { 
    getToken, 
    hasToken, 
    clearAuthData, 
    getTeamId,
    setAuthCookie,
    removeAuthCookie
} from './storage';

// User
export { 
    getCurrentUser, 
    getUserRole, 
    getUserName, 
    getUserEmail, 
    getUserAvatar, 
    updateLocalUser 
} from './user';

// Permissions
export { 
    hasRole, 
    hasAnyRole, 
    hasPermission, 
    getRoleLevel, 
    canAccessTeam,
    hasToken as hasTokenPermission
} from './permissions';

// Admin
export { 
    isAdmin, 
    isSuperAdmin, 
    isAnyAdmin, 
    isManagerOrAbove, 
    isEditorOrAbove 
} from './admin';

// Token
export { setAuthCookie as setAuthToken, removeAuthCookie as removeAuthToken } from './token';