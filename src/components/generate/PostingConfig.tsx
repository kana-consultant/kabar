import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Send, Calendar, AlertCircle, Settings2, Clock, ImageIcon } from "lucide-react";
import { TargetProducts } from "./TargetProducts";
import type { Product } from "@/types/product";

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
    isPosting : boolean;
}

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
    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <Settings2 className="h-4 w-4" />
                    Konfigurasi Posting
                </CardTitle>
                <CardDescription>
                    Atur jadwal dan target posting konten
                </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
                {/* Mode Posting */}
                <div>
                    <Label className="mb-2 block">Mode Posting</Label>
                    <div className="grid grid-cols-3 gap-2">
                        <Button
                            type="button"
                            variant={postMode === "instant" ? "default" : "outline"}
                            className="flex items-center gap-2"
                            onClick={() => setPostMode("instant")}
                        >
                            <Send className="h-3 w-3" />
                            Langsung
                        </Button>
                        <Button
                            type="button"
                            variant={postMode === "scheduled" ? "default" : "outline"}
                            className="flex items-center gap-2"
                            onClick={() => setPostMode("scheduled")}
                        >
                            <Calendar className="h-3 w-3" />
                            Terjadwal
                        </Button>
                        <Button
                            type="button"
                            variant={postMode === "draft" ? "default" : "outline"}
                            className="flex items-center gap-2"
                            onClick={() => setPostMode("draft")}
                        >
                            <AlertCircle className="h-3 w-3" />
                            Draft
                        </Button>
                    </div>
                </div>

                {/* Schedule Settings */}
                {postMode === "scheduled" && (
                    <div className="rounded-lg border p-3 space-y-3">
                        <div className="flex items-center justify-between">
                            <Label className="flex items-center gap-2">
                                <Clock className="h-3 w-3" />
                                Posting Berulang (Daily)
                            </Label>
                            <Switch
                                checked={dailySchedule}
                                onCheckedChange={setDailySchedule}
                            />
                        </div>

                        {dailySchedule ? (
                            <div>
                                <Label className="mb-1 block text-sm">Waktu Posting Harian</Label>
                                <Input
                                    type="time"
                                    value={dailyTime}
                                    onChange={(e) => setDailyTime(e.target.value)}
                                    className="w-full"
                                />
                                <p className="text-xs text-slate-500 mt-1">
                                    Konten akan diposting setiap hari jam {dailyTime}
                                </p>
                            </div>
                        ) : (
                            <div className="grid grid-cols-2 gap-3">
                                <div>
                                    <Label className="mb-1 block text-sm">Tanggal</Label>
                                    <Input
                                        type="date"
                                        value={scheduleDate}
                                        onChange={(e) => setScheduleDate(e.target.value)}
                                    />
                                </div>
                                <div>
                                    <Label className="mb-1 block text-sm">Waktu</Label>
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

                {/* Auto Generate Image */}
                {/* <div className="flex items-center justify-between">
                    <Label className="flex items-center gap-2">
                        <ImageIcon className="h-3 w-3" />
                        Auto-generate gambar
                    </Label>
                    <Switch
                        checked={autoGenerateImage}
                        onCheckedChange={setAutoGenerateImage}
                    />
                </div> */}

                {/* Post Button */}
                <Button
                    className="w-full mt-4"
                    onClick={onPost}
                    disabled={
                        selectedProducts.length === 0 ||
                        !article ||
                        isPosting
                    }
                >
                    {isPosting ? (
                        <>
                            <svg
                                className="mr-2 h-4 w-4 animate-spin"
                                viewBox="0 0 24 24"
                                fill="none"
                            >
                                <circle
                                    cx="12"
                                    cy="12"
                                    r="10"
                                    stroke="currentColor"
                                    strokeWidth="4"
                                    opacity="0.25"
                                />
                                <path
                                    d="M22 12a10 10 0 00-10-10"
                                    stroke="currentColor"
                                    strokeWidth="4"
                                />
                            </svg>

                            Memproses Posting...
                        </>
                    ) : (
                        <>
                            <Send className="mr-2 h-4 w-4" />

                            {postMode === "instant" &&
                                "Post Sekarang"}

                            {postMode === "scheduled" &&
                                (
                                    dailySchedule
                                        ? "Jadwalkan Harian"
                                        : "Jadwalkan"
                                )
                            }

                            {postMode === "draft" &&
                                "Simpan Draft"}

                            {" ke "}
                            {selectedProducts.length}
                            {" produk"}
                        </>
                    )}
                </Button>

                {!article && (
                    <p className="text-xs text-amber-600 text-center">
                        ⚠️ Generate artikel terlebih dahulu sebelum posting
                    </p>
                )}
            </CardContent>
        </Card>
    );
}