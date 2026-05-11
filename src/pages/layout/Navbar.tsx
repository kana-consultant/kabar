import { useState, useEffect } from "react";
import { Link, useLocation, useNavigate } from "@tanstack/react-router";
import { Menu, X, Rocket, LayoutDashboard, FileText, Package, History, Settings, ChevronLeft, ChevronRight, Sun, Moon, FileStack, Calendar, LogOut, User, Shield } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ThemeSwitch } from "@/components/switch";
import { toast } from "sonner";
import { removeAuthCookie } from "@/services/api";

const menuItems = [
    { title: "Dashboard", href: "/", icon: LayoutDashboard },
    { title: "Generate Konten", href: "/generate", icon: FileText },
    { title: "Produk", href: "/products", icon: Package },
    { title: "Draft", href: "/drafts", icon: FileStack },
    { title: "Schedule", href: "/schedule", icon: Calendar },
    { title: "History", href: "/history", icon: History },
    { title: "Settings", href: "/settings", icon: Settings },
];

interface NavbarProps {
    onToggleSidebar: () => void;
    sidebarOpen: boolean;
}

export function Navbar({ onToggleSidebar, sidebarOpen }: NavbarProps) {
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
    const [isScrolled, setIsScrolled] = useState(false);
    const [dropdownOpen, setDropdownOpen] = useState(false);
    const location = useLocation();
    const navigate = useNavigate();
    const currentPath = location.pathname;

    const currentTitle = menuItems.find((item) => item.href === currentPath)?.title || "Dashboard";

    // Deteksi scroll
    useEffect(() => {
        const handleScroll = () => {
            setIsScrolled(window.scrollY > 10);
        };
        window.addEventListener("scroll", handleScroll);
        return () => window.removeEventListener("scroll", handleScroll);
    }, []);

    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            const target = event.target as HTMLElement;
            if (!target.closest('.avatar-dropdown')) {
                setDropdownOpen(false);
            }
        };
        document.addEventListener('click', handleClickOutside);
        return () => document.removeEventListener('click', handleClickOutside);
    }, []);

    const handleLogout = () => {
        removeAuthCookie();
        toast.success("Berhasil logout", {
            description: "Anda telah keluar dari aplikasi",
        });
        navigate({ to: "/login" });
    };

    const handleSettings = () => {
        setDropdownOpen(false);
        navigate({ to: "/settings" });
    };



    // Get user info from cookie or localStorage
    const getUserInitial = () => {
        try {
            const userStr = localStorage.getItem('user') || document.cookie.match('user=([^;]+)')?.[1];
            if (userStr) {
                const user = JSON.parse(decodeURIComponent(userStr));
                return user.name ? user.name.charAt(0).toUpperCase() : 'A';
            }
        } catch (error) {
            console.error('Failed to get user info:', error);
        }
        return 'A';
    };

    const getUserName = () => {
        try {
            const userStr = localStorage.getItem('user') || document.cookie.match('user=([^;]+)')?.[1];
            if (userStr) {
                const user = JSON.parse(decodeURIComponent(userStr));
                return user.name || user.email || 'Admin';
            }
        } catch (error) {
            console.error('Failed to get user name:', error);
        }
        return 'Admin';
    };

    const getUserEmail = () => {
        try {
            const userStr = localStorage.getItem('user') || document.cookie.match('user=([^;]+)')?.[1];
            if (userStr) {
                const user = JSON.parse(decodeURIComponent(userStr));
                return user.email || 'admin@example.com';
            }
        } catch (error) {
            console.error('Failed to get user email:', error);
        }
        return 'admin@example.com';
    };

    return (
        <>
            <header
                className={`sticky top-0 z-40 flex h-16 items-center gap-4 border-b px-4 ${
                    isScrolled
                        ? "bg-white/80 backdrop-blur-lg dark:bg-black/80"
                        : "bg-white dark:bg-black"
                } border-slate-200 dark:border-slate-800`}
            >
                {/* Tombol toggle sidebar (desktop) */}
                <Button
                    variant="ghost"
                    size="icon"
                    className="hidden lg:flex shrink-0"
                    onClick={onToggleSidebar}
                >
                    {sidebarOpen ? (
                        <ChevronLeft className="h-5 w-5" />
                    ) : (
                        <ChevronRight className="h-5 w-5" />
                    )}
                </Button>

                {/* Tombol menu mobile */}
                <Button
                    variant="ghost"
                    size="icon"
                    className="lg:hidden shrink-0"
                    onClick={() => setMobileMenuOpen(true)}
                >
                    <Menu className="h-5 w-5" />
                </Button>

                <div className="flex flex-1 items-center justify-between">
                    <h1 className="text-lg font-semibold truncate">{currentTitle}</h1>
                    <div className="flex items-center gap-3">
                        {/* Theme Toggle Switch */}
                        <div className="flex items-center gap-2">
                            <Sun className="h-4 w-4 text-slate-500 dark:text-slate-400" />
                            <ThemeSwitch />
                            <Moon className="h-4 w-4 text-slate-500 dark:text-slate-400" />
                        </div>
                        
                        {/* Avatar with Dropdown */}
                        <div className="relative avatar-dropdown">
                            <button
                                onClick={(e) => {
                                    e.stopPropagation();
                                    setDropdownOpen(!dropdownOpen);
                                }}
                                className="flex h-8 w-8 items-center justify-center rounded-full bg-gradient-to-r from-blue-600 to-blue-500 text-sm font-medium text-white shadow-sm hover:shadow-md transition-all focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-slate-900"
                            >
                                {getUserInitial()}
                            </button>

                            {/* Dropdown Menu */}
                            {dropdownOpen && (
                                <div className="absolute right-0 mt-2 w-56 rounded-lg bg-white shadow-lg ring-1 ring-black ring-opacity-5 dark:bg-slate-800 dark:ring-slate-700 z-50">
                                    <div className="px-4 py-3 border-b border-slate-100 dark:border-slate-700">
                                        <p className="text-sm font-medium text-slate-900 dark:text-white">
                                            {getUserName()}
                                        </p>
                                        <p className="text-xs text-slate-500 dark:text-slate-400 truncate">
                                            {getUserEmail()}
                                        </p>
                                    </div>
                                    <div className="py-1">
                                        {/* <button
                                            onClick={handleProfile}
                                            className="flex w-full items-center gap-3 px-4 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-700 transition-colors"
                                        >
                                            <User className="h-4 w-4" />
                                            Profile
                                        </button> */}
                                        <button
                                            onClick={handleSettings}
                                            className="flex w-full items-center gap-3 px-4 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-700 transition-colors"
                                        >
                                            <Settings className="h-4 w-4" />
                                            Settings
                                        </button>
                                        {getUserEmail() === 'admin@example.com' && (
                                            <button
                                                className="flex w-full items-center gap-3 px-4 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-700 transition-colors"
                                            >
                                                <Shield className="h-4 w-4" />
                                                Admin Panel
                                            </button>
                                        )}
                                    </div>
                                    <div className="border-t border-slate-100 dark:border-slate-700 py-1">
                                        <button
                                            onClick={handleLogout}
                                            className="flex w-full items-center gap-3 px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-950 transition-colors"
                                        >
                                            <LogOut className="h-4 w-4" />
                                            Logout
                                        </button>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </header>

            {/* Mobile sidebar */}
            {mobileMenuOpen && (
                <div className="fixed inset-0 z-50 lg:hidden">
                    <div
                        className="fixed inset-0 bg-black/50 backdrop-blur-sm"
                        onClick={() => setMobileMenuOpen(false)}
                    />
                    <div className="fixed inset-y-0 left-0 z-50 w-64 bg-white shadow-xl dark:bg-slate-950">
                        <div className="flex h-16 items-center justify-between border-b px-4 dark:border-slate-800">
                            <div className="flex items-center gap-2">
                                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600">
                                    <Rocket className="h-4 w-4 text-white" />
                                </div>
                                <span className="font-bold text-slate-800 dark:text-white">SEO Multi-Post</span>
                            </div>
                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => setMobileMenuOpen(false)}
                                className="h-8 w-8"
                            >
                                <X className="h-4 w-4" />
                            </Button>
                        </div>
                        <nav className="space-y-1 p-4">
                            {menuItems.map((item) => {
                                const isActive = currentPath === item.href;
                                return (
                                    <Link
                                        key={item.href}
                                        to={item.href}
                                        onClick={() => setMobileMenuOpen(false)}
                                        className={`flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-all ${
                                            isActive
                                                ? "bg-blue-50 text-blue-600 dark:bg-blue-950 dark:text-blue-400"
                                                : "text-slate-600 hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-800"
                                        }`}
                                    >
                                        <item.icon className="h-4 w-4" />
                                        {item.title}
                                    </Link>
                                );
                            })}
                            
                            {/* Mobile logout button */}
                            <div className="pt-4 mt-4 border-t border-slate-200 dark:border-slate-800">
                                <button
                                    onClick={() => {
                                        setMobileMenuOpen(false);
                                        handleLogout();
                                    }}
                                    className="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-950 transition-colors"
                                >
                                    <LogOut className="h-4 w-4" />
                                    Logout
                                </button>
                            </div>
                        </nav>
                    </div>
                </div>
            )}
        </>
    );
}