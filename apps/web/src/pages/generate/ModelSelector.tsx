import {
    Select, SelectContent, SelectItem,
    SelectTrigger, SelectValue,
} from "@kana-consultant/ui-kit";
import { Loader2, RefreshCw } from "lucide-react";
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

    const filteredModels = (models ?? []).filter((model) =>
        filterByService === "all" ? true : model.service === filterByService
    );
    const selectedModel = filteredModels.find(m => m.id === selectedModelId);

    const label = filterByService === "text" ? "Model Teks" : filterByService === "image" ? "Model Gambar" : "AI Model";

    return (
        <div className="space-y-1.5">
            <div className="flex items-center justify-between">
                <label className="text-[11px] font-medium uppercase tracking-wide text-slate-400 dark:text-slate-600">
                    {label}
                </label>
                <button
                    type="button"
                    onClick={() => refetch()}
                    title="Refresh models"
                    className="text-slate-400 hover:text-green-600 dark:hover:text-purple-400 disabled:opacity-40"
                    disabled={loadingModels}
                >
                    <RefreshCw className={cn("h-3 w-3", `${loadingModels && "animate-spin"}`)} />
                </button>
            </div>

            {loadingModels ? (
                <div className="flex h-8 items-center gap-2 rounded-lg border border-slate-200/80 bg-slate-50/60 px-3 text-xs text-slate-400 dark:border-white/[0.08] dark:bg-white/[0.02] dark:text-slate-600">
                    <Loader2 className="h-3 w-3 animate-spin" /> Memuat model...
                </div>
            ) : isError ? (
                <div className="flex h-8 items-center justify-between rounded-lg border border-red-200/60 bg-red-50 px-3 text-xs text-red-500 dark:border-red-500/20 dark:bg-red-500/10 dark:text-red-400">
                    <span className="truncate">{error?.message || "Gagal memuat model"}</span>
                    <button onClick={() => refetch()} className="ml-2 shrink-0 underline">Retry</button>
                </div>
            ) : filteredModels.length === 0 ? (
                <div className="flex h-8 items-center rounded-lg border border-slate-200/80 bg-slate-50/60 px-3 text-xs text-slate-400 dark:border-white/[0.08] dark:bg-white/[0.02] dark:text-slate-600">
                    Belum ada model dikonfigurasi
                </div>
            ) : (
                <>
                    <Select value={selectedModelId} onValueChange={onModelChange}>
                        <SelectTrigger className={cn(
                            "h-9 text-sm rounded-lg overflow-hidden",
                            "border-slate-200/80 bg-white",
                            "dark:border-white/[0.08] dark:bg-white/[0.03]",
                            "focus:outline-none focus:ring-1 focus:ring-green-500/40 focus:border-green-400/60",
                            "dark:focus:ring-purple-500/40 dark:focus:border-purple-400/40"
                        )}>
                            <SelectValue placeholder="Pilih AI Model">
                                <span className="block truncate text-left">
                                    {selectedModel?.model_name}
                                </span>
                            </SelectValue>
                        </SelectTrigger>
                        <SelectContent className="w-[320px] max-w-[95vw]">
                            {filteredModels.map((model) => (
                                <SelectItem key={model.id} value={model.id} className="py-2">
                                    <div className="flex items-center gap-1.5 min-w-0">
                                        <span className="text-sm font-medium truncate">{model.model_name}</span>
                                        {filterByService === "all" && (
                                            <span className={cn(
                                                "inline-flex items-center rounded-md border px-1.5 py-0.5 text-[10px] font-medium shrink-0",
                                                getServiceBadgeColor(model.service)
                                            )}>
                                                {model.service === "text" ? "Text" : "Image"}
                                            </span>
                                        )}
                                        <span className={cn(
                                            "inline-flex items-center rounded-md px-1.5 py-0.5 text-[10px] font-medium shrink-0",
                                            getProviderColor(model.provider_name)
                                        )}>
                                            {model.provider_display_name}
                                        </span>
                                    </div>
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>

                    {selectedModelId && selectedModel && (
                        <div className="flex items-center gap-1.5 text-[11px] text-slate-400 dark:text-slate-600 min-w-0">
                            <span className="relative flex h-1.5 w-1.5 shrink-0">
                                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-500 opacity-75" />
                                <span className="relative inline-flex h-1.5 w-1.5 rounded-full bg-green-500" />
                            </span>
                            {filterByService === "all" && (
                                <span className={cn(
                                    "inline-flex items-center rounded-md border px-1.5 py-0.5 text-[10px] font-medium shrink-0",
                                    getServiceBadgeColor(selectedModel.service)
                                )}>
                                    {selectedModel.service === "text" ? "Text" : "Image"}
                                </span>
                            )}
                            <span className={cn(
                                "inline-flex items-center rounded-md px-1.5 py-0.5 text-[10px] font-medium shrink-0",
                                getProviderColor(selectedModel.provider_name)
                            )}>
                                {selectedModel.provider_display_name}
                            </span>
                            <span className="truncate">siap digunakan</span>
                        </div>
                    )}
                </>
            )}
        </div>
    );
}