import { useState, useEffect } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@kana-consultant/ui-kit";
import { ImageIcon, FileText, Settings2 } from "lucide-react";
import { cn } from "@/lib/utils";
import type { PreviewSectionProps } from "./preview-section/types";
import { ArticleEditorPanel } from "./preview-section/ArticleEditorPanel";
import { ImageUploadPanel } from "./preview-section/ImageUploadPanel";
import { ConfigSummaryPanel } from "./preview-section/ConfigSummaryPanel";

export function PreviewSection({
    article,
    keywords,
    excerpt,
    slug,
    imageUrl,
    hasImage,
    postMode,
    dailySchedule,
    dailyTime,
    scheduleDate,
    scheduleTime,
    selectedProductsCount,
    autoGenerateImage,
    onArticleUpdate,
    onImageUpload,
}: PreviewSectionProps) {
    // uploadedImage diangkat ke sini karena dipakai bareng oleh tab Gambar & tab Ringkasan
    const [uploadedImage, setUploadedImage] = useState<string | null>(null);
    console.log(keywords,slug,excerpt)

    useEffect(() => {
        if (imageUrl) {
            setUploadedImage(imageUrl);
        }
    }, [imageUrl]);

    return (
        <Tabs defaultValue="article" className="w-full">
            <TabsList
                className={cn(
                    "h-8 gap-0.5 rounded-lg border p-0.5",
                    "bg-slate-50 border-slate-200/80",
                    "dark:bg-white/[0.02] dark:border-white/[0.06]"
                )}
            >
                {[
                    { value: "article", label: "Preview Artikel", icon: FileText },
                    { value: "image", label: "Preview Gambar", icon: ImageIcon },
                    { value: "config", label: "Ringkasan", icon: Settings2 },
                ].map(({ value, label, icon: Icon }) => (
                    <TabsTrigger
                        key={value}
                        value={value}
                        className={cn(
                            "h-7 gap-1.5 rounded-md px-3 text-xs data-[state=active]:shadow-sm",
                            "data-[state=active]:bg-white data-[state=active]:text-slate-800",
                            "dark:data-[state=active]:bg-white/[0.08] dark:data-[state=active]:text-white"
                        )}
                    >
                        <Icon className="h-3 w-3" />
                        {label}
                    </TabsTrigger>
                ))}
            </TabsList>

            <TabsContent value="article" className="mt-3">
                <ArticleEditorPanel article={article} onArticleUpdate={onArticleUpdate} />
            </TabsContent>

            <TabsContent value="image" className="mt-3">
                <ImageUploadPanel
                    hasImage={hasImage}
                    uploadedImage={uploadedImage}
                    onUploadedImageChange={setUploadedImage}
                    onImageUpload={onImageUpload}
                />
            </TabsContent>

            <TabsContent value="config" className="mt-3">
                <ConfigSummaryPanel
                    postMode={postMode}
                    dailySchedule={dailySchedule}
                    dailyTime={dailyTime}
                    scheduleDate={scheduleDate}
                    scheduleTime={scheduleTime}
                    selectedProductsCount={selectedProductsCount}
                    autoGenerateImage={autoGenerateImage}
                    uploadedImage={uploadedImage}
                    imageUrl={imageUrl}
                    slug={slug}
                    keywords={keywords}
                    excerpt={excerpt}
                />
            </TabsContent>
        </Tabs>
    );
}