import { GenerateHeader } from "./GenerateHeader";
import { TopicInput } from "./TopicInput";
import { PostingConfig } from "./PostingConfig";
import { PreviewSection } from "./PreviewSection";
import { PublishResultDialog } from "../draft/PublishResultDialog";
import { useGenerate } from "@/hooks/useGenerate";
import { ModelSelector } from "./ModelSelector";

export default function Generate() {
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

    console.log(postMode);
    
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
                    isError={isError}
                    results={results}
                    errorData={errorData}
                    onRetry={onRetry}
                    onCloseError={onCloseError}
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