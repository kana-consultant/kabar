import { useState } from "react";

export function useGenerateState() {
    const [topic, setTopic] = useState("");
    const [article, setArticle] = useState("");
    const [articleResponse, setArticleResponse] = useState<any>(null);
    const [imageUrl, setImageUrl] = useState("");
    const [loadingArticle, setLoadingArticle] = useState(false);
    const [loadingImage, setLoadingImage] = useState(false);
    const [selectedProducts, setSelectedProducts] = useState<string[]>([]);
    const [postMode, setPostMode] = useState<"instant" | "scheduled" | "draft">("instant");
    const [scheduleTime, setScheduleTime] = useState("09:00");
    const [scheduleDate, setScheduleDate] = useState(() => {
        const tomorrow = new Date();
        tomorrow.setDate(tomorrow.getDate() + 1);
        return tomorrow.toISOString().split("T")[0];
    });
    const [dailySchedule, setDailySchedule] = useState(false);
    const [dailyTime, setDailyTime] = useState("08:00");
    const [autoGenerateImage, setAutoGenerateImage] = useState(false);
    const [postToAll, setPostToAll] = useState(false);
    const [currentDraftId, setCurrentDraftId] = useState<string | null>(null);
    const [publishing, setPublishing] = useState(false);
    const [showResultDialog, setShowResultDialog] = useState(false);
    const [publishResults, setPublishResults] = useState<any>(null);
    const [models, setModels] = useState<any[]>([]);
    const [selectedModelId, setSelectedModelId] = useState("");
    const [loadingModels, setLoadingModels] = useState(true);
    const [tone, setTone] = useState<"professional" | "casual" | "friendly" | "formal">("professional");
    const [articleLength, setArticleLength] = useState<"short" | "medium" | "long">("medium");
    const [keywords, setKeywords] = useState<string[]>([]);
    const [keywordInput, setKeywordInput] = useState("");
    const [language, setLanguage] = useState<"id" | "en">("id");
    const [seoScore, setSeoScore] = useState<number | null>(null);
    const [readabilityScore, setReadabilityScore] = useState<number | null>(null);
    const [wordCount, setWordCount] = useState<number | null>(null);
    const [products, setProducts] = useState<any[]>([]);
    const [productNames, setProductNames] = useState<string[]>([]);
    const [productsLoading, setProductsLoading] = useState(true);
    const [productsError, setProductsError] = useState<string | null>(null);
    const [isPosting, setIsPosting] = useState(false);

    return {
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
    };
}