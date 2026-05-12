import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Send, Calendar, FileText, Settings2, Clock, Loader2, AlertTriangle } from "lucide-react";
import { TargetProducts } from "./TargetProducts";
import type { Product } from "@/types/product";
import { cn } from "@/lib/utils";

interface PostingConfigProps {
    postMode: "instant" | "scheduled" | "draft";
    setPostMode: (mode: "instant" | "scheduled" | "draft") => void;
    scheduleDate: string;
    setScheduleDate: (value: string) => void;
    scheduleTime: string;
    setScheduleTime: (value: string) => void;
    dailySchedule: boolean;
    setDailySchedule: (value: boolean) => void;
    dailyTime: string;
    setDailyTime: (value: string) => void;
    autoGenerateImage: boolean;
    setAutoGenerateImage: (value: boolean) => void;
    products: Product[];
    selectedProducts: string[];
    postToAll: boolean;
    onToggleProduct: (product: string) => void;
    onSelectAll: () => void;
    article: string;
    onPost: () => void;
    isPosting: boolean;
}

const modes = [
    { value: "instant", label: "Langsung", icon: Send },
    { value: "scheduled", label: "Terjadwal", icon: Calendar },
    { value: "draft", label: "Draft", icon: FileText },
] as const;

export function PostingConfig({
    postMode, setPostMode,
    scheduleDate, setScheduleDate,
    scheduleTime, setScheduleTime,
    dailySchedule, setDailySchedule,
    dailyTime, setDailyTime,
    products, selectedProducts, postToAll,
    onToggleProduct, onSelectAll,
    article, onPost, isPosting,
    autoGenerateImage, setAutoGenerateImage,
}: PostingConfigProps) {

    const postLabel = postMode === "instant"
        ? "Post Sekarang"
        : postMode === "scheduled"
            ? dailySchedule ? "Jadwalkan Harian" : "Jadwalkan"
            : "Simpan Draft";

    return (
        <div className={cn(
            "overflow-hidden rounded-2xl border",
            "bg-white border-slate-200/80",
            "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
        )}>
            {/* Header */}
            <div className={cn(
                "flex items-center gap-3 px-5 py-4 border-b",
                "border-slate-100 bg-slate-50/60",
                "dark:border-white/[0.05] dark:bg-white/[0.02]"
            )}>
                <div className="flex h-8 w-8 items-center justify-center rounded-lg ring-1 bg-blue-50 text-blue-600 ring-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:ring-blue-500/20">
                    <Settings2 className="h-3.5 w-3.5" />
                </div>
                <div>
                    <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                        Konfigurasi Posting
                    </p>
                    <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                        Atur jadwal dan target posting konten
                    </p>
                </div>
            </div>

            <div className="p-5 space-y-5">
                {/* Mode selector */}
                <div>
                    <p className="text-[11px] font-medium uppercase tracking-wide text-slate-400 dark:text-slate-600 mb-2.5">
                        Mode Posting
                    </p>
                    <div className={cn(
                        "flex gap-1 rounded-lg border p-1",
                        "bg-slate-50 border-slate-200/80",
                        "dark:bg-white/[0.02] dark:border-white/[0.06]"
                    )}>
                        {modes.map(({ value, label, icon: Icon }) => (
                            <button
                                key={value}
                                type="button"
                                onClick={() => setPostMode(value)}
                                className={cn(
                                    "flex flex-1 items-center justify-center gap-1.5 rounded-md py-1.5 text-xs font-medium transition-all",
                                    postMode === value
                                        ? [
                                            "bg-white text-slate-800 shadow-sm border border-slate-200/80",
                                            "dark:bg-white/[0.08] dark:text-white dark:border-white/[0.10]"
                                        ]
                                        : "text-slate-500 hover:text-slate-700 dark:text-slate-500 dark:hover:text-slate-300"
                                )}
                            >
                                <Icon className="h-3 w-3" />
                                {label}
                            </button>
                        ))}
                    </div>
                </div>

                {/* Schedule options */}
                {postMode === "scheduled" && (
                    <div className={cn(
                        "rounded-xl border p-4 space-y-3",
                        "bg-slate-50/60 border-slate-200/60",
                        "dark:bg-white/[0.02] dark:border-white/[0.05]"
                    )}>
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                                <Clock className="h-3.5 w-3.5 text-slate-400" />
                                <p className="text-xs font-medium text-slate-700 dark:text-slate-300">
                                    Posting Berulang (Daily)
                                </p>
                            </div>
                            <Switch
                                checked={dailySchedule}
                                onCheckedChange={setDailySchedule}
                                className="data-[state=checked]:bg-green-600 dark:data-[state=checked]:bg-purple-600"
                            />
                        </div>

                        {dailySchedule ? (
                            <div className="space-y-1.5">
                                <label className="text-[10px] uppercase tracking-wide text-slate-400">
                                    Waktu Posting Harian
                                </label>
                                <Input
                                    type="time"
                                    value={dailyTime}
                                    onChange={(e) => setDailyTime(e.target.value)}
                                    className="h-8 text-sm rounded-lg border-slate-200/80 dark:border-white/[0.08]"
                                />
                                <p className="text-[10px] text-slate-400">
                                    Setiap hari jam {dailyTime}
                                </p>
                            </div>
                        ) : (
                            <div className="grid grid-cols-2 gap-2">
                                <div className="space-y-1.5">
                                    <label className="text-[10px] uppercase tracking-wide text-slate-400">
                                        Tanggal
                                    </label>
                                    <Input
                                        type="date"
                                        value={scheduleDate}
                                        onChange={(e) => setScheduleDate(e.target.value)}
                                        className="h-8 text-sm rounded-lg border-slate-200/80 dark:border-white/[0.08]"
                                    />
                                </div>
                                <div className="space-y-1.5">
                                    <label className="text-[10px] uppercase tracking-wide text-slate-400">
                                        Waktu
                                    </label>
                                    <Input
                                        type="time"
                                        value={scheduleTime}
                                        onChange={(e) => setScheduleTime(e.target.value)}
                                        className="h-8 text-sm rounded-lg border-slate-200/80 dark:border-white/[0.08]"
                                    />
                                </div>
                            </div>
                        )}
                    </div>
                )}

                {/* Target products */}
                <TargetProducts
                    products={products}
                    selectedProducts={selectedProducts}
                    postToAll={postToAll}
                    onToggleProduct={onToggleProduct}
                    onSelectAll={onSelectAll}
                />

                {/* Post button */}
                <div className="space-y-2 pt-1">
                    <Button
                        onClick={onPost}
                        disabled={selectedProducts.length === 0 || !article || isPosting}
                        className={cn(
                            "w-full h-9 gap-2 rounded-lg text-sm font-medium",
                            "bg-green-600 hover:bg-green-700 text-white shadow-sm",
                            "dark:bg-purple-600 dark:hover:bg-purple-700",
                            "disabled:opacity-40"
                        )}
                    >
                        {isPosting
                            ? <><Loader2 className="h-3.5 w-3.5 animate-spin" /> Memproses...</>
                            : <><Send className="h-3.5 w-3.5" /> {postLabel} ke {selectedProducts.length} produk</>
                        }
                    </Button>

                    {!article && (
                        <p className="flex items-center justify-center gap-1.5 text-[11px] text-amber-600 dark:text-amber-400">
                            <AlertTriangle className="h-3 w-3" />
                            Generate artikel terlebih dahulu sebelum posting
                        </p>
                    )}
                </div>
            </div>
        </div>
    );
}