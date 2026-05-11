import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Send, Calendar, AlertCircle, Settings2, Clock, ImageIcon, Loader2 } from "lucide-react";
import { TargetProducts } from "./TargetProducts";
import type { Product } from "@/services/product";
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

const postModeConfig = {
    instant: { icon: Send, label: "Langsung", color: "default" },
    scheduled: { icon: Calendar, label: "Terjadwal", color: "default" },
    draft: { icon: AlertCircle, label: "Draft", color: "default" },
};

export function PostingConfig({
    postMode,
    setPostMode,
    scheduleDate,
    setScheduleDate,
    scheduleTime,
    setScheduleTime,
    dailySchedule,
    setDailySchedule,
    dailyTime,
    setDailyTime,
    autoGenerateImage,
    setAutoGenerateImage,
    products,
    selectedProducts,
    postToAll,
    onToggleProduct,
    onSelectAll,
    article,
    onPost,
    isPosting
}: PostingConfigProps) {
    const getPostButtonText = () => {
        const productCount = selectedProducts.length;
        const baseText = {
            instant: "Post Sekarang",
            scheduled: dailySchedule ? "Jadwalkan Harian" : "Jadwalkan",
            draft: "Simpan Draft"
        }[postMode];
        
        return `${baseText} ke ${productCount} ${productCount === 1 ? 'produk' : 'produk'}`;
    };

    const isPostDisabled = () => {
        return selectedProducts.length === 0 || !article || isPosting;
    };

    return (
        <Card className="sticky top-4">
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <Settings2 className="h-4 w-4" />
                    Konfigurasi Posting
                </CardTitle>
                <CardDescription>
                    Atur jadwal dan target posting konten
                </CardDescription>
            </CardHeader>
            
            <CardContent className="space-y-6">
                {/* Mode Posting */}
                <div className="space-y-2">
                    <Label>Mode Posting</Label>
                    <div className="grid grid-cols-3 gap-2">
                        {Object.entries(postModeConfig).map(([mode, { icon: Icon, label }]) => (
                            <Button
                                key={mode}
                                type="button"
                                variant={postMode === mode ? "default" : "outline"}
                                className={cn(
                                    "flex items-center gap-2 transition-all",
                                    postMode === mode && "shadow-md"
                                )}
                                onClick={() => setPostMode(mode as any)}
                            >
                                <Icon className="h-3.5 w-3.5" />
                                {label}
                            </Button>
                        ))}
                    </div>
                </div>

                {/* Schedule Settings */}
                {postMode === "scheduled" && (
                    <div className="rounded-lg border bg-muted/20 p-4 space-y-4">
                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label className="flex items-center gap-2">
                                    <Clock className="h-3.5 w-3.5" />
                                    Posting Berulang (Daily)
                                </Label>
                                <p className="text-xs text-muted-foreground">
                                    Jadwalkan posting otomatis setiap hari
                                </p>
                            </div>
                            <Switch
                                checked={dailySchedule}
                                onCheckedChange={setDailySchedule}
                            />
                        </div>

                        {dailySchedule ? (
                            <div className="space-y-2">
                                <Label className="text-sm">Waktu Posting Harian</Label>
                                <Input
                                    type="time"
                                    value={dailyTime}
                                    onChange={(e) => setDailyTime(e.target.value)}
                                    className="w-full"
                                />
                                <p className="text-xs text-muted-foreground">
                                    Konten akan diposting setiap hari jam {dailyTime || "---"}
                                </p>
                            </div>
                        ) : (
                            <div className="grid grid-cols-2 gap-3">
                                <div className="space-y-2">
                                    <Label className="text-sm">Tanggal</Label>
                                    <Input
                                        type="date"
                                        value={scheduleDate}
                                        onChange={(e) => setScheduleDate(e.target.value)}
                                    />
                                </div>
                                <div className="space-y-2">
                                    <Label className="text-sm">Waktu</Label>
                                    <Input
                                        type="time"
                                        value={scheduleTime}
                                        onChange={(e) => setScheduleTime(e.target.value)}
                                    />
                                </div>
                            </div>
                        )}
                    </div>
                )}

                {/* Target Products */}
                <TargetProducts
                    products={products}
                    selectedProducts={selectedProducts}
                    postToAll={postToAll}
                    onToggleProduct={onToggleProduct}
                    onSelectAll={onSelectAll}
                />

                {/* Auto Generate Image - Uncomment if needed */}
                {false && (
                    <div className="flex items-center justify-between rounded-lg border p-4">
                        <div className="space-y-0.5">
                            <Label className="flex items-center gap-2">
                                <ImageIcon className="h-3.5 w-3.5" />
                                Auto-generate gambar
                            </Label>
                            <p className="text-xs text-muted-foreground">
                                Generate gambar otomatis menggunakan AI
                            </p>
                        </div>
                        <Switch
                            checked={autoGenerateImage}
                            onCheckedChange={setAutoGenerateImage}
                        />
                    </div>
                )}

                {/* Error/Warning Message */}
                {!article && (
                    <div className="rounded-lg bg-amber-50 border border-amber-200 p-3">
                        <p className="text-xs text-amber-700 text-center flex items-center justify-center gap-2">
                            <AlertCircle className="h-3.5 w-3.5" />
                            Generate artikel terlebih dahulu sebelum posting
                        </p>
                    </div>
                )}

                {/* Post Button */}
                <Button
                    className="w-full mt-6"
                    size="lg"
                    onClick={onPost}
                    disabled={isPostDisabled()}
                >
                    {isPosting ? (
                        <>
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                            Memproses Posting...
                        </>
                    ) : (
                        <>
                            <Send className="mr-2 h-4 w-4" />
                            {getPostButtonText()}
                        </>
                    )}
                </Button>
            </CardContent>
        </Card>
    );
}