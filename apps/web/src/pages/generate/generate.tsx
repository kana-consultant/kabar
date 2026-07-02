// components/Generate.tsx
import { useState } from "react"; // 👈 PENTING: Import useState
import { GenerateHeader } from "./GenerateHeader";
import { TopicInput } from "./TopicInput";
import { PostingConfig } from "./PostingConfig";
import { PreviewSection } from "./PreviewSection";
import { PublishResultDialog } from "../draft/PublishResultDialog";
import { useGenerate } from "@/hooks/useGenerate";
import { ModelSelector } from "./ModelSelector";
import { cn } from "@/lib/utils";

export default function Generate() {
    const {
        topic, setTopic,
        article,
        imageUrl,
        loadingArticle, loadingImage,
        selectedProducts, 
        setSelectedProducts,
        postMode, setPostMode,
        scheduleTime, setScheduleTime,
        scheduleDate, setScheduleDate,
        dailySchedule, setDailySchedule,
        dailyTime, setDailyTime,
        autoGenerateImage, setAutoGenerateImage,
        products,
        productsLoading,
        productsError,
        currentDraftId,
        generateArticle,
        handlePost,
        publishResults,
        showResultDialog,
        closeResultDialog,
        isPosting,
        selectedModelId,
        setSelectedModelId,
        isError,
        results,
        errorData,
        onRetry,
        onCloseError,
    } = useGenerate();

    // 👇 STATE LOKAL: model khusus untuk generate gambar (step kedua, muncul jika autoGenerateImage true)
    const [selectedImageModelId, setSelectedImageModelId] = useState("");
    const [imageModelError, setImageModelError] = useState(false);

    // 👇 STATE KHUSUS UNTUK SINGLE SELECT

    // 👇 HANDLER UNTUK SINGLE SELECT
    const handleSelectProduct = (productId: string | null) => {
        if (productId === null) {
            setSelectedProducts([]);
        } else {
            setSelectedProducts([productId]);
        }
    };

    // 👇 VALIDASI: kalau auto-generate gambar aktif, model image wajib sudah dipilih
    const handleGenerateArticle = () => {
        if (autoGenerateImage && !selectedImageModelId) {
            setImageModelError(true);
            return;
        }
        setImageModelError(false);
        // 👇 selectedImageModelId & autoGenerateImage diteruskan supaya generateArticle
        // bisa langsung generate gambar-gambar di dalam artikel pakai model ini.
        // NOTE: signature generateArticle() di useGenerate.ts perlu disesuaikan
        // untuk menerima & menggunakan kedua argumen ini.
        generateArticle(autoGenerateImage, selectedImageModelId);
    };

    // 👇 Toggle auto-generate gambar, sekaligus bersihkan error kalau dimatikan
    const handleToggleAutoGenerateImage = (value: boolean) => {
        setAutoGenerateImage(value);
        if (!value) setImageModelError(false);
    };

    // 👇 Pilih model image, otomatis hapus error begitu user pilih sesuatu
    const handleSelectImageModel = (modelId: string) => {
        setSelectedImageModelId(modelId);
        if (modelId) setImageModelError(false);
    };

    if (productsLoading) {
        return (
            <div className="flex min-h-[400px] items-center justify-center">
                <div className="text-center">
                    <div className="mb-4 h-8 w-8 animate-spin rounded-full border-4 border-cyan-500 border-t-transparent mx-auto" />
                    <p className="text-zinc-600 dark:text-zinc-400">Memuat data produk...</p>
                </div>
            </div>
        );
    }

    if (productsError) {
        return (
            <div className="flex min-h-[400px] items-center justify-center">
                <div className="text-center">
                    <p className="mb-2 text-red-600 dark:text-red-400">{productsError}</p>
                    <button
                        onClick={() => window.location.reload()}
                        className="rounded-lg bg-cyan-500 px-4 py-2 text-white hover:bg-cyan-600 transition-colors"
                    >
                        Coba Lagi
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {/* Draft indicator */}
            {currentDraftId && (
                <div className="rounded-lg bg-cyan-50 p-3 text-sm text-cyan-700 dark:bg-cyan-950/50 dark:text-cyan-300 border border-cyan-200 dark:border-cyan-800">
                    Sedang mengedit draft. Simpan akan memperbarui draft yang ada.
                </div>
            )}

            <GenerateHeader />

            {/* Main content grid */}
            <div className="grid gap-6 lg:grid-cols-2">
                {/* Left column */}
                <div className="space-y-6">
                    <TopicInput
                        topic={topic}
                        setTopic={setTopic}
                        loadingArticle={loadingArticle}
                        onGenerateArticle={handleGenerateArticle}
                        autoGenerateImage={autoGenerateImage}
                        setAutoGenerateImage={handleToggleAutoGenerateImage}
                        article={article}
                    />

                    <div className="grid gap-4 sm:grid-cols-2">
                        <div className={cn(
                            "rounded-2xl border border-slate-200/80 bg-white p-4 dark:bg-[#0f0d1a] dark:border-white/[0.06]",
                            `${!autoGenerateImage && "sm:col-span-2"}`
                        )}>
                            <ModelSelector
                                selectedModelId={selectedModelId}
                                onModelChange={setSelectedModelId}
                                filterByService="text"
                            />
                        </div>

                        {/* Step kedua: muncul hanya kalau auto-generate gambar aktif */}
                        {autoGenerateImage && (
                            <div className="rounded-2xl border border-slate-200/80 bg-white p-4 dark:bg-[#0f0d1a] dark:border-white/[0.06]">
                                <ModelSelector
                                    selectedModelId={selectedImageModelId}
                                    onModelChange={handleSelectImageModel}
                                    filterByService="image"
                                />
                                {imageModelError && (
                                    <p className="mt-1.5 text-xs text-red-500 dark:text-red-400">
                                        Pilih model gambar dulu sebelum generate artikel.
                                    </p>
                                )}
                            </div>
                        )}
                    </div>
                </div>

                {/* Right column */}
                <div className="space-y-6">
                    <PostingConfig
                        postMode={postMode}
                        setPostMode={setPostMode}
                        scheduleDate={scheduleDate}
                        setScheduleDate={setScheduleDate}
                        scheduleTime={scheduleTime}
                        setScheduleTime={setScheduleTime}
                        dailySchedule={dailySchedule}
                        setDailySchedule={setDailySchedule}
                        dailyTime={dailyTime}
                        setDailyTime={setDailyTime}
                        autoGenerateImage={autoGenerateImage}
                        setAutoGenerateImage={handleToggleAutoGenerateImage}
                        products={products}
                        selectedProduct={selectedProducts[0]}
                        onSelectProduct={handleSelectProduct}
                        article={article}
                        onPost={handlePost}
                        isPosting={isPosting}
                        isError={isError}
                        results={results}
                        errorData={errorData}
                        onRetry={onRetry}
                        onCloseError={onCloseError}
                    />
                </div>
            </div>

            {/* Preview Section */}
            <PreviewSection
                article={article}
                imageUrl={imageUrl}
                hasImage={!loadingImage}
                postMode={postMode}
                dailySchedule={dailySchedule}
                dailyTime={dailyTime}
                scheduleDate={scheduleDate}
                scheduleTime={scheduleTime}
                selectedProductsCount={selectedProducts ? 1 : 0}
                autoGenerateImage={autoGenerateImage}
            />

            {/* Publish Result Dialog */}
            <PublishResultDialog
                open={showResultDialog}
                onOpenChange={closeResultDialog}
                results={publishResults as any}
            />
        </div>
    );
}