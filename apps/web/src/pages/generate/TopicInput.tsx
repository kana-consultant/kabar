import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Loader2, Sparkles, Image as ImageIcon, Pencil } from "lucide-react";
import { cn } from "@/lib/utils";

interface TopicInputProps {
    topic: string;
    setTopic: (value: string) => void;
    loadingArticle: boolean;
    loadingImage: boolean;
    onGenerateArticle: () => void;
    onGenerateImage: () => void;
    autoGenerateImage: boolean;
    setAutoGenerateImage: (value: boolean) => void;
    article: string;
}

export function TopicInput({
    topic, setTopic,
    loadingArticle, loadingImage,
    onGenerateArticle, onGenerateImage,
    autoGenerateImage, article,
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
                        onChange={(e) => setTopic(e.target.value)}
                        className={cn(
                            "h-8 text-sm rounded-lg",
                            "border-slate-200/80 bg-white placeholder:text-slate-400",
                            "dark:border-white/[0.08] dark:bg-white/[0.03] dark:placeholder:text-slate-600",
                            "focus-visible:ring-1 focus-visible:ring-green-500/40 focus-visible:border-green-400/60",
                            "dark:focus-visible:ring-purple-500/40 dark:focus-visible:border-purple-400/40"
                        )}
                    />
                </div>

                {/* Action buttons */}
                <div className="flex gap-2">
                    <Button
                        onClick={onGenerateArticle}
                        disabled={!topic || loadingArticle}
                        className={cn(
                            "flex-1 h-9 gap-2 rounded-lg text-sm font-medium",
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

                    <Button
                        onClick={onGenerateImage}
                        disabled={!article || loadingImage}
                        variant="outline"
                        className={cn(
                            "h-9 gap-2 rounded-lg text-xs font-medium",
                            "border-slate-200/80 text-slate-500",
                            "hover:text-green-600 hover:border-green-300/60 hover:bg-green-50/50",
                            "dark:border-white/[0.08] dark:text-slate-400",
                            "dark:hover:text-purple-400 dark:hover:border-purple-500/30 dark:hover:bg-purple-500/5",
                            "disabled:opacity-40"
                        )}
                    >
                        {loadingImage
                            ? <><Loader2 className="h-3.5 w-3.5 animate-spin" /> Generating...</>
                            : <><ImageIcon className="h-3.5 w-3.5" /> Generate Gambar</>
                        }
                    </Button>
                </div>

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
                                : "Klik 'Generate Gambar' untuk menambahkan ilustrasi."
                            }
                        </span>
                    </div>
                )}
            </div>
        </div>
    );
}