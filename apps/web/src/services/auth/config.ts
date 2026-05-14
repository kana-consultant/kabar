// Cookie configuration
export const COOKIE_OPTIONS = {
    expires: 7, // 7 days
    secure: import.meta.env.PROD, // true in production
    sameSite: 'strict' as const,
    path: '/'
};