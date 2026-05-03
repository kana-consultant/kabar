import { useLocation } from "@tanstack/react-router";
import { useGenerateState } from "./useGenerateState";
import { useGenerateData } from "./useGenerateData";
import { useLoadDraft } from "./useLoadDraft";
import { generateArticleContent } from "./useGenerateArticle";
import { generateImageManually } from "./useGenerateImage";
import { saveAsDraft, saveAsSchedule, postInstant } from "./usePublishActions";
import { handleAddKeyword, handleRemoveKeyword, handleProductToggle, handleSelectAll, resetForm } from "./useFormManagement";
import { quickGenerate } from "./useQuickGenerate";
import { closeResultDialog } from "./useDialogState";
import { toast } from "sonner"; 

export function useGenerate() {
    const location = useLocation();
    const searchParams = new URLSearchParams(location.search);
    const editId = searchParams.get("edit") || undefined;
    const TopicId = searchParams.get("topic") || undefined;

    const {
        topic, setTopic,
        article, setArticle,
        articleResponse, setArticleResponse,
        imageUrl, setImageUrl,
        loadingArticle, setLoadingArticle,
        loadingImage, setLoadingImage,
        selectedProducts, setSelectedProducts,
        postMode, setPostMode,
        scheduleTime, setScheduleTime,
        scheduleDate, setScheduleDate,
        dailySchedule, setDailySchedule,
        dailyTime, setDailyTime,
        autoGenerateImage, setAutoGenerateImage,
        postToAll, setPostToAll,
        currentDraftId, setCurrentDraftId,
        publishing, setPublishing,
        showResultDialog, setShowResultDialog,
        publishResults, setPublishResults,
        models, setModels,
        selectedModelId, setSelectedModelId,
        loadingModels, setLoadingModels,
        tone, setTone,
        articleLength, setArticleLength,
        keywords, setKeywords,
        keywordInput, setKeywordInput,
        language, setLanguage,
        seoScore, setSeoScore,
        readabilityScore, setReadabilityScore,
        wordCount, setWordCount,
        products, setProducts,
        productNames, setProductNames,
        productsLoading, setProductsLoading,
        productsError, setProductsError,
        isPosting, setIsPosting,
    } = useGenerateState();

    useGenerateData(
        setProducts, setProductNames, setProductsLoading, setProductsError,
        setModels, setLoadingModels, setSelectedModelId
    );

    useLoadDraft(
        editId, TopicId,
        setTopic, setArticle, setImageUrl, setSelectedProducts, setCurrentDraftId
    );

    const generateArticle = () => generateArticleContent(
        topic, selectedModelId, tone, articleLength , language, 
        setLoadingArticle, setArticleResponse, setArticle,
        setSeoScore, setReadabilityScore, setWordCount
    );

    const generateImage = () => generateImageManually(
        articleResponse, topic, selectedModelId, setLoadingImage, setImageUrl
    );

    const handleSaveAsDraft = () => saveAsDraft(
        article, topic, imageUrl, selectedProducts, currentDraftId, setCurrentDraftId
    );

    const handleSaveAsSchedule = () => saveAsSchedule(
        article, topic, imageUrl, selectedProducts, currentDraftId,
        dailySchedule, dailyTime, scheduleDate, scheduleTime,
        setPublishing, setPublishResults, setShowResultDialog,
        () => resetForm(
            setTopic, setArticle, setArticleResponse, setImageUrl, setSelectedProducts,
            setPostMode, setCurrentDraftId, setKeywords, setKeywordInput,
            setSeoScore, setReadabilityScore, setWordCount, setTone, setArticleLength, setLanguage
        )
    );

    const handlePostInstant = () => postInstant(
        article, topic, imageUrl, selectedProducts, currentDraftId,
        setPublishing, setPublishResults, setShowResultDialog,
        () => resetForm(
            setTopic, setArticle, setArticleResponse, setImageUrl, setSelectedProducts,
            setPostMode, setCurrentDraftId, setKeywords, setKeywordInput,
            setSeoScore, setReadabilityScore, setWordCount, setTone, setArticleLength, setLanguage
        )
    );

    const handlePost = async () => {
        try {
            if (selectedProducts.length === 0) {
                toast.error("Pilih minimal 1 produk terlebih dahulu");
                return;
            }
            if (!article) {
                toast.error("Generate artikel terlebih dahulu");
                return;
            }

            setIsPosting(true);
            if (postMode === "draft") {
                await handleSaveAsDraft();
            } else if (postMode === "scheduled") {
                await handleSaveAsSchedule();
            } else {
                await handlePostInstant();
            }
        } finally {
            setIsPosting(false);
        }
    };

    const onResetForm = () => resetForm(
        setTopic, setArticle, setArticleResponse, setImageUrl, setSelectedProducts,
        setPostMode, setCurrentDraftId, setKeywords, setKeywordInput,
        setSeoScore, setReadabilityScore, setWordCount, setTone, setArticleLength, setLanguage
    );

    const onQuickGenerate = () => quickGenerate(
        topic, selectedProducts, article, imageUrl, setPublishing, setCurrentDraftId, onResetForm
    );

    return {
        topic, setTopic,
        article, setArticle,
        articleResponse,
        imageUrl, setImageUrl,
        loadingArticle, loadingImage,
        publishing,
        selectedProducts, setSelectedProducts,
        postMode, setPostMode,
        scheduleTime, setScheduleTime,
        scheduleDate, setScheduleDate,
        dailySchedule, setDailySchedule,
        dailyTime, setDailyTime,
        autoGenerateImage, setAutoGenerateImage,
        postToAll, setPostToAll,
        currentDraftId,
        models,
        selectedModelId, setSelectedModelId,
        loadingModels,
        products, productNames, productsLoading, productsError,
        tone, setTone,
        articleLength, setArticleLength,
        keywords, keywordInput, setKeywordInput,
        language, setLanguage,
        seoScore, readabilityScore, wordCount,
        showResultDialog, publishResults,
        closeResultDialog: () => closeResultDialog(setShowResultDialog, setPublishResults),
        generateArticle,
        generateImage,
        handleProductToggle: (product: string) => handleProductToggle(product, selectedProducts, postToAll, setSelectedProducts),
        handleSelectAll: () => handleSelectAll(postToAll, productNames, setPostToAll, setSelectedProducts),
        handlePost,
        saveAsDraft: handleSaveAsDraft,
        resetForm: onResetForm,
        handleAddKeyword: () => handleAddKeyword(keywordInput, keywords, setKeywords, setKeywordInput),
        handleRemoveKeyword: (keyword: string) => handleRemoveKeyword(keyword, keywords, setKeywords),
        isPosting, setIsPosting,
        quickGenerate: onQuickGenerate,
    };
}