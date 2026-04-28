// src/components/generate/TopicInput.tsx
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Loader2, Sparkles, Image as ImageIcon } from "lucide-react";

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
    topic,
    setTopic,
    loadingArticle,
    loadingImage,
    onGenerateArticle,
    onGenerateImage,
    autoGenerateImage,
    setAutoGenerateImage,
    article,
}: TopicInputProps) {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Topik Artikel</CardTitle>
                <CardDescription>Masukkan topik yang ingin dibuatkan artikelnya</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
                <div>
                    <Label>Topik</Label>
                    <Input
                        placeholder="Contoh: Cara Memilih Sepatu Lari yang Tepat"
                        value={topic}
                        onChange={(e) => setTopic(e.target.value)}
                        className="mt-1"
                    />
                </div>

                {/* <div className="flex items-center justify-between rounded-lg border p-3">
                    <div className="space-y-0.5">
                        <Label className="text-sm">Auto Generate Gambar</Label>
                        <p className="text-xs text-slate-500">
                            Generate gambar otomatis setelah artikel selesai
                        </p>
                    </div>
                    <Switch
                        checked={autoGenerateImage}
                        onCheckedChange={setAutoGenerateImage}
                    />
                </div> */}

                <div className="flex gap-3">
                    <Button
                        onClick={onGenerateArticle}
                        disabled={!topic || loadingArticle}
                        className="flex-1"
                    >
                        {loadingArticle ? (
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        ) : (
                            <Sparkles className="mr-2 h-4 w-4" />
                        )}
                        {loadingArticle ? "Mengenerate..." : "Generate Artikel"}
                    </Button>

                    <Button
                        variant="outline"
                        onClick={onGenerateImage}
                        disabled={!article || loadingImage}
                    >
                        {loadingImage ? (
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        ) : (
                            <ImageIcon className="mr-2 h-4 w-4" />
                        )}
                        {loadingImage ? "Generating..." : "Generate Gambar"}
                    </Button>
                </div>

                {article && (
                    <div className="mt-4 rounded-lg bg-green-50 p-3 text-xs text-green-700 dark:bg-green-950 dark:text-green-300">
                        ✅ Artikel siap. {autoGenerateImage ? "Gambar otomatis akan digenerate." : "Klik 'Generate Gambar' untuk menambahkan ilustrasi."}
                    </div>
                )}
            </CardContent>
        </Card>
    );
}