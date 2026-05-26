import {
    Select, SelectContent, SelectItem,
    SelectTrigger, SelectValue,
} from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Loader2, Cpu, Key, RefreshCw } from "lucide-react";
import { useModels } from "@/hooks/useGenerate/useModel";
import { useEffect } from "react";
import { cn } from "@/lib/utils";

interface ModelSelectorProps {
    selectedModelId: string;
    onModelChange: (modelId: string) => void;
    filterByService?: "text" | "image" | "all";
}

export function ModelSelector({
    selectedModelId, onModelChange, filterByService = "all",
}: ModelSelectorProps) {
    const {
        models, loadingModels, isError, error,
        refetch, getProviderColor, getServiceBadgeColor, getTextAndImageModels,
    } = useModels();

    useEffect(() => { getTextAndImageModels(); }, []);

    const filteredModels = models ?? [];
    const selectedModel = filteredModels.find(m => m.id === selectedModelId);

    const cardCls = cn(
        "overflow-hidden rounded-2xl border",
        "bg-white border-slate-200/80",
        "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
    );

    const hdrCls = cn(
        "flex items-center justify-between px-5 py-4 border-b",
        "border-slate-100 bg-slate-50/60",
        "dark:border-white/[0.05] dark:bg-white/[0.02]"
    );

    if (loadingModels) return (
        <div className={cardCls}>
            <div className={hdrCls}>
                <div className="flex items-center gap-3">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg ring-1 bg-violet-50 text-violet-600 ring-violet-200/60 dark:bg-violet-500/10 dark:text-violet-400 dark:ring-violet-500/20">
                        <Cpu className="h-3.5 w-3.5" />
                    </div>
                    <p className="text-sm font-medium text-slate-800 dark:text-slate-100">AI Model</p>
                </div>
            </div>
            <div className="flex items-center justify-center py-8">
                <Loader2 className="h-5 w-5 animate-spin text-slate-400" />
            </div>
        </div>
    );

    if (isError) return (
        <div className={cardCls}>
            <div className={hdrCls}>
                <div className="flex items-center gap-3">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg ring-1 bg-red-50 text-red-500 ring-red-200/60 dark:bg-red-500/10 dark:text-red-400 dark:ring-red-500/20">
                        <Cpu className="h-3.5 w-3.5" />
                    </div>
                    <p className="text-sm font-medium text-slate-800 dark:text-slate-100">AI Model</p>
                </div>
            </div>
            <div className="flex flex-col items-center gap-2 py-8">
                <p className="text-xs text-red-500 dark:text-red-400">
                    {error?.message || "Gagal memuat model"}
                </p>
                <Button
                    variant="outline" size="sm"
                    onClick={() => refetch()}
                    className="h-7 gap-1.5 px-3 text-xs"
                >
                    <RefreshCw className="h-3 w-3" /> Retry
                </Button>
            </div>
        </div>
    );

    return (
        <div className={cardCls}>
            {/* Header */}
            <div className={hdrCls}>
                <div className="flex items-center gap-3">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg ring-1 bg-violet-50 text-violet-600 ring-violet-200/60 dark:bg-violet-500/10 dark:text-violet-400 dark:ring-violet-500/20">
                        <Cpu className="h-3.5 w-3.5" />
                    </div>
                    <div>
                        <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                            AI Model
                            {filterByService !== "all" && (
                                <span className="ml-1.5 text-xs font-normal text-slate-400">
                                    ({filterByService === "text" ? "Text" : "Image"})
                                </span>
                            )}
                        </p>
                        <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                            Model yang digunakan untuk generate
                        </p>
                    </div>
                </div>
                <Button
                    variant="ghost" size="icon"
                    onClick={() => refetch()}
                    className="h-7 w-7 rounded-lg text-slate-400 hover:text-green-600 hover:bg-green-50 dark:hover:text-purple-400 dark:hover:bg-purple-500/10"
                    title="Refresh models"
                >
                    <RefreshCw className="h-3.5 w-3.5" />
                </Button>
            </div>

            {/* Body */}
            <div className="p-5">
                {filteredModels.length === 0 ? (
                    <div className="flex flex-col items-center py-8 gap-2">
                        <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-slate-100 text-slate-400 dark:bg-white/[0.04] dark:text-slate-600">
                            <Key className="h-5 w-5" />
                        </div>
                        <p className="text-sm font-medium text-slate-600 dark:text-slate-400">
                            Belum ada model dikonfigurasi
                        </p>
                        <p className="text-xs text-slate-400 dark:text-slate-600">
                            Tambahkan API key di halaman Settings
                        </p>
                    </div>
                ) : (
                    <div className="space-y-3">
                        <Select value={selectedModelId} onValueChange={onModelChange}>
                            <SelectTrigger className={cn(
                                "h-8 text-sm rounded-lg",
                                "border-slate-200/80 bg-white",
                                "dark:border-white/[0.08] dark:bg-white/[0.03]",
                                "focus:ring-1 focus:ring-green-500/40 dark:focus:ring-purple-500/40"
                            )}>
                                <SelectValue placeholder="Pilih AI Model" />
                            </SelectTrigger>
                            <SelectContent className="w-[580px] max-w-[95vw]">
                                {filteredModels.map((model) => (
                                    <SelectItem key={model.id} value={model.id} className="py-2.5">
                                        <div className="flex flex-wrap items-center gap-1.5">
                                            <span className="text-sm font-medium">{model.modelDisplayName}</span>
                                            <span className={cn(
                                                "inline-flex items-center rounded-md border px-1.5 py-0.5 text-[10px] font-medium",
                                                getServiceBadgeColor(model.service)
                                            )}>
                                                {model.service === "text" ? "Text" : "Image"}
                                            </span>
                                            <span className={cn(
                                                "inline-flex items-center rounded-md px-1.5 py-0.5 text-[10px] font-medium",
                                                getProviderColor(model.providerDisplayName)
                                            )}>
                                                {model.providerDisplayName}
                                            </span>
                                            {model.isActive && (
                                                <span className="inline-flex items-center rounded-md border px-1.5 py-0.5 text-[10px] font-medium bg-green-50 text-green-700 border-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:border-green-500/20">
                                                    Active
                                                </span>
                                            )}
                                        </div>
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>

                        {selectedModelId && selectedModel && (
                            <div className="flex items-center gap-2 text-[11px] text-slate-400 dark:text-slate-600">
                                <span className="relative flex h-1.5 w-1.5">
                                    <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-500 opacity-75" />
                                    <span className="relative inline-flex h-1.5 w-1.5 rounded-full bg-green-500" />
                                </span>
                                {selectedModel.modelDisplayName} siap digunakan
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}