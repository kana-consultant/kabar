import { useState } from "react";
import type { Tone, ArticleLength, Language, PostMode } from "./types";

export function useFormState() {
    const [topic, setTopic] = useState("");
    const [article, setArticle] = useState("");
    const [imageUrl, setImageUrl] = useState("");
    const [selectedProducts, setSelectedProducts] = useState<string[]>([]);
    const [postMode, setPostMode] = useState<PostMode>("instant");
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
    const [tone, setTone] = useState<Tone>("professional");
    const [articleLength, setArticleLength] = useState<ArticleLength>("medium");
    const [keywords, setKeywords] = useState<string[]>([]);
    const [keywordInput, setKeywordInput] = useState("");
    const [language, setLanguage] = useState<Language>("id");
    const [isPosting, setIsPosting] = useState(false);
    
    // Result states
    const [showResultDialog, setShowResultDialog] = useState(false);
    const [publishResults, setPublishResults] = useState<any>(null);
    
    // Loading states
    const [loadingArticle, setLoadingArticle] = useState(false);
    const [loadingImage, setLoadingImage] = useState(false);
    const [publishing, setPublishing] = useState(false);
    
    // SEO states
    const [seoScore, setSeoScore] = useState<number | null>(null);
    const [readabilityScore, setReadabilityScore] = useState<number | null>(null);
    const [wordCount, setWordCount] = useState<number | null>(null);
    
    const [articleResponse, setArticleResponse] = useState<any>(null);
    const [selectedModelId, setSelectedModelId] = useState("");

    return {
        // Topic & Content
        topic, setTopic,
        article, setArticle,
        articleResponse, setArticleResponse,
        imageUrl, setImageUrl,
        
        // Selection
        selectedProducts, setSelectedProducts,
        postMode, setPostMode,
        
        // Schedule
        scheduleTime, setScheduleTime,
        scheduleDate, setScheduleDate,
        dailySchedule, setDailySchedule,
        dailyTime, setDailyTime,
        
        // Options
        autoGenerateImage, setAutoGenerateImage,
        postToAll, setPostToAll,
        currentDraftId, setCurrentDraftId,
        
        // Article config
        tone, setTone,
        articleLength, setArticleLength,
        keywords, setKeywords,
        keywordInput, setKeywordInput,
        language, setLanguage,
        
        // Loading states
        loadingArticle, setLoadingArticle,
        loadingImage, setLoadingImage,
        publishing, setPublishing,
        isPosting, setIsPosting,
        
        // SEO
        seoScore, setSeoScore,
        readabilityScore, setReadabilityScore,
        wordCount, setWordCount,
        
        // Dialog
        showResultDialog, setShowResultDialog,
        publishResults, setPublishResults,
        
        // Model
        selectedModelId, setSelectedModelId,
        
        // Actions
        resetForm: () => {
            setTopic("");
            setArticle("");
            setArticleResponse(null);
            setImageUrl("");
            setSelectedProducts([]);
            setPostMode("instant");
            setCurrentDraftId(null);
            setKeywords([]);
            setKeywordInput("");
            setSeoScore(null);
            setReadabilityScore(null);
            setWordCount(null);
            setTone("professional");
            setArticleLength("medium");
            setLanguage("id");
        }
    };
}