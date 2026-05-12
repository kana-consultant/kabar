import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ImageIcon, FileText, Settings2 } from "lucide-react";
import { cn } from "@/lib/utils";

interface PreviewSectionProps {
    article: string;
    imageUrl: string;
    hasImage: boolean;
    postMode: "instant" | "scheduled" | "draft";
    dailySchedule: boolean;
    dailyTime: string;
    scheduleDate: string;
    scheduleTime: string;
    selectedProductsCount: number;
    autoGenerateImage: boolean;
}

export function PreviewSection({
    article, imageUrl, hasImage,
    postMode, dailySchedule, dailyTime,
    scheduleDate, scheduleTime,
    selectedProductsCount, autoGenerateImage,
}: PreviewSectionProps) {
    return (
        <Tabs defaultValue="article" className="w-full">
            <TabsList className={cn(
                "h-8 gap-0.5 rounded-lg border p-0.5",
                "bg-slate-50 border-slate-200/80",
                "dark:bg-white/[0.02] dark:border-white/[0.06]"
            )}>
                {[
                    { value: "article", label: "Preview Artikel", icon: FileText },
                    { value: "image", label: "Preview Gambar", icon: ImageIcon },
                    { value: "config", label: "Ringkasan", icon: Settings2 },
                ].map(({ value, label, icon: Icon }) => (
                    <TabsTrigger
                        key={value}
                        value={value}
                        className={cn(
                            "h-7 gap-1.5 rounded-md px-3 text-xs data-[state=active]:shadow-sm",
                            "data-[state=active]:bg-white data-[state=active]:text-slate-800",
                            "dark:data-[state=active]:bg-white/[0.08] dark:data-[state=active]:text-white"
                        )}
                    >
                        <Icon className="h-3 w-3" />
                        {label}
                    </TabsTrigger>
                ))}
            </TabsList>

            {/* Article */}
            <TabsContent value="article" className="mt-3">
                <div className={cn(
                    "rounded-xl border p-5 min-h-[200px]",
                    "bg-white border-slate-200/80",
                    "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
                )}>
                    {article ? (
                        <div
                            className="prose dark:prose-invert max-w-none prose-sm prose-headings:font-semibold prose-p:text-slate-600 dark:prose-p:text-slate-400"
                            dangerouslySetInnerHTML={{ __html: article }}
                        />
                    ) : (
                        <div className="flex flex-col items-center justify-center py-12 text-slate-400">
                            <FileText className="h-8 w-8 mb-2 opacity-30" />
                            <p className="text-sm">Belum ada artikel. Klik "Generate Artikel" dulu.</p>
                        </div>
                    )}
                </div>
            </TabsContent>

            {/* Image */}
            <TabsContent value="image" className="mt-3">
                <div className={cn(
                    "rounded-xl border p-5 min-h-[200px]",
                    "bg-white border-slate-200/80",
                    "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
                )}>
                    {imageUrl ? (
                        <div className="space-y-3">
                            <img
                                src={imageUrl}
                                alt="Preview"
                                className="w-full rounded-lg border border-slate-200/80 object-cover max-h-[400px] dark:border-white/[0.06]"
                                onError={(e) => {
                                    (e.target as HTMLImageElement).src = "https://placehold.co/800x400?text=Image+Failed";
                                }}
                            />
                            <p className="text-center text-xs text-slate-400">
                                Klik kanan → Save Image As untuk menyimpan
                            </p>
                        </div>
                    ) : (
                        <div className="flex flex-col items-center justify-center py-12 text-slate-400">
                            <ImageIcon className="h-8 w-8 mb-2 opacity-30" />
                            <p className="text-sm">{hasImage ? "Generate gambar terlebih dahulu" : "Belum ada gambar"}</p>
                        </div>
                    )}
                </div>
            </TabsContent>

            {/* Config summary */}
            <TabsContent value="config" className="mt-3">
                <div className={cn(
                    "rounded-xl border p-5",
                    "bg-white border-slate-200/80",
                    "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
                )}>
                    <p className="text-sm font-medium text-slate-800 dark:text-slate-100 mb-3">
                        Ringkasan Konfigurasi
                    </p>
                    <div className="space-y-2">
                        {[
                            {
                                label: "Mode",
                                value: postMode === "instant" ? "Langsung posting"
                                    : postMode === "scheduled" ? "Terjadwal"
                                    : "Simpan sebagai draft"
                            },
                            ...(postMode === "scheduled" ? [{
                                label: "Jadwal",
                                value: dailySchedule
                                    ? `Setiap hari jam ${dailyTime}`
                                    : `${scheduleDate} jam ${scheduleTime}`
                            }] : []),
                            { label: "Target", value: `${selectedProductsCount} produk` },
                            { label: "Auto gambar", value: autoGenerateImage ? "Ya" : "Tidak" },
                            ...(imageUrl ? [{ label: "Gambar", value: "Tersedia" }] : []),
                        ].map(({ label, value }) => (
                            <div key={label} className="flex items-center justify-between py-1.5 border-b border-slate-100 dark:border-white/[0.04] last:border-0">
                                <span className="text-xs text-slate-400 dark:text-slate-600">{label}</span>
                                <span className="text-xs font-medium text-slate-700 dark:text-slate-300">{value}</span>
                            </div>
                        ))}
                    </div>
                </div>
            </TabsContent>
        </Tabs>
    );
}