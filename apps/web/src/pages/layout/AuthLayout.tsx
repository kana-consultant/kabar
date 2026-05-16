// src/layouts/AuthLayout.tsx
import { useEffect } from 'react';
import { Outlet, useLocation, useNavigate } from '@tanstack/react-router';
import { useAuth } from '@/contexts/AuthContext';
import { hasToken } from '@/services/auth';

// Public routes yang bisa diakses tanpa login
const PUBLIC_ROUTES = ['/login', '/register', '/forgot-password', '/reset-password'];

export function AuthLayout() {
    const navigate = useNavigate();
    const location = useLocation();
    const { isAuthenticated, isLoading } = useAuth();

    // Handle redirect berdasarkan auth status
    useEffect(() => {
        if (isLoading) return;

        const currentPath = location.pathname;
        const isPublicRoute = PUBLIC_ROUTES.some(route => currentPath.startsWith(route));

        // Jika tidak login dan bukan di public route, redirect ke login
        if (!hasToken()) {
            navigate({
                to: '/login',
                replace: true,
                state: { from: currentPath }
            });
            return;
        }

        // Jika sudah login dan mencoba akses halaman auth (login/register)
        if (hasToken()) {
            // Redirect ke dashboard atau intended destination
            const from = (location.state as any)?.from || '/dashboard';
            navigate({
                to: from,
                replace: true
            });
            return;
        }
    }, [isAuthenticated, isLoading, location.pathname, location.state, navigate]);

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