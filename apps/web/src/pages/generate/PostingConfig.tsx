// components/PostingConfig.tsx
import { Input } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Switch } from "@kana-consultant/ui-kit";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter
} from "@kana-consultant/ui-kit";
import {
    Send,
    Calendar,
    FileText,
    Settings2,
    Clock,
    Loader2,
    AlertTriangle,
    X,
    CheckCircle2,
    RotateCcw,
    AlertCircle,
    XCircle
} from "lucide-react";
import { cn } from "@/lib/utils";
import type { Product } from "@/services/product";
import { useState } from "react";
import { ScrollArea } from "./scroll-area";

// ============ INTERFACE ============
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
    selectedProduct: string; // 👈 SINGLE STRING
    onSelectProduct: (productId: string | null) => void; // 👈 SINGLE SELECT HANDLER
    article: string;
    onPost: () => void;
    isPosting: boolean;
    isError?: boolean;
    results?: any[];
    errorData?: {
        title: string;
        message: string;
        errors: string[];
        results: any[];
        someFailed: boolean;
        allFailed: boolean;
    } | null;
    onRetry?: () => void;
    onCloseError?: () => void;
}

// ============ CONSTANTS ============
const modes = [
    { value: "instant", label: "Langsung", icon: Send },
    { value: "scheduled", label: "Terjadwal", icon: Calendar },
    { value: "draft", label: "Draft", icon: FileText },
] as const;

// ============ ERROR MODAL ============
function ErrorModal({
    open,
    onOpenChange,
    title,
    message,
    errors,
    results,
    someFailed,
    allFailed,
    onRetry,
    onClose,
}: {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    title: string;
    message: string;
    errors: string[];
    results: any[];
    someFailed: boolean;
    allFailed: boolean;
    onRetry: () => void;
    onClose: () => void;
}) {
    const totalProducts = results.length || 0;
    const successCount = results.filter((r: any) => r?.success === true).length || 0;
    const failedCount = results.filter((r: any) => r?.success === false).length || 0;

    const groupedErrors = results
        .filter((r: any) => r?.success === false)
        .reduce((acc: Record<string, Array<{ node: string; error: string; statusCode: number }>>, r: any) => {
            const productName = r?.product || 'Unknown Product';
            if (!acc[productName]) acc[productName] = [];
            acc[productName].push({
                node: r?.node || 'Unknown Node',
                error: r?.error || 'Unknown error',
                statusCode: r?.statusCode || 0,
            });
            return acc;
        }, {});

    const errorCount = Object.keys(groupedErrors).length;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className={cn(
                "sm:max-w-[640px] max-h-[90vh] p-0 overflow-hidden",
                "bg-white dark:bg-[#0f0d1a]",
                "border border-red-200/80 dark:border-red-500/20",
                "shadow-2xl"
            )}>
                <DialogHeader className={cn(
                    "px-6 py-4 border-b",
                    allFailed
                        ? "bg-red-50/60 border-red-200/60 dark:bg-red-500/10 dark:border-red-500/20"
                        : "bg-amber-50/60 border-amber-200/60 dark:bg-amber-500/10 dark:border-amber-500/20"
                )}>
                    <div className="flex items-start gap-3">
                        <div className={cn(
                            "flex h-10 w-10 items-center justify-center rounded-full flex-shrink-0",
                            allFailed
                                ? "bg-red-100 text-red-600 dark:bg-red-500/20 dark:text-red-400"
                                : "bg-amber-100 text-amber-600 dark:bg-amber-500/20 dark:text-amber-400"
                        )}>
                            <AlertTriangle className="h-5 w-5" />
                        </div>
                        <div className="flex-1 min-w-0">
                            <DialogTitle className="text-base font-semibold text-slate-800 dark:text-slate-100">
                                {title}
                            </DialogTitle>
                            <p className="text-sm text-slate-500 dark:text-slate-400 mt-0.5">
                                {message}
                            </p>
                        </div>
                        <button
                            onClick={onClose}
                            className="text-slate-400 hover:text-slate-600 dark:text-slate-500 dark:hover:text-slate-300 transition-colors"
                        >
                            <X className="h-4 w-4" />
                        </button>
                    </div>
                </DialogHeader>

                <ScrollArea className="flex-1 max-h-[420px] px-6 py-4">
                    <div className="space-y-4">
                        <div className="grid grid-cols-3 gap-2">
                            <div className={cn(
                                "rounded-lg p-3 text-center",
                                "bg-slate-50 border border-slate-200/60",
                                "dark:bg-white/[0.02] dark:border-white/[0.05]"
                            )}>
                                <p className="text-xl font-bold text-slate-800 dark:text-slate-200">
                                    {totalProducts}
                                </p>
                                <p className="text-[10px] uppercase tracking-wider text-slate-400">
                                    Total
                                </p>
                            </div>
                            <div className={cn(
                                "rounded-lg p-3 text-center",
                                "bg-green-50 border border-green-200/60",
                                "dark:bg-green-500/10 dark:border-green-500/20"
                            )}>
                                <p className="text-xl font-bold text-green-600 dark:text-green-400">
                                    {successCount}
                                </p>
                                <p className="text-[10px] uppercase tracking-wider text-green-600/70 dark:text-green-400/60">
                                    Sukses
                                </p>
                            </div>
                            <div className={cn(
                                "rounded-lg p-3 text-center",
                                "bg-red-50 border border-red-200/60",
                                "dark:bg-red-500/10 dark:border-red-500/20"
                            )}>
                                <p className="text-xl font-bold text-red-600 dark:text-red-400">
                                    {failedCount}
                                </p>
                                <p className="text-[10px] uppercase tracking-wider text-red-600/70 dark:text-red-400/60">
                                    Gagal
                                </p>
                            </div>
                        </div>

                        {errorCount > 0 && (
                            <div className="space-y-2">
                                <div className="flex items-center justify-between">
                                    <p className="text-[11px] font-medium uppercase tracking-wider text-slate-400 dark:text-slate-500">
                                        Detail Error ({errorCount} produk gagal)
                                    </p>
                                    <span className="text-[10px] px-2 py-0.5 rounded-full bg-red-100 text-red-700 dark:bg-red-500/20 dark:text-red-400">
                                        {failedCount} error
                                    </span>
                                </div>

                                {Object.entries(groupedErrors).map(([productName, productErrors]) => (
                                    <div
                                        key={productName}
                                        className={cn(
                                            "rounded-lg border p-3",
                                            "bg-red-50/50 border-red-200/60",
                                            "dark:bg-red-500/5 dark:border-red-500/20"
                                        )}
                                    >
                                        <div className="flex items-center gap-2 mb-2">
                                            <XCircle className="h-3.5 w-3.5 text-red-500 dark:text-red-400 flex-shrink-0" />
                                            <p className="text-xs font-medium text-red-800 dark:text-red-300 truncate">
                                                {productName}
                                            </p>
                                            <span className="ml-auto text-[10px] px-1.5 py-0.5 rounded bg-red-200/60 text-red-700 dark:bg-red-500/20 dark:text-red-400">
                                                {productErrors.length} error
                                            </span>
                                        </div>

                                        <div className="space-y-1.5 ml-5">
                                            {productErrors.map((err, idx) => (
                                                <div
                                                    key={idx}
                                                    className="text-[11px] text-red-700/80 dark:text-red-400/80"
                                                >
                                                    <span className="font-medium text-red-800 dark:text-red-300">
                                                        {err.node}:
                                                    </span>
                                                    <span className="ml-1">{err.error}</span>
                                                    {err.statusCode > 0 && (
                                                        <span className="ml-1.5 px-1.5 py-0.5 rounded text-[9px] bg-red-200/60 text-red-800/80 dark:bg-red-500/20 dark:text-red-400">
                                                            {err.statusCode}
                                                        </span>
                                                    )}
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                ))}

                                {errors && errors.length > 0 && Object.keys(groupedErrors).length === 0 && (
                                    <div className={cn(
                                        "rounded-lg border p-3",
                                        "bg-red-50/50 border-red-200/60",
                                        "dark:bg-red-500/5 dark:border-red-500/20"
                                    )}>
                                        {errors.map((err, idx) => (
                                            <div key={idx} className="text-xs text-red-700 dark:text-red-400 flex items-start gap-2">
                                                <AlertCircle className="h-3 w-3 mt-0.5 flex-shrink-0" />
                                                <span>{err}</span>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </div>
                        )}

                        {successCount > 0 && (
                            <div className={cn(
                                "rounded-lg border p-3",
                                "bg-green-50/50 border-green-200/60",
                                "dark:bg-green-500/5 dark:border-green-500/20"
                            )}>
                                <div className="flex items-center gap-2">
                                    <CheckCircle2 className="h-3.5 w-3.5 text-green-500 dark:text-green-400 flex-shrink-0" />
                                    <p className="text-xs text-green-700 dark:text-green-300">
                                        {successCount} produk berhasil diposting
                                    </p>
                                </div>
                            </div>
                        )}

                        {totalProducts === 0 && (
                            <div className="text-center py-8">
                                <AlertCircle className="h-12 w-12 text-slate-300 dark:text-slate-600 mx-auto mb-3" />
                                <p className="text-sm text-slate-500 dark:text-slate-400">
                                    Tidak ada data yang tersedia
                                </p>
                            </div>
                        )}
                    </div>
                </ScrollArea>

                <DialogFooter className={cn(
                    "px-6 py-4 border-t",
                    "bg-slate-50/60 border-slate-200/60",
                    "dark:bg-white/[0.02] dark:border-white/[0.05]"
                )}>
                    <div className="flex items-center gap-2 w-full sm:w-auto">
                        <Button
                            onClick={onClose}
                            variant="outline"
                            className="flex-1 sm:flex-none h-9 px-4 text-sm"
                        >
                            <X className="h-3.5 w-3.5 mr-1.5" />
                            Tutup
                        </Button>
                        {someFailed && !allFailed && (
                            <Button
                                onClick={onRetry}
                                className="flex-1 sm:flex-none h-9 px-4 text-sm bg-blue-600 hover:bg-blue-700 text-white"
                            >
                                <RotateCcw className="h-3.5 w-3.5 mr-1.5" />
                                Retry Gagal
                            </Button>
                        )}
                        {allFailed && (
                            <Button
                                onClick={() => {
                                    onClose();
                                    onRetry();
                                }}
                                className="flex-1 sm:flex-none h-9 px-4 text-sm bg-red-600 hover:bg-red-700 text-white"
                            >
                                <RotateCcw className="h-3.5 w-3.5 mr-1.5" />
                                Retry Semua
                            </Button>
                        )}
                    </div>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}

// ============ MAIN COMPONENT ============
export function PostingConfig({
    postMode, setPostMode,
    scheduleDate, setScheduleDate,
    scheduleTime, setScheduleTime,
    dailySchedule, setDailySchedule,
    dailyTime, setDailyTime,
    products,
    selectedProduct, // 👈 SINGLE STRING
    onSelectProduct, // 👈 SINGLE SELECT HANDLER
    article, onPost, isPosting,
    isError = false,
    results = [],
    errorData = null,
    onRetry,
    onCloseError,
}: PostingConfigProps) {

    const [errorModalOpen, setErrorModalOpen] = useState(false);

    const postLabel = postMode === "instant"
        ? "Post Sekarang"
        : postMode === "scheduled"
            ? dailySchedule ? "Jadwalkan Harian" : "Jadwalkan"
            : "Simpan Draft";

    const hasProductError = (productId: string) => {
        return results.some(r => {
            const rProductId = r?.product?.toString() || '';
            return rProductId === productId && r?.success === false;
        });
    };

    const selectedProductData = products.find(p => p.id?.toString() === selectedProduct);

    const handleShowError = () => {
        if (errorData) setErrorModalOpen(true);
    };

    const handleRetry = () => {
        setErrorModalOpen(false);
        if (onRetry) onRetry();
    };

    const handleCloseError = () => {
        setErrorModalOpen(false);
        if (onCloseError) onCloseError();
    };

    return (
        <>
            <div className={cn(
                "overflow-hidden rounded-2xl border",
                "bg-white border-slate-200/80",
                `${isError && "border-red-300 dark:border-red-500/30"}`,
                "dark:bg-[#0f0d1a] dark:border-white/[0.06]",
            )}>
                {/* Header */}
                <div className={cn(
                    "flex items-center gap-3 px-5 py-4 border-b",
                    isError
                        ? "border-red-200 bg-red-50/60 dark:border-red-500/20 dark:bg-red-500/5"
                        : "border-slate-100 bg-slate-50/60 dark:border-white/[0.05] dark:bg-white/[0.02]"
                )}>
                    <div className={cn(
                        "flex h-8 w-8 items-center justify-center rounded-lg ring-1",
                        isError
                            ? "bg-red-100 text-red-600 ring-red-200/60 dark:bg-red-500/20 dark:text-red-400 dark:ring-red-500/20"
                            : "bg-blue-50 text-blue-600 ring-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:ring-blue-500/20"
                    )}>
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
                    {isError && (
                        <button
                            onClick={handleShowError}
                            className="ml-auto flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-red-100 text-red-700 hover:bg-red-200 dark:bg-red-500/20 dark:text-red-400 dark:hover:bg-red-500/30 transition-colors"
                        >
                            <AlertTriangle className="h-3 w-3" />
                            <span className="text-[10px] font-medium">Lihat Error</span>
                        </button>
                    )}
                </div>

                {/* Content */}
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
                                            ? "bg-white text-slate-800 shadow-sm border border-slate-200/80 dark:bg-white/[0.08] dark:text-white dark:border-white/[0.10]"
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
                                        onChange={(e: any) => setDailyTime(e.target.value)}
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
                                            onChange={(e: any) => setScheduleDate(e.target.value)}
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
                                            onChange={(e: any) => setScheduleTime(e.target.value)}
                                            className="h-8 text-sm rounded-lg border-slate-200/80 dark:border-white/[0.08]"
                                        />
                                    </div>
                                </div>
                            )}
                        </div>
                    )}
                    <div>

                        {/* Products - SINGLE SELECT GRID */}
                        <div>
                            <div className="flex items-center justify-between mb-2.5">
                                <p className="text-[11px] font-medium uppercase tracking-wide text-slate-400 dark:text-slate-600">
                                    Pilih Produk
                                </p>
                                <span className="text-[10px] text-slate-400 dark:text-slate-500">
                                    {selectedProduct ? '1 produk dipilih' : 'Pilih 1 produk'}
                                </span>
                            </div>

                            <div className="h-[180px] overflow-y-auto rounded-xl border border-slate-200/80 dark:border-white/[0.06]">
                                <div className="p-2">
                                    {products.length === 0 ? (
                                        <div className="py-8 text-center">
                                            <p className="text-xs text-slate-400 dark:text-slate-500">
                                                Tidak ada produk tersedia
                                            </p>
                                        </div>
                                    ) : (
                                        <div className="grid grid-cols-2 gap-1.5">
                                            {products.map((product) => {
                                                const productId = product.id?.toString() || '';
                                                const isSelected = selectedProduct === productId;
                                                const isErrorProduct = hasProductError(productId);

                                                return (
                                                    <button
                                                        key={productId}
                                                        type="button"
                                                        onClick={() => {
                                                            if (isSelected) {
                                                                onSelectProduct(null);
                                                            } else {
                                                                onSelectProduct(productId);
                                                            }
                                                        }}
                                                        className={cn(
                                                            "rounded-lg px-3 py-2.5 transition-all cursor-pointer text-left border",
                                                            `${isSelected && !isErrorProduct && "border-blue-500 bg-blue-50 dark:bg-blue-500/10"}`,
                                                            `${isErrorProduct && isSelected && "border-red-400 bg-red-50 dark:bg-red-500/10"}`,
                                                            `${!isSelected && "border-slate-200 dark:border-white/[0.08]"}`
                                                        )}
                                                    >
                                                        <p className={cn(
                                                            "text-xs font-medium truncate",
                                                            isSelected && !isErrorProduct
                                                                ? "text-blue-700 dark:text-blue-300"
                                                                : isErrorProduct && isSelected
                                                                    ? "text-red-700 dark:text-red-300"
                                                                    : "text-slate-700 dark:text-slate-300",
                                                        )}>
                                                            {product.name || 'Unnamed Product'}
                                                        </p>

                                                        {/* Selected indicator - subtle top-right accent */}
                                                        {isSelected && !isErrorProduct && (
                                                            <div className="absolute top-0 right-0 w-0 h-0 border-t-[8px] border-r-[8px] border-t-blue-500 border-r-transparent rounded-tr-lg" />
                                                        )}

                                                        {isErrorProduct && isSelected && (
                                                            <div className="absolute top-0 right-0 w-0 h-0 border-t-[8px] border-r-[8px] border-t-red-400 border-r-transparent rounded-tr-lg" />
                                                        )}
                                                    </button>
                                                );
                                            })}
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>

                        <div className="mt-2 flex items-center justify-between">
                            <p className="text-[10px] text-slate-400 dark:text-slate-500">
                                <span className="font-medium text-blue-600 dark:text-blue-400">●</span> Pilih 1 produk untuk diposting
                            </p>
                            {selectedProduct && (
                                <button
                                    type="button"
                                    onClick={() => onSelectProduct(null)}
                                    className="text-[10px] font-medium text-red-500 hover:text-red-600 dark:text-red-400 dark:hover:text-red-300 transition-colors cursor-pointer"
                                >
                                    Hapus Pilihan
                                </button>
                            )}
                        </div>
                    </div>

                    {/* Selected Product Box */}
                    {selectedProductData && (
                        <div className={cn(
                            "rounded-xl border p-3.5",
                            isError && hasProductError(selectedProduct)
                                ? "bg-red-50/60 border-red-200/60 dark:bg-red-500/5 dark:border-red-500/20"
                                : "bg-blue-50/60 border-blue-200/60 dark:bg-blue-500/[0.04] dark:border-blue-500/20"
                        )}>
                            <div className="flex items-start gap-3">
                                <div className={cn(
                                    "flex h-8 w-8 items-center justify-center rounded-lg ring-1 flex-shrink-0",
                                    isError && hasProductError(selectedProduct)
                                        ? "bg-red-100 text-red-600 ring-red-300/60 dark:bg-red-500/10 dark:text-red-400 dark:ring-red-500/20"
                                        : "bg-blue-100 text-blue-600 ring-blue-300/60 dark:bg-blue-500/10 dark:text-blue-400 dark:ring-blue-500/20"
                                )}>
                                    {isError && hasProductError(selectedProduct) ? (
                                        <AlertTriangle className="h-4 w-4" />
                                    ) : (
                                        <CheckCircle2 className="h-4 w-4" />
                                    )}
                                </div>
                                <div className="flex-1 min-w-0">
                                    <p className={cn(
                                        "text-sm font-medium",
                                        isError && hasProductError(selectedProduct)
                                            ? "text-red-800 dark:text-red-300"
                                            : "text-blue-800 dark:text-blue-300"
                                    )}>
                                        {selectedProductData.name || 'Produk Terpilih'}
                                    </p>
                                </div>
                                {isError && hasProductError(selectedProduct) ? (
                                    <span className="text-[10px] px-2 py-0.5 rounded-full bg-red-100 text-red-700 dark:bg-red-500/20 dark:text-red-400 flex-shrink-0">
                                        Gagal
                                    </span>
                                ) : (
                                    <span className="text-[10px] px-2 py-0.5 rounded-full bg-green-100 text-green-700 dark:bg-green-500/20 dark:text-green-400 flex-shrink-0">
                                        Siap Posting
                                    </span>
                                )}
                            </div>
                        </div>
                    )}

                    {/* Empty state */}
                    {!selectedProductData && (
                        <div className={cn(
                            "rounded-xl border p-3.5 text-center",
                            "bg-slate-50/60 border-slate-200/60",
                            "dark:bg-white/[0.02] dark:border-white/[0.05]"
                        )}>
                            <p className="text-xs text-slate-400 dark:text-slate-500">
                                Pilih 1 produk untuk diposting
                            </p>
                        </div>
                    )}

                    {/* Post button */}
                    <div className="space-y-2 pt-1">
                        <Button
                            onClick={onPost}
                            disabled={!selectedProduct || !article || isPosting}
                            className={cn(
                                "w-full h-9 gap-2 rounded-lg text-sm font-medium",
                                isError
                                    ? "bg-red-600 hover:bg-red-700 text-white shadow-sm dark:bg-red-600 dark:hover:bg-red-700"
                                    : "bg-green-600 hover:bg-green-700 text-white shadow-sm dark:bg-purple-600 dark:hover:bg-purple-700",
                                "disabled:opacity-40"
                            )}
                        >
                            {isPosting
                                ? <><Loader2 className="h-3.5 w-3.5 animate-spin" /> Memproses...</>
                                : <>
                                    {isError ? (
                                        <RotateCcw className="h-3.5 w-3.5" />
                                    ) : (
                                        <Send className="h-3.5 w-3.5" />
                                    )}
                                    {isError ? 'Retry Posting' : postLabel}
                                    {selectedProductData && ` ke ${selectedProductData.name}`}
                                </>
                            }
                        </Button>

                        {!article && (
                            <p className="flex items-center justify-center gap-1.5 text-[11px] text-amber-600 dark:text-amber-400">
                                <AlertTriangle className="h-3 w-3" />
                                Generate artikel terlebih dahulu sebelum posting
                            </p>
                        )}

                        {selectedProduct && !isError && article && (
                            <p className="text-center text-[10px] text-slate-400 dark:text-slate-500">
                                Akan diposting ke 1 produk
                            </p>
                        )}
                    </div>
                </div>
            </div>

            {/* Error Modal */}
            {errorData && (
                <ErrorModal
                    open={errorModalOpen}
                    onOpenChange={setErrorModalOpen}
                    title={errorData.title}
                    message={errorData.message}
                    errors={errorData.errors}
                    results={errorData.results}
                    someFailed={errorData.someFailed}
                    allFailed={errorData.allFailed}
                    onRetry={handleRetry}
                    onClose={handleCloseError}
                />
            )}
        </>
    );
}