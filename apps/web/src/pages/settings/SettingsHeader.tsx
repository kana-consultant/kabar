import { cn } from "@/lib/utils";

interface SettingsHeaderProps {
    title: string;
    description: string;
}

export function SettingsHeader({ title, description }: SettingsHeaderProps) {
    return (
        <div>
            <h2 className="text-xl font-semibold tracking-tight text-slate-900 dark:text-white">
                {title}
            </h2>
            <p className="mt-0.5 text-sm text-slate-400 dark:text-slate-500">
                {description}
            </p>
        </div>
    );
}