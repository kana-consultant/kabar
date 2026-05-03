import { useFormState } from "./useFormState";
import { useProducts } from "./useProducts";
import { useModels } from "./useModels";
import { useArticleGeneration } from "./useArticleGeneration";
import { useDraftManagement } from "./useDraftManagement";
import { usePublishing } from "./usePublishing";
import { createDraft as createDraftService, publishDraft as publishDraftService } from "@/services/draft";
import { generateImage as generateImageService } from "@/services/generateService";
import type { UseGenerateReturn } from "./types";
import { toast } from "sonner";
import { useEffect } from "react";

export function useGenerate(): UseGenerateReturn {
    // Form State
    const formState = useFormState();
    
    // Products
    const productsState = useProducts();
    
    // Models
    const modelsState = useModels();
    
    // Draft Management
    const draftManagement = useDraftManagement({
        setTopic: formState.setTopic,
        setArticle: formState.setArticle,
        setImageUrl: formState.setImageUrl,
        setSelectedProducts: formState.setSelectedProducts,
        setCurrentDraftId: formState.setCurrentDraftId,
    });
    
    // Article Generation
    const articleGeneration = useArticleGeneration({
        topic: formState.topic,
        selectedModelId: modelsState.selectedModelId,
        tone: formState.tone,
        articleLength: formState.articleLength,
        keywords: formState.keywords,
        language: formState.language,
        autoGenerateImage: formState.autoGenerateImage,
        setArticle: formState.setArticle,
        setArticleResponse: formState.setArticleResponse,
        setImageUrl: formState.setImageUrl,
        setSeoScore: formState.setSeoScore,
        setReadabilityScore: formState.setReadabilityScore,
        setWordCount: formState.setWordCount,
        setLoadingArticle: formState.setLoadingArticle,
        setLoadingImage: formState.setLoadingImage,
    });
    
    // Publishing
    const publishing = usePublishing({
        setShowResultDialog: formState.setShowResultDialog,
        setPublishResults: formState.setPublishResults,
        resetForm: formState.resetForm,
        setPublishing: formState.setPublishing,
    });
    
    // Handlers
    const handleProductToggle = (product: string) => {
        if (formState.postToAll) return;
        formState.setSelectedProducts(prev =>
            prev.includes(product) ? prev.filter(p => p !== product) : [...prev, product]
        );
    };
    
    const handleSelectAll = () => {
        if (formState.postToAll) {
            formState.setPostToAll(false);
            formState.setSelectedProducts([]);
            toast.info("Semua produk dibatalkan");
        } else {
            formState.setPostToAll(true);
            formState.setSelectedProducts([...productsState.productNames]);
            toast.success("Semua produk dipilih");
        }
    };
    
    const handleAddKeyword = () => {
        if (formState.keywordInput.trim() && !formState.keywords.includes(formState.keywordInput.trim())) {
            formState.setKeywords([...formState.keywords, formState.keywordInput.trim()]);
            formState.setKeywordInput("");
        }
    };
    
    const handleRemoveKeyword = (keyword: string) => {
        formState.setKeywords(formState.keywords.filter(k => k !== keyword));
    };
    
    const handlePost = async () => {
        try {
            if (formState.selectedProducts.length === 0) {
                toast.error("Pilih minimal 1 produk terlebih dahulu");
                return;
            }
            if (!formState.article) {
                toast.error("Generate artikel terlebih dahulu");
                return;
            }
            
            formState.setIsPosting(true);
            
            if (formState.postMode === "draft") {
                await draftManagement.saveAsDraft(
                    formState.topic,
                    formState.article,
                    formState.imageUrl,
                    formState.selectedProducts,
                    formState.currentDraftId,
                    formState.setCurrentDraftId
                );
            } else if (formState.postMode === "scheduled") {
                await publishing.saveAsSchedule(
                    formState.topic,
                    formState.article,
                    formState.imageUrl,
                    formState.selectedProducts,
                    formState.currentDraftId,
                    formState.scheduleDate,
                    formState.scheduleTime,
                    formState.dailySchedule,
                    formState.dailyTime
                );
            } else {
                await publishing.postInstant(
                    formState.topic,
                    formState.article,
                    formState.imageUrl,
                    formState.selectedProducts,
                    formState.currentDraftId,
                    formState.resetForm
                );
            }
        } finally {
            formState.setIsPosting(false);
        }
    };
    
    const closeResultDialog = () => {
        formState.setShowResultDialog(false);
        formState.setPublishResults(null);
    };
    
    const quickGenerate = () => {
        return publishing.quickGenerate(
            formState.topic,
            formState.selectedProducts,
            modelsState.getTextAndImageModels,
            async (modelId: string) => {
                // Wrapper for generate article
                const response = await articleGeneration.generateArticleContent();
                return response;
            },
            async (prompt: string, modelId: string) => {
                const response = await generateImageService({ prompt, modelId });
                return response;
            },
            createDraftService,
            publishDraftService
        );
    };
    
    // Effect for initial product selection
    useEffect(() => {
        if (productsState.products.length > 0 && !formState.selectedProducts.length) {
            formState.setSelectedProducts(productsState.products[0]?.id as any);
        }
    }, [productsState.products]);
    
    return {
        // States
        topic: formState.topic,
        setTopic: formState.setTopic,
        article: formState.article,
        setArticle: formState.setArticle,
        articleResponse: formState.articleResponse,
        imageUrl: formState.imageUrl,
        setImageUrl: formState.setImageUrl,
        loadingArticle: formState.loadingArticle,
        loadingImage: formState.loadingImage,
        publishing: formState.publishing,
        selectedProducts: formState.selectedProducts,
        setSelectedProducts: formState.setSelectedProducts,
        postMode: formState.postMode,
        setPostMode: formState.setPostMode,
        scheduleTime: formState.scheduleTime,
        setScheduleTime: formState.setScheduleTime,
        scheduleDate: formState.scheduleDate,
        setScheduleDate: formState.setScheduleDate,
        dailySchedule: formState.dailySchedule,
        setDailySchedule: formState.setDailySchedule,
        dailyTime: formState.dailyTime,
        setDailyTime: formState.setDailyTime,
        autoGenerateImage: formState.autoGenerateImage,
        setAutoGenerateImage: formState.setAutoGenerateImage,
        postToAll: formState.postToAll,
        setPostToAll: formState.setPostToAll,
        currentDraftId: formState.currentDraftId,
        models: modelsState.models,
        selectedModelId: modelsState.selectedModelId,
        setSelectedModelId: modelsState.setSelectedModelId,
        loadingModels: modelsState.loadingModels,
        products: productsState.products,
        productNames: productsState.productNames,
        productsLoading: productsState.productsLoading,
        productsError: productsState.productsError,
        tone: formState.tone,
        setTone: formState.setTone,
        articleLength: formState.articleLength,
        setArticleLength: formState.setArticleLength,
        keywords: formState.keywords,
        keywordInput: formState.keywordInput,
        setKeywordInput: formState.setKeywordInput,
        language: formState.language,
        setLanguage: formState.setLanguage,
        seoScore: formState.seoScore,
        readabilityScore: formState.readabilityScore,
        wordCount: formState.wordCount,
        showResultDialog: formState.showResultDialog,
        publishResults: formState.publishResults,
        isPosting: formState.isPosting,
        setIsPosting: formState.setIsPosting,
        
        // Actions
        generateArticle: articleGeneration.generateArticleContent,
        generateImage: () => articleGeneration.generateImageManually(formState.articleResponse, formState.topic),
        handleProductToggle,
        handleSelectAll,
        handlePost,
        saveAsDraft: () => draftManagement.saveAsDraft(
            formState.topic,
            formState.article,
            formState.imageUrl,
            formState.selectedProducts,
            formState.currentDraftId,
            formState.setCurrentDraftId
        ),
        resetForm: formState.resetForm,
        handleAddKeyword,
        handleRemoveKeyword,
        closeResultDialog,
        quickGenerate,
    };
}