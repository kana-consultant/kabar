import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent } from "@/components/ui/card";
import { ImageIcon } from "lucide-react";

interface PreviewSectionProps {
    article: string;
    imageUrl: string;
    hasImage: boolean;
    postMode: "instant" | "scheduled" | "draft";
    dailySchedule: boolean;
    dailyTime: string;
    scheduleDate: string;
    scheduleTime: string;
    selectedProductsCount: number;
    autoGenerateImage: boolean;
}

export function PreviewSection({
    article,
    imageUrl,
    hasImage,
    postMode,
    dailySchedule,
    dailyTime,
    scheduleDate,
    scheduleTime,
    selectedProductsCount,
    autoGenerateImage,
}: PreviewSectionProps) {
    return (
        <Tabs defaultValue="article" className="w-full">
            <TabsList>
                <TabsTrigger value="article">Preview Artikel</TabsTrigger>
                <TabsTrigger value="image">Preview Gambar</TabsTrigger>
                <TabsTrigger value="config">Ringkasan Konfigurasi</TabsTrigger>
            </TabsList>

            <TabsContent value="article">
                <Card>
                    <CardContent className="p-6">
                        {article ? (
                            <div 
                                className="prose max-w-none prose-headings:font-bold prose-h1:text-2xl prose-h2:text-xl prose-h3:text-lg prose-p:my-2 prose-ul:my-2 prose-li:my-1"
                                dangerouslySetInnerHTML={{ __html: article }}
                            />
                        ) : (
                            <p className="py-8 text-center text-slate-500">
                                Belum ada artikel. Klik "Generate Artikel" dulu.
                            </p>
                        )}
                    </CardContent>
                </Card>
            </TabsContent>

            <TabsContent value="image">
                <Card>
                    <CardContent className="p-6">
                        {imageUrl ? (
                            <div className="space-y-4">
                                <img
                                    src={imageUrl}
                                    alt="Preview"
                                    className="w-full rounded-lg border object-cover max-h-[400px]"
                                    onError={(e) => {
                                        (e.target as HTMLImageElement).src = "https://placehold.co/800x400?text=Image+Failed+to+Load";
                                    }}
                                />
                                <p className="text-center text-sm text-slate-500">
                                    Klik kanan pada gambar → Save Image As untuk menyimpan
                                </p>
                            </div>
                        ) : (
                            <div className="flex min-h-[200px] flex-col items-center justify-center rounded-lg border-2 border-dashed">
                                <ImageIcon className="h-8 w-8 text-slate-400" />
                                <p className="mt-2 text-slate-500">
                                    {hasImage ? "Generate gambar terlebih dahulu" : "Belum ada gambar"}
                                </p>
                            </div>
                        )}
                    </CardContent>
                </Card>
            </TabsContent>

            <TabsContent value="config">
                <Card>
                    <CardContent className="p-6 space-y-2">
                        <h4 className="font-medium">Ringkasan Konfigurasi Posting</h4>
                        <div className="text-sm space-y-1 text-slate-600">
                            <p>📝 Mode: {
                                postMode === "instant" ? "Langsung posting" :
                                postMode === "scheduled" ? "Terjadwal" : "Simpan sebagai draft"
                            }</p>
                            {postMode === "scheduled" && (
                                <p>📅 Jadwal: {dailySchedule ? `Setiap hari jam ${dailyTime}` : `${scheduleDate} jam ${scheduleTime}`}</p>
                            )}
                            <p>🎯 Target: {selectedProductsCount} produk</p>
                            <p>🖼️ Auto gambar: {autoGenerateImage ? "Ya" : "Tidak"}</p>
                            {imageUrl && <p>🖼️ Gambar: ✓ Ada</p>}
                        </div>
                    </CardContent>
                </Card>
            </TabsContent>
        </Tabs>
    );
}