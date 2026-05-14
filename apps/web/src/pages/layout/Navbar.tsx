import { useState, useEffect } from "react";
import { useLocation, useNavigate } from "@tanstack/react-router";
import { Settings, LogOut, Shield, Sun, Moon, ChevronRight, Home } from "lucide-react";
import { ThemeSwitch } from "@/components/switch";
import { toast } from "sonner";
import { removeAuthCookie } from "@/services/api";
import { cn } from "@/lib/utils";

const menuItems = [
    { title: "Dashboard", href: "/" },
    { title: "Generate Konten", href: "/generate" },
    { title: "Produk", href: "/products" },
    { title: "Draft", href: "/drafts" },
    { title: "Schedule", href: "/schedule" },
    { title: "History", href: "/history" },
    { title: "Settings", href: "/settings" },
];

function getUserData() {
    try {
        const raw = localStorage.getItem("user") ||
            document.cookie.match("user=([^;]+)")?.[1];
        if (raw) return JSON.parse(decodeURIComponent(raw));
    } catch {
        return
    }
    return null;
}

function useThemeToggle() {
    const [isDark, setIsDark] = useState(
        () => document.documentElement.classList.contains("dark")
    );

    const toggle = () => {
        const html = document.documentElement;
        if (isDark) {
            html.classList.remove("dark");
            localStorage.setItem("theme", "light");
        } else {
            html.classList.add("dark");
            localStorage.setItem("theme", "dark");
        }
        setIsDark(!isDark);
    };

    return { isDark, toggle };
}

export function Navbar() {
    const [scrolled, setScrolled] = useState(false);
    const [dropdownOpen, setDropdownOpen] = useState(false);
   
    const location = useLocation();
    const navigate = useNavigate();
    const currentPath = location.pathname;

    const currentPage = menuItems.find(i => i.href === currentPath);
    const currentTitle = currentPage?.title ?? "Dashboard";
    const segments = currentPath.split("/").filter(Boolean);

    const user = getUserData();
    const initial = user?.name?.charAt(0).toUpperCase() ?? "A";
    const userName = (user?.name || user?.email || "Admin").split(" ")[0]; // first name only
    const userEmail = user?.email || "admin@example.com";
    const isAdmin = userEmail === "admin@example.com";

    useEffect(() => {
        const onScroll = () => setScrolled(window.scrollY > 10);
        window.addEventListener("scroll", onScroll);
        return () => window.removeEventListener("scroll", onScroll);
    }, []);

    useEffect(() => {
        const handler = (e: MouseEvent) => {
            if (!(e.target as HTMLElement).closest(".avatar-dropdown"))
                setDropdownOpen(false);
        };
        document.addEventListener("click", handler);
        return () => document.removeEventListener("click", handler);
    }, []);

    const handleLogout = () => {
        removeAuthCookie();
        toast.success("Berhasil logout");
        navigate({ to: "/login" });
    };

    return (
        <header className={cn(
            "sticky top-0 z-40 flex h-14 items-center gap-4 border-b px-5 transition-all duration-150",
            scrolled
                ? "bg-white/80 backdrop-blur-lg dark:bg-[#080612]/80"
                : "bg-white dark:bg-[#080612]",
            "border-slate-200/80 dark:border-white/[0.06]"
        )}>

            {/* ── Left: title + breadcrumb stacked ── */}
            <div className="flex flex-col justify-center flex-1 min-w-0">
                <span className="text-sm font-semibold text-slate-800 dark:text-white leading-tight">
                    {currentTitle}
                </span>
                <div className="flex items-center gap-1 mt-0.5">
                    <Home className="h-2.5 w-2.5 shrink-0 text-slate-400 dark:text-slate-600" />
                    <span className="text-[10px] text-slate-400 dark:text-slate-600">Home</span>
                    {segments.map((seg, i) => (
                        <span key={i} className="flex items-center gap-1">
                            <ChevronRight className="h-2.5 w-2.5 text-slate-300 dark:text-slate-700 shrink-0" />
                            <span className={cn(
                                "text-[10px]",
                                i === segments.length - 1
                                    ? "text-slate-500 dark:text-slate-500 font-medium"
                                    : "text-slate-400 dark:text-slate-600"
                            )}>
                                {menuItems.find(m => m.href === "/" + seg)?.title ?? seg}
                            </span>
                        </span>
                    ))}
                </div>
            </div>

            {/* ── Right controls ── */}
            <div className="flex items-center gap-2 shrink-0">

                {/* Theme toggle — single icon button */}
                <div className="flex items-center gap-1.5">
    <Sun className="h-3.5 w-3.5 text-slate-400 dark:text-slate-600" />
    <ThemeSwitch />
    <Moon className="h-3.5 w-3.5 text-slate-400 dark:text-slate-600" />
</div>

                {/* Divider */}
                <div className="h-5 w-px bg-slate-200 dark:bg-white/[0.08]" />

                {/* Avatar pill */}
                <div className="relative avatar-dropdown">
                    <button
                        onClick={(e) => {
                            e.stopPropagation();
                            setDropdownOpen(v => !v);
                        }}
                        className={cn(
                            "flex items-center gap-2 rounded-full border pl-1 pr-2.5 py-1 transition-all",
                            "border-slate-200/80 bg-white hover:bg-slate-50",
                            "dark:border-white/[0.08] dark:bg-white/[0.03] dark:hover:bg-white/[0.06]",
                            dropdownOpen && "border-green-300/60 dark:border-purple-500/40"
                        )}
                    >
                        {/* Mini avatar */}
                        <div className={cn(
                            "flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-[11px] font-semibold text-white",
                            "bg-gradient-to-br from-green-500 to-emerald-600",
                            "dark:from-purple-500 dark:to-violet-600"
                        )}>
                            {initial}
                        </div>

                        {/* Name */}
                        <span className="text-xs font-medium text-slate-700 dark:text-slate-200 max-w-[80px] truncate">
                            {userName}
                        </span>

                        {/* Chevron */}
                        <ChevronRight className={cn(
                            "h-3 w-3 text-slate-400 dark:text-slate-600 transition-transform duration-150",
                            dropdownOpen ? "rotate-90" : "rotate-0"
                        )} />
                    </button>

                    {/* Dropdown */}
                    {dropdownOpen && (
                        <div className={cn(
                            "absolute right-0 top-full mt-2 w-56 rounded-xl border shadow-lg z-50 overflow-hidden",
                            "bg-white border-slate-200/80",
                            "dark:bg-[#0f0d1a] dark:border-white/[0.08]"
                        )}>
                            {/* User info header */}
                            <div className={cn(
                                "flex items-center gap-2.5 px-4 py-3 border-b",
                                "border-slate-100 bg-slate-50/60",
                                "dark:border-white/[0.05] dark:bg-white/[0.02]"
                            )}>
                                <div className={cn(
                                    "flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-xs font-semibold text-white",
                                    "bg-gradient-to-br from-green-500 to-emerald-600",
                                    "dark:from-purple-500 dark:to-violet-600"
                                )}>
                                    {initial}
                                </div>
                                <div className="min-w-0">
                                    <p className="text-xs font-medium text-slate-800 dark:text-slate-100 truncate">
                                        {user?.name || "Admin"}
                                    </p>
                                    <p className="text-[10px] text-slate-400 dark:text-slate-600 truncate">
                                        {userEmail}
                                    </p>
                                </div>
                            </div>

                            {/* Menu */}
                            <div className="p-1">
                                <button
                                    onClick={() => { setDropdownOpen(false); navigate({ to: "/settings" }); }}
                                    className={cn(
                                        "flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-medium transition-colors",
                                        "text-slate-600 hover:bg-slate-100/80 hover:text-slate-800",
                                        "dark:text-slate-400 dark:hover:bg-white/[0.05] dark:hover:text-slate-200"
                                    )}
                                >
                                    <Settings className="h-3.5 w-3.5" />
                                    Settings
                                </button>

                                {isAdmin && (
                                    <button className={cn(
                                        "flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-medium transition-colors",
                                        "text-slate-600 hover:bg-slate-100/80 hover:text-slate-800",
                                        "dark:text-slate-400 dark:hover:bg-white/[0.05] dark:hover:text-slate-200"
                                    )}>
                                        <Shield className="h-3.5 w-3.5" />
                                        Admin Panel
                                    </button>
                                )}
                            </div>

                            {/* Logout */}
                            <div className="border-t p-1 border-slate-100 dark:border-white/[0.05]">
                                <button
                                    onClick={handleLogout}
                                    className={cn(
                                        "flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-medium transition-colors",
                                        "text-red-500 hover:bg-red-50 hover:text-red-600",
                                        "dark:text-red-400 dark:hover:bg-red-500/10"
                                    )}
                                >
                                    <LogOut className="h-3.5 w-3.5" />
                                    Logout
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </header>
    );
}