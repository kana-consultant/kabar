// components/PostingConfig.tsx
import { Input } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Switch } from "@kana-consultant/ui-kit";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
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
    XCircle,
    Copy,
    ExternalLink
} from "lucide-react";
import { cn } from "@/lib/utils";
import type { Product } from "@/services/product";
import { useMemo, useState } from "react";
import { ScrollArea } from "./scroll-area";

// ============ INTERFACES ============
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
    selectedProduct: string;
    onSelectProduct: (productId: string | null) => void;
    article: string;
    onPost: () => void;
    isPosting: boolean;
    isError?: boolean;
    results?: PostingResult[];
    errorData?: ErrorData | null;
    onRetry?: () => void;
    onCloseError?: () => void;
}

interface PostingResult {
    success: boolean;
    product: string;
    node?: string;
    error?: string;
    statusCode?: number;
}

interface ErrorData {
    title: string;
    message: string;
    errors: string[];
    results: PostingResult[];
}

interface GroupedError {
    node: string;
    error: string;
    statusCode: number;
}

interface ErrorModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    title?: string;
    message?: string;
    results?: PostingResult[];
    errors?: string[];
    onRetry?: () => void;
    onClose?: () => void;
}

// ============ CONSTANTS ============
const modes = [
    { value: "instant", label: "Langsung", icon: Send },
    { value: "scheduled", label: "Terjadwal", icon: Calendar },
    { value: "draft", label: "Draft", icon: FileText },
] as const;

type ParsedFailure = {
    success: false;
    product: string;
    node: string;
    error: string;
    statusCode: number;
};

// ============ UTILITY FUNCTIONS ============
function parseBackendErrorMessage(message: string): ParsedFailure[] {
    if (!message) return [];

    const results: ParsedFailure[] = [];

    // Split by "Product " to get individual product blocks
    const blocks = message.split(/Product\s+/);

    for (const block of blocks) {
        if (!block.trim()) continue;

        // Extract product ID
        const productMatch = block.match(/^([\w-]+)\s+failed:/);
        if (!productMatch) continue;

        const productId = productMatch[1];

        // Extract node errors - improved regex
        const nodePattern = /Node\s+([\w-]+):\s*request failed with status\s+(\d+)[^:]*:\s*(.+?)(?=(?:Node\s+[\w-]+:)|$)/gs;
        let nodeMatch;
        let hasNodes = false;

        while ((nodeMatch = nodePattern.exec(block)) !== null) {
            hasNodes = true;
            const nodeId = nodeMatch[1];
            const statusCode = Number(nodeMatch[2]);
            let errorText = nodeMatch[3].trim();

            // Clean up error message
            // Remove HTTP prefix and attempt info
            errorText = errorText.replace(/^HTTP\s+\d+\s*(?:\(attempt\s*\d+\/\d+\))?:\s*/, '').trim();

            // Extract JSON message if present
            if (errorText.includes('{')) {
                try {
                    const jsonMatch = errorText.match(/\{[^}]*"message"\s*:\s*"([^"]+)"[^}]*\}/);
                    if (jsonMatch) {
                        errorText = jsonMatch[1];
                    }
                } catch {
                    // Keep original if parsing fails
                }
            }

            results.push({
                success: false,
                product: productId,
                node: nodeId,
                error: errorText,
                statusCode,
            });
        }

        // If no nodes found, extract what we can
        if (!hasNodes) {
            const anyNodeMatch = block.match(/Node\s+([\w-]+)/);
            const statusMatch = block.match(/status\s+(\d+)/);
            const messageMatch = block.match(/"message"\s*:\s*"([^"]+)"/);

            results.push({
                success: false,
                product: productId,
                node: anyNodeMatch ? anyNodeMatch[1] : 'Unknown Node',
                error: messageMatch ? messageMatch[1] : 'Unknown error',
                statusCode: statusMatch ? Number(statusMatch[1]) : 0,
            });
        }
    }

    return results.length > 0 ? results : [];
}

// ============ ERROR MODAL COMPONENT ============
export function ErrorModal({
    open,
    onOpenChange,
    title = 'Posting Gagal',
    message = '',
    results = [],
    errors = [],
    onRetry = () => { },
    onClose = () => { },
}: ErrorModalProps) {
    const [copiedId, setCopiedId] = useState<string | null>(null);

    // Derive structured rows from results or fallback parsing
    const effectiveResults = useMemo(() => {
        if (results && results.length > 0) {
            return results;
        }

        const fromMessage = parseBackendErrorMessage(message);
        if (fromMessage.length > 0) {
            return fromMessage;
        }

        if (errors && errors.length > 0) {
            const fromErrors = errors.flatMap((e) => parseBackendErrorMessage(e));
            if (fromErrors.length > 0) {
                return fromErrors;
            }
        }

        return results || [];
    }, [results, message, errors]);

    const totalProducts = new Set(effectiveResults.map((r) => r?.product)).size;
    const successCount = effectiveResults.filter((r) => r?.success === true).length || 0;
    const failedCount = effectiveResults.filter((r) => r?.success === false).length || 0;

    // Derived booleans based on actual data
    const someFailed = failedCount > 0;
    const allFailed = failedCount > 0 && successCount === 0 && totalProducts > 0;

    const groupedErrors = effectiveResults
        .filter((r) => r?.success === false)
        .reduce<Record<string, GroupedError[]>>((acc, r) => {
            const productName = r?.product || 'Unknown Product';
            if (!acc[productName]) {
                acc[productName] = [];
            }
            acc[productName].push({
                node: r?.node || 'Unknown Node',
                error: r?.error || 'Unknown error',
                statusCode: r?.statusCode || 0,
            });
            return acc;
        }, {});

    const errorCount = Object.keys(groupedErrors).length;
    const showRawErrorsFallback = errorCount === 0 && errors && errors.length > 0;

    const handleOpenChange = (newOpen: boolean) => {
        if (!newOpen) {
            onClose();
        }
        onOpenChange(newOpen);
    };

    const handleRetry = () => {
        onRetry();
        if (allFailed) {
            onClose();
        }
    };

    const handleCopyError = async (text: string, id: string) => {
        await navigator.clipboard.writeText(text);
        setCopiedId(id);
        setTimeout(() => setCopiedId(null), 2000);
    };

    const getStatusColor = (code: number) => {
        if (code >= 500) return "text-red-600 dark:text-red-400 bg-red-100 dark:bg-red-500/20";
        if (code >= 400) return "text-orange-600 dark:text-orange-400 bg-orange-100 dark:bg-orange-500/20";
        if (code >= 300) return "text-yellow-600 dark:text-yellow-400 bg-yellow-100 dark:bg-yellow-500/20";
        return "text-slate-600 dark:text-slate-400 bg-slate-100 dark:bg-slate-500/20";
    };

    return (
        <Dialog open={open} onOpenChange={handleOpenChange}>
            <DialogContent
                className={cn(
                    'sm:max-w-[680px] max-h-[90vh] p-0 overflow-hidden',
                    'bg-white dark:bg-[#0f0d1a]',
                    'border shadow-2xl',
                    allFailed
                        ? 'border-red-300 dark:border-red-500/30'
                        : 'border-orange-200/80 dark:border-orange-500/20'
                )}
            >
                {/* Header */}
                <DialogHeader className={cn(
                    "px-6 pt-6 pb-4 border-b",
                    allFailed
                        ? "border-red-200/60 dark:border-red-500/20"
                        : "border-orange-200/60 dark:border-orange-500/20"
                )}>
                    <div className="flex items-start gap-3">
                        <div className={cn(
                            "flex h-10 w-10 items-center justify-center rounded-xl ring-1 flex-shrink-0",
                            allFailed
                                ? "bg-red-100 text-red-600 ring-red-200/60 dark:bg-red-500/20 dark:text-red-400 dark:ring-red-500/20"
                                : "bg-orange-100 text-orange-600 ring-orange-200/60 dark:bg-orange-500/20 dark:text-orange-400 dark:ring-orange-500/20"
                        )}>
                            <XCircle className="h-5 w-5" />
                        </div>
                        <div className="flex-1 min-w-0">
                            <DialogTitle className={cn(
                                "text-lg font-semibold",
                                allFailed ? "text-red-600 dark:text-red-400" : "text-orange-600 dark:text-orange-400"
                            )}>
                                {title}
                            </DialogTitle>

                        </div>
                    </div>

                    {/* Stats */}

                </DialogHeader>

                {/* Error Details */}
                <ScrollArea className="flex-1 max-h-[400px]">
                    <div className="px-6 py-4 space-y-4">
                        {errorCount > 0 && (
                            <div className="space-y-3">
                                <div className="flex items-center gap-2">
                                    <p className="text-[11px] font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-500">
                                        Detail Error
                                    </p>
                                    <span className="text-[10px] px-1.5 py-0.5 rounded-full bg-red-100 text-red-700 dark:bg-red-500/20 dark:text-red-400 font-medium">
                                        {errorCount} produk
                                    </span>
                                </div>

                                {Object.entries(groupedErrors).map(([productName, productErrors]) => (
                                    <div
                                        key={productName}
                                        className={cn(
                                            'rounded-xl border overflow-hidden',
                                            'bg-white dark:bg-white/[0.02]',
                                            'border-red-200/60 dark:border-red-500/20',
                                            'shadow-sm'
                                        )}
                                    >
                                        {/* Product Header */}
                                        <div className={cn(
                                            "flex items-center gap-2 px-4 py-2.5",
                                            "bg-red-50/60 dark:bg-red-500/5",
                                            "border-b border-red-200/40 dark:border-red-500/10"
                                        )}>
                                            <XCircle className="h-3.5 w-3.5 text-red-500 dark:text-red-400 flex-shrink-0" />
                                            <div className="flex-1 min-w-0">
                                                <p className="text-xs font-semibold text-red-800 dark:text-red-300 truncate">
                                                    {productName}
                                                </p>
                                            </div>
                                            <span className="text-[10px] px-2 py-0.5 rounded-full bg-red-200/60 text-red-700 dark:bg-red-500/20 dark:text-red-400 font-medium">
                                                {productErrors.length} error{productErrors.length > 1 ? 's' : ''}
                                            </span>
                                        </div>

                                        {/* Error Items */}
                                        <div className="divide-y divide-red-100/60 dark:divide-red-500/10">
                                            {productErrors.map((err, idx) => {
                                                const copyId = `${productName}-${idx}`;
                                                return (
                                                    <div
                                                        key={idx}
                                                        className="px-4 py-3 hover:bg-red-50/30 dark:hover:bg-red-500/5 transition-colors"
                                                    >
                                                        <div className="flex items-start gap-3">
                                                            <div className="flex-1 min-w-0 space-y-1.5">
                                                                {/* Node ID */}
                                                                <div className="flex items-center gap-2">
                                                                    <span className="text-[10px] font-mono text-slate-400 dark:text-slate-500">
                                                                        Node:
                                                                    </span>
                                                                    <code className="text-[11px] font-mono text-slate-700 dark:text-slate-300 bg-slate-100 dark:bg-slate-800 px-1.5 py-0.5 rounded truncate">
                                                                        {err.node}
                                                                    </code>
                                                                    <button
                                                                        onClick={() => handleCopyError(err.node, copyId)}
                                                                        className="flex-shrink-0 p-0.5 hover:bg-slate-200 dark:hover:bg-slate-700 rounded transition-colors"
                                                                        title="Copy node ID"
                                                                    >
                                                                        {copiedId === copyId ? (
                                                                            <CheckCircle2 className="h-3 w-3 text-green-500" />
                                                                        ) : (
                                                                            <Copy className="h-3 w-3 text-slate-400" />
                                                                        )}
                                                                    </button>
                                                                </div>

                                                                {/* Error Message */}
                                                                <div className="pl-4 border-l-2 border-red-200 dark:border-red-500/20">
                                                                    <p className="text-xs text-red-700 dark:text-red-400 leading-relaxed">
                                                                        {err.error}
                                                                    </p>
                                                                </div>

                                                                {/* Status Code */}
                                                                {err.statusCode > 0 && (
                                                                    <div className="flex items-center gap-2">
                                                                        <span className="text-[10px] text-slate-400 dark:text-slate-500">
                                                                            Status:
                                                                        </span>
                                                                        <span className={cn(
                                                                            "text-[10px] px-1.5 py-0.5 rounded font-mono font-medium",
                                                                            getStatusColor(err.statusCode)
                                                                        )}>
                                                                            {err.statusCode}
                                                                        </span>
                                                                        <span className="text-[10px] text-slate-400 dark:text-slate-500">
                                                                            {err.statusCode >= 500 ? 'Server Error' :
                                                                                err.statusCode >= 400 ? 'Client Error' :
                                                                                    err.statusCode >= 300 ? 'Redirect' : 'Unknown'}
                                                                        </span>
                                                                    </div>
                                                                )}
                                                            </div>
                                                        </div>
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}

                        {/* Raw errors fallback */}
                        {showRawErrorsFallback && (
                            <div className={cn(
                                'rounded-xl border p-4 space-y-3',
                                'bg-red-50/50 border-red-200/60',
                                'dark:bg-red-500/5 dark:border-red-500/20'
                            )}>
                                <p className="text-[11px] font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-500">
                                    Error Messages
                                </p>
                                {errors.map((err, idx) => (
                                    <div
                                        key={idx}
                                        className="flex items-start gap-2 text-xs text-red-700 dark:text-red-400 bg-white dark:bg-white/[0.02] rounded-lg p-3 border border-red-100 dark:border-red-500/10"
                                    >
                                        <AlertCircle className="h-3.5 w-3.5 mt-0.5 flex-shrink-0" />
                                        <span className="break-words leading-relaxed">{err}</span>
                                    </div>
                                ))}
                            </div>
                        )}

                        {/* Success summary */}
                        {successCount > 0 && (
                            <div className={cn(
                                'rounded-xl border p-4',
                                'bg-green-50/50 border-green-200/60',
                                'dark:bg-green-500/5 dark:border-green-500/20'
                            )}>
                                <div className="flex items-center gap-2">
                                    <CheckCircle2 className="h-4 w-4 text-green-500 dark:text-green-400 flex-shrink-0" />
                                    <p className="text-sm font-medium text-green-700 dark:text-green-300">
                                        {successCount} produk berhasil diposting
                                    </p>
                                </div>
                            </div>
                        )}

                        {/* Empty state */}
                        {totalProducts === 0 && errorCount === 0 && !showRawErrorsFallback && (
                            <div className="text-center py-12">
                                <AlertCircle className="h-12 w-12 text-slate-300 dark:text-slate-600 mx-auto mb-3" />
                                <p className="text-sm text-slate-500 dark:text-slate-400">
                                    Tidak ada detail error yang tersedia
                                </p>
                                <p className="text-xs text-slate-400 dark:text-slate-500 mt-1">
                                    Silakan coba lagi atau hubungi administrator
                                </p>
                            </div>
                        )}
                    </div>
                </ScrollArea>

                {/* Footer */}
                <DialogFooter
                    className={cn(
                        'px-6 py-4 border-t',
                        'bg-slate-50/60 border-slate-200/60',
                        'dark:bg-white/[0.02] dark:border-white/[0.05]'
                    )}
                >
                    <div className="flex items-center gap-2 w-full">
                        <Button
                            onClick={onClose}
                            variant="outline"
                            className="flex-1 h-9 px-4 text-sm"
                        >
                            <X className="h-3.5 w-3.5 mr-1.5" />
                            Tutup
                        </Button>
                        {someFailed && !allFailed && (
                            <Button
                                onClick={onRetry}
                                className="flex-1 h-9 px-4 text-sm bg-blue-600 hover:bg-blue-700 text-white shadow-sm"
                            >
                                <RotateCcw className="h-3.5 w-3.5 mr-1.5" />
                                Retry yang Gagal
                            </Button>
                        )}
                        {allFailed && (
                            <Button
                                onClick={handleRetry}
                                className="flex-1 h-9 px-4 text-sm bg-red-600 hover:bg-red-700 text-white shadow-sm"
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
    autoGenerateImage, setAutoGenerateImage,
    products,
    selectedProduct,
    onSelectProduct,
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
    const isProductError = selectedProduct ? hasProductError(selectedProduct) : false;

    const handleShowError = () => {
        console.log(errorData)
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
                                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setDailyTime(e.target.value)}
                                        className="h-8 text-sm rounded-lg border-slate-200/80 dark:border-white/[0.08]"
                                    />
                                    <p className="text-[10px] text-slate-400">
                                        Setiap hari jam {dailyTime || '00:00'}
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
                                            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setScheduleDate(e.target.value)}
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
                                            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setScheduleTime(e.target.value)}
                                            className="h-8 text-sm rounded-lg border-slate-200/80 dark:border-white/[0.08]"
                                        />
                                    </div>
                                </div>
                            )}
                        </div>
                    )}

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
                                            const productHasError = hasProductError(productId);

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
                                                        "relative rounded-lg px-3 py-2.5 transition-all cursor-pointer text-left border",
                                                        `${isSelected ? "bg-blue-50 dark:bg-blue-500/10" : "bg-white dark:bg-white/[0.02]"}`,
                                                        `${isSelected && productHasError ? "border-red-400 bg-red-50 dark:bg-red-500/10" : ""}`,
                                                        `${!isSelected ? "border-slate-200 dark:border-white/[0.08] hover:border-slate-300 dark:hover:border-white/[0.15]" : ""}`
                                                    )}
                                                >
                                                    <p className={cn(
                                                        "text-xs font-medium truncate",
                                                        `${isSelected && !productHasError ? "text-blue-700 dark:text-blue-300" : ""}`,
                                                        `${isSelected && productHasError ? "text-red-700 dark:text-red-300" : ""}`
                                                    )}>
                                                        {product.name || 'Unnamed Product'}
                                                    </p>

                                                    {isSelected && !productHasError && (
                                                        <div className="absolute top-0 right-0 w-0 h-0 border-t-[8px] border-r-[8px] border-t-blue-500 border-r-transparent rounded-tr-lg" />
                                                    )}

                                                    {isSelected && productHasError && (
                                                        <div className="absolute top-0 right-0 w-0 h-0 border-t-[8px] border-r-[8px] border-t-red-400 border-r-transparent rounded-tr-lg" />
                                                    )}
                                                </button>
                                            );
                                        })}
                                    </div>
                                )}
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
                            isProductError
                                ? "bg-red-50/60 border-red-200/60 dark:bg-red-500/5 dark:border-red-500/20"
                                : "bg-blue-50/60 border-blue-200/60 dark:bg-blue-500/[0.04] dark:border-blue-500/20"
                        )}>
                            <div className="flex items-start gap-3">
                                <div className={cn(
                                    "flex h-8 w-8 items-center justify-center rounded-lg ring-1 flex-shrink-0",
                                    isProductError
                                        ? "bg-red-100 text-red-600 ring-red-300/60 dark:bg-red-500/10 dark:text-red-400 dark:ring-red-500/20"
                                        : "bg-blue-100 text-blue-600 ring-blue-300/60 dark:bg-blue-500/10 dark:text-blue-400 dark:ring-blue-500/20"
                                )}>
                                    {isProductError ? (
                                        <AlertTriangle className="h-4 w-4" />
                                    ) : (
                                        <CheckCircle2 className="h-4 w-4" />
                                    )}
                                </div>
                                <div className="flex-1 min-w-0">
                                    <p className={cn(
                                        "text-sm font-medium",
                                        isProductError
                                            ? "text-red-800 dark:text-red-300"
                                            : "text-blue-800 dark:text-blue-300"
                                    )}>
                                        {selectedProductData.name || 'Produk Terpilih'}
                                    </p>
                                </div>
                                {isProductError ? (
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
                    onRetry={handleRetry}
                    onClose={handleCloseError}
                />
            )}
        </>
    );
}