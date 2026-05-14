import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import type { User } from "@/types/user";

const roleConfig: Record<string, { label: string; class: string }> = {
    admin: {
        label: "Administrator",
        class: "bg-red-50 text-red-700 border-red-200/60 dark:bg-red-500/10 dark:text-red-400 dark:border-red-500/20",
    },
    manager: {
        label: "Manager",
        class: "bg-blue-50 text-blue-700 border-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:border-blue-500/20",
    },
    editor: {
        label: "Editor",
        class: "bg-green-50 text-green-700 border-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:border-green-500/20",
    },
    viewer: {
        label: "Viewer",
        class: "bg-slate-100 text-slate-600 border-slate-200/60 dark:bg-white/[0.05] dark:text-slate-400 dark:border-white/[0.08]",
    },
};

interface ProfileTabProps {
    currentUser: User | null;
}

export function ProfileTab({ currentUser }: ProfileTabProps) {
    const role = currentUser?.role ?? "viewer";
    const cfg = roleConfig[role] ?? roleConfig.viewer;
    const initials = currentUser?.name?.charAt(0)?.toUpperCase() ?? "?";

    return (
        <div className={cn(
            "overflow-hidden rounded-2xl border",
            "bg-white border-slate-200/80",
            "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
        )}>
            {/* Card header */}
            <div className={cn(
                "px-6 py-4 border-b",
                "border-slate-100 bg-slate-50/60",
                "dark:border-white/[0.05] dark:bg-white/[0.02]"
            )}>
                <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                    Profile Information
                </p>
                <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                    Informasi akun Anda
                </p>
            </div>

            {/* Card body */}
            <div className="px-6 py-5 space-y-5">
                {/* Avatar row */}
                <div className={cn(
                    "flex items-center gap-4 pb-5 border-b",
                    "border-slate-100 dark:border-white/[0.05]"
                )}>
                    <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl bg-gradient-to-br from-green-500 to-emerald-600 dark:from-purple-600 dark:to-violet-700 text-white text-2xl font-semibold shadow-sm">
                        {initials}
                    </div>
                    <div>
                        <p className="text-sm font-semibold text-slate-900 dark:text-white">
                            {currentUser?.name}
                        </p>
                        <p className="text-xs text-slate-400 dark:text-slate-500 mt-0.5">
                            {currentUser?.email}
                        </p>
                        <span className={cn(
                            "mt-2 inline-flex items-center rounded-md border px-2 py-0.5 text-[10px] font-medium",
                            cfg.class
                        )}>
                            {cfg.label}
                        </span>
                    </div>
                </div>

                {/* Fields */}
                <div className="grid gap-4 sm:grid-cols-2">
                    {[
                        { label: "Name", value: currentUser?.name },
                        { label: "Email", value: currentUser?.email },
                    ].map(({ label, value }) => (
                        <div key={label} className="space-y-1.5">
                            <p className="text-[11px] font-medium uppercase tracking-wide text-slate-400 dark:text-slate-600">
                                {label}
                            </p>
                            <Input
                                defaultValue={value}
                                disabled
                                className={cn(
                                    "h-8 text-sm",
                                    "border-slate-200/80 bg-slate-50/80 text-slate-600",
                                    "dark:border-white/[0.06] dark:bg-white/[0.02] dark:text-slate-400",
                                    "disabled:opacity-100 disabled:cursor-default"
                                )}
                            />
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}