import { useState } from "react";
import { Link, useLocation } from "@tanstack/react-router";
import { Button } from "@kana-consultant/ui-kit";
import {
    LayoutDashboard, FileText, Settings, History,
    Package, FileStack, Calendar, Menu, X,
    Sparkles, Rocket, ChevronLeft, ChevronRight
} from "lucide-react";
import { cn } from "@/lib/utils";

const menuGroups = [
    {
        label: "MENU",
        items: [
            { title: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
            { title: "Generate Konten", href: "/generate", icon: FileText },
            { title: "Produk", href: "/products", icon: Package },
            { title: "Draft", href: "/drafts", icon: FileStack },
        ],
    },
    {
        label: "MANAJEMEN",
        items: [
            { title: "Schedule", href: "/schedule", icon: Calendar },
            { title: "History", href: "/history", icon: History },
        ],
    },
    {
        label: "GENERAL",
        items: [
            { title: "Settings", href: "/settings", icon: Settings },
        ],
    },
];

function NavItems({
    expanded,
    currentPath,
    onClose,
}: {
    expanded: boolean;
    currentPath: string;
    onClose?: () => void;
}) {
    return (
        <nav className="flex-1 overflow-y-auto px-2 py-3 space-y-4">
            {menuGroups.map((group) => (
                <div key={group.label}>
                    <div className={cn(
                        "overflow-hidden transition-all duration-200",
                        expanded ? "max-h-6 opacity-100 mb-1" : "max-h-0 opacity-0 mb-0"
                    )}>
                        <p className="px-2.5 text-[10px] font-semibold tracking-widest text-slate-400 dark:text-slate-600">
                            {group.label}
                        </p>
                    </div>

                    <div className="space-y-0.5">
                        {group.items.map((item) => {
                            const isActive = currentPath === item.href;
                            return (
                                <Link
                                    key={item.href}
                                    to={item.href}
                                    onClick={onClose}
                                    title={!expanded ? item.title : undefined}
                                    className={cn(
                                        "group relative flex items-center rounded-lg text-sm font-medium transition-all duration-150 overflow-hidden",
                                        expanded
                                            ? "gap-2.5 px-2.5 py-1.5"
                                            : "justify-center w-9 h-9 mx-auto",
                                        isActive
                                            ? "bg-green-50 text-green-700 dark:bg-purple-500/[0.12] dark:text-purple-300"
                                            : "text-slate-500 hover:text-slate-800 hover:bg-slate-100/80 dark:text-slate-400 dark:hover:text-slate-100 dark:hover:bg-white/[0.05]"
                                    )}
                                >
                                    {isActive && expanded && (
                                        <span className="absolute left-0 top-1 bottom-1 w-[3px] rounded-r-full bg-green-600 dark:bg-purple-500" />
                                    )}

                                    <item.icon className={cn(
                                        "shrink-0 transition-colors",
                                        expanded ? "h-4 w-4" : "h-[18px] w-[18px]",
                                        isActive
                                            ? "text-green-600 dark:text-purple-400"
                                            : "text-slate-400 group-hover:text-slate-600 dark:text-slate-500 dark:group-hover:text-slate-300"
                                    )} />

                                    <span className={cn(
                                        "whitespace-nowrap overflow-hidden transition-all duration-200 leading-none",
                                        expanded ? "max-w-[160px] opacity-100" : "max-w-0 opacity-0 w-0"
                                    )}>
                                        {item.title}
                                    </span>
                                </Link>
                            );
                        })}
                    </div>
                </div>
            ))}
        </nav>
    );
}

interface SidebarProps {
    expanded: boolean;
    onExpandedChange: (value: boolean) => void;
}

export function Sidebar({ expanded, onExpandedChange }: SidebarProps) {
    const location = useLocation();
    const currentPath = location.pathname;
    const [mobileOpen, setMobileOpen] = useState(false);

    return (
        <>
            {/* ── Desktop sidebar ── */}
            <aside className={cn(
                "fixed inset-y-0 left-0 z-40 hidden lg:flex flex-col",
                "border-r transition-all duration-200 ease-in-out",
                "bg-white border-slate-200/80",
                "dark:bg-[#0a0812] dark:border-white/[0.06]",
                expanded ? "w-56" : "w-[60px]"
            )}>
                <div className={cn(
                    "flex h-14 shrink-0 items-center border-b px-3",
                    "border-slate-100 dark:border-white/[0.05]",
                    expanded ? "justify-between" : "justify-center"
                )}>
                    <div className={cn(
                        "flex items-center gap-2 overflow-hidden transition-all duration-200",
                        expanded ? "max-w-[120px] opacity-100" : "max-w-0 opacity-0 w-0"
                    )}>
                        <div className="relative shrink-0">
                            <Rocket className="h-5 w-5 text-green-600 dark:text-purple-400" />
                            <Sparkles className="absolute -top-1 -right-1 h-2.5 w-2.5 text-amber-400 animate-pulse" />
                        </div>
                        <span className="text-base font-bold tracking-tight text-slate-800 dark:text-white whitespace-nowrap">
                            KABAR
                        </span>
                    </div>

                    {!expanded && (
                        <div className="relative shrink-0">
                            <Rocket className="h-5 w-5 text-green-600 dark:text-purple-400" />
                            <Sparkles className="absolute -top-1 -right-1 h-2.5 w-2.5 text-amber-400 animate-pulse" />
                        </div>
                    )}

                    <button
                        onClick={() => onExpandedChange(!expanded)}
                        title={expanded ? "Collapse" : "Expand"}
                        className={cn(
                            "flex h-7 w-7 shrink-0 items-center justify-center rounded-lg transition-colors",
                            "text-slate-400 hover:text-slate-700 hover:bg-slate-100",
                            "dark:text-slate-600 dark:hover:text-slate-300 dark:hover:bg-white/[0.06]",
                            !expanded ? "absolute opacity-0 pointer-events-none" : ""
                        )}
                    >
                        <ChevronLeft className="h-4 w-4" />
                    </button>
                </div>

                {!expanded && (
                    <div className="flex justify-center pt-2 px-2">
                        <button
                            onClick={() => onExpandedChange(true)}
                            title="Expand"
                            className={cn(
                                "flex h-8 w-8 items-center justify-center rounded-lg transition-colors",
                                "text-slate-400 hover:text-slate-700 hover:bg-slate-100",
                                "dark:text-slate-600 dark:hover:text-slate-300 dark:hover:bg-white/[0.06]"
                            )}
                        >
                            <ChevronRight className="h-4 w-4" />
                        </button>
                    </div>
                )}

                <NavItems expanded={expanded} currentPath={currentPath} />
            </aside>

            {/* ── Mobile trigger ── */}
            <Button
                variant="ghost"
                size="icon"
                onClick={() => setMobileOpen(true)}
                className={cn(
                    "fixed left-3 top-3 z-50 h-8 w-8 lg:hidden",
                    "bg-white border border-slate-200/80 shadow-sm",
                    "dark:bg-[#0a0812] dark:border-white/[0.08]"
                )}
            >
                <Menu className="h-4 w-4 text-slate-600 dark:text-slate-300" />
            </Button>

            {/* ── Mobile sidebar ── */}
            {mobileOpen && (
                <div className="fixed inset-0 z-50 lg:hidden">
                    <div
                        className="fixed inset-0 bg-black/40 backdrop-blur-sm"
                        onClick={() => setMobileOpen(false)}
                    />
                    <aside className={cn(
                        "fixed inset-y-0 left-0 z-50 flex w-56 flex-col",
                        "border-r bg-white border-slate-200/80 shadow-xl",
                        "dark:bg-[#0a0812] dark:border-white/[0.06]"
                    )}>
                        <div className="flex h-14 shrink-0 items-center justify-between border-b px-4 border-slate-100 dark:border-white/[0.05]">
                            <div className="flex items-center gap-2">
                                <div className="relative">
                                    <Rocket className="h-5 w-5 text-green-600 dark:text-purple-400" />
                                    <Sparkles className="absolute -top-1 -right-1 h-2.5 w-2.5 text-amber-400 animate-pulse" />
                                </div>
                                <span className="text-base font-bold tracking-tight text-slate-800 dark:text-white">
                                    KABAR
                                </span>
                            </div>
                            <button
                                onClick={() => setMobileOpen(false)}
                                className="flex h-7 w-7 items-center justify-center rounded-lg text-slate-400 hover:bg-slate-100 dark:hover:bg-white/[0.05]"
                            >
                                <X className="h-4 w-4" />
                            </button>
                        </div>
                        <NavItems
                            expanded={true}
                            currentPath={currentPath}
                            onClose={() => setMobileOpen(false)}
                        />
                    </aside>
                </div>
            )}
        </>
    );
}