import type { APIKey } from "@/services/apiKey";
import type {  PublishResponse } from "@/services/draft";
import type { Product } from "@/types/product";

export type Tone = "professional" | "casual" | "friendly" | "formal";
export type ArticleLength = "short" | "medium" | "long";
export type Language = "id" | "en";
export type PostMode = "instant" | "scheduled" | "draft";

export interface GenerateArticleResponse {
    title: string;
    content: string;
    excerpt: string;
    keywords: string[];
    imagePrompt: string;
    imageUrl: string;
    wordCount: number;
    readabilityScore: number;
    seoScore: number;
}

export interface UseGenerateState {
    topic: string;
    article: string;
    articleResponse: GenerateArticleResponse | null;
    imageUrl: string;
    loadingArticle: boolean;
    loadingImage: boolean;
    publishing: boolean;
    selectedProducts: string[];
    postMode: PostMode;
    scheduleTime: string;
    scheduleDate: string;
    dailySchedule: boolean;
    dailyTime: string;
    autoGenerateImage: boolean;
    postToAll: boolean;
    currentDraftId: string | null;
    selectedModelId: string;
    tone: Tone;
    articleLength: ArticleLength;
    keywords: string[];
    keywordInput: string;
    language: Language;
    seoScore: number | null;
    readabilityScore: number | null;
    wordCount: number | null;
    isPosting: boolean;
}

export interface UseGenerateReturn extends UseGenerateState {
    setTopic: (topic: string) => void;
    setArticle: (article: string) => void;
    setImageUrl: (url: string) => void;
    setSelectedProducts: (products: string[]) => void;
    setPostMode: (mode: PostMode) => void;
    setScheduleTime: (time: string) => void;
    setScheduleDate: (date: string) => void;
    setDailySchedule: (enabled: boolean) => void;
    setDailyTime: (time: string) => void;
    setAutoGenerateImage: (enabled: boolean) => void;
    setPostToAll: (enabled: boolean) => void;
    setSelectedModelId: (id: string) => void;
    setTone: (tone: Tone) => void;
    setArticleLength: (length: ArticleLength) => void;
    setKeywordInput: (input: string) => void;
    setLanguage: (lang: Language) => void;
    setIsPosting: (posting: boolean) => void;
    
    // Computed states
    models: APIKey[];
    loadingModels: boolean;
    products: Product[];
    productNames: string[];
    productsLoading: boolean;
    productsError: string | null;
    showResultDialog: boolean;
    publishResults: PublishResponse | null;
    
    // Actions
    generateArticle: () => Promise<void>;
    generateImage: () => Promise<void>;
    handleProductToggle: (product: string) => void;
    handleSelectAll: () => void;
    handlePost: () => Promise<void>;
    saveAsDraft: () => Promise<boolean>;
    resetForm: () => void;
    handleAddKeyword: () => void;
    handleRemoveKeyword: (keyword: string) => void;
    closeResultDialog: () => void;
    quickGenerate: () => Promise<string | undefined>;
}