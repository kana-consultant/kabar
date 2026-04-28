interface SettingsHeaderProps {
    title: string;
    description: string;
}

export function SettingsHeader({ title, description }: SettingsHeaderProps) {
    return (
        <div>
            <h2 className="text-2xl font-bold tracking-tight">{title}</h2>
            <p className="text-slate-500">{description}</p>
        </div>
    );
}