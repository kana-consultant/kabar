// src/layouts/AuthLayout.tsx
import { useEffect } from 'react';
import { Outlet, useLocation, useNavigate } from '@tanstack/react-router';
import { useAuth } from '@/contexts/AuthContext';
import { hasToken } from '@/services/auth';

// Public routes yang bisa diakses tanpa login
const PUBLIC_ROUTES = ['/login', '/register', '/forgot-password', '/reset-password', '/invite'];

export function AuthLayout() {
    const navigate = useNavigate();
    const location = useLocation();
    const { isAuthenticated, isLoading } = useAuth();

    // Handle redirect berdasarkan auth status
    useEffect(() => {
       
        const currentPath = location.pathname;
        const isPublicRoute = PUBLIC_ROUTES.some(route => currentPath.startsWith(route));
        const isLoggedIn = isAuthenticated && hasToken();

        // CASE 1: User sudah login dan mencoba akses public route (login/register/etc)
        if (isLoggedIn && isPublicRoute) {
            // Redirect ke dashboard
            navigate({
                to: '/dashboard',
                replace: true
            });
            return;
        }

        // CASE 2: User belum login dan mencoba akses protected route
        if (!isLoggedIn && !isPublicRoute) {
            // Redirect ke login, simpan intended destination
            navigate({
                to: '/login',
                replace: true,
                state: { from: currentPath }
            });
            return;
        }

        // CASE 3: User sudah login di protected route -> allow access
        // CASE 4: User belum login di public route -> allow access
        
    }, [isAuthenticated, isLoading, location.pathname, navigate]);

    // Show loading spinner while checking auth
    if (isLoading) {
        return (
            <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900 flex items-center justify-center">
                <div className="flex flex-col items-center gap-4">
                    <div className="h-12 w-12 animate-spin rounded-full border-4 border-primary border-t-transparent" />
                    <p className="text-sm text-muted-foreground">Loading...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900">
            <Outlet />
        </div>
    );
}