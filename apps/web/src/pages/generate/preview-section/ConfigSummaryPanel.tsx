import { cn } from "@/lib/utils";

interface ConfigSummaryPanelProps {
    postMode: "instant" | "scheduled" | "draft";
    dailySchedule: boolean;
    dailyTime: string;
    scheduleDate: string;
    scheduleTime: string;
    selectedProductsCount: number;
    autoGenerateImage: boolean;
    uploadedImage: string | null;
    imageUrl: string;
}

export function ConfigSummaryPanel({
    postMode,
    dailySchedule,
    dailyTime,
    scheduleDate,
    scheduleTime,
    selectedProductsCount,
    autoGenerateImage,
    uploadedImage,
    imageUrl,
}: ConfigSummaryPanelProps) {
    const rows = [
        {
            label: "Mode",
            value: postMode === "instant" ? "Langsung posting" : postMode === "scheduled" ? "Terjadwal" : "Simpan sebagai draft",
        },
        ...(postMode === "scheduled" ? [
            { label: "Jadwal", value: dailySchedule ? `Setiap hari jam ${dailyTime}` : `${scheduleDate} jam ${scheduleTime}` },
        ] : []),
        { label: "Target", value: `${selectedProductsCount} produk` },
        { label: "Auto gambar", value: autoGenerateImage ? "Ya" : "Tidak" },
        ...(uploadedImage ? [{ label: "Gambar", value: "Tersedia (Upload)" }] : []),
        ...(!uploadedImage && imageUrl ? [{ label: "Gambar", value: "Tersedia (Generated)" }] : []),
    ];

    return (
        <div
            className={cn(
                "rounded-xl border p-5",
                "bg-white border-slate-200/80",
                "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
            )}
        >
            <p className="text-sm font-medium text-slate-800 dark:text-slate-100 mb-3">
                Ringkasan Konfigurasi
            </p>
            <div className="space-y-2">
                {rows.map(({ label, value }) => (
                    <div
                        key={label}
                        className="flex items-center justify-between py-1.5 border-b border-slate-100 dark:border-white/[0.04] last:border-0"
                    >
                        <span className="text-xs text-slate-400 dark:text-slate-600">
                            {label}
                        </span>
                        <span className="text-xs font-medium text-slate-700 dark:text-slate-300">
                            {value}
                        </span>
                    </div>
                ))}
            </div>
        </div>
    );
}
