import { Input } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Loader2, Sparkles, Pencil, Wand2 } from "lucide-react";
import { cn } from "@/lib/utils";

interface TopicInputProps {
    topic: string;
    setTopic: (value: string) => void;
    loadingArticle: boolean;
    onGenerateArticle: () => void;
    autoGenerateImage: boolean;
    setAutoGenerateImage: (value: boolean) => void;
    article: string;
}

export function TopicInput({
    topic, setTopic,
    loadingArticle,
    onGenerateArticle,
    autoGenerateImage, setAutoGenerateImage, article,
}: TopicInputProps) {
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
                <div className="flex h-8 w-8 items-center justify-center rounded-lg ring-1 bg-green-50 text-green-600 ring-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:ring-green-500/20">
                    <Pencil className="h-3.5 w-3.5" />
                </div>
                <div>
                    <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                        Topik Artikel
                    </p>
                    <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                        Masukkan topik yang ingin dibuatkan artikelnya
                    </p>
                </div>
            </div>

            {/* Body */}
            <div className="p-5 space-y-4">
                <div className="space-y-1.5">
                    <label className="text-[11px] font-medium uppercase tracking-wide text-slate-400 dark:text-slate-600">
                        Topik / Keyword
                    </label>
                    <Input
                        placeholder="Contoh: Cara Memilih Sepatu Lari yang Tepat"
                        value={topic}
                        onChange={(e: any) => {
                            let value = e.target.value;

                            // 🔒 Sanitasi
                            let sanitized = value
                                .replace(/<[^>]*>/g, '')
                                .replace(/```[\s\S]*?```/g, '')
                                .replace(/`[^`]*`/g, '')
                                .replace(/[<>{}|\\^~\[\]]/g, '');

                            // 🔒 Hapus semua simbol di depan sampai ketemu huruf
                            sanitized = sanitized.replace(/^[^a-zA-Z]+/, '');

                            // Batasi 50 karakter
                            sanitized = sanitized.substring(0, 50);

                            setTopic(sanitized);
                        }}
                        maxLength={50}
                        className={cn(
                            "h-8 text-sm rounded-lg",
                            "border-slate-200/80 bg-white placeholder:text-slate-400",
                            "dark:border-white/[0.08] dark:bg-white/[0.03] dark:placeholder:text-slate-600",
                            "focus-visible:ring-1 focus-visible:ring-green-500/40 focus-visible:border-green-400/60",
                            "dark:focus-visible:ring-purple-500/40 dark:focus-visible:border-purple-400/40"
                        )}
                    />
                    <div className="flex justify-end">
                        <p className={cn(
                            "text-[10px]",
                            topic.length >= 50 ? "text-red-500" : "text-slate-400 dark:text-slate-600"
                        )}>
                            {topic.length}/50
                        </p>
                    </div>
                </div>

                {/* Auto-generate image toggle */}
                <div className={cn(
                    "flex items-center justify-between gap-3 rounded-lg border px-3 py-2.5",
                    "border-slate-200/80 bg-slate-50/60",
                    "dark:border-white/[0.06] dark:bg-white/[0.02]"
                )}>
                    <div className="flex items-center gap-2.5 min-w-0">
                        <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md ring-1 bg-amber-50 text-amber-600 ring-amber-200/60 dark:bg-amber-500/10 dark:text-amber-400 dark:ring-amber-500/20">
                            <Wand2 className="h-3.5 w-3.5" />
                        </div>
                        <div className="min-w-0">
                            <p className="text-xs font-medium text-slate-700 dark:text-slate-200 truncate">
                                Generate Gambar Otomatis
                            </p>
                            <p className="text-[11px] text-slate-400 dark:text-slate-600 truncate">
                                Gambar dibuat otomatis setelah artikel selesai
                            </p>
                        </div>
                    </div>

                    <button
                        type="button"
                        role="switch"
                        aria-checked={autoGenerateImage}
                        aria-label="Toggle generate gambar otomatis"
                        onClick={() => setAutoGenerateImage(!autoGenerateImage)}
                        className={cn(
                            "relative inline-flex h-5 w-9 shrink-0 items-center rounded-full transition-colors duration-200",
                            "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2",
                            autoGenerateImage
                                ? "bg-green-600 dark:bg-purple-600 focus-visible:ring-green-500/40 dark:focus-visible:ring-purple-500/40"
                                : "bg-slate-200 dark:bg-white/[0.08] focus-visible:ring-slate-400/40"
                        )}
                    >
                        <span
                            className={cn(
                                "inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow transition-transform duration-200",
                                autoGenerateImage ? "translate-x-[19px]" : "translate-x-[3px]"
                            )}
                        />
                    </button>
                </div>

                {/* Action button */}
                <Button
                    onClick={onGenerateArticle}
                    disabled={!topic || loadingArticle}
                    className={cn(
                        "w-full h-9 gap-2 rounded-lg text-sm font-medium",
                        "bg-green-600 hover:bg-green-700 text-white shadow-sm",
                        "dark:bg-purple-600 dark:hover:bg-purple-700",
                        "disabled:opacity-50"
                    )}
                >
                    {loadingArticle
                        ? <><Loader2 className="h-3.5 w-3.5 animate-spin" /> Mengenerate...</>
                        : <><Sparkles className="h-3.5 w-3.5" /> Generate Artikel</>
                    }
                </Button>

                {/* Article ready indicator */}
                {article && (
                    <div className={cn(
                        "flex items-center gap-2 rounded-lg border px-3 py-2",
                        "bg-green-50 border-green-200/60 text-green-700",
                        "dark:bg-green-500/10 dark:border-green-500/20 dark:text-green-400"
                    )}>
                        {/* Ping dot */}
                        <span className="relative flex h-1.5 w-1.5 shrink-0">
                            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-500 opacity-75" />
                            <span className="relative inline-flex h-1.5 w-1.5 rounded-full bg-green-500" />
                        </span>
                        <span className="text-xs">
                            Artikel siap.{" "}
                            {autoGenerateImage
                                ? "Gambar otomatis akan digenerate."
                                : "Tidak ada gambar yang akan ditambahkan."
                            }
                        </span>
                    </div>
                )}
            </div>
        </div>
    );
}