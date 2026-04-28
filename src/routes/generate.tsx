// src/routes/generate.tsx
import { createFileRoute } from "@tanstack/react-router";
import { GenerateHeader } from "@/components/generate/GenerateHeader";
import { TopicInput } from "@/components/generate/TopicInput";
import { PostingConfig } from "@/components/generate/PostingConfig";
import { PreviewSection } from "@/components/generate/PreviewSection";
import { PublishResultDialog } from "@/components/draft/PublishResultDialog";
import { useGenerate } from "@/hooks/useGenerate";
import { ModelSelector } from '@/components/generate/ModelSelector';

export const Route = createFileRoute("/generate")({
    component: Generate,
});

export function Generate() {
    const {
        topic, setTopic,
        article,
        imageUrl,
        loadingArticle, loadingImage,
        selectedProducts,
        postMode, setPostMode,
        scheduleTime, setScheduleTime,
        scheduleDate, setScheduleDate,
        dailySchedule, setDailySchedule,
        dailyTime, setDailyTime,
        autoGenerateImage, setAutoGenerateImage,
        postToAll,
        products,
        productsLoading,
        productsError,
        currentDraftId,
        generateArticle,
        generateImage,
        handleProductToggle,
        handleSelectAll,
        handlePost,
        seoScore,
        readabilityScore,
        wordCount,
        publishResults,
        showResultDialog,
        closeResultDialog,
        isPosting,
        selectedModelId, 
        setSelectedModelId
    } = useGenerate();

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
    console.log(postMode)
    return (
        <div className="space-y-6">
            {/* Draft indicator */}
            {currentDraftId && (
                <div className="rounded-lg bg-cyan-50 p-3 text-sm text-cyan-700 dark:bg-cyan-950/50 dark:text-cyan-300 border border-cyan-200 dark:border-cyan-800">
                    Sedang mengedit draft. Simpan akan memperbarui draft yang ada.
                </div>
            )}

            <GenerateHeader />

            {/* SEO Metrics Display */}
            {(seoScore || readabilityScore || wordCount) && (
                <div className="rounded-lg bg-zinc-50 dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 p-4">
                    <h3 className="mb-2 font-semibold text-zinc-900 dark:text-zinc-100">SEO Metrics</h3>
                    <div className="flex flex-wrap gap-4 text-sm">
                        {seoScore && (
                            <div className="flex items-center gap-1 text-zinc-600 dark:text-zinc-400">
                                <span>📈 SEO Score:</span>
                                <span className="font-medium text-cyan-600 dark:text-cyan-400">{seoScore}</span>
                            </div>
                        )}
                        {readabilityScore && (
                            <div className="flex items-center gap-1 text-zinc-600 dark:text-zinc-400">
                                <span>📖 Readability:</span>
                                <span className="font-medium text-emerald-600 dark:text-emerald-400">{readabilityScore}</span>
                            </div>
                        )}
                        {wordCount && (
                            <div className="flex items-center gap-1 text-zinc-600 dark:text-zinc-400">
                                <span>📝 Word Count:</span>
                                <span className="font-medium text-purple-600 dark:text-purple-400">{wordCount}</span>
                            </div>
                        )}
                    </div>
                </div>
            )}

            {/* Main content grid */}
            <div className="grid gap-6 lg:grid-cols-2">
                {/* Left column */}
                <div className="space-y-6">
                    <TopicInput
                        topic={topic}
                        setTopic={setTopic}
                        loadingArticle={loadingArticle}
                        loadingImage={loadingImage}
                        onGenerateArticle={generateArticle}
                        onGenerateImage={generateImage}
                        autoGenerateImage={autoGenerateImage}
                        setAutoGenerateImage={setAutoGenerateImage}
                        article={article}
                    />

                    <ModelSelector
                        selectedModelId={selectedModelId}
                        onModelChange={setSelectedModelId}
                    />
                </div>

                {/* Right column */}
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
                    setAutoGenerateImage={setAutoGenerateImage}
                    products={products}
                    selectedProducts={selectedProducts}
                    postToAll={postToAll}
                    onToggleProduct={handleProductToggle}
                    onSelectAll={handleSelectAll}
                    article={article}
                    onPost={handlePost}
                    isPosting={isPosting}
                />
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
                selectedProductsCount={selectedProducts.length}
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