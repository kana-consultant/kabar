import { useState } from "react";
import { Link, useLocation } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { 
    LayoutDashboard, 
    FileText, 
    Settings, 
    History, 
    Package, 
    Rocket,
    FileStack,
    Calendar,
    Menu,
    X,
    Sparkles
} from "lucide-react";

const menuItems = [
    { title: "Dashboard", href: "/", icon: LayoutDashboard },
    { title: "Generate Konten", href: "/generate", icon: FileText },
    { title: "Produk", href: "/products", icon: Package },
    { title: "Draft", href: "/drafts", icon: FileStack },
    { title: "Schedule", href: "/schedule", icon: Calendar },
    { title: "History", href: "/history", icon: History },
    { title: "Settings", href: "/settings", icon: Settings },
];

interface SidebarProps {
    isOpen: boolean;
}

export function Sidebar({ isOpen }: SidebarProps) {
    const location = useLocation();
    const currentPath = location.pathname;
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

    const DesktopSidebar = () => (
        <aside
            className={`fixed inset-y-0 left-0 z-50 hidden border-r backdrop-blur-xl  lg:block ${
                isOpen ? "w-64" : "w-20"
            } ${
                isOpen 
                    ? "bg-white/80 dark:bg-black/30 border-slate-200 dark:border-white/20" 
                    : "bg-white/60 dark:bg-black/20 border-slate-200 dark:border-white/10"
            }`}
        >
            <div className={`flex h-16 items-center gap-2 border-b px-6 transition-all duration-300 ${
                isOpen ? "border-slate-200 dark:border-white/20" : "border-slate-200 dark:border-white/10"
            }`}>
                <div className="relative">
                    <Rocket className="h-6 w-6 text-emerald-600 dark:text-purple-400" />
                    <Sparkles className="absolute -top-1 -right-1 h-3 w-3 text-amber-500 dark:text-yellow-400 animate-pulse" />
                </div>
                {isOpen && (
                    <span className="text-lg font-bold text-slate-800 dark:bg-gradient-to-r dark:from-purple-400 dark:to-cyan-400 dark:bg-clip-text dark:text-transparent">
                        KABAR
                    </span>
                )}
            </div>
            
            <nav className="space-y-1 p-4">
                {menuItems.map((item) => {
                    const isActive = currentPath === item.href;
                    return (
                        <Link
                            key={item.href}
                            to={item.href}
                            className={`group relative flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-all duration-300 ${
                                isActive
                                    ? // Active item: background solid
                                      "bg-emerald-100 text-emerald-800 dark:bg-purple-900/40 dark:text-purple-300 shadow-md"
                                    : // Normal item: transparan, hover berbeda
                                      "text-slate-600 hover:bg-slate-100 hover:text-slate-900 dark:text-slate-300 dark:hover:bg-white/10 dark:hover:text-white"
                            } ${!isOpen && "justify-center"}`}
                            title={!isOpen ? item.title : undefined}
                        >
                            {isActive && (
                                <div className="absolute left-0 h-8 w-0.5 rounded-r-full bg-emerald-500 dark:bg-purple-500" />
                            )}
                            
                            <item.icon className={`h-5 w-5 transition-all duration-300 ${
                                isActive 
                                    ? "text-emerald-700 dark:text-purple-400" 
                                    : "text-slate-500 group-hover:text-slate-700 dark:text-slate-400 dark:group-hover:text-white"
                            }`} />
                            
                            {isOpen && (
                                <span className="transition-all duration-300">
                                    {item.title}
                                </span>
                            )}
                            
                            {!isOpen && !isActive && (
                                <span className="absolute left-full ml-2 hidden rounded-md bg-slate-800 px-2 py-1 text-xs text-white opacity-0 transition-opacity group-hover:opacity-100 group-hover:block dark:bg-slate-800">
                                    {item.title}
                                </span>
                            )}
                        </Link>
                    );
                })}
            </nav>
            
            <div className="absolute bottom-0 left-0 right-0 h-32 bg-gradient-to-t from-emerald-500/5 to-transparent pointer-events-none dark:from-purple-500/5" />
        </aside>
    );

    const MobileSidebar = () => (
        <>
            <Button
                variant="ghost"
                size="icon"
                className="fixed left-4 top-3 z-50 lg:hidden bg-white/80 backdrop-blur-lg border border-slate-200 hover:bg-white dark:bg-black/30 dark:border-white/20 dark:hover:bg-white/10"
                onClick={() => setMobileMenuOpen(true)}
            >
                <Menu className="h-5 w-5 text-slate-700 dark:text-white" />
            </Button>

            {mobileMenuOpen && (
                <div className="fixed inset-0 z-50 lg:hidden">
                    <div
                        className="fixed inset-0 bg-black/50 backdrop-blur-sm"
                        onClick={() => setMobileMenuOpen(false)}
                    />
                    <aside className="fixed inset-y-0 left-0 z-50 w-64 bg-white shadow-xl dark:bg-gradient-to-b dark:from-slate-900 dark:to-slate-950 dark:border-r dark:border-white/10">
                        <div className="flex h-16 items-center justify-between border-b px-4 border-slate-200 dark:border-white/10">
                            <div className="flex items-center gap-2">
                                <div className="relative">
                                    <Rocket className="h-6 w-6 text-emerald-600 dark:text-purple-400" />
                                    <Sparkles className="absolute -top-1 -right-1 h-3 w-3 text-amber-500 dark:text-yellow-400 animate-pulse" />
                                </div>
                                <span className="font-bold text-slate-800 dark:bg-gradient-to-r dark:from-purple-400 dark:to-cyan-400 dark:bg-clip-text dark:text-transparent">
                                    KABAR
                                </span>
                            </div>
                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => setMobileMenuOpen(false)}
                                className="h-8 w-8 text-slate-600 hover:bg-slate-100 dark:text-white dark:hover:bg-white/10"
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
                                        className={`flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-all duration-300 ${
                                            isActive
                                                ? "bg-emerald-100 text-emerald-800 dark:bg-purple-900/40 dark:text-purple-300"
                                                : "text-slate-600 hover:bg-slate-100 hover:text-slate-900 dark:text-slate-300 dark:hover:bg-white/10 dark:hover:text-white"
                                        }`}
                                    >
                                        <item.icon className="h-4 w-4" />
                                        <span>{item.title}</span>
                                    </Link>
                                );
                            })}
                        </nav>
                    </aside>
                </div>
            )}
        </>
    );

    return (
        <>
            <DesktopSidebar />
            <MobileSidebar />
        </>
    );
}