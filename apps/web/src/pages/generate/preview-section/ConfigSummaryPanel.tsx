import { cn } from "@/lib/utils";

interface ConfigSummaryPanelProps {
    postMode: "instant" | "scheduled" | "draft";
    keywords: string[];
    excerpt: string;
    slug: string;
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
    keywords,
    excerpt,
    slug,
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
        <div className="space-y-4">
            {/* SEO Information Section */}
            <div
                className={cn(
                    "rounded-xl border p-5",
                    "bg-white border-slate-200/80",
                    "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
                )}
            >
                <p className="text-sm font-medium text-slate-800 dark:text-slate-100 mb-3">
                    Informasi SEO
                </p>
                <div className="space-y-3">
                    {/* Slug */}
                    <div>
                        <span className="text-xs text-slate-400 dark:text-slate-600 block mb-1">
                            Slug
                        </span>
                        <p className="text-sm text-slate-700 dark:text-slate-300 font-mono bg-slate-50 dark:bg-white/[0.02] rounded-md px-3 py-2 border border-slate-200/80 dark:border-white/[0.06]">
                            {slug || "Belum ada slug"}
                        </p>
                    </div>

                    {/* Keywords */}
                    <div>
                        <span className="text-xs text-slate-400 dark:text-slate-600 block mb-1">
                            Keywords
                        </span>
                        <div className="flex flex-wrap gap-1.5">
                            {keywords?.length > 0 ? (
                                keywords.map((keyword, index) => (
                                    <span
                                        key={index}
                                        className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-50 text-purple-700 dark:bg-purple-950/30 dark:text-purple-300 border border-purple-200/60 dark:border-purple-800/50"
                                    >
                                        {keyword}
                                    </span>
                                ))
                            ) : (
                                <p className="text-sm text-slate-500 dark:text-slate-500 italic">
                                    Belum ada keywords
                                </p>
                            )}
                        </div>
                    </div>

                    {/* Excerpt */}
                    <div>
                        <span className="text-xs text-slate-400 dark:text-slate-600 block mb-1">
                            Excerpt
                        </span>
                        <p className="text-sm text-slate-700 dark:text-slate-300 leading-relaxed bg-slate-50 dark:bg-white/[0.02] rounded-md px-3 py-2 border border-slate-200/80 dark:border-white/[0.06]">
                            {excerpt || "Belum ada excerpt"}
                        </p>
                    </div>
                </div>
            </div>

            {/* Configuration Summary Section */}
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
        </div>
    );
}