import { useState, useEffect } from "react";
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import { Loader2, Sparkles, Zap } from "lucide-react";
import { useGenerate } from "@/hooks/useGenerate";
import { useNavigate } from "@tanstack/react-router";
import { cn } from "@/lib/utils";

export function QuickGenerate() {
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();
    const { products, topic, setSelectedProducts, selectedProducts, setTopic, quickGenerate } = useGenerate();

    useEffect(() => {
        if (products.length > 0 && selectedProducts.length === 0) {
            setSelectedProducts([products[0].id]);
        }
    }, [products]);

    const handleGenerate = async () => {
        setLoading(true);
        try {
            const draftId = await quickGenerate();
            if (draftId) navigate({ to: "/history" });
        } finally {
            setLoading(false);
        }
    };

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
                <div className={cn(
                    "flex h-8 w-8 items-center justify-center rounded-lg ring-1",
                    "bg-green-50 text-green-600 ring-green-200/60",
                    "dark:bg-purple-500/10 dark:text-purple-400 dark:ring-purple-500/20"
                )}>
                    <Zap className="h-3.5 w-3.5" />
                </div>
                <div>
                    <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                        Quick Generate
                    </p>
                    <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                        Generate artikel + gambar otomatis
                    </p>
                </div>
            </div>

            {/* Form */}
            <div className="p-5 space-y-4">
                <div className="space-y-1.5">
                    <label className="text-[11px] font-medium uppercase tracking-wide text-slate-400 dark:text-slate-600">
                        Topik / Keyword
                    </label>
                    <Input
                        id="topic"
                        value={topic}
                        placeholder="Contoh: Cara Memilih Sepatu Gunung untuk Pemula"
                        onChange={(e : any) => setTopic(e.target.value)}
                        className={cn(
                            "h-8 text-sm rounded-lg",
                            "border-slate-200/80 bg-white placeholder:text-slate-400",
                            "dark:border-white/[0.08] dark:bg-white/[0.03] dark:placeholder:text-slate-600",
                            "focus-visible:ring-1 focus-visible:ring-green-500/40 focus-visible:border-green-400/60",
                            "dark:focus-visible:ring-purple-500/40 dark:focus-visible:border-purple-400/40"
                        )}
                    />
                </div>

                <div className="space-y-1.5">
                    <label className="text-[11px] font-medium uppercase tracking-wide text-slate-400 dark:text-slate-600">
                        Target Produk
                    </label>
                    <select
                        value={selectedProducts[0] || ""}
                        onChange={(e : any) => setSelectedProducts([e.target.value])}
                        className={cn(
                            "w-full h-8 rounded-lg border px-2.5 text-sm appearance-none",
                            "border-slate-200/80 bg-white text-slate-700",
                            "dark:border-white/[0.08] dark:bg-white/[0.03] dark:text-slate-300",
                            "focus:outline-none focus:ring-1 focus:ring-green-500/40",
                            "dark:focus:ring-purple-500/40"
                        )}
                    >
                        {products.map((item) => (
                            <option key={item.id} value={item.id}>{item.name}</option>
                        ))}
                    </select>
                </div>

                <Button
                    onClick={handleGenerate}
                    disabled={loading || !topic}
                    className={cn(
                        "w-full h-9 gap-2 rounded-lg text-sm font-medium",
                        "bg-green-600 hover:bg-green-700 text-white shadow-sm",
                        "dark:bg-purple-600 dark:hover:bg-purple-700",
                        "disabled:opacity-50"
                    )}
                >
                    {loading
                        ? <><Loader2 className="h-3.5 w-3.5 animate-spin" /> Generating...</>
                        : <><Sparkles className="h-3.5 w-3.5" /> Generate Konten</>
                    }
                </Button>
            </div>
        </div>
    );
}