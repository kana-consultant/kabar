import { Outlet, useNavigate } from '@tanstack/react-router';
import { useState, useEffect } from 'react';
import { Sidebar } from './Sidebar';
import { Navbar } from './Navbar';
import { getToken } from '@/services/auth';

export function Layout() {
    const navigate = useNavigate();
    const [sidebarOpen, setSidebarOpen] = useState(true);
    const [isAuth, setIsAuth] = useState<boolean | null>(null);

    useEffect(() => {
        const token = getToken();
        setIsAuth(!!token);
        
        // Redirect ke login jika tidak ada token
        if (!token) {
            navigate({ to: '/login' });
        }
    }, [navigate]);

    if (isAuth === null) {
        return (
            <div className="flex h-screen items-center justify-center">
                <div className="h-8 w-8 animate-spin rounded-full border-4 border-blue-600 border-t-transparent" />
            </div>
        );
    }

    if (!isAuth) {
        return null; // Akan redirect oleh useEffect
    }

    const toggleSidebar = () => {
        setSidebarOpen(!sidebarOpen);
    };

    return (
        <div className="min-h-screen bg-white dark:bg-black">
            <Sidebar isOpen={sidebarOpen} />
            <div className={`transition-all duration-300 ${sidebarOpen ? "lg:ml-64" : "lg:ml-20"}`}>
                <Navbar onToggleSidebar={toggleSidebar} sidebarOpen={sidebarOpen} />
                <main className="p-6">
                    <Outlet />
                </main>
            </div>
        </div>
    );
}