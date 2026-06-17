import { Outlet, useNavigate } from '@tanstack/react-router';
import { useState, useEffect } from 'react';
import { Sidebar } from './Sidebar';
import { Navbar } from './Navbar';
import { getToken } from '@/services/auth';

export function Layout() {
    const navigate = useNavigate();
    const [sidebarExpanded, setSidebarExpanded] = useState(false);
    const [isAuth, setIsAuth] = useState<boolean | null>(null);

    useEffect(() => {
        const token = getToken();
        setIsAuth(!!token);
        if (!token) navigate({ to: '/login' });
    }, [navigate]);

    if (isAuth === null) {
        return (
            <div className="flex h-screen items-center justify-center">
                <div className="h-8 w-8 animate-spin rounded-full border-4 border-green-600 border-t-transparent dark:border-purple-500" />
            </div>
        );
    }

    if (!isAuth) return null;

    return (
        <div className="min-h-screen bg-slate-50 dark:bg-[#080612]">
            <Sidebar
                expanded={sidebarExpanded}
                onExpandedChange={setSidebarExpanded}
            />
            <div className={`transition-all duration-200 ease-in-out ${
                sidebarExpanded ? "lg:ml-56" : "lg:ml-[60px]"
            }`}>
                <Navbar />
                <main className="p-6">
                    <Outlet />
                </main>
            </div>
        </div>
    );
}