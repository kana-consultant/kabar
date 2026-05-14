import { Outlet } from '@tanstack/react-router';

export function AuthLayout() {
    return (
        <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900">
            <Outlet />
        </div>
    );
}