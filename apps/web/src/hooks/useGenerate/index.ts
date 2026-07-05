// hooks/useGenerate.ts
import { useLocation } from "@tanstack/react-router";
import { useToast } from "@/hooks/use-toast";
import { useGenerateState } from "./useGenerateState";
import { useGenerateData } from "./useGenerateData";
import { useLoadDraft } from "./useLoadDraft";
import { generateArticleContent } from "./useGenerateArticle";
import { generateImageManually } from "./useGenerateImage";
import { saveAsDraft, saveAsSchedule, postInstant } from "./usePublishActions";
import { handleAddKeyword, handleRemoveKeyword, handleProductToggle, handleSelectAll, resetForm } from "./useFormManagement";
import { quickGenerate } from "./useQuickGenerate";
import { closeResultDialog } from "./useDialogState";
import { useState } from "react";

export function useGenerate() {
    const toast = useToast();
    const location = useLocation();
    const searchParams = new URLSearchParams(location.search);
    const editId = searchParams.get("edit") || undefined;
    const TopicId = searchParams.get("topic") || undefined;

    // State untuk error handling di PostingConfig
    const [isError, setIsError] = useState(false);
    const [results, setResults] = useState<any[]>([]);
    const [errorData, setErrorData] = useState<any>(null);

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
        slug, setSlug,
        excerpt, setExcerpt
    } = useGenerateState();

    useGenerateData(
        setProducts, setProductNames, setProductsLoading, setProductsError,
        setModels, setLoadingModels, setSelectedModelId, toast
    );

    useLoadDraft(
        editId, TopicId,
        setTopic, setArticle, setImageUrl, setSelectedProducts, setCurrentDraftId,
        setSlug, setKeywords,setExcerpt,toast
    );

    const generateArticle = (autoGenerateImage?: boolean, imageModelId?: string) => generateArticleContent(
        topic, selectedModelId, tone, articleLength, language,
        setLoadingArticle, setArticleResponse, setArticle,
        setSeoScore, setReadabilityScore, setWordCount, setSlug, setKeywords, setExcerpt, toast,
        autoGenerateImage, imageModelId  
    );

    const generateImage = () => generateImageManually(
        articleResponse, topic, selectedModelId, setLoadingImage, setImageUrl,
        toast
    );

    // Reset error state
    const resetErrorState = () => {
        setIsError(false);
        setResults([]);
        setErrorData(null);
    };



    // Handle reset form setelah publish
    const handleResetAfterPublish = () => {
        resetForm(
            setTopic, setArticle, setArticleResponse, setImageUrl, setSelectedProducts,
            setPostMode, setCurrentDraftId, setKeywords, setKeywordInput,
            setSeoScore, setReadabilityScore, setWordCount, setTone, setArticleLength, setLanguage,
            toast
        );
    };

    const handleSaveAsDraft = () => saveAsDraft(
        article, topic, imageUrl, selectedProducts, currentDraftId, setCurrentDraftId,
        slug as string, keywords as string[], excerpt as string, toast
    );

    const handleSaveAsSchedule = () => {
        resetErrorState();
        return saveAsSchedule(
            article, topic, imageUrl, selectedProducts, currentDraftId,
            dailySchedule, dailyTime, scheduleDate, scheduleTime,
            setPublishing, setPublishResults, setShowResultDialog,
            (result: any) => {
                // onSuccess callback
                if (result) {
                    console.log("Hasil")
                    console.log(result.data)
                    setResults(result.results || []);
                    setIsError(true);
                  
                }
            },
            () => {
                // onReset callback
                resetForm(
                    setTopic, setArticle, setArticleResponse, setImageUrl, setSelectedProducts,
                    setPostMode, setCurrentDraftId, setKeywords, setKeywordInput,
                    setSeoScore, setReadabilityScore, setWordCount, setTone, setArticleLength, setLanguage,
                    toast
                );
            },
            slug as string, keywords as string[],
            toast
        );
    };

    const handlePostInstant = () => {
        resetErrorState();
        return postInstant(
            article, topic, imageUrl, selectedProducts, currentDraftId,
            setPublishing, setPublishResults, setShowResultDialog,
            (result: any) => {
                // onSuccess callback
                console.log("Hasil")
                console.log(result.data)
            },
            () => {
                // onReset callback
                resetForm(
                    setTopic, setArticle, setArticleResponse, setImageUrl, setSelectedProducts,
                    setPostMode, setCurrentDraftId, setKeywords, setKeywordInput,
                    setSeoScore, setReadabilityScore, setWordCount, setTone, setArticleLength, setLanguage,
                    toast
                );
            },
            slug as string, keywords as string[],
            excerpt as string,
            toast
        );
    };

    const handlePost = async () => {
        try {
            console.log(selectedProducts)
            if (selectedProducts.length === 0) {
                toast.error("Pilih minimal 1 produk terlebih dahulu");
                return;
            }
            if (!article) {
                toast.error("Generate artikel terlebih dahulu");
                return;
            }

            setIsPosting(true);
            resetErrorState();

            let result: any = null;

            if (postMode === "draft") {
                await handleSaveAsDraft();
                toast.success('Draft berhasil disimpan');
            } else if (postMode === "scheduled") {
                result = await handleSaveAsSchedule();
            } else {
                result = await handlePostInstant();
            }
        } catch (error: any) {
            console.log("Error saat posting:", error.response?.data || error.message || error);
            setIsError(true);
            setErrorData(error.response?.data.results);
            toast.error("Terjadi kesalahan saat memposting. Silakan coba lagi.");
        } finally {
            setIsPosting(false);
        }
    };

    const onResetForm = () => {
        resetForm(
            setTopic, setArticle, setArticleResponse, setImageUrl, setSelectedProducts,
            setPostMode, setCurrentDraftId, setKeywords, setKeywordInput,
            setSeoScore, setReadabilityScore, setWordCount, setTone, setArticleLength, setLanguage,
            toast
        );
        resetErrorState();
    };

    const onQuickGenerate = () => quickGenerate(
        topic, selectedProducts, article, imageUrl, setPublishing, setCurrentDraftId, onResetForm, excerpt as string,
        toast
    );

    // Handle retry
    const handleRetry = () => {
        resetErrorState();
        handlePost();
    };

    // Handle close error
    const handleCloseError = () => {
        resetErrorState();
        closeResultDialog(setShowResultDialog, setPublishResults);
    };

   

    return {
        topic, setTopic,
        article, setArticle,
        excerpt,slug,
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
        handleSelectAll: () => handleSelectAll(postToAll, productNames, setPostToAll, setSelectedProducts, toast),
        handlePost,
        saveAsDraft: handleSaveAsDraft,
        resetForm: onResetForm,
        handleAddKeyword: () => handleAddKeyword(keywordInput, keywords, setKeywords, setKeywordInput),
        handleRemoveKeyword: (keyword: string) => handleRemoveKeyword(keyword, keywords, setKeywords),
        isPosting, setIsPosting,
        quickGenerate: onQuickGenerate,
        // Props untuk PostingConfig error handling
        isError,
        results,
        errorData,
        onRetry: handleRetry,
        onCloseError: handleCloseError,
    };
}