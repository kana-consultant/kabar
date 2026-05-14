import { useState } from "react";
import { Switch } from "@/components/ui/switch";
import { Moon, Bell, Save } from "lucide-react";
import { cn } from "@/lib/utils";

const preferences = [
    {
        key: "darkMode",
        icon: Moon,
        label: "Dark Mode",
        desc: "Tampilan gelap untuk aplikasi",
    },
    {
        key: "emailNotifications",
        icon: Bell,
        label: "Email Notifications",
        desc: "Terima notifikasi via email",
    },
    {
        key: "autoSave",
        icon: Save,
        label: "Auto-save Draft",
        desc: "Simpan draft secara otomatis",
    },
] as const;

type PrefKey = (typeof preferences)[number]["key"];

export function PreferencesTab() {
    const [prefs, setPrefs] = useState<Record<PrefKey, boolean>>({
        darkMode: false,
        emailNotifications: true,
        autoSave: true,
    });

    const toggle = (key: PrefKey) =>
        setPrefs((p) => ({ ...p, [key]: !p[key] }));

    return (
        <div className={cn(
            "overflow-hidden rounded-2xl border",
            "bg-white border-slate-200/80",
            "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
        )}>
            {/* Header */}
            <div className={cn(
                "px-6 py-4 border-b",
                "border-slate-100 bg-slate-50/60",
                "dark:border-white/[0.05] dark:bg-white/[0.02]"
            )}>
                <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                    Preferences
                </p>
                <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                    Pengaturan tampilan dan notifikasi
                </p>
            </div>

            {/* Rows */}
            <div className="divide-y divide-slate-100 dark:divide-white/[0.05]">
                {preferences.map(({ key, icon: Icon, label, desc }) => (
                    <div
                        key={key}
                        className="flex items-center justify-between px-6 py-4 transition-colors hover:bg-slate-50/60 dark:hover:bg-white/[0.01]"
                    >
                        <div className="flex items-center gap-3">
                            <div className={cn(
                                "flex h-8 w-8 shrink-0 items-center justify-center rounded-lg",
                                prefs[key]
                                    ? "bg-green-50 text-green-600 ring-1 ring-green-200/60 dark:bg-purple-500/10 dark:text-purple-400 dark:ring-purple-500/20"
                                    : "bg-slate-100 text-slate-400 dark:bg-white/[0.04] dark:text-slate-600"
                            )}>
                                <Icon className="h-3.5 w-3.5" />
                            </div>
                            <div>
                                <p className="text-sm font-medium text-slate-800 dark:text-slate-200">
                                    {label}
                                </p>
                                <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                                    {desc}
                                </p>
                            </div>
                        </div>

                        <Switch
                            checked={prefs[key]}
                            onCheckedChange={() => toggle(key)}
                            className={cn(
                                "data-[state=checked]:bg-green-600",
                                "dark:data-[state=checked]:bg-purple-600"
                            )}
                        />
                    </div>
                ))}
            </div>
        </div>
    );
}